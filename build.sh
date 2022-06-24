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
  ld_flags='-s -w'
  if [ "${GOOS}" == "windows" ]; then
    output_name+=".exe"
  elif [ "${GOOS}" == "linux" ] && [ "${GOARCH}" == "amd64" ]; then
    ld_flags+=' -linkmode external -extldflags "-static"'
  fi
  env GOOS="${GOOS}" GOARCH="${GOARCH}" CC="/home/sinc/.nix-profile/bin/musl-gcc" go build -ldflags "${ld_flags}" -o "builds/${output_name}"
done
