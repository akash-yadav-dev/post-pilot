#!/usr/bin/env bash
set -euo pipefail

echo "=============================="
echo "🚀 Running Local CI Pipeline"
echo "=============================="

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# Ensure Go bin tools are available
export PATH="$PATH:$(go env GOPATH)/bin"

echo ""
echo "📦 Verifying Go modules..."
go mod tidy
git diff --exit-code go.mod go.sum

echo ""
echo "🔍 Running go vet..."
go vet ./...

echo ""
echo "🧹 Running Go linter..."
bash ./scripts/lint.sh


echo ""
echo "🔐 Running gosec security scan..."

GOBIN="$(go env GOPATH)/bin"

if ! command -v gosec &> /dev/null; then
    echo "Installing gosec..."
    go install github.com/securego/gosec/v2/cmd/gosec@latest
fi

"$GOBIN/gosec" ./...


echo ""
echo "🛡 Running govulncheck..."

if ! command -v govulncheck &> /dev/null; then
    echo "Installing govulncheck..."
    go install golang.org/x/vuln/cmd/govulncheck@latest
fi

"$GOBIN/govulncheck" ./...


echo ""
echo "🧪 Running Go tests..."
if [[ "$(go env CGO_ENABLED)" == "1" ]]; then
    go test -v -race -coverprofile=coverage.out ./apps/... ./packages/...
else
    echo "CGO is disabled; running tests without -race"
    go test -v -coverprofile=coverage.out ./apps/... ./packages/...
fi

echo ""
echo "⚙️ Building Go binaries..."
mkdir -p build

go build -v -o build/api ./apps/api/cmd/server
go build -v -o build/worker ./apps/worker/cmd/worker


echo ""
echo "🌐 Running Frontend Checks..."
cd apps/web

echo "Installing dependencies..."
npm ci

echo "Running frontend lint..."
npm run lint

echo "Building frontend..."
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080 npm run build

cd "$ROOT_DIR"

echo ""
echo "✅ Local CI Passed Successfully!"
echo "=============================="