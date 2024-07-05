#!/usr/bin/env bash

function checkExcept () {
    if [[ $branchName =~ ^($except)$ ]]
    then
        exit 0
    fi
}


function checkBranchName () {
    if [[ $branchName =~ ^($feature)\/$code+\-+$number+\_+.* ]]
    then
        exit 0
    else
        IFS='|' read -r -a features <<< "$feature"
        {
            echo "💩 wrong branch name"
            echo "👉 using pattern: <branch-type>/<ISSUE-CODE>_<task-title>"
        }
        for element in "${features[@]}"
        do
            echo "- $element"
        done
        exit 1
    fi
}

