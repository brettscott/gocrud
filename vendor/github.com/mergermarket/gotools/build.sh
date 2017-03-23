#!/bin/bash -e

set -e
set -u

source "./build-functions.sh"

buildCode
fmtCode
lintCode
runTests
