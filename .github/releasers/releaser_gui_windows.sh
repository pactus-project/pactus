#!/bin/bash

set -e

ROOT_DIR="$(pwd)"
VERSION="$(echo `git -C ${ROOT_DIR} describe --abbrev=0 --tags` | sed 's/^.//')" # "v1.2.3" -> "1.2.3"
BUILD_DIR="${ROOT_DIR}/build"
PACKAGE_NAME="pactus-gui_${VERSION}"
PACKAGE_DIR="${ROOT_DIR}/${PACKAGE_NAME}"

# Check the architecture
ARC="$(uname -m)"

if [ "${ARC}" = "x86_64" ]; then
    FILE_NAME="${PACKAGE_NAME}_windows_amd64"
    INNO_PATH="/c/Program Files (x86)/Inno Setup 6"
elif [ "${ARC}" = "arm64" ]; then
    FILE_NAME="${PACKAGE_NAME}_windows_arm64"
    INNO_PATH="/c/Program Files/Inno Setup 6"
else
    echo "Unsupported architecture: ${ARC}"
    exit 1
fi

#!/bin/bash

########################################
# Function to sign an app using SignPath
# Parameters:
#   $1 - Path to unsigned app
#   $2 - Path to save signed app
#   $3 - App Name
sign_app() {
    local unsigned_app_path="$1"
    local signed_app_path="$2"
    local name="$3"

    # Check if all parameters are provided
    if [[ -z "${unsigned_app_path}" || -z "${signed_app_path}" || -z "${name}" ]]; then
        echo "Error: Missing parameters. Usage: sign_app_with_signpath <unsigned_app_path> <signed_app_path> <name>"
        return 1
    fi

    # Check if unsigned_app_path exists
    if [[ ! -f "${unsigned_app_path}" ]]; then
        echo "Error: Unsigned app not found at $unsigned_app_path"
        return 1
    fi

    # Check if SIGNPATH_PROJECT_SLUG is set
    if [[ -z "${SIGNPATH_PROJECT_SLUG}" ]]; then
        echo "Error: SIGNPATH_PROJECT_SLUG environment variable not set"
        return 1
    fi

    # Check if SIGNPATH_ORGANIZATION_ID is set
    if [[ -z "${SIGNPATH_ORGANIZATION_ID}" ]]; then
        echo "Error: SIGNPATH_ORGANIZATION_ID environment variable not set"
        return 1
    fi

    # Check if SIGNPATH_SIGNING_POLICY_SLUG is set
    if [[ -z "${SIGNPATH_SIGNING_POLICY_SLUG}" ]]; then
        echo "Error: SIGNPATH_SIGNING_POLICY_SLUG environment variable not set"
        return 1
    fi

    # Check if SIGNPATH_API_TOKEN is set
    if [[ -z "${SIGNPATH_API_TOKEN}" ]]; then
        echo "Error: SIGNPATH_API_TOKEN environment variable not set"
        return 1
    fi

    base_url="https://app.signpath.io/API/v1/${SIGNPATH_ORGANIZATION_ID}"

    # Submit the signing request
    echo "Submitting signing request to SignPath for ${name}..."
    local submit_response
    submit_response=$(curl -v -X POST \
        -H "Authorization: Bearer ${SIGNPATH_API_TOKEN}" \
        -F "ProjectSlug=${SIGNPATH_PROJECT_SLUG}" \
        -F "SigningPolicySlug=${SIGNPATH_SIGNING_POLICY_SLUG}" \
        -F "Artifact=@${unsigned_app_path}" \
        -F "Description=$name version ${version}" \
        -F "Parameters[productVersion]=${version}" \
        -F "Parameters[productName]=${name}" \
        "${base_url}/SigningRequests")

    # Check if submission was successful
    local signing_request_id
    signing_request_id=$(echo "${submit_response}" | jq -r '.signingRequestId')
    if [[ -z "${signing_request_id}" || "${signing_request_id}" == "null" ]]; then
        echo "Error: Failed to submit signing request"
        echo "Response: ${submit_response}"
        return 1
    fi

    echo "Signing request submitted. ID: ${signing_request_id}"

    # Wait for signing to complete
    local status
    local attempts=0
    local max_attempts=3

    while [[ $attempts -lt $max_attempts ]]; do
        sleep 5
        attempts=$((attempts + 1))

        echo "Checking signing status (attempt $attempts/$max_attempts)..."
        local status_response
        status_response=$(curl -v -X GET \
            -H "Authorization: Bearer ${SIGNPATH_API_TOKEN}" \
            "${base_url}/SigningRequests/${signing_request_id}")

        echo "status_response: ${status_response}"
        status=$(echo "${status_response}" | jq -r '.status')

        if [[ "${status}" == "Failed" ]]; then
            echo "Error: Signing failed"
            echo "Status response: $status_response"
            return 1
        elif [[ "${status}" == "Completed" ]]; then
            break
        fi
    done

    if [[ "${status}" != "Completed" ]]; then
        echo "Error: Signing timed out after $max_attempts attempts"
        return 1
    fi

    # Download the signed artifact
    echo "Downloading signed artifact..."
    if ! curl -v -X GET \
        -H "Authorization: Bearer ${SIGNPATH_API_TOKEN}" \
        -o "${signed_app_path}" \
        "${base_url}/SigningRequests/$signing_request_id/SignedArtifact"; then
        echo "Error: Failed to download signed artifact"
        return 1
    fi

    # Verify the signed artifact was downloaded
    if [[ ! -f "${signed_app_path}" ]]; then
        echo "Error: Signed artifact was not saved to $signed_app_path"
        return 1
    fi

    echo "Successfully signed and saved artifact to $signed_app_path"
    return 0
}
########################################

mkdir ${PACKAGE_DIR}

echo "Building the binaries for Windows ${ARC} architecture"

# This fixes a bug in pkgconfig: invalid flag in pkg-config --libs: -Wl,-luuid
sed -i -e 's/-Wl,-luuid/-luuid/g' /mingw64/lib/pkgconfig/gdk-3.0.pc

CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/pactus-daemon_unsigned.exe        ./cmd/daemon
CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/pactus-wallet_unsigned.exe        ./cmd/wallet
CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ${BUILD_DIR}/pactus-shell_unsigned.exe         ./cmd/shell
go build -ldflags "-s -w -H windowsgui" -trimpath -tags gtk -o ${BUILD_DIR}/pactus-gui_unsigned.exe ./cmd/gtk

sign_app ${BUILD_DIR}/pactus-daemon_unsigned.exe  ${BUILD_DIR}/pactus-daemon.exe "Pactus Daemon"
sign_app "${BUILD_DIR}/pactus-wallet_unsigned.exe"  "${BUILD_DIR}/pactus-wallet.exe" "Pactus Wallet"
sign_app "${BUILD_DIR}/pactus-shell_unsigned.exe"   "${BUILD_DIR}/pactus-shell.exe"  "Pactus Shell"
sign_app "${BUILD_DIR}/pactus-gui_unsigned.exe"     "${BUILD_DIR}/pactus-gui.exe"    "Pactus GUI"

# Copying the necessary libraries
echo "Creating GUI directory"
GUI_DIR="${PACKAGE_DIR}/pactus-gui"
mkdir ${GUI_DIR}

echo "Copying required DLLs and EXEs from ${MINGW_PREFIX}/bin"

cd "${MINGW_PREFIX}/bin"
cp -p \
  "gdbus.exe" \
  "gspawn-win64-helper.exe" \
  "gspawn-win64-helper-console.exe" \
  "libatk-1.0-0.dll" \
  "libbrotlicommon.dll" \
  "libbrotlidec.dll" \
  "libbz2-1.dll" \
  "libcairo-2.dll" \
  "libcairo-gobject-2.dll" \
  "libdatrie-1.dll" \
  "libdeflate.dll" \
  "libepoxy-0.dll" \
  "libexpat-1.dll" \
  "libffi-8.dll" \
  "libfontconfig-1.dll" \
  "libfreetype-6.dll" \
  "libfribidi-0.dll" \
  "libgcc_s_seh-1.dll" \
  "libgdk_pixbuf-2.0-0.dll" \
  "libgdk-3-0.dll" \
  "libgio-2.0-0.dll" \
  "libglib-2.0-0.dll" \
  "libgmodule-2.0-0.dll" \
  "libgobject-2.0-0.dll" \
  "libgomp-1.dll" \
  "libgraphite2.dll" \
  "libgtk-3-0.dll" \
  "libharfbuzz-0.dll" \
  "libiconv-2.dll" \
  "libintl-8.dll" \
  "libjbig-0.dll" \
  "libjpeg-8.dll" \
  "libLerc.dll" \
  "liblzma-5.dll" \
  "libpango-1.0-0.dll" \
  "libpangocairo-1.0-0.dll" \
  "libpangoft2-1.0-0.dll" \
  "libpangowin32-1.0-0.dll" \
  "libpcre2-16-0.dll" \
  "libpcre2-32-0.dll" \
  "libpcre2-8-0.dll" \
  "libpcre2-posix-3.dll" \
  "libpixman-1-0.dll" \
  "libpng16-16.dll" \
  "librsvg-2-2.dll" \
  "libstdc++-6.dll" \
  "libsystre-0.dll" \
  "libthai-0.dll" \
  "libtiff-6.dll" \
  "libtre-5.dll" \
  "libwebp-7.dll" \
  "libwinpthread-1.dll" \
  "libxml2-2.dll" \
  "libzstd.dll" \
  "zlib1.dll" \
  "libsharpyuv-0.dll" \
  "${GUI_DIR}"


echo "Copying GDK pixbuf from ${MINGW_PREFIX}/lib/gdk-pixbuf-2.0"
mkdir -p "${GUI_DIR}/lib/gdk-pixbuf-2.0"
cp -rp "${MINGW_PREFIX}/lib/gdk-pixbuf-2.0" "${GUI_DIR}/lib"

#########
### Based on this tutorial: https://www.gtk.org/docs/installations/windows#building-and-distributing-your-application

echo "Step 1: gtk-3.0/"
mkdir -p "${GUI_DIR}/share/themes/Windows10/gtk-3.0/"
cp -rp "${MINGW_PREFIX}/share/gtk-3.0" "${GUI_DIR}/share/themes/Windows10/"

echo "Step 2: Adwaita icons"
mkdir -p "${GUI_DIR}/share/icons/Adwaita"
cp -rp "${MINGW_PREFIX}/share/icons/Adwaita" "${GUI_DIR}/share/icons"

echo "Step 3: hicolor icons"
mkdir -p "${GUI_DIR}/share/icons/hicolor"
cp -rp "${MINGW_PREFIX}/share/icons/hicolor" "${GUI_DIR}/share/icons"

echo "Step 4: settings.ini"
mkdir "${GUI_DIR}/share/gtk-3.0/"
echo "[Settings]
gtk-theme-name=Windows10" > "${GUI_DIR}/share/gtk-3.0/settings.ini"

echo "Step 5: glib-compile-schemas"
mkdir -p "${GUI_DIR}/share/glib-2.0/schemas"
cp -rp "${MINGW_PREFIX}/share/glib-2.0/schemas/gschemas.compiled" "${GUI_DIR}/share/glib-2.0/schemas"


# Moving binaries to package directory
cd ${ROOT_DIR}
echo "Moving binaries"
mv ${BUILD_DIR}/pactus-daemon.exe  ${PACKAGE_DIR}/pactus-daemon.exe
mv ${BUILD_DIR}/pactus-wallet.exe  ${PACKAGE_DIR}/pactus-wallet.exe
mv ${BUILD_DIR}/pactus-shell.exe   ${PACKAGE_DIR}/pactus-shell.exe
mv ${BUILD_DIR}/pactus-gui.exe     ${PACKAGE_DIR}/pactus-gui/pactus-gui.exe

echo "Archiving the package"
7z a ${ROOT_DIR}/${FILE_NAME}.zip ${PACKAGE_DIR}

echo "Creating installer"
cat << EOF > ${ROOT_DIR}/inno.iss
[Setup]
AppId=Pactus
AppName=Pactus
AppVersion=${VERSION}
AppPublisher=Pactus
AppPublisherURL=https://pactus.org/
DefaultDirName={autopf}/Pactus
DefaultGroupName=Pactus
SetupIconFile=${ROOT_DIR}/.github/releasers/pactus.ico
LicenseFile=${ROOT_DIR}/LICENSE
Uninstallable=yes
UninstallDisplayIcon={app}\\pactus-gui\\pactus-gui.exe

[Files]
Source:"${PACKAGE_NAME}/*"; DestDir:"{app}"; Flags: recursesubdirs

[Icons]
Name:"{group}\\Pactus"; Filename:"{app}\\pactus-gui\\pactus-gui.exe"
Name:"{commondesktop}\\Pactus"; Filename:"{app}\\pactus-gui\\pactus-gui.exe"

[Run]
Filename:"{app}\\pactus-gui\\pactus-gui.exe"; Description:"Launch Pactus"; Flags: postinstall nowait
EOF

cd ${ROOT_DIR}

INNO_DIR=$(cygpath -w -s "${INNO_PATH}")
"${INNO_DIR}/ISCC.exe" "${ROOT_DIR}/inno.iss"
mv "Output/mysetup.exe" "${ROOT_DIR}/${FILE_NAME}_installer_unsigned.exe"

sign_app ${ROOT_DIR}/${FILE_NAME}_installer_unsigned.exe  ${ROOT_DIR}/${FILE_NAME}_installer.exe "Pactus Installer"
