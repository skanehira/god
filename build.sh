#!/bin/bash

# build MacOS
GOOS=darwin GOARCH=amd64 go build
zip MacOS.zip ./god && rm -rf ./god

# build Linux
GOOS=linux GOARCH=amd64 go build
zip Linux.zip ./god && rm -rf ./god

