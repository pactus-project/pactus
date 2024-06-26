import argparse
import os
import shutil
import subprocess
import hashlib
import json
import logging
import zipfile
from datetime import datetime


def setup_logging():
    logging.basicConfig(
        format='[%(asctime)s] %(message)s',
        datefmt='%Y-%m-%d-%H:%M',
        level=logging.INFO
    )


def get_timestamp_str():
    return datetime.now().strftime("%Y%m%d%H%M%S")


def get_current_time_iso():
    return datetime.now().isoformat()


class Metadata:
    @staticmethod
    def sha256(file_path):
        hash_sha = hashlib.sha256()
        with open(file_path, "rb") as f:
            for chunk in iter(lambda: f.read(4096), b""):
                hash_sha.update(chunk)
        return hash_sha.hexdigest()

    @staticmethod
    def update_metadata_file(snapshot_path, snapshot_metadata):
        metadata_file = os.path.join(snapshot_path, 'snapshots', 'metadata.json')
        if os.path.isfile(metadata_file):
            logging.info(f"Updating existing metadata file '{metadata_file}'")
            with open(metadata_file, 'r') as f:
                metadata = json.load(f)
        else:
            logging.info(f"Creating new metadata file '{metadata_file}'")
            metadata = []

        formatted_metadata = {
            "name": snapshot_metadata["name"],
            "created_at": snapshot_metadata["created_at"],
            "compress": snapshot_metadata["compress"],
            "data": snapshot_metadata["data"]
        }

        metadata.append(formatted_metadata)

        with open(metadata_file, 'w') as f:
            json.dump(metadata, f, indent=4)

    @staticmethod
    def update_metadata_after_removal(snapshots_dir, removed_snapshots):
        metadata_file = os.path.join(snapshots_dir, 'metadata.json')
        if not os.path.isfile(metadata_file):
            return

        logging.info(f"Updating metadata file '{metadata_file}' after snapshot removal")
        with open(metadata_file, 'r') as f:
            metadata = json.load(f)

        updated_metadata = [entry for entry in metadata if entry["name"] not in removed_snapshots]

        with open(metadata_file, 'w') as f:
            json.dump(updated_metadata, f, indent=4)

    @staticmethod
    def create_snapshot_json(data_dir, snapshot_subdir):
        files = []
        for root, _, filenames in os.walk(data_dir):
            for filename in filenames:
                file_path = os.path.join(root, filename)
                rel_path = os.path.relpath(file_path, data_dir)
                snapshot_rel_path = os.path.join(snapshot_subdir, rel_path).replace('\\', '/')
                file_info = {
                    "name": filename,
                    "path": snapshot_rel_path,
                    "sha": Metadata.sha256(file_path)
                }
                files.append(file_info)

        return {"data": files}

    @staticmethod
    def create_compressed_snapshot_json(compressed_file, rel_path):
        file_info = {
            "name": os.path.basename(compressed_file),
            "path": rel_path,
            "sha": Metadata.sha256(compressed_file)
        }

        return {"data": file_info}


class DaemonManager:
    @staticmethod
    def is_daemon_running(daemon_path):
        process_name = os.path.basename(daemon_path)
        if os.name == 'nt':
            tasklist = subprocess.run(['tasklist'], capture_output=True, text=True).stdout
            return process_name in tasklist
        else:
            pgrep_command = ['pgrep', '-f', daemon_path]
            pgrep_result = subprocess.run(pgrep_command, capture_output=True, text=True)
            pgrep_output = pgrep_result.stdout.strip()
            logging.debug(f"Running command: {' '.join(pgrep_command)}")
            logging.debug(f"Pgrep output: '{pgrep_output}'")
            return bool(pgrep_output)

    @staticmethod
    def stop_daemon(daemon_path):
        if DaemonManager.is_daemon_running(daemon_path):
            logging.info(f"Daemon process found for '{daemon_path}', attempting to stop it")
            process_name = os.path.basename(daemon_path)
            try:
                if os.name == 'nt':
                    subprocess.run(['taskkill', '/F', '/IM', process_name], timeout=10)
                else:
                    pids = subprocess.run(['pgrep', '-f', daemon_path], capture_output=True,
                                          text=True).stdout.strip().split()
                    current_pid = os.getpid()
                    pids_to_kill = [pid for pid in pids if pid != str(current_pid)]
                    logging.debug(f"PIDs to kill: {pids_to_kill}")
                    for pid in pids_to_kill:
                        subprocess.run(['kill', pid], timeout=10)
                logging.info(f"Daemon process '{process_name}' stopped successfully")
            except subprocess.TimeoutExpired:
                logging.error(f"Failed to stop daemon process '{process_name}' within the timeout period")
        else:
            logging.info(f"No running daemon process found for '{daemon_path}'")

    @staticmethod
    def start_node(daemon_path):
        logging.info(f"Starting daemon with command '{daemon_path} start'")
        try:
            subprocess.run([daemon_path, 'start'], capture_output=True, text=True, check=True)
            logging.info(f"Daemon '{daemon_path}' started successfully")
        except subprocess.CalledProcessError as e:
            logging.error(f"Failed to start daemon '{daemon_path}': {e}")
            logging.error(f"stdout: {e.stdout}")
            logging.error(f"stderr: {e.stderr}")


class SnapshotManager:
    def __init__(self, args):
        self.args = args

    def manage_snapshots(self):
        snapshots_dir = os.path.join(self.args.snapshot_path, 'snapshots')
        logging.info(f"Managing snapshots in '{snapshots_dir}'")

        if not os.path.exists(snapshots_dir):
            logging.info(f"Snapshots directory '{snapshots_dir}' does not exist. Creating it.")
            os.makedirs(snapshots_dir)

        snapshots = sorted([s for s in os.listdir(snapshots_dir) if s != 'metadata.json'])

        logging.info(f"Found snapshots: {snapshots}")
        logging.info(f"Retention policy is to keep {self.args.retention} snapshots")

        if len(snapshots) >= self.args.retention:
            num_to_remove = len(snapshots) - self.args.retention + 1
            to_remove = snapshots[:num_to_remove]
            logging.info(f"Snapshots to remove: {to_remove}")
            for snapshot in to_remove:
                snapshot_path = os.path.join(snapshots_dir, snapshot)
                logging.info(f"Removing old snapshot '{snapshot_path}'")
                shutil.rmtree(snapshot_path)

            Metadata.update_metadata_after_removal(snapshots_dir, to_remove)

    def create_snapshot(self):
        timestamp_str = get_timestamp_str()
        snapshot_dir = os.path.join(self.args.snapshot_path, 'snapshots', timestamp_str)
        logging.info(f"Creating snapshot directory '{snapshot_dir}'")
        os.makedirs(snapshot_dir, exist_ok=True)

        data_dir = os.path.join(snapshot_dir, 'data')
        if self.args.compress == 'none':
            logging.info(f"Copying data from '{self.args.data_path}' to '{data_dir}'")
            shutil.copytree(self.args.data_path, data_dir)
            snapshot_metadata = Metadata.create_snapshot_json(data_dir, timestamp_str)
        elif self.args.compress == 'zip':
            zip_file = os.path.join(snapshot_dir, 'data.zip')
            rel = os.path.relpath(zip_file, snapshot_dir)
            meta_path = os.path.join(timestamp_str, rel)
            logging.info(f"Creating ZIP archive '{zip_file}'")
            with zipfile.ZipFile(zip_file, 'w', zipfile.ZIP_DEFLATED) as zipf:
                for root, _, files in os.walk(self.args.data_path):
                    for file in files:
                        full_path = os.path.join(root, file)
                        rel_path = os.path.relpath(full_path, self.args.data_path)
                        zipf.write(full_path, os.path.join('data', rel_path))
            snapshot_metadata = Metadata.create_compressed_snapshot_json(zip_file, meta_path)
        elif self.args.compress == 'tar':
            tar_file = os.path.join(snapshot_dir, 'data.tar.gz')
            rel = os.path.relpath(tar_file, snapshot_dir)
            meta_path = os.path.join(timestamp_str, rel)
            logging.info(f"Creating TAR.GZ archive '{tar_file}'")
            subprocess.run(['tar', '-czvf', tar_file, '-C', self.args.data_path, '.'])
            snapshot_metadata = Metadata.create_compressed_snapshot_json(tar_file, meta_path)

        snapshot_metadata["name"] = timestamp_str
        snapshot_metadata["created_at"] = get_current_time_iso()
        snapshot_metadata["compress"] = self.args.compress

        Metadata.update_metadata_file(self.args.snapshot_path, snapshot_metadata)


class Validation:
    @staticmethod
    def validate_args(args):
        logging.info('Validating arguments')

        if not os.path.isfile(args.daemon_path):
            raise ValueError(f"Daemon path '{args.daemon_path}' does not exist.")
        logging.info(f"Daemon path '{args.daemon_path}' exists")

        if not os.path.isdir(args.data_path):
            raise ValueError(f"Data path '{args.data_path}' does not exist.")
        logging.info(f"Data path '{args.data_path}' exists")

        if not os.access(args.data_path, os.W_OK):
            raise PermissionError(f"No permission to access data path '{args.data_path}'.")
        logging.info(f"Permission to access data path '{args.data_path}' confirmed")

        if args.compress == 'zip' and not shutil.which('zip'):
            raise EnvironmentError("The 'zip' command is not available.")
        elif args.compress == 'zip':
            logging.info("The 'zip' command is available")

        if args.compress == 'tar' and not shutil.which('tar'):
            raise EnvironmentError("The 'tar' command is not available.")
        elif args.compress == 'tar':
            logging.info("The 'tar' command is available")

        if args.retention <= 0:
            raise ValueError("Retention value must be greater than 0.")
        logging.info(f"Retention value is set to {args.retention}")

        if not os.access(args.snapshot_path, os.W_OK):
            raise PermissionError(f"No permission to access snapshot path '{args.snapshot_path}'.")
        logging.info(f"Permission to access snapshot path '{args.snapshot_path}' confirmed")

        snapshots_dir = os.path.join(args.snapshot_path, 'snapshots')
        if not os.path.isdir(snapshots_dir):
            logging.info("Snapshots directory does not exist, creating it")
            os.makedirs(snapshots_dir)
        else:
            logging.info("Snapshots directory exists")


class ProcessBackup:
    def __init__(self, args):
        self.args = args

    def run(self):
        Validation.validate_args(self.args)
        DaemonManager.stop_daemon(self.args.daemon_path)
        snapshot_manager = SnapshotManager(self.args)
        snapshot_manager.manage_snapshots()
        snapshot_manager.create_snapshot()
        DaemonManager.start_node(self.args.daemon_path)


def parse_args():
    user_home = os.path.expanduser("~")
    default_data_path = os.path.join(user_home, 'pactus')

    parser = argparse.ArgumentParser(description='Pactus Blockchain Backup Tool')
    parser.add_argument('--daemon_path', required=True, help='Path to daemon executable')
    parser.add_argument('--data_path', default=default_data_path, help='Path to data directory')
    parser.add_argument('--compress', choices=['none', 'zip', 'tar'], default='none', help='Compression type')
    parser.add_argument('--retention', type=int, default=3, help='Number of snapshots to retain')
    parser.add_argument('--snapshot_path', default=os.getcwd(), help='Path to store snapshots')

    return parser.parse_args()


def main():
    setup_logging()
    args = parse_args()
    process_backup = ProcessBackup(args)
    process_backup.run()


if __name__ == "__main__":
    main()
