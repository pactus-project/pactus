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

        metadata.append(snapshot_metadata)

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


def run_command(command):
    logging.info(f"Running command: {' '.join(command)}")
    try:
        result = subprocess.run(command, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
        logging.info(f"Command output: {result.stdout.strip()}")
        if result.stderr.strip():
            logging.error(f"Command error: {result.stderr.strip()}")
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        logging.error(f"Command failed with error: {e.stderr.strip()}")
        return f"Error: {e.stderr.strip()}"


def get_service_name(service_path):
    base_name = os.path.basename(service_path)
    service_name = os.path.splitext(base_name)[0]
    return service_name


class DaemonManager:
    @staticmethod
    def start_service(service_path):
        sv = get_service_name(service_path)
        logging.info(f"Starting '{sv}' service")
        return run_command(['sudo', 'systemctl', 'start', sv])

    @staticmethod
    def stop_service(service_path):
        sv = get_service_name(service_path)
        logging.info(f"Stopping '{sv}' service")
        return run_command(['sudo', 'systemctl', 'stop', sv])


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

        data_dir = self.args.data_path
        snapshot_metadata = {"name": timestamp_str, "created_at": get_current_time_iso(),
                             "compress": self.args.compress, "total_size": 0, "data": []}

        for root, _, files in os.walk(data_dir):
            for file in files:
                file_path = os.path.join(root, file)
                file_name, file_ext = os.path.splitext(file)
                compressed_file_name = f"{file_name}.{self.args.compress}"
                compressed_file_path = os.path.join(snapshot_dir, compressed_file_name)
                rel_path = os.path.relpath(compressed_file_path, self.args.snapshot_path)

                if rel_path.startswith('snapshots' + os.path.sep):
                    rel_path = rel_path[len('snapshots' + os.path.sep):]

                if self.args.compress == 'zip':
                    logging.info(f"Creating ZIP archive '{compressed_file_path}'")
                    with zipfile.ZipFile(compressed_file_path, 'w', zipfile.ZIP_DEFLATED) as zipf:
                        zipf.write(file_path, file)
                elif self.args.compress == 'tar':
                    logging.info(f"Creating TAR archive '{compressed_file_path}'")
                    subprocess.run(['tar', '-cvf', compressed_file_path, '-C', os.path.dirname(file_path), file])

                compressed_file_size = os.path.getsize(compressed_file_path)
                snapshot_metadata["total_size"] += compressed_file_size

                file_info = {
                    "name": file_name,
                    "path": rel_path,
                    "sha": Metadata.sha256(compressed_file_path),
                    "size": compressed_file_size
                }
                snapshot_metadata["data"].append(file_info)

        Metadata.update_metadata_file(self.args.snapshot_path, snapshot_metadata)


class Validation:
    @staticmethod
    def validate_args(args):
        logging.info('Validating arguments')

        if not os.path.isfile(args.service_path):
            raise ValueError(f"Service file '{args.service_path}' does not exist.")
        logging.info(f"Service file '{args.service_path}' exists")

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

    @staticmethod
    def validate():
        if os.name == "nt":
            raise EnvironmentError("Windows not supported.")
        if os.geteuid() != 0:
            raise PermissionError("This script requires sudo/root access. Please run with sudo.")


class ProcessBackup:
    def __init__(self, args):
        self.args = args

    def run(self):
        Validation.validate()
        Validation.validate_args(self.args)
        DaemonManager.stop_service(self.args.service_path)
        snapshot_manager = SnapshotManager(self.args)
        snapshot_manager.manage_snapshots()
        snapshot_manager.create_snapshot()
        DaemonManager.start_service(self.args.service_path)


def parse_args():
    user_home = os.path.expanduser("~")
    default_data_path = os.path.join(user_home, 'pactus')

    parser = argparse.ArgumentParser(description='Pactus Blockchain Snapshot Tool')
    parser.add_argument('--service_path', required=True, help='Path to pactus systemctl service')
    parser.add_argument('--data_path', default=default_data_path, help='Path to data directory')
    parser.add_argument('--compress', choices=['zip', 'tar'], default='zip', help='Compression type')
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
