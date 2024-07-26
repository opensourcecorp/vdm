#!/usr/bin/env bash
set -euo pipefail

failures=()

printf '> Running CI checks\n'

printf '>> Go vet\n'
if ! go vet ./... ; then
  printf '>>> Failed go-vet check\n'
  failures+=('go-vet')
fi

printf '>> Go linter\n'
if ! staticcheck ./... ; then
  printf '>>> Failed go-lint\n'
  failures+=('go-lint')
fi

printf '>> Go test\n'
if ! go test -cover -coverprofile=./cover.out ./... ; then
  printf '>>> Failed go-test check\n'
  failures+=('go-test')
fi

if [[ "${#failures[@]}" -gt 0 ]] ; then
  printf '> One or more checks failed, see output above\n'
  exit 1
fi

printf '> All checks passed!\n'
