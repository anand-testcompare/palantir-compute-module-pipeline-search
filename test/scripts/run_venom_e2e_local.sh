#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

VENOM_VERSION="v1.3.0"
VENOM_SHA256="89832ec25e820c605cf0d3c09122e60bad43d13c1724aa6d375ef7109fbfe201"

VENOM_BIN="${ROOT_DIR}/bin/venom"

mkdir -p "${ROOT_DIR}/bin" "${ROOT_DIR}/out/venom"

if [[ ! -x "${VENOM_BIN}" ]]; then
  curl -fsSL -o "${VENOM_BIN}" "https://github.com/ovh/venom/releases/download/${VENOM_VERSION}/venom.linux-amd64"
  echo "${VENOM_SHA256}  ${VENOM_BIN}" | sha256sum -c -
  chmod +x "${VENOM_BIN}"
fi

rm -rf "${ROOT_DIR}/out/venom"
mkdir -p "${ROOT_DIR}/out/venom"

"${VENOM_BIN}" run "${ROOT_DIR}/test/venom/enricher_e2e.yml" --format=xml --output-dir="${ROOT_DIR}/out/venom" --html-report

echo "wrote ${ROOT_DIR}/out/venom/test_results.html"

