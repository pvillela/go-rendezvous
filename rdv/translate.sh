#!/bin/bash

#
#  Copyright © 2021 Paulo Villela. All rights reserved.
#  Use of this source code is governed by the Apache 2.0 license
#  that can be found in the LICENSE file.
#

MYDIR="$(dirname "$0")"

${MYDIR}/../util/translate.sh
go tool go2go translate ${MYDIR}/*.go2
