#!/bin/bash

MYDIR="$(dirname "$0")"

${MYDIR}/../rdv/clean.sh
rm ${MYDIR}/*.go
