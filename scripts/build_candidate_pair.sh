#!/usr/bin/env bash
# Build and validate a local Linux pair. Never logs in or publishes registry tags.
set -euo pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$ROOT"
SOURCE_COMMIT=$(git rev-parse HEAD)
[[ -z "$(git status --porcelain --untracked-files=all)" ]] || {
  echo 'Commit the complete candidate tree before building' >&2; exit 1;
}
OUTPUT_DIR="${1:?usage: bash scripts/build_candidate_pair.sh /absolute/evidence-directory}"
[[ "$OUTPUT_DIR" = /* ]] || { echo 'Evidence directory must be absolute' >&2; exit 1; }
case "$OUTPUT_DIR/" in "$ROOT/"*) echo 'Evidence must be outside the source tree' >&2; exit 1;; esac
# Check the local build tool before downloading dependencies or images.
DOCKER=(docker)
[[ -z "${CANDIDATE_DOCKER_CONTEXT:-}" ]] || DOCKER+=(--context "$CANDIDATE_DOCKER_CONTEXT")
"${DOCKER[@]}" buildx version >/dev/null || {
  echo 'Docker Buildx is required; no compilation has started' >&2; exit 1;
}
# Ignore inherited publishing switches entirely; this script invokes Buildx directly.
GO_VERSION=$(awk '$1 == "go" { print $2; exit }' go.mod)
source scripts/lib_go_module_cache.sh
prepare_go_module_cache "$ROOT"
mkdir -p "$OUTPUT_DIR"
# Resolve exact base digests once, so node and plugin use the same Linux toolchain.
resolve_image() {
  "${DOCKER[@]}" pull --platform "$1" "$2" >&2
  "${DOCKER[@]}" image inspect "$2" --format '{{index .RepoDigests 0}}'
}
BUILD_ARCH=$("${DOCKER[@]}" info --format '{{.Architecture}}')
case "$BUILD_ARCH" in aarch64|arm64) BUILD_ARCH=arm64;; x86_64|amd64) BUILD_ARCH=amd64;;
  *) echo "Unsupported builder: $BUILD_ARCH" >&2; exit 1;; esac
GO_IMAGE=$(resolve_image "linux/$BUILD_ARCH" "golang:$GO_VERSION-bookworm")
RUNTIME_IMAGE=$(resolve_image linux/amd64 debian:12-slim)
IMAGE="qa53-candidate:$SOURCE_COMMIT"
"${DOCKER[@]}" buildx build --load --platform linux/amd64 \
  --build-context "gomodcache=$(go_module_proxy_cache)" \
  --build-arg "GO_IMAGE=$GO_IMAGE" --build-arg "RUNTIME_IMAGE=$RUNTIME_IMAGE" \
  --build-arg "SOURCE_COMMIT=$SOURCE_COMMIT" \
  --metadata-file "$OUTPUT_DIR/build-metadata.json" \
  -f Dockerfile.candidate -t "$IMAGE" .
RUN=("${DOCKER[@]}" run --rm --platform linux/amd64 --network none --read-only
     --cap-drop ALL --security-opt no-new-privileges --memory 512m --cpus 1 --pids-limit 64
     --tmpfs /tmp:rw,noexec,nosuid,size=16m)
"${RUN[@]}" "$IMAGE" --version > "$OUTPUT_DIR/node-version.txt"
"${RUN[@]}" --entrypoint /plugins-template/subnet-evm "$IMAGE" --version > "$OUTPUT_DIR/plugin-version.txt"
"${RUN[@]}" --entrypoint /usr/local/bin/candidate-handshake "$IMAGE" /plugins-template/subnet-evm > "$OUTPUT_DIR/handshake.txt"
python3 scripts/candidate/verify_versions.py "$SOURCE_COMMIT" "$OUTPUT_DIR"
"${DOCKER[@]}" image inspect "$IMAGE" \
  --format '{{json .}}' > "$OUTPUT_DIR/image-inspect.json"
printf '%s\n' "$GO_IMAGE" > "$OUTPUT_DIR/go-base.txt"
printf '%s\n' "$RUNTIME_IMAGE" > "$OUTPUT_DIR/runtime-base.txt"
printf 'PASS: Linux candidate pair %s; evidence %s\n' "$IMAGE" "$OUTPUT_DIR"
