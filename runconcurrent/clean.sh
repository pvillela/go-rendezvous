#!/bin/bash

MYDIR="$(dirname "$0")"

${MYDIR}/../errorof/clean.sh
rm ${MYDIR}/*.go
