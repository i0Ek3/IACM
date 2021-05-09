#!/bin/bash

usage() {
    echo "usage:"
    echo "  $0 project_name [init|tidy]"
    echo "  $0 big0 init"
}

if [ $1 -o $2 -e ] 
then
    usage
else
    if [ $2 == "init" ] 
    then
        go mod init $2
    else
        go mod tidy 
    fi
fi
