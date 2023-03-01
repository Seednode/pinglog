#!/usr/bin/env bash

package_name="pinglog"
mkdir -p builds

platforms=(
  "android/arm64"
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
  "netbsd/386"
  "netbsd/amd64"
  "netbsd/arm"
  "netbsd/arm64"
  "openbsd/386"
  "openbsd/amd64"
  "openbsd/arm"
  "openbsd/arm64"
  "windows/386"
  "windows/amd64"
  "windows/arm"
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
  elif [ "${GOOS}" == "linux" ] && [ "${GOARCH}" == "amd64" ]; then
    ld_flags+=' -linkmode external -extldflags "-static"'
  fi
  env GOOS="${GOOS}" GOARCH="${GOARCH}" CC="musl-gcc" go build -ldflags "${ld_flags}" -o "builds/${output_name}"
done
