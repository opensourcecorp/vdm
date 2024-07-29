#!/usr/bin/env bash
set -euo pipefail

################################################################################
# A singular, consistent version needs to match across several places across the
# codebase tree, such as the Debian Control file & man pages, the embedded
# version listed for the CLI command, etc. This script is used to bump version
# identifiers in any locations listed in the loop at the bottom.
################################################################################

root="$(git rev-parse --show-toplevel)"

old_version="${1:-}"
new_version="$(grep -v '#' "${root}"/VERSION)"

if [[ -z "${old_version:-}" ]] ; then
  printf 'ERROR: you must specify old_version as the first script argument\n'
  exit 1
fi

if [[ -z "${new_version:-}" ]] ; then
  printf 'ERROR: version unable to be determined from ./VERSION file; possible malformed?\n'
  exit 1
fi

for f in \
  dist/debian/vdm/DEBIAN/control \
  dist/man/* \
  cmd/root.go \
; do
  if grep -q "${old_version}" "${root}/${f}" ; then
    printf 'Updating version in %s: %s -> %s\n' "${f}" "${old_version}" "${new_version}"
    sed -i "s/${old_version}/${new_version}/g" "${root}/${f}"
  fi
done
