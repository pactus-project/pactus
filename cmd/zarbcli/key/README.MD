# zarb key

`zarbcli key` is a simple command-line tool to work with zarb key files.

## Usage

### Generate new key

Generate new key file. This key file is encrypted with the password you provide while creating a new account.

Example:

```bash
zarbcli key generate
```

### Recovering a key

You can recover a key if you have the seed phrase and passphrase.

Example:

```bash
zarbcli key recover
```

### Inspect KeyFile

Print various information about the given key file.
Private key information can be printed using the `-e` flag. Make sure you use this feature with great caution.

Example:

```bash
zarbcli key inspect <PATH_TO_KEYFILE>
zarbcli key inspect -e <PATH_TO_KEYFILE>
```

### Signing a message

Signs the message with the private key. User can directly specify the private key to sign the message or specify the a key file and password to decrypt the private key and sign the message.

Example:

```bash
zarbcli key sign -m hello -k <PATH_TO_KEYFILE>
```

### Verifying signature

Verify the signature of a message.

Example:

```bash
zarbcli key verify -m hello <PUBLIC_KEY> <SIGNATURE>
```

### Change KeyFile password

Change the password of the key file which is used for encryption and decryption.

Example:

```bash
zarbcli key change-auth <PATH_TO_KEYFILE>
```
