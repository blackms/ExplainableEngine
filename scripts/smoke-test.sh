#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${1:?Usage: $0 <base-url>}"

echo "Smoke testing ${BASE_URL}..."

# Health
echo -n "GET /health... "
curl -sf "${BASE_URL}/health" | grep -q '"healthy"' && echo "OK" || { echo "FAIL"; exit 1; }

# POST explain
echo -n "POST /api/v1/explain... "
RESP=$(curl -sf -X POST "${BASE_URL}/api/v1/explain" \
  -H "Content-Type: application/json" \
  -d '{"target":"smoke_test","value":0.72,"components":[{"name":"a","value":0.8,"weight":0.4,"confidence":0.9},{"name":"b","value":0.5,"weight":0.3,"confidence":0.7},{"name":"c","value":0.6,"weight":0.3,"confidence":0.85}]}')
ID=$(echo "$RESP" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
[ -n "$ID" ] && echo "OK (id=$ID)" || { echo "FAIL"; exit 1; }

# GET explain
echo -n "GET /api/v1/explain/${ID}... "
curl -sf "${BASE_URL}/api/v1/explain/${ID}" | grep -q '"smoke_test"' && echo "OK" || { echo "FAIL"; exit 1; }

# GET graph (mermaid)
echo -n "GET /api/v1/explain/${ID}/graph?format=mermaid... "
curl -sf "${BASE_URL}/api/v1/explain/${ID}/graph?format=mermaid" | grep -q 'graph LR' && echo "OK" || { echo "FAIL"; exit 1; }

# GET narrative
echo -n "GET /api/v1/explain/${ID}/narrative... "
curl -sf "${BASE_URL}/api/v1/explain/${ID}/narrative" | grep -q '"narrative"' && echo "OK" || { echo "FAIL"; exit 1; }

# POST what-if
echo -n "POST /api/v1/explain/${ID}/what-if... "
curl -sf -X POST "${BASE_URL}/api/v1/explain/${ID}/what-if" \
  -H "Content-Type: application/json" \
  -d '{"modifications":[{"component":"a","new_value":0.95}]}' | grep -q '"delta_value"' && echo "OK" || { echo "FAIL"; exit 1; }

echo ""
echo "All smoke tests passed!"
