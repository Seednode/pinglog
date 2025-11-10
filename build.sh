#!/usr/bin/env bash

package_name="pinglog"
mkdir -p builds

platforms=(
  "aix/ppc64"
  "darwin/amd64"
  "darwin/arm64"
  "dragonfly/amd64"
  "freebsd/386"
  "freebsd/amd64"
  "freebsd/arm"
  "freebsd/arm64"
  "linux/386"
  "linux/amd64"
  "linux/arm"
  "linux/arm64"
  "linux/loong64"
  "linux/mips"
  "linux/mips64"
  "linux/mips64le"
  "linux/mipsle"
  "linux/ppc64"
  "linux/ppc64le"
  "linux/riscv64"
  "linux/s390x"
  "netbsd/386"
  "netbsd/amd64"
  "netbsd/arm"
  "netbsd/arm64"
  "windows/386"
  "windows/amd64"
  "windows/arm64"
)

for platform in "${platforms[@]}"; do
  IFS=" " read -r -a platform_split <<< "${platform//\// }"
  GOOS="${platform_split[0]}"
  GOARCH="${platform_split[1]}"
  output_name="${package_name}-${GOOS}-${GOARCH}"
  ld_flags='-s -w'
  if [ "${GOOS}" == "windows" ]; then
    output_name+=".exe"
  fi
  env GOOS="${GOOS}" GOARCH="${GOARCH}" CC="musl-gcc" CGO_ENABLED=0 go build -trimpath -ldflags "${ld_flags}" -tags "netgo timetzdata" -o "builds/${output_name}"
done
