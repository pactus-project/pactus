#!/bin/bash

# The 'set -e' command causes the script to immediately exit
# if any command returns a non-zero exit status (i.e., an error).
set -e

replace_in_place() {
  if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "$1" "$2"
  else
    sed -i "$1" "$2"
  fi
}

ROOT_DIR="$(pwd)"
PACKAGE_DIR="${ROOT_DIR}/packages"
GEN_DIR="${ROOT_DIR}/www/grpc/gen"
GENERATOR_DIR="${PACKAGE_DIR}/js/generator"
VERSION="$(echo `git -C ${ROOT_DIR} describe --abbrev=0 --tags` | sed 's/^.//')" # "v1.2.3" -> "1.2.3"

if [[ -z "$VERSION" ]]; then
  echo "❌ Error: Version tag not found."
  exit 1
fi

echo "Packing Version:" ${VERSION}

rm -rf ${PACKAGE_DIR}
mkdir -p ${PACKAGE_DIR}
mkdir -p ${PACKAGE_DIR}/js/{pactus-grpc,pactus-jsonrpc}/src

echo "== Building pactus-grpc package for JavaScript"
cp -R ${ROOT_DIR}/.github/packager/js/grpc/package.json ${PACKAGE_DIR}/js/pactus-grpc
cp -R ${ROOT_DIR}/.github/packager/js/index.js ${PACKAGE_DIR}/js/pactus-grpc/src
cp -R ${GEN_DIR}/js/* ${PACKAGE_DIR}/js/pactus-grpc/src
cp ${ROOT_DIR}/LICENSE ${PACKAGE_DIR}/js/pactus-grpc
cp ${ROOT_DIR}/README.md ${PACKAGE_DIR}/js/pactus-grpc
replace_in_place "s/{{ VERSION }}/$VERSION/g" "${PACKAGE_DIR}/js/pactus-grpc/package.json"

echo "== Building pactus-jsonrpc package for JavaScript"
rm -rf "$GENERATOR_DIR"
git clone https://github.com/pactus-project/generator.git "$GENERATOR_DIR" && cd "$GENERATOR_DIR"
npm i && npm run build
cd "$ROOT_DIR" && $GENERATOR_DIR/build/cli.js generate \
  -t client \
  -l typescript \
  -n pactusClientTs \
  -d "www/grpc/gen/open-rpc/pactus-openrpc.json" \
  -o "$GENERATOR_DIR/gen"
cp -R ${ROOT_DIR}/.github/packager/js/jsonrpc/package.json ${PACKAGE_DIR}/js/pactus-jsonrpc
cp -R $GENERATOR_DIR/gen/client/typescript/src/* ${PACKAGE_DIR}/js/pactus-jsonrpc/src
replace_in_place "s/{{ VERSION }}/$VERSION/g" "${PACKAGE_DIR}/js/pactus-jsonrpc/package.json"

echo "== Building pactus-grpc package for Python"
cp -R ${ROOT_DIR}/.github/packager/python ${PACKAGE_DIR}/python
cp ${GEN_DIR}/python/* ${PACKAGE_DIR}/python/pactus_grpc
cp ${ROOT_DIR}/LICENSE ${PACKAGE_DIR}/python/
cp ${ROOT_DIR}/README.md ${PACKAGE_DIR}/python/
replace_in_place "s/{{ VERSION }}/$VERSION/g" ${PACKAGE_DIR}/python/setup.py
