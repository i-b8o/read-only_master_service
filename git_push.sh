#!/bin/bash
echo "supreme push"
git add . && \
git add -u && \
git commit -m $1 && \
git push origin master
