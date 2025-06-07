#!/bin/bash

cleanup(){

    echo "${BOLD}${RED} CLEANING UP: ${NC} Merged branches into main"
    git branch --merged main | grep -v '^\*' | grep -v main | xargs -n 1 git branch -d
    echo -e "\t${TICK} Branches Cleaned."
}