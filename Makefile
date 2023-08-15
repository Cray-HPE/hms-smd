# MIT License
#
# (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP
#
# Permission is hereby granted, free of charge, to any person obtaining a
# copy of this software and associated documentation files (the "Software"),
# to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included
# in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
# THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.

# Service
NAME ?= cray-smd
VERSION ?= $(shell cat .version)

all: image unittest ct snyk ct_image

.PHONY : all image unittest snyk ct ct_image binaries coverage

image:
	docker build ${NO_CACHE} --pull ${DOCKER_ARGS} --tag '${NAME}:${VERSION}' -f Dockerfile .

unittest:
	./runUnitTest.sh

snyk:
	./runSnyk.sh

ct:
	./runCT.sh

ct_image:
	docker build --no-cache -f test/ct/Dockerfile test/ct/ --tag hms-smd-test:${VERSION} 

binaries: smd smd-init smd-loader



BUILD := `git rev-parse HEAD --short`
VERSION := `git describe --tags --abbrev=0`
LDFLAGS=-ldflags "-X=$(GIT)build.Build=$(BUILD) -X=$(GIT)build.Version=$(VERSION)"

smd:
	go build -v $(LDFLAGS) ./cmd/smd

smd-init:
	go build -v $(LDFLAGS) ./cmd/smd-init

smd-loader:
	go build -v $(LDFLAGS) ./cmd/smd-loader

coverage:
	go test -cover -v -tags musl ./cmd/* ./internal/* ./pkg/*