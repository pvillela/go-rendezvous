#!/bin/bash

MYDIR="$(dirname "$0")"

${MYDIR}/../async/translate.sh
# go tool go2go translate ${MYDIR}/*.go2
