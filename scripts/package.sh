#!/usr/bin/env bash
set -euo pipefail

# cd-jumps because it makes logs cleaner, not sorry
mkdir -p dist/zipped
cd build || exit 1
for built in * ; do
  printf 'Packaging for %s into dist/zipped/\n' "${built}"
  cd "${built}" || exit 1
  # Windows might like .zips better, otherwise make .tar.gzs
  if [[ "${built}" =~ 'windows' ]] ; then
    zip -r9 ../../dist/zipped/vdm_"${built}".zip ./*
  else
    tar -czf ../../dist/zipped/vdm_"${built}".tar.gz ./*
  fi
  cd - > /dev/null
done
