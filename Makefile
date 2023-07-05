SHELL = /usr/bin/env bash -euo pipefail

BINNAME := vdm

all: test clean

.PHONY: %

test: clean
	go vet ./...
	go test -cover -coverprofile=./cover.out ./...
	staticcheck ./...
	make -s clean

test-coverage: test
	go tool cover -html=./cover.out -o cover.html
	xdg-open ./cover.html

build: clean
	@mkdir -p build/$$(go env GOOS)-$$(go env GOARCH)
	@go build -o build/$$(go env GOOS)-$$(go env GOARCH)/$(BINNAME) -ldflags "-s -w"

xbuild: clean
	@for target in \
		darwin-amd64 \
		linux-amd64 \
		linux-arm \
		linux-arm64 \
		windows-amd64 \
	; \
	do \
		GOOS=$$(echo "$${target}" | cut -d'-' -f1) ; \
		GOARCH=$$(echo "$${target}" | cut -d'-' -f2) ; \
		export GOOS GOARCH ; \
		outdir=build/"$${GOOS}-$${GOARCH}" ; \
		mkdir -p "$${outdir}" ; \
		printf "Building for %s-%s into build/ ...\n" "$${GOOS}" "$${GOARCH}" ; \
		go build -o "$${outdir}"/$(BINNAME) -ldflags "-s -w" ; \
	done

package: xbuild
	@mkdir -p dist
	@cd build || exit 1; \
	for built in * ; do \
		printf 'Packaging for %s into dist/ ...\n' "$${built}" ; \
		cd $${built} && tar -czf ../../dist/$(BINNAME)_$${built}.tar.gz * && cd - >/dev/null ; \
	done

clean:
	@rm -rf \
		/tmp/$(BINNAME)-tests \
		*cache* \
		.*cache* \
		./build/ \
		./dist/*.gz
# TODO: until I sort out the tests to write test data consistently, these deps/
# directories etc. can kind of show up anywhere
	@find . -type d -name '*deps*' -exec rm -rf {} +
	@find . -type f -name '*VDMMETA*' -delete

pre-commit-hook:
	cp ./scripts/ci.sh ./.git/hooks/pre-commit

# Some targets that help set up local workstations with rhad tooling. Assumes
# ~/.local/bin is on $PATH
add-local-symlinks:
	@mkdir -p "$${HOME}"/.local/bin
	@ln -fs $$(realpath build/$$(go env GOOS)-$$(go env GOARCH)/$(BINNAME)) "$${HOME}"/.local/bin/$(BINNAME)
	@printf 'Symlinked vdm to %s\n' "$${HOME}"/.local/bin/$(BINNAME)
