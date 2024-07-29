SHELL = /usr/bin/env bash -euo pipefail

BINNAME := vdm

.PHONY: %

all: ci package package-debian

ci: clean
	@bash ./scripts/ci.sh

# test is just an alias for ci
test: ci

ci-container:
	@docker build -f ./Containerfile -t vdm-test:latest .

test-coverage: test
	go tool cover -html=./cover.out -o cover.html
	xdg-open ./cover.html

# builds for the current platform only
build: clean
	@go build -buildmode=pie -o build/$$(go env GOOS)-$$(go env GOARCH)/$(BINNAME) -ldflags "-s -w"

xbuild: clean
	@bash ./scripts/xbuild.sh

package: xbuild
	@bash ./scripts/package.sh

package-debian: build
	@bash ./scripts/package-debian.sh

clean:
	@rm -rf \
		/tmp/$(BINNAME)-tests \
		*cache* \
		.*cache* \
		./build/ \
		./dist/zipped/*.tar.gz \
		./dist/zipped/*.zip \
		./dist/debian/vdm.deb \
		*.out
	@sudo rm -rf ./dist/debian/vdm/usr
# TODO: until I sort out the tests to write test data consistently, these deps/
# directories etc. can kind of show up anywhere
	@find . -type d -name '*deps*' -exec rm -rf {} +
	@find . -type f -name '*VDMMETA*' -delete

bump-versions: clean
	@bash ./scripts/bump-versions.sh "$${old_version:-}"

tag-release: clean
	@bash ./scripts/tag-release.sh

pre-commit-hook:
	cp ./scripts/ci.sh ./.git/hooks/pre-commit

# Some targets that help set up local workstations with rhad tooling. Assumes
# ~/.local/bin is on $PATH
add-local-symlinks:
	@mkdir -p "$${HOME}"/.local/bin
	@ln -fs $$(realpath build/$$(go env GOOS)-$$(go env GOARCH)/$(BINNAME)) "$${HOME}"/.local/bin/$(BINNAME)
	@printf 'Symlinked vdm to %s\n' "$${HOME}"/.local/bin/$(BINNAME)
