#!/usr/bin/env bash
set -euo pipefail

# cd-jumps because it makes logs cleaner, not sorry
mkdir -p dist/compressed
cd build || exit 1
for built in * ; do
  printf 'Packaging for %s into dist/compressed/\n' "${built}"
  cd "${built}" || exit 1
  # Windows might like .zips better, otherwise make .tar.gzs
  if [[ "${built}" =~ 'windows' ]] ; then
    zip -r9 ../../dist/compressed/vdm_"${built}".zip ./*
  else
    tar -czf ../../dist/compressed/vdm_"${built}".tar.gz ./*
  fi
  cd - > /dev/null
done
