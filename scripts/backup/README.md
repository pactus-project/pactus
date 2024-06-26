# Pactus Blockchain Backup tool

This script first stops the Pactus node if it is running, then creates a backup of the blockchain data by copying
or compressing it based on the specified options. The backup is stored in a timestamped snapshot directory along with
a `metadata.json` file that contains detailed information about the snapshot, including file paths and checksums. Finally,
the script manages the retention of snapshots, ensuring only a specified number of recent backups are kept.

## Arguments

- `--daemon_path`: This argument specifies the path to the `pactus-daemon` file to manage and
check the process.
- `--data_path`: This argument specifies the path to the Pactus data folder to create snapshots.
   - Windows: `C:\Users\{user}\pactus\data`
   - Linux or Mac: `/home/{user}/pactus/data`
- `--compress`: This argument specifies the compression method based on your choice ['none', 'zip', 'tar'],
with 'none' being without compression.
- `--retention`: This argument sets the number of snapshots to keep.
- `--snapshot_path`: This argument sets a custom path for snapshots, with the default being the current
working directory of the script.

## How to run?

For create snapshots just run this command:

```shell
python3 backup.py --daemon_path /path/pactus/pactus-daemon --data_path /home/{user}/pactus/data --compress zip --retention 3
```
