#!/bin/bash

# Required parameters:
# @raycast.schemaVersion 1
# @raycast.title Expand URL
# @raycast.mode silent

# Optional parameters:
# @raycast.icon 🤖

# Documentation:
# @raycast.description Follow URL redirects to final location
# @raycast.author William W. Marx
# @raycast.authorURL https://marx.sh

curl -sIL $(pbpaste) | rg -Po "(?<=location: ).+" | tail -1 | pbcopy
