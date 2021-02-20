#!/bin/bash

MYDIR="$(dirname "$0")"

${MYDIR}/../rdv/translate.sh
go tool go2go translate ${MYDIR}/*.go2
