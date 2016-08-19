#!/bin/sh

if [ ! -z "$(git status --porcelain)" ]
then
  echo "Git is not clean"
  git status
fi
