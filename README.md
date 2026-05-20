
# 矿山智能调度系统 (Intelligent Mine Dispatch System)

[![CI](https://github.com/alvis-ai/intelligent-mine-dispatch-system/actions/workflows/ci.yml/badge.svg)](https://github.com/alvis-ai/intelligent-mine-dispatch-system/actions/workflows/ci.yml)

生产级 **Go 微服务** 矿山智能调度平台，面向大型露天矿 / 地下矿场景，支持自动派车、实时追踪、路径规划和 AI 优化调度。

---

## 架构概览

```
┌─────────────┐     ┌──────────────┐     ┌─────────────────┐
│  Web Admin  │────▶│   Gateway    │────▶│  gRPC Services  │
│ (React/TS)  │     │  (go-zero)   │     │  (Go 微服务)    │
└─────────────┘     └──────┬───────┘     └────────┬────────┘
                           │                      │
                    ┌──────┴──────┐        ┌───────┴────────┐
                    │  WebSocket  │        │  PostgreSQL    │
                    │  + Redis    │        │  + Redis       │
                    └─────────────┘        └────────────────┘
```

## 技术栈

| 类型 | 技术 |
|------|------|
| 主语言 | Go 1.25 |
| 微服务框架 | go-zero + gRPC |
| 数据库 | PostgreSQL 16 |
| 缓存 / 消息 | Redis 7 (Pub/Sub) |
| 前端 | React 19 + TypeScript + Ant Design + Vite |
| 实时通信 | WebSocket + Redis Pub/Sub |
| 部署 | Docker Compose |
| CI/CD | GitHub Actions |
| API协议 | gRPC + REST (Gateway 代理) |

## 微服务

| 服务 | 端口 | 职责 |
|------|------|------|
| **gateway** | 8080 | API 网关，统一 REST 入口 |
| **user-service** | 8081 | 用户 CRUD 与管理 |
| **auth-service** | 8082 | JWT 登录认证与鉴权 |
| **vehicle-service** | 8083 | 车辆管理 (矿卡/挖机/铲车) |
| **telemetry-service** | 8084 | GPS 实时定位 & WebSocket 广播 |
| **dispatch-service** | 8085 | 调度引擎 (多算法) |
| **ai-service** | 8086 | AI 拥堵预测 & 智能调度建议 |
| **alarm-service** | 8087 | 电子围栏 & 实时告警 |
| **route-service** | 8088 | 道路网络 & 最短路径规划 |
| **device-service** | 8089 | IoT 设备管理（GPS/传感器/摄像头） |

## 调度算法

- **FIFO** — 先进先出
- **Weighted Round Robin** — 加权轮询
- **Nearest First** — 基于道路距离的最近优先（集成路线服务）
- **Genetic Algorithm** — 遗传算法批量优化调度
- **AI Suggest** — 基于拥堵/负载/距离多因子评分的 AI 调度建议
- **Dijkstra / A\*** — 矿区道路最短路径（路线服务）
- **AI Route** — 拥堵感知的 AI 加权路线推荐

## 前端页面 (9个)

| 页面 | 功能 |
|------|------|
| Dashboard | 调度看板，实时统计 |
| Login | JWT 登录 |
| Map | 实时地图追踪，车辆定位 |
| Tasks | 调度任务管理，创建/完成/取消 |
| Vehicles | 车辆 CRUD，编辑/删除 |
| VehicleTypes | 车型管理 |
| LoadingPoints | 装载点/卸载点管理 |
| Alarms | 告警中心，告警事件管理 |
| Geofences | 电子围栏管理 |

## 快速开始

### 前置要求

- Go 1.25+
- Docker & Docker Compose
- Node.js 20+ (前端开发)

### 本地开发

```bash
# 1. 编译所有服务二进制 (Linux/ARM64)
cd services/dispatch-service && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ../../deploy/docker/build/dispatch-service ./cmd
# (其他服务同理)

# 2. 启动所有服务
cd deploy/docker
docker compose up -d --build

# 3. 验证部署
curl http://localhost:8080/api/v1/vehicles
```

### 运行 API 测试

```bash
bash scripts/api_test.sh
```

### 运行单元测试

```bash
# 测试指定服务
cd services/dispatch-service && go test ./... -v -count=1
```

## 项目结构

```
├── gateway/                 # API 网关
│   └── internal/handler/    # REST 路由处理器
├── services/                # 微服务
│   ├── ai-service/
│   ├── alarm-service/
│   ├── auth-service/
│   ├── device-service/
│   ├── dispatch-service/
│   ├── route-service/
│   ├── telemetry-service/
│   ├── user-service/
│   └── vehicle-service/
├── proto/                   # gRPC proto 定义
│   ├── dispatch/v1/
│   ├── route/v1/
│   ├── telemetry/
│   └── ...
├── web-admin/               # React 前端
│   └── src/pages/           # 页面组件
├── deploy/
│   ├── docker/              # Docker Compose 部署
│   └── k8s/                 # Kubernetes 配置
├── scripts/
│   └── api_test.sh          # API 端到端测试
├── pkg/                     # 共享工具包
└── .github/workflows/       # CI/CD
```

## 数据流

```
设备 GPS → Telemetry Service → Redis → WebSocket → 前端地图
                                   ↓
                             Dispatch Service → 调度算法决策 (含 AI)
                                   ↓
                             AI Service → 拥堵预测 / 需求预测 / 智能建议
                                   ↓
                             Route Service → 道路距离计算
                                   ↓
                             Alarm Service → 电子围栏检查
```

## 路线规划

路线服务 (`route-service`) 提供完整的矿山道路网络管理和路径规划：

- 道路节点 / 边 CRUD
- Dijkstra 最短路径算法
- A\* 最短路径算法（Haversine 距离作为启发函数）
- 批量距离计算
- 行驶时间估算（基于速度限制）

### 路线 API

```bash
# 计算两点间道路距离
curl -X POST http://localhost:8080/api/v1/route/calculate \
  -H "Content-Type: application/json" \
  -d '{"from_lat":39.895,"from_lon":116.405,"to_lat":39.908,"to_lon":116.402,"algorithm":"dijkstra"}'
```

## 生产部署

### 使用预编译二进制 (推荐)

```bash
# 1. 交叉编译所有服务二进制
for svc in gateway user-service auth-service vehicle-service telemetry-service dispatch-service alarm-service route-service ai-service device-service; do
  if [ "$svc" = "gateway" ]; then
    cd gateway && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ../deploy/docker/build/gateway ./cmd && cd ..
  else
    cd services/$svc && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ../../deploy/docker/build/$svc ./cmd && cd ../..
  fi
done

# 2. 构建并启动
cd deploy/docker && docker compose up -d --build
```

### 从源码构建 (Docker 多阶段)

```bash
cd deploy/docker && docker compose -f docker-compose.source.yaml up -d --build
```

## 路线图

### Phase 1 — MVP ✅
- [x] 用户系统 & JWT 鉴权
- [x] 车辆 CRUD 管理
- [x] GPS 实时定位 & WebSocket
- [x] 基础调度 (FIFO)

### Phase 2 — 生产基础版 ✅
- [x] 遗传算法调度优化
- [x] 路线规划服务 (Dijkstra/A*)
- [x] NearestFirst 集成道路距离
- [x] 电子围栏 & 告警系统
- [x] 前端功能补全
- [x] CI/CD 流水线
- [x] Docker 容器化部署

### Phase 3 — 企业版 🚧
- [ ] AI 调度 / 拥堵预测 (ai-service)
- [ ] IoT 设备管理 (device-service)
- [ ] BI 报表分析 (report-service)
- [ ] 语音调度 (voice-service)
- [ ] 多租户 / 多矿区

### Phase 4 — 生产级
- [ ] Kubernetes 编排
- [ ] 自动扩缩容
- [ ] Prometheus + Grafana 监控
- [ ] 分布式链路追踪
- [ ] 灰度发布

## 许可证

MIT
