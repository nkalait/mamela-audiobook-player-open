#!/bin/bash

RPATH="./lib/mac"
TARGET="tester"

go test -c -o $TARGET

install_name_tool -change @loader_path/libbass.dylib @loader_path/lib/mac/libbass.dylib $TARGET

./$TARGET