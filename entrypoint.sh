#!/bin/bash

# Leverage the default env variables as described in:
# https://docs.github.com/en/actions/reference/environment-variables#default-environment-variables
if [[ $GITHUB_ACTIONS != "true" ]]
then
  /usr/bin/sato "$@"
  exit $?
fi

flags=""

echo "running command:"
echo sato parse -f "$INPUT_FILE" "$flags"

/usr/bin/sato parse -f "$INPUT_FILE" "$flags"
export sato_EXIT_CODE=$?
