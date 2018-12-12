#!/bin/bash

set -e

export GOPATH="$(mktemp -d)"

# Runc is built from the 17.06 branch in https://github.com/docker/runc
git clone git://github.com/docker/runc.git "$GOPATH/src/github.com/opencontainers/runc"
cd "$GOPATH/src/github.com/opencontainers/runc"
git checkout -q "$RUNC_COMMIT"
make BUILDTAGS="seccomp apparmor selinux"
sudo make install
