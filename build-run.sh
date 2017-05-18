#!/bin/bash -e

set -o errexit
set -o nounset
set -o pipefail

go install
/go/bin/gocrud