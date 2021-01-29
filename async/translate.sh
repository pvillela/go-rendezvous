#!/bin/bash

MYDIR="$(dirname "$0")"

go tool go2go translate ${MYDIR}/../errorof/*.go2
go tool go2go translate ${MYDIR}/../tuple/*.go2
go tool go2go translate ${MYDIR}/*.go2
