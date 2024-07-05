#!/usr/bin/env bash

. tools/pre-commit/branch-name.sh
. tools/pre-commit/var.sh

checkExcept
checkBranchName

