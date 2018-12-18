#! /bin/bash

set -e -o pipefail

IMAGE=${1-gluster/anthill}

# This sets the version variable to (hopefully) a semver compatible string. We
# expect released versions to have a tag of vX.Y.Z (with Y & Z optional), so we
# only look for those tags. For version info on non-release commits, we want to
# include the git commit info as a "build" suffix ("+stuff" at the end). There
# is also special casing here for when no tags match.
VERSION_GLOB="v[:digit:]*"
# Get the nearest "version" tag if one exists. If not, this returns the full
# git hash
NEAREST_TAG="$(git describe --always --tags --match "$VERSION_GLOB" --abbrev=0)"
# Full output of git describe for us to parse: TAG-<N>-g<hash>-<dirty>
FULL_DESCRIBE="$(git describe --always --tags --match "$VERSION_GLOB" --dirty)"
# If full matches against nearest, we found a valid tag earlier
if [[ $FULL_DESCRIBE =~ ${NEAREST_TAG}-(.*) ]]; then
        # Build suffix is the last part of describe w/ "-" replaced by "."
        version="$NEAREST_TAG+${BASH_REMATCH[1]//-/.}"
else
        # We didn't find a valid tag, so assume version 0 and everything ends up
        # in build suffix.
        version="0.0.0+g${FULL_DESCRIBE//-/.}"
fi
builddate="$(date -u '+%Y-%m-%dT%H:%M:%S.%NZ')"

operator-sdk build "$IMAGE"
docker tag "$IMAGE" "$IMAGE:base-image"
docker build \
    --build-arg from="$IMAGE:base-image" \
    --build-arg version="$version" \
    --build-arg builddate="$builddate" \
    -t "$IMAGE" \
    -f build/Dockerfile.stage2 \
    .
