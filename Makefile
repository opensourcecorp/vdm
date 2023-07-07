SHELL = /usr/bin/env bash -euo pipefail

BINNAME := vdm

# DO NOT TOUCH -- use `make bump-version oldversion=X newversion=Y`!
VERSION := 0.4.0

all: test

.PHONY: test
test: clean
	go vet ./...
	go test -cover -coverprofile=./cover.out ./...
	staticcheck ./...
	make -s clean

.PHONY: test-coverage
test-coverage: test
	go tool cover -html=./cover.out -o cover.html
	xdg-open ./cover.html

# builds for the current platform only
build: clean
	@go build -buildmode=pie -o build/$$(go env GOOS)-$$(go env GOARCH)/$(BINNAME) -ldflags "-s -w"

xbuild: clean
	@bash ./scripts/xbuild.sh

package: xbuild
	bash ./scripts/package.sh

package-debian: build
	@bash ./scripts/package-debian.sh

.PHONY: clean
clean:
	@rm -rf \
		/tmp/$(BINNAME)-tests \
		*cache* \
		.*cache* \
		./build/ \
		./dist/*.gz \
		./dist/debian/vdm.deb
	@sudo rm -rf ./dist/debian/vdm/usr
# TODO: until I sort out the tests to write test data consistently, these deps/
# directories etc. can kind of show up anywhere
	@find . -type d -name '*deps*' -exec rm -rf {} +
	@find . -type f -name '*VDMMETA*' -delete

bump-version: clean
	@if [[ -z "$(oldversion)" ]] || [[ -z "$(newversion)" ]] ; then printf 'ERROR: "oldversion" and "newversion" must be provided\n' && exit 1 ; fi
	find . \
		-type f \
		-not -path './go.*' \
		-exec sed -i 's/$(oldversion)/$(newversion)/g' {} \;

pre-commit-hook:
	cp ./scripts/ci.sh ./.git/hooks/pre-commit

# Some targets that help set up local workstations with rhad tooling. Assumes
# ~/.local/bin is on $PATH
add-local-symlinks:
	@mkdir -p "$${HOME}"/.local/bin
	@ln -fs $$(realpath build/$$(go env GOOS)-$$(go env GOARCH)/$(BINNAME)) "$${HOME}"/.local/bin/$(BINNAME)
	@printf 'Symlinked vdm to %s\n' "$${HOME}"/.local/bin/$(BINNAME)
