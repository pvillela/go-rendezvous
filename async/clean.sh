#!/bin/bash

MYDIR="$(dirname "$0")"

${MYDIR}/../errorof/clean.sh
${MYDIR}/../tuple/clean.sh
rm ${MYDIR}/*.go
