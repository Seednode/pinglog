#!/usr/bin/env bash
# build, tag, and push docker images

# exit if a command fails
set -o errexit

# go docker image tag to use
tag="${TAG:-latest}"

# if no registry is provided, tag image as "local" registry
registry="${REGISTRY:-local}"

# set image name
image_name="pinglog"

# set image version
image_version="$(grep "ReleaseVersion" main.go | head -n1 | awk '{print $4}' | sed 's/\"//g')"

# platforms to build for
platforms="linux/amd64"
platforms+=",linux/arm"
platforms+=",linux/arm64"
platforms+=",linux/ppc64le"

# copy native image to local image repository
docker buildx build \
                    --build-arg TAG="${tag}" \
                    -t "${registry}/${image_name}:${image_version}" \
                    $(if [ "${LATEST}" == "yes" ]; then echo "-t ${registry}/${image_name}:latest"; fi) \
                    -f docker/Dockerfile . \
                    --load

# push image to remote registry
docker buildx build --platform "${platforms}" \
                    --build-arg TAG="${tag}" \
                    -t "${registry}/${image_name}:${image_version}" \
                    $(if [ "${LATEST}" == "yes" ]; then echo "-t ${registry}/${image_name}:latest"; fi) \
                    -f docker/Dockerfile . \
                    --push

# copy debug image to local image repository
docker buildx build \
                    --build-arg TAG="${tag}" \
                    -t "${registry}/${image_name}:${image_version}-debug" \
                    $(if [ "${LATEST}" == "yes" ]; then echo "-t ${registry}/${image_name}:debug"; fi) \
                    -f docker/Dockerfile.debug . \
                    --load

# push debug image to remote registry
docker buildx build --platform "${platforms}" \
                    --build-arg TAG="${tag}" \
                    -t "${registry}/${image_name}:${image_version}-debug" \
                    $(if [ "${LATEST}" == "yes" ]; then echo "-t ${registry}/${image_name}:debug"; fi) \
                    -f docker/Dockerfile.debug . \
                    --push
