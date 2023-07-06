#!/usr/bin/env bash
set -euo pipefail

# TODO: figure out cross-arch distribution for Debian at some point
printf 'Note: only supporting Debian builds for amd64 at the moment\n'

# We need several directories to exist, and lintian says they need to be owned
# by root, so instead of keeping them in the repo, we create them (and populate
# their contents) on the fly
mkdir -p \
  ./dist/debian/vdm/usr/bin \
  ./dist/debian/vdm/usr/share/doc/vdm \
  ./dist/debian/vdm/usr/share/man/man1
sudo chown -R 0:0 ./dist/debian/vdm/usr

sudo cp ./build/linux-amd64/vdm ./dist/debian/vdm/usr/bin/vdm
sudo cp ./LICENSE ./dist/debian/vdm/usr/share/doc/vdm/copyright
sudo sh -c 'gzip -9 -n -c ./CHANGELOG > ./dist/debian/vdm/usr/share/doc/vdm/changelog.gz'
sudo sh -c 'pandoc ./dist/man/man.1.md -s -t man | gzip -9 -n -c - > ./dist/debian/vdm/usr/share/man/man1/vdm.1.gz'

# Actually build the debfile
dpkg-deb --build ./dist/debian/vdm

# Ask lintian (the Debian package linter) to scream as loudly as it can about
# anything it finds
lintian --info --tag-display-limit=0 ./dist/debian/vdm.deb
