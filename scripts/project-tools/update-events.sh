#!/bin/bash

currentPath=$(pwd)
echo "$currentPath"
for f in ./events/*; do
    if [ -d "${f}" ]; then
        if [ -d "${f}/.git" ]; then
          printf "%s\n" "${f}"
          cd "${f}" && git pull && cd "$currentPath" || exit
        fi;
    fi
done