#!/usr/bin/env bash

package_name="pinglog"
mkdir -p builds

platforms=(
  "darwin/amd64"
  "darwin/arm64"
  "linux/386"
  "linux/amd64"
  "linux/arm"
  "linux/arm64"
  "windows/386"
  "windows/amd64"
)

for platform in ${platforms[@]}; do
  platform_split=(${platform//\// })
  GOOS="${platform_split[0]}"
  GOARCH="${platform_split[1]}"
  output_name="${package_name}-${GOOS}-${GOARCH}"
  if [ "${GOOS}" == "windows" ]; then
    output_name+=".exe"
  fi
  env GOOS="${GOOS}" GOARCH="${GOARCH}" go build -ldflags "-s -w" -o "builds/${output_name}"
  upx -q --brute "builds/${output_name}"
done
