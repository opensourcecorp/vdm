SHELL = /usr/bin/env bash -euo pipefail

BINNAME := vdm

all: test clean

.PHONY: %

test: clean
	@go vet ./...
	@go test -cover ./...

build: clean
	@mkdir -p build/$$(go env GOOS)-$$(go env GOARCH)
	@go build -o build/$$(go env GOOS)-$$(go env GOARCH)/$(BINNAME)

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
		outdir=build/"$${GOOS}-$${GOARCH}" ; \
		mkdir -p "$${outdir}" ; \
		printf "Building for %s-%s into build/ ...\n" "$${GOOS}" "$${GOARCH}" ; \
		GOOS="$${GOOS}" GOARCH="$${GOARCH}" go build -o "$${outdir}"/$(BINNAME) ; \
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
		./dist/ \
		./deps/ \
		./testdata/deps/

# Some targets that help set up local workstations with rhad tooling. Assumes
# ~/.local/bin is on $PATH
add-local-symlinks:
	@mkdir -p "$${HOME}"/.local/bin
	@ln -fs $$(realpath build/$$(go env GOOS)-$$(go env GOARCH)/$(BINNAME)) "$${HOME}"/.local/bin/$(BINNAME)
	@printf 'Symlinked vdm to %s\n' "$${HOME}"/.local/bin/$(BINNAME)
