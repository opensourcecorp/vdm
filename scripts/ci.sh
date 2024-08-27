#!/usr/bin/env bash
set -euo pipefail

failures=()

printf '> Running CI checks\n'

printf '>> Go vet\n'
if ! go vet ./... ; then
  printf '>>> Failed go-vet check\n' > /dev/stderr
  failures+=('go-vet')
fi

printf '>> Go linter (staticcheck)\n'
if ! go run honnef.co/go/tools/cmd/staticcheck@latest ./... ; then
  printf '>>> Failed go-lint-staticcheck\n' > /dev/stderr
  failures+=('go-lint-staticcheck')
fi

printf '>> Go linter (revive)\n'
if ! go run github.com/mgechev/revive@latest --set_exit_status ./... ; then
  printf '>>> Failed go-lint-revive\n' > /dev/stderr
  failures+=('go-lint-revive')
fi

printf '>> Go linter (errcheck)\n'
if ! go run github.com/kisielk/errcheck@latest ./... ; then
  printf '>>> Failed go-lint-errcheck\n' > /dev/stderr
  failures+=('go-lint-errcheck')
fi

printf '>> Go test\n'
go clean -testcache
if ! go test -cover -coverprofile=./cover.out ./... ; then
  printf '>>> Failed go-test check\n' > /dev/stderr
  failures+=('go-test')
fi

printf '>> Packaging checker\n'
if ! make -s package ; then
  printf '>>> Failed packaging check\n' > /dev/stderr
  failures+=('packaging')
fi

if [[ "${#failures[@]}" -gt 0 ]] ; then
  printf '> The following checks failed -- logs for each are above\n' > /dev/stderr
  printf '%s\n' "${failures[@]}"
  exit 1
fi

printf '> All checks passed!\n'
