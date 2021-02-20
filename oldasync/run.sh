#!/bin/bash

MYDIR="$(dirname "$0")"

${MYDIR}/translate.sh
go tool go2go run ${MYDIR}/cmd/*.go2
