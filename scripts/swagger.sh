#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ROUTES_FILE="$ROOT_DIR/apps/api/routes/routes.go"
OUTPUT_DIR="$ROOT_DIR/docs/swagger"
OUTPUT_FILE="$OUTPUT_DIR/swagger.json"
BASE_URL="${SWAGGER_BASE_URL:-http://localhost:8080}"

if [[ ! -f "$ROUTES_FILE" ]]; then
  echo "routes file not found: $ROUTES_FILE" >&2
  exit 1
fi

mkdir -p "$OUTPUT_DIR"

tmp_routes="$(mktemp)"
trap 'rm -f "$tmp_routes"' EXIT

awk '
BEGIN {
  prefix["router"] = ""
}
{
  if (match($0, /([A-Za-z0-9_]+)[[:space:]]*:=[[:space:]]*([A-Za-z0-9_]+)\.Group\("([^"]*)"\)/, m)) {
    parent = m[2]
    segment = m[3]
    base = (parent in prefix) ? prefix[parent] : ""
    prefix[m[1]] = base segment
  }

  if (match($0, /([A-Za-z0-9_]+)\.(GET|POST|PUT|PATCH|DELETE|OPTIONS|HEAD)\("([^"]*)"/, m)) {
    group = m[1]
    method = tolower(m[2])
    path = m[3]
    base = (group in prefix) ? prefix[group] : ""
    route = base path
    gsub(/\/\/+/, "/", route)
    print route "|" method
  }
}
' "$ROUTES_FILE" >"$tmp_routes"

# Include health endpoint from router bootstrap.
echo "/health|get" >>"$tmp_routes"

mapfile -t endpoints < <(sort -u "$tmp_routes")

if [[ ${#endpoints[@]} -eq 0 ]]; then
  echo "no endpoints found while parsing routes" >&2
  exit 1
fi

{
  printf '{\n'
  printf '  "openapi": "3.0.3",\n'
  printf '  "info": {\n'
  printf '    "title": "Post Pilot API",\n'
  printf '    "version": "dev",\n'
  printf '    "description": "Development-generated endpoint catalog from Gin route declarations."\n'
  printf '  },\n'
  printf '  "servers": [\n'
  printf '    { "url": "%s" }\n' "$BASE_URL"
  printf '  ],\n'
  printf '  "paths": {\n'

  for i in "${!endpoints[@]}"; do
    IFS='|' read -r raw_path method <<<"${endpoints[$i]}"
    openapi_path="$(sed -E 's/:([A-Za-z0-9_]+)/{\1}/g' <<<"$raw_path")"
    operation_id="$(tr '/:{}-' '_' <<<"${method}_${openapi_path}" | sed -E 's/[^A-Za-z0-9_]+/_/g; s/_+/_/g; s/^_|_$//g')"

    if [[ $i -gt 0 ]]; then
      printf ',\n'
    fi

    printf '    "%s": {\n' "$openapi_path"
    printf '      "%s": {\n' "$method"
    printf '        "summary": "%s %s",\n' "${method^^}" "$openapi_path"
    printf '        "operationId": "%s",\n' "$operation_id"
    printf '        "responses": {\n'
    printf '          "200": { "description": "Success" }\n'
    printf '        }\n'
    printf '      }\n'
    printf '    }'
  done

  printf '\n  }\n'
  printf '}\n'
} >"$OUTPUT_FILE"

echo "Swagger endpoint catalog generated: $OUTPUT_FILE"
echo "Detected endpoints:"
for endpoint in "${endpoints[@]}"; do
  IFS='|' read -r path method <<<"$endpoint"
  printf '  - %-6s %s\n' "${method^^}" "$path"
done
