#!/bin/bash

MYDIR="$(dirname "$0")"

${MYDIR}/run.sh
${MYDIR}/clean.sh
${MYDIR}/../errorof/clean.sh
