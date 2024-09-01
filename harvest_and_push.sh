#!/bin/bash

export $(grep -v '^#' /app/.env | xargs)

cd /app

./harvester

git config --global user.name "$GIT_COMMITTER_NAME"
git config --global user.email "$GIT_COMMITTER_EMAIL"

cd public

git init .

git remote add origin https://$GITHUB_TOKEN@github.com/$GITHUB_USERNAME/$GITHUB_REPOSITORY.git
git remote set-url origin https://$GITHUB_TOKEN@github.com/$GITHUB_USERNAME"/$GITHUB_REPOSITORY.git

git checkout -b main

git add .

git commit -m "Automated data harvest $(date)"

git push --force origin main

exit 0
