#!/bin/bash

MYDIR="$(dirname "$0")"

go tool go2go translate ${MYDIR}/*.go2
