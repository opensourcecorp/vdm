#!/usr/bin/env bash
set -euo pipefail

# cd-jumps because it makes logs cleaner, not sorry
mkdir -p dist
cd build || exit 1
for built in * ; do
  printf 'Packaging for %s into dist/ ...\n' "${built}"
  cd "${built}" && tar -czf ../../dist/vdm_"${built}".tar.gz ./*
  cd - > /dev/null
done
