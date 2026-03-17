#!/usr/bin/env bash

set -euo pipefail

if [[ $# -lt 1 ]]; then
	echo "Usage: $0 <command> [--steps N] [--version N] [--yes]"
	echo "Commands: up | down | steps | goto | force | version"
	exit 1
fi

COMMAND="$1"
shift

echo "[migrate] running command=${COMMAND} args=$*"

go run ./apps/migrator/cmd/migrator -command "${COMMAND}" "$@"

echo "[migrate] completed"
