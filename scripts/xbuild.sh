#!/usr/bin/env bash
set -euo pipefail

targets=$(go tool dist list | grep -E 'linux|windows|darwin' | grep -E 'amd64|arm64')
printf 'Will build for:\n'
while read -r line ; do
  printf '\t%s\n' "${line}"
done <<< "${targets}"

for target in ${targets} ; do
  GOOS=$(echo "${target}" | cut -d'/' -f1)
  GOARCH=$(echo "${target}" | cut -d'/' -f2)
  export GOOS GOARCH
  outdir=build/"${GOOS}-${GOARCH}"
  mkdir -p "${outdir}"
  printf "Building for %s-%s into build/ ...\n" "${GOOS}" "${GOARCH}"
  go build -buildmode=pie -o "${outdir}"/vdm -ldflags "-s -w"
done
