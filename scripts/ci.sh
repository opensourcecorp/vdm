#!/usr/bin/env bash
set -euo pipefail

failures=()

printf '> Running CI checks\n'

printf '>> Go vet\n'
if ! go vet ./... ; then
  printf '>>> Failed go-vet check\n' > /dev/stderr
  failures+=('go-vet')
fi

printf '>> Go linter\n'
if ! go run github.com/mgechev/revive@latest --set_exit_status ./... ; then
  printf '>>> Failed go-lint\n' > /dev/stderr
  failures+=('go-lint')
fi

printf '>> Go error checker\n'
if ! go run github.com/kisielk/errcheck@latest ./... ; then
  printf '>>> Failed go-error-check\n' > /dev/stderr
  failures+=('go-error-check')
fi

printf '>> Go test\n'
if ! go test -cover -coverprofile=./cover.out ./... ; then
  printf '>>> Failed go-test check\n' > /dev/stderr
  failures+=('go-test')
fi

printf '>> Packaging checker\n'
if ! make -s package ; then
  printf '>>> Failed packaging check\n' > /dev/stderr
  failures+=('packaging')
fi

printf '>> Version tag checker\n'
if ! make -s tag-release ; then
  printf '>>> Failed verion tag checker\n' > /dev/stderr
  failures+=('version-tag')
fi

if [[ "${#failures[@]}" -gt 0 ]] ; then
  printf '> One or more checks failed, see output above\n' > /dev/stderr
  exit 1
fi

printf '> All checks passed!\n'
