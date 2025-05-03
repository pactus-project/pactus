# Pactus Blockchain Snapshot tool
#
# This script first stops the Pactus node if it is running, then creates a backup of the blockchain data by copying
# or compressing it based on the specified options. The backup is stored in a timestamped snapshot directory along
# with a `metadata.json` file that contains detailed information about the snapshot, including file paths and
# checksums. Finally, the script manages the retention of snapshots, ensuring only a specified number of recent
# backups are kept.
#
# Arguments
#
# - `--service_path`: This argument specifies the path to the `pactus` service file to manage systemctl service.
# - `--data_path`: This argument specifies the path to the Pactus data folder to create snapshots.
#    - Windows: `C:\Users\{user}\pactus\data`
#    - Linux or Mac: `/home/{user}/pactus/data`
# - `--compress`: This argument specifies the compression method based on your choice ['none', 'zip', 'tar'],
# with 'none' being without compression.
# - `--retention`: This argument sets the number of snapshots to keep.
# - `--snapshot_path`: This argument sets a custom path for snapshots, with the default being the current
# working directory of the script.
#
# How to run?
#
# For create snapshots just run this command:
#
# sudo python3 snapshot.py --service_path /etc/systemd/system/pactus.service --data_path /home/{user}/pactus/data
# --compress zip --retention 3


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
        format="[%(asctime)s] %(message)s", datefmt="%Y-%m-%d-%H:%M", level=logging.INFO
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
        metadata_file = os.path.join(snapshot_path, "snapshots", "metadata.json")
        if os.path.isfile(metadata_file):
            logging.info(f"Updating existing metadata file '{metadata_file}'")
            with open(metadata_file, "r") as f:
                metadata = json.load(f)
        else:
            logging.info(f"Creating new metadata file '{metadata_file}'")
            metadata = []

        formatted_metadata = {
            "name": snapshot_metadata["name"],
            "created_at": snapshot_metadata["created_at"],
            "compress": snapshot_metadata["compress"],
            "data": snapshot_metadata["data"],
        }

        metadata.append(formatted_metadata)

        with open(metadata_file, "w") as f:
            json.dump(metadata, f, indent=4)

    @staticmethod
    def update_metadata_after_removal(snapshots_dir, removed_snapshots):
        metadata_file = os.path.join(snapshots_dir, "metadata.json")
        if not os.path.isfile(metadata_file):
            return

        logging.info(f"Updating metadata file '{metadata_file}' after snapshot removal")
        with open(metadata_file, "r") as f:
            metadata = json.load(f)

        updated_metadata = [
            entry for entry in metadata if entry["name"] not in removed_snapshots
        ]

        with open(metadata_file, "w") as f:
            json.dump(updated_metadata, f, indent=4)

    @staticmethod
    def create_snapshot_json(data_dir, snapshot_subdir):
        files = []
        for root, _, filenames in os.walk(data_dir):
            for filename in filenames:
                file_path = os.path.join(root, filename)
                rel_path = os.path.relpath(file_path, data_dir)
                snapshot_rel_path = os.path.join(snapshot_subdir, rel_path).replace(
                    "\\", "/"
                )
                file_info = {
                    "name": filename,
                    "path": snapshot_rel_path,
                    "sha": Metadata.sha256(file_path),
                }
                files.append(file_info)

        return {"data": files}

    @staticmethod
    def create_compressed_snapshot_json(compressed_file, rel_path):
        compressed_file_size = os.path.getsize(compressed_file)
        file_info = {
            "name": os.path.basename(compressed_file),
            "path": rel_path,
            "sha": Metadata.sha256(compressed_file),
            "size": compressed_file_size,
        }

        return {"data": file_info}


def run_command(command):
    logging.info(f"Running command: {' '.join(command)}")
    try:
        result = subprocess.run(
            command,
            check=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
        )

        if result.stdout.strip():
            logging.info(f"Command output: {result.stdout.strip()}")

        if result.stderr.strip():
            # Downgrade from error to info for successful commands
            logging.info(f"Command stderr: {result.stderr.strip()}")

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
    def start_service(service_path=None, docker_compose_path=None, docker_service_name=None):
        if docker_compose_path and docker_service_name:
            logging.info(f"Starting Docker Compose service '{docker_service_name}' at '{docker_compose_path}'")
            return run_command(["docker", "compose", "-f", docker_compose_path, "start", docker_service_name])
        elif service_path:
            sv = get_service_name(service_path)
            logging.info(f"Starting systemctl service '{sv}'")
            return run_command(["sudo", "systemctl", "start", sv])
        return None

    @staticmethod
    def stop_service(service_path=None, docker_compose_path=None, docker_service_name=None):
        if docker_compose_path and docker_service_name:
            logging.info(f"Stopping Docker Compose service '{docker_service_name}' at '{docker_compose_path}'")
            return run_command(["docker", "compose", "-f", docker_compose_path, "stop", docker_service_name])
        elif service_path:
            sv = get_service_name(service_path)
            logging.info(f"Stopping systemctl service '{sv}'")
            return run_command(["sudo", "systemctl", "stop", sv])
        return None


class SnapshotManager:
    def __init__(self, args):
        self.args = args

    def manage_snapshots(self):
        snapshots_dir = os.path.join(self.args.snapshot_path, "snapshots")
        logging.info(f"Managing snapshots in '{snapshots_dir}'")

        if not os.path.exists(snapshots_dir):
            logging.info(
                f"Snapshots directory '{snapshots_dir}' does not exist. Creating it."
            )
            os.makedirs(snapshots_dir)

        snapshots = sorted(
            [s for s in os.listdir(snapshots_dir) if s != "metadata.json"]
        )

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
        snapshot_dir = os.path.join(self.args.snapshot_path, "snapshots", timestamp_str)
        logging.info(f"Creating snapshot directory '{snapshot_dir}'")
        os.makedirs(snapshot_dir, exist_ok=True)

        data_dir = os.path.join(snapshot_dir, "data")
        if self.args.compress == "none":
            logging.info(f"Copying data from '{self.args.data_path}' to '{data_dir}'")
            shutil.copytree(self.args.data_path, data_dir)
            snapshot_metadata = Metadata.create_snapshot_json(data_dir, timestamp_str)
        elif self.args.compress == "zip":
            zip_file = os.path.join(snapshot_dir, "data.zip")
            rel = os.path.relpath(zip_file, snapshot_dir)
            meta_path = os.path.join(timestamp_str, rel)
            logging.info(f"Creating ZIP archive '{zip_file}'")
            with zipfile.ZipFile(zip_file, "w", zipfile.ZIP_DEFLATED) as zipf:
                for root, _, files in os.walk(self.args.data_path):
                    for file in files:
                        full_path = os.path.join(root, file)
                        rel_path = os.path.relpath(full_path, self.args.data_path)
                        zipf.write(full_path, os.path.join("data", rel_path))
            snapshot_metadata = Metadata.create_compressed_snapshot_json(
                zip_file, meta_path
            )
        elif self.args.compress == "tar":
            tar_file = os.path.join(snapshot_dir, "data.tar.gz")
            rel = os.path.relpath(tar_file, snapshot_dir)
            meta_path = os.path.join(timestamp_str, rel)
            logging.info(f"Creating TAR.GZ archive '{tar_file}'")
            subprocess.run(["tar", "-czvf", tar_file, "-C", self.args.data_path, "."])
            snapshot_metadata = Metadata.create_compressed_snapshot_json(
                tar_file, meta_path
            )

        snapshot_metadata["name"] = timestamp_str
        snapshot_metadata["created_at"] = get_current_time_iso()
        snapshot_metadata["compress"] = self.args.compress

        Metadata.update_metadata_file(self.args.snapshot_path, snapshot_metadata)


class Validation:
    @staticmethod
    def validate_args(args):
        logging.info("Validating arguments")

        # Ensure at least one service method is provided
        if not args.service_path and not args.docker_compose_path:
            raise ValueError("Either --service_path or --docker_compose_path must be provided.")

        # Validate systemctl service path if provided
        if args.service_path:
            if not os.path.isfile(args.service_path):
                raise ValueError(f"Service file '{args.service_path}' does not exist.")
            logging.info(f"Service file '{args.service_path}' exists")

        # Validate docker-compose if provided
        if args.docker_compose_path:
            if not os.path.isfile(args.docker_compose_path):
                raise ValueError(f"Docker Compose file '{args.docker_compose_path}' does not exist.")
            if not args.docker_service_name:
                raise ValueError("--docker_service_name is required when using --docker_compose_path")
            logging.info(f"Docker Compose file '{args.docker_compose_path}' exists")
            logging.info(f"Docker service name is '{args.docker_service_name}'")

        # Common validations
        if not os.path.isdir(args.data_path):
            raise ValueError(f"Data path '{args.data_path}' does not exist.")
        logging.info(f"Data path '{args.data_path}' exists")

        if not os.access(args.data_path, os.W_OK):
            raise PermissionError(
                f"No permission to access data path '{args.data_path}'."
            )
        logging.info(f"Permission to access data path '{args.data_path}' confirmed")

        if args.compress == "zip" and not shutil.which("zip"):
            raise EnvironmentError("The 'zip' command is not available.")
        elif args.compress == "zip":
            logging.info("The 'zip' command is available")

        if args.compress == "tar" and not shutil.which("tar"):
            raise EnvironmentError("The 'tar' command is not available.")
        elif args.compress == "tar":
            logging.info("The 'tar' command is available")

        if args.retention <= 0:
            raise ValueError("Retention value must be greater than 0.")
        logging.info(f"Retention value is set to {args.retention}")

        if not os.access(args.snapshot_path, os.W_OK):
            raise PermissionError(
                f"No permission to access snapshot path '{args.snapshot_path}'."
            )
        logging.info(
            f"Permission to access snapshot path '{args.snapshot_path}' confirmed"
        )

        snapshots_dir = os.path.join(args.snapshot_path, "snapshots")
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
        DaemonManager.stop_service(
            service_path=self.args.service_path,
            docker_compose_path=self.args.docker_compose_path,
            docker_service_name=self.args.docker_service_name
        )
        snapshot_manager = SnapshotManager(self.args)
        snapshot_manager.manage_snapshots()
        snapshot_manager.create_snapshot()
        DaemonManager.start_service(
            service_path=self.args.service_path,
            docker_compose_path=self.args.docker_compose_path,
            docker_service_name=self.args.docker_service_name
        )


def parse_args():
    user_home = os.path.expanduser("~")
    default_data_path = os.path.join(user_home, "pactus")

    parser = argparse.ArgumentParser(description="Pactus Blockchain Snapshot Tool")
    parser.add_argument(
        "--service_path", help="Path to pactus systemctl service"
    )
    parser.add_argument(
        "--docker_compose_path",
        help="Path to docker-compose.yml file to manage Docker-based service"
    )
    parser.add_argument(
        "--docker_service_name",
        help="Name of the Docker service in the Compose file"
    )
    parser.add_argument(
        "--data_path", default=default_data_path, help="Path to data directory"
    )
    parser.add_argument(
        "--compress",
        choices=["none", "zip", "tar"],
        default="none",
        help="Compression type",
    )
    parser.add_argument(
        "--retention", type=int, default=3, help="Number of snapshots to retain"
    )
    parser.add_argument(
        "--snapshot_path", default=os.getcwd(), help="Path to store snapshots"
    )

    return parser.parse_args()


def main():
    setup_logging()
    args = parse_args()
    process_backup = ProcessBackup(args)
    process_backup.run()


if __name__ == "__main__":
    main()
