#!/bin/bash
# 矿山智能调度系统 - 完整 API 自动化测试
GATEWAY="http://localhost:8080"
PASS=0
FAIL=0

GREEN='\033[0;32m'; RED='\033[0;31m'; NC='\033[0m'

api() { curl -s -o /tmp/api_resp.json -w '%{http_code}' "$@"; }
ok() { echo -e "  ${GREEN}✓${NC} $1"; PASS=$((PASS+1)); }
fail() { echo -e "  ${RED}✗${NC} $1"; FAIL=$((FAIL+1)); }
assert_code() { local c=$(api "$@"); [ "$c" = "200" ] && ok "$2" || fail "$2 (HTTP $c)"; }

# Usage: assert_json <expected_code> <test_name> <curl_args...>
assert_json() {
  local expected_code="$1"
  local test_name="$2"
  shift 2
  local http_code
  http_code=$(curl -s -o /tmp/api_resp.json -w '%{http_code}' "$@")
  local code
  code=$(python3 -c "import sys,json; print(json.load(open('/tmp/api_resp.json')).get('code',0))" 2>/dev/null || echo "parse_error")
  [ "$code" = "$expected_code" ] && ok "$test_name" || fail "$test_name (code=$code, http=$http_code)"
}

echo "============================================"
echo "  矿山智能调度系统 - API 自动化测试"
echo "============================================"
echo ""

# ── 1. CORS ──
echo "── 1. CORS 跨域 ──"
CORS=$(curl -s -D - -X OPTIONS "$GATEWAY/api/v1/users" 2>/dev/null | grep -i "access-control" | tr -d '\r')
echo "$CORS" | grep -qi "allow-origin" && ok "OPTIONS 预检返回 CORS 头" || fail "CORS 预检失败"

# ── 2. 用户 CRUD ──
echo ""
echo "── 2. 用户服务 ──"

assert_json 0 "创建用户 alice" \
  -X POST "$GATEWAY/api/v1/users" \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pass123","realName":"测试用户A","mineId":1}'

USER_ID=$(python3 -c "import json; print(json.load(open('/tmp/api_resp.json'))['data']['id'])" 2>/dev/null)
[ -n "$USER_ID" ] && ok "用户 ID 已生成: $USER_ID" || fail "用户 ID 为空"

assert_json 0 "创建用户 bob" \
  -X POST "$GATEWAY/api/v1/users" \
  -H "Content-Type: application/json" \
  -d '{"username":"bob","password":"pass456","realName":"测试用户B","mineId":1}'

assert_json 500 "重复 username 拒绝" \
  -X POST "$GATEWAY/api/v1/users" \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"","realName":"","mineId":0}'

LIST_RESP=$(curl -s "$GATEWAY/api/v1/users")
LIST_CNT=$(echo "$LIST_RESP" | python3 -c "import sys,json; print(len(json.load(sys.stdin).get('data',[])))" 2>/dev/null)
[ "$LIST_CNT" -ge 2 ] && ok "用户列表 ≥2 条: $LIST_CNT" || fail "用户列表条数不足: $LIST_CNT"

# ── 3. 登录 ──
echo ""
echo "── 3. 登录认证 ──"

assert_json 0 "alice 正确密码登录" \
  -X POST "$GATEWAY/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pass123"}'

TOKEN=$(python3 -c "import json; print(json.load(open('/tmp/api_resp.json')).get('data',{}).get('token',''))" 2>/dev/null)
[ -n "$TOKEN" ] && ok "Token 已签发: ${TOKEN:0:40}..." || fail "Token 为空"

assert_json 401 "错误密码被拒绝" \
  -X POST "$GATEWAY/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"wrong"}'

assert_json 401 "不存在的用户被拒绝" \
  -X POST "$GATEWAY/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"nonexist","password":"x"}'

assert_json 0 "有效 token 验证通过" \
  -X POST "$GATEWAY/api/v1/auth/validate" \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\"}"

assert_json 401 "无效 token 验证拒绝" \
  -X POST "$GATEWAY/api/v1/auth/validate" \
  -H "Content-Type: application/json" \
  -d '{"token":"bad.token.value"}'

# ── 4. 车辆 ──
echo ""
echo "── 4. 车辆管理 ──"

assert_json 0 "创建矿用卡车 A001" \
  -X POST "$GATEWAY/api/v1/vehicles" \
  -H "Content-Type: application/json" \
  -d '{"plate":"矿卡-A001","type":1,"mineId":1}'

V_ID=$(python3 -c "import json; print(json.load(open('/tmp/api_resp.json')).get('data',{}).get('id',''))" 2>/dev/null)
[ -n "$V_ID" ] && ok "车辆 ID: $V_ID" || fail "车辆 ID 为空"

assert_json 0 "创建挖掘机 B001" \
  -X POST "$GATEWAY/api/v1/vehicles" \
  -H "Content-Type: application/json" \
  -d '{"plate":"挖机-B001","type":2,"mineId":1}'

assert_json 500 "重复车牌拒绝" \
  -X POST "$GATEWAY/api/v1/vehicles" \
  -H "Content-Type: application/json" \
  -d '{"plate":"矿卡-A001","type":1,"mineId":1}'

assert_json 0 "创建装载机 C001" \
  -X POST "$GATEWAY/api/v1/vehicles" \
  -H "Content-Type: application/json" \
  -d '{"plate":"铲车-C001","type":3,"mineId":1}'

assert_json 0 "创建矿卡 A002 (矿区2)" \
  -X POST "$GATEWAY/api/v1/vehicles" \
  -H "Content-Type: application/json" \
  -d '{"plate":"矿卡-A002","type":1,"mineId":2}'

VEH_LIST=$(curl -s "$GATEWAY/api/v1/vehicles")
VEH_CNT=$(echo "$VEH_LIST" | python3 -c "import sys,json; print(len(json.load(sys.stdin).get('data',[])))" 2>/dev/null)
[ "$VEH_CNT" -ge 3 ] && ok "车辆列表 ≥3 条: $VEH_CNT" || fail "车辆列表不足: $VEH_CNT"

# ── 5. 调度 ──
echo ""
echo "── 5. 调度服务 ──"

assert_json 0 "创建调度任务 (FIFO)" \
  -X POST "$GATEWAY/api/v1/dispatch/assign" \
  -H "Content-Type: application/json" \
  -d "{\"vehicle_id\":$V_ID,\"load_point_id\":1,\"dump_point_id\":2,\"algorithm\":\"fifo\"}"

TASK_ID=$(python3 -c "import json; print(json.load(open('/tmp/api_resp.json')).get('data',{}).get('id',''))" 2>/dev/null)
[ -n "$TASK_ID" ] && ok "调度任务 ID: $TASK_ID" || fail "调度任务 ID 为空"

assert_json 0 "创建调度任务 (最近优先)" \
  -X POST "$GATEWAY/api/v1/dispatch/assign" \
  -H "Content-Type: application/json" \
  -d "{\"vehicle_id\":$V_ID,\"load_point_id\":1,\"dump_point_id\":3,\"algorithm\":\"nearest_first\"}"

TASK_LIST=$(curl -s "$GATEWAY/api/v1/dispatch/tasks")
TASK_CNT=$(echo "$TASK_LIST" | python3 -c "import sys,json; print(len(json.load(sys.stdin).get('data',[])))" 2>/dev/null)
[ "$TASK_CNT" -ge 1 ] && ok "调度任务列表 ≥1 条: $TASK_CNT" || fail "调度任务列表为空"

# ── 6. 实时位置 ──
echo ""
echo "── 6. 实时定位 ──"

assert_json 0 "上报车辆位置" \
  -X POST "$GATEWAY/api/v1/telemetry/location" \
  -H "Content-Type: application/json" \
  -d "{\"location\":{\"vehicle_id\":$V_ID,\"latitude\":39.9042,\"longitude\":116.4074,\"speed\":35,\"heading\":180,\"altitude\":100}}"

# ── 7. 告警服务 ──
echo ""
echo "── 7. 告警服务 ──"

assert_json 0 "创建圆形电子围栏" \
  -X POST "$GATEWAY/api/v1/geofences" \
  -H "Content-Type: application/json" \
  -d '{"name":"测试围栏","shape":"circle","center_lat":39.9042,"center_lon":116.4074,"radius_m":500,"fence_type":"restricted","max_speed_kmh":40,"enabled":true}'

FENCE_ID=$(python3 -c "import json; print(json.load(open('/tmp/api_resp.json')).get('data',{}).get('id',''))" 2>/dev/null)
[ -n "$FENCE_ID" ] && ok "电子围栏 ID: $FENCE_ID" || fail "围栏 ID 为空"

assert_json 0 "创建告警规则" \
  -X POST "$GATEWAY/api/v1/alarm-rules" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"测试禁区规则\",\"rule_type\":\"geofence\",\"geofence_id\":$FENCE_ID,\"severity\":\"critical\",\"description\":\"测试\",\"enabled\":true}"

RULE_ID=$(python3 -c "import json; print(json.load(open('/tmp/api_resp.json')).get('data',{}).get('id',''))" 2>/dev/null)
[ -n "$RULE_ID" ] && ok "告警规则 ID: $RULE_ID" || fail "规则 ID 为空"

assert_json 0 "获取围栏列表" \
  -X GET "$GATEWAY/api/v1/geofences"

assert_json 0 "获取告警规则列表" \
  -X GET "$GATEWAY/api/v1/alarm-rules"

assert_json 0 "获取告警事件列表" \
  -X GET "$GATEWAY/api/v1/alarms"

assert_json 0 "获取告警统计" \
  -X GET "$GATEWAY/api/v1/alarms/stats"

assert_json 0 "位置检查 (安全位置)" \
  -X POST "$GATEWAY/api/v1/alarms/check-position" \
  -H "Content-Type: application/json" \
  -d "{\"latitude\":39.9000,\"longitude\":116.4200,\"speed\":10,\"vehicle_id\":1}"

assert_json 0 "位置检查 (闯入禁区)" \
  -X POST "$GATEWAY/api/v1/alarms/check-position" \
  -H "Content-Type: application/json" \
  -d "{\"latitude\":39.9042,\"longitude\":116.4074,\"speed\":50,\"vehicle_id\":1}"

echo ""

# ── 8. 结果 ──
echo "============================================"
echo -e "  通过: ${GREEN}$PASS${NC}  |  失败: ${RED}$FAIL${NC}  |  总计: $((PASS+FAIL))"
echo "============================================"

[ "$FAIL" -gt 0 ] && exit 1 || exit 0
