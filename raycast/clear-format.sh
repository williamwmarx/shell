#!/bin/bash

# Required parameters:
# @raycast.schemaVersion 1
# @raycast.title Clear format
# @raycast.mode silent

# Optional parameters:
# @raycast.icon 🤖
# @raycast.packageName TextTools

# Documentation:
# @raycast.description Clears the format of text while maintaining spacing
# @raycast.author William W. Marx
# @raycast.authorURL https://marx.sh

pbpaste | pbcopy
