#!/usr/bin/env bash

set -euo pipefail

echo "[lint] go version: $(go version)"

echo "[lint] running go vet"
go vet ./...

echo "[lint] installing staticcheck"
go install honnef.co/go/tools/cmd/staticcheck@latest

echo "[lint] running staticcheck"
"$(go env GOPATH)/bin/staticcheck" ./...

echo "[lint] completed successfully"