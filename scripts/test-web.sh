#!/usr/bin/env bash
set -euo pipefail

WEB="${WEB_URL:-http://localhost:3000}"
API="${API_URL:-http://localhost:8080}"
SUFFIX=$(date +%s)

echo "Web UI test against $WEB (API $API)"

for page in /login /register; do
  code=$(curl -s -o /dev/null -w "%{http_code}" "$WEB$page")
  if [ "$code" != "200" ]; then
    echo "FAIL: $page returned $code"
    exit 1
  fi
  echo "  $page OK ($code)"
done

register=$(curl -sf -X POST "$API/api/v1/auth/register" \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"web-$SUFFIX@test.com\",\"username\":\"web_$SUFFIX\",\"password\":\"password123\"}")
token=$(echo "$register" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['access_token'])")

daily=$(curl -sf "$API/api/v1/questions/daily" -H "Authorization: Bearer $token")
qcount=$(echo "$daily" | python3 -c "import sys,json; print(len(json.load(sys.stdin)['data']['questions']))")
streak=$(curl -sf "$API/api/v1/streaks/me" -H "Authorization: Bearer $token")
progress=$(curl -sf "$API/api/v1/progress/me" -H "Authorization: Bearer $token")

echo "  dashboard data: questions=$qcount streak=$(echo "$streak" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['current_streak'])") xp=$(echo "$progress" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['total_xp'])")"
echo "Web test passed"
