#!/bin/bash

MYDIR="$(dirname "$0")"

${MYDIR}/../util/clean.sh
rm ${MYDIR}/*.go
