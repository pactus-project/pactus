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
PROTO_GEN_DIR="${ROOT_DIR}/www/grpc/gen"

if [[ -z "$VERSION" ]]; then
  echo "❌ Error: Version tag not found."
  exit 1
fi

echo "Packing Version:" ${VERSION}

rm -rf ${PACKAGE_DIR}
mkdir -p ${PACKAGE_DIR}
mkdir -p ${PACKAGE_DIR}/js/{pactus-grpc,pactus-jsonrpc}
mkdir -p ${PACKAGE_DIR}/python/{pactus-grpc,pactus-jsonrpc}

echo "== Building pactus-grpc package for JavaScript"
cp -R ${ROOT_DIR}/.github/packager/js/grpc/package.json ${PACKAGE_DIR}/js/pactus-grpc
cp -R ${PROTO_GEN_DIR}/js/* ${PACKAGE_DIR}/js/pactus-grpc
cp ${ROOT_DIR}/LICENSE ${PACKAGE_DIR}/js/pactus-grpc
cp ${ROOT_DIR}/README.md ${PACKAGE_DIR}/js/pactus-grpc
replace_in_place "s/{{ VERSION }}/$VERSION/g" "${PACKAGE_DIR}/js/pactus-grpc/package.json"

echo "== Building pactus-jsonrpc package for JavaScript"
GENERATOR_DIR="${PACKAGE_DIR}/js/generator"
git clone https://github.com/pactus-project/generator.git "$GENERATOR_DIR" && cd "$GENERATOR_DIR"
npm install && npm run build
cd "$ROOT_DIR" && $GENERATOR_DIR/build/cli.js generate \
  -t client \
  -l typescript \
  -n pactus-jsonrpc \
  -d "${ROOT_DIR}/www/grpc/gen/open-rpc/pactus-openrpc.json" \
  -o "$GENERATOR_DIR/gen"
cd "$GENERATOR_DIR/gen/client/typescript"
npm install && tsc
cp $GENERATOR_DIR/gen/client/typescript/build/index.d.ts ${PACKAGE_DIR}/js/pactus-jsonrpc
cp $GENERATOR_DIR/gen/client/typescript/build/index.js ${PACKAGE_DIR}/js/pactus-jsonrpc
cp $GENERATOR_DIR/gen/client/typescript/build/index.js.map ${PACKAGE_DIR}/js/pactus-jsonrpc
cp ${ROOT_DIR}/.github/packager/js/jsonrpc/package.json ${PACKAGE_DIR}/js/pactus-jsonrpc
cp ${ROOT_DIR}/LICENSE ${PACKAGE_DIR}/js/pactus-jsonrpc
cp ${ROOT_DIR}/README.md ${PACKAGE_DIR}/js/pactus-jsonrpc
replace_in_place "s/{{ VERSION }}/$VERSION/g" "${PACKAGE_DIR}/js/pactus-jsonrpc/package.json"


echo "== Building pactus-grpc package for Python"
cp -R ${ROOT_DIR}/.github/packager/python/grpc/* ${PACKAGE_DIR}/python/pactus-grpc
cp ${PROTO_GEN_DIR}/python/* ${PACKAGE_DIR}/python/pactus-grpc/pactus_grpc
cp ${ROOT_DIR}/LICENSE ${PACKAGE_DIR}/python/pactus-grpc
cp ${ROOT_DIR}/README.md ${PACKAGE_DIR}/python/pactus-grpc
replace_in_place "s/{{ VERSION }}/$VERSION/g" ${PACKAGE_DIR}/python/pactus-grpc/setup.py

echo "== Building pactus-jsonrpc package for Python"
pip install openrpcclientgenerator
ORPC_DIR="${PACKAGE_DIR}/python/orpc"
mkdir -p ${ORPC_DIR}
cp "${ROOT_DIR}/www/grpc/gen/open-rpc/pactus-openrpc.json" ${ORPC_DIR}/openrpc.json
cd ${ORPC_DIR}
orpc python example.com ./out
cp -R ${ROOT_DIR}/.github/packager/python/jsonrpc/* ${PACKAGE_DIR}/python/pactus-jsonrpc
cp ${ORPC_DIR}/out/python/pactus-open-rpc-http-client/pactus_open_rpc_http_client/client.py ${PACKAGE_DIR}/python/pactus-jsonrpc/pactus_jsonrpc/client.py
cp ${ORPC_DIR}/out/python/pactus-open-rpc-http-client/pactus_open_rpc_http_client/models.py ${PACKAGE_DIR}/python/pactus-jsonrpc/pactus_jsonrpc/models.py
cp ${ROOT_DIR}/LICENSE ${PACKAGE_DIR}/python/pactus-jsonrpc
cp ${ROOT_DIR}/README.md ${PACKAGE_DIR}/python/pactus-jsonrpc
replace_in_place "s/{{ VERSION }}/$VERSION/g" ${PACKAGE_DIR}/python/pactus-jsonrpc/setup.py
