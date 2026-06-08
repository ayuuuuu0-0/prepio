#!/usr/bin/env bash
set -euo pipefail

API="${API_URL:-http://localhost:8080}"
SUFFIX=$(date +%s)

echo "E2E test against $API"

register=$(curl -sf -X POST "$API/api/v1/auth/register" \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"e2e-$SUFFIX@test.com\",\"username\":\"e2e_$SUFFIX\",\"password\":\"password123\"}")

token=$(echo "$register" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['access_token'])")
echo "Registered and got token"

daily=$(curl -sf "$API/api/v1/questions/daily" -H "Authorization: Bearer $token")
session=$(echo "$daily" | python3 -c "import sys,json; d=json.load(sys.stdin)['data']; print(d['session_id'])")
qid=$(echo "$daily" | python3 -c "import sys,json; d=json.load(sys.stdin)['data']; print(d['questions'][0]['id'])")
echo "Daily paper: session=$session question=$qid"

submit=$(curl -sf -X POST "$API/api/v1/questions/$qid/submit" \
  -H "Authorization: Bearer $token" \
  -H 'Content-Type: application/json' \
  -d "{\"session_id\":\"$session\",\"answer\":\"hash map approach with O(n) time and O(n) space complexity\",\"time_spent_seconds\":60}")
correct=$(echo "$submit" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['correct'])")
echo "Submit correct=$correct"

sleep 2

streak=$(curl -sf "$API/api/v1/streaks/me" -H "Authorization: Bearer $token")
current=$(echo "$streak" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['current_streak'])")
echo "Streak=$current"

progress=$(curl -sf "$API/api/v1/progress/me" -H "Authorization: Bearer $token")
xp=$(echo "$progress" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['total_xp'])")
echo "XP=$xp"

if [ "$correct" != "True" ] && [ "$correct" != "true" ]; then
  echo "FAIL: answer not marked correct"
  exit 1
fi

echo "E2E API test passed"
