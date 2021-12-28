#!/bin/sh
: > .gitignore
echo ".env" >> .gitignore
gibo dump Go Docker >> .gitignore
