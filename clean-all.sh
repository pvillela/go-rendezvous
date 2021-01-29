#!/bin/bash

MYDIR="$(dirname "$0")"

${MYDIR}/async/clean.sh
${MYDIR}/runconcurrent/clean.sh
${MYDIR}/util/clean.sh
