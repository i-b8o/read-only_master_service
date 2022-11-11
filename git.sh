#!/bin/bash
read -p "Commit description: " m
git add . && \
git add -u && \
git commit -m "$m" && \
git push origin master
