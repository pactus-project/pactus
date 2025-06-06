#!/bin/sh

# Function to check if a command is available
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Download function
download_file() {
    file_name=$1
    url=$2
    echo "== Downloading ${file_name}..."
    if ! curl --fail --proto '=https' --tlsv1.2 -# -L -o "${file_name}" "${url}" ; then
        echo "Failed to download ${file_name}. Check the URL or your internet connection."
        exit 1
    fi
}


# Check if sha256sum is installed
if ! command_exists sha256sum; then
    echo "sha256sum is not installed. Please install sha256sum and try again."
    exit 1
fi

# Detect the operating system and CPU type
OS=$(uname -s)
ARCH=$(uname -m)
EXT=".tar.gz"

case $OS in
    ######################
    Android)
        OS="android"
        case $ARCH in
            aarch64 | arm64)
                ARCH="arm64"
                ;;
            *)
                echo "Unrecognized CPU type: $ARCH"
                exit 1
                ;;
        esac
        ;;

    ######################
    Linux)
        OS="linux"
        case $ARCH in
            x86_64)
                ARCH="amd64"
                ;;
            aarch64 | arm64)
                ARCH="arm64"
                ;;
            *)
                echo "Unrecognized CPU type: $ARCH"
                exit 1
                ;;
        esac
        ;;

    ######################
    FreeBSD)
        OS="freebsd"
        case $ARCH in
            amd64)
                ARCH="amd64"
                ;;
            *)
                echo "Unrecognized CPU type: $ARCH"
                exit 1
                ;;
        esac
        ;;

    ######################
    Darwin)
        OS="darwin"
        case $ARCH in
            x86_64)
                ARCH="amd64"
                ;;
            arm64)
                ARCH="arm64"
                ;;
            *)
                echo "Unrecognized CPU type: $ARCH"
                exit 1
                ;;
        esac
        ;;

    ######################
    MINGW* | MSYS* | CYGWIN*)
        if ! command_exists unzip; then
            echo "unzip is not installed. Please install unzip and try again."
            exit 1
        fi
        OS="windows"
        EXT=".zip"
        case $ARCH in
            i686 | i386)
                ARCH="386"
                ;;
            x86_64)
                ARCH="amd64"
                ;;
            aarch64 | arm64)
                ARCH="arm64"
                ;;
            *)
                echo "Unrecognized CPU type: $ARCH"
                exit 1
                ;;
        esac
        ;;

    *)
        echo "Unrecognized OS type: $OS"
        exit 1
        ;;
esac

# Version number
VERSION="__VERSION__"

# Set the server URL where the binary and checksum files are hosted
SERVER_URL="https://github.com/pactus-project/pactus/releases/download/v${VERSION}"

# Set the filename of the binary and checksum
CHECKSUM_FILE="SHA256SUMS"

# Set the filename of the binary and checksum
FILE_NAME="pactus-cli_${VERSION}_${OS}_${ARCH}${EXT}"

# Destination directory that is the current directory
DEST_DIR="$(pwd)"
EXTRACTED_DIR="${DEST_DIR}/pactus-cli_${VERSION}"

# Create a temporary directory for downloads
DOWN_DIR="$(mktemp -d)"

# Check if extractopn folder exists and print an error if it does
if [ -e "${EXTRACTED_DIR}" ]; then
    echo "Destination directory '${EXTRACTED_DIR}' already exists."
    exit 1
fi

cd ${DOWN_DIR}

# Download the files using the download_file function
download_file "${FILE_NAME}" "${SERVER_URL}/${FILE_NAME}"
download_file "${CHECKSUM_FILE}" "${SERVER_URL}/${CHECKSUM_FILE}"

# Verify the checksum
echo "== Verify the checksum..."
sha256sum --ignore-missing -c "${CHECKSUM_FILE}"
if [ $? -ne 0 ]; then
    echo "Checksum verification failed for ${FILE_NAME}!"
    exit 1
fi

# Extracting
echo "== Extracting ${FILE_NAME}..."
if [ "${OS}" = "windows" ]; then
    unzip -n "${FILE_NAME}" -d "${DEST_DIR}" || {
        echo "Error: Extraction failed."
        exit 1
    }
else
    tar -xzf "${FILE_NAME}" -C "${DEST_DIR}" || {
        echo "Error: Extraction failed."
        exit 1
    }
fi
echo "Extracted at ${EXTRACTED_DIR}"


echo ""
echo "Installation completed."
