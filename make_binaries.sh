#!/usr/bin/env bash
set -euo pipefail

ARTIFACT="bart"
RELEASE_INFO="release_info.md"
GIT_REVISION=$(git rev-parse HEAD)
GO_VERSION=$(go version)

rm -f \
  "./$RELEASE_INFO" \
  "./$ARTIFACT" \
  "./$ARTIFACT.exe" \
  "./${ARTIFACT}_macos_x64.tar.gz" \
  "./${ARTIFACT}_linux_x64.tar.gz" \
  "./${ARTIFACT}_windows_x64.zip"

# $1 File to be hashed
# $2 Target filename
function generate_digests {
local SHA_256_DIGEST=$(sha256sum $1)
local SHA_512_DIGEST=$(sha512sum $1)
cat >> "$2" <<EOF
SHA256 digest:
\`\`\`
$SHA_256_DIGEST
\`\`\`
SHA512 digest:
\`\`\`
$SHA_512_DIGEST
\`\`\`
EOF
}

cat >> "$RELEASE_INFO" <<EOF
# \`$ARTIFACT\` binaries

git revision:

\`\`\`
$GIT_REVISION
\`\`\`

Go compiler version:

\`\`\`
$GO_VERSION
\`\`\`

## Linux x64
EOF
make xcompile_linux
generate_digests "$ARTIFACT" "$RELEASE_INFO"
tar -cvzf "${ARTIFACT}_linux_x64.tar.gz" "$ARTIFACT"

printf "\n## macOS x64\n\n" >> "$RELEASE_INFO"
make xcompile_mac
generate_digests "$ARTIFACT" "$RELEASE_INFO"
tar -cvzf "${ARTIFACT}_macos_x64.tar.gz" "$ARTIFACT"

printf "\n## Windows x64\n\n" >> "$RELEASE_INFO"
make xcompile_win
generate_digests "$ARTIFACT.exe" "$RELEASE_INFO"
zip "${ARTIFACT}_windows_x64.zip" "$ARTIFACT.exe"
