#!/usr/bin/env bash

len=`expr "$msg" : '.*'`

. tools/pre-commit/var.sh
. tools/pre-commit/msg.sh
. tools/pre-commit/branch-name.sh

checkExcept
checkFirstLetter
checkLengthMsg

