#!/usr/bin/env bash

# checking the first letter of message
function checkFirstLetter () {
    if [[ $msg =~ ^[A-Z]{1}.*[[:alnum:]]$ ]]
    then
        echo "pass"
    else
        echo "❌ message must be start with uppercase"
        exit 1
    fi
}
# checking length of message
function checkLengthMsg () {
    if [ $len -lt 100 ]
    then
        echo "pass"
        exit 0
    fi
    echo "❌ message must be less than 100 characters"
    echo "❌ message length $len"
    exit 1
}
