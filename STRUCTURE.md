# 双生花目录结构

> 2026-05-17 | v1 最终版 | 79 个目录，9 个顶层域

---

## 顶层总览

| 目录 | 英文 | 中文 |
|------|------|------|
| `docs/` | Documentation | 文档 |
| `root/` | Root — Model Intelligence | 根 — 模型能力 |
| `stem/` | Stem — Shared Tools | 茎 — 通用工具 |
| `flowers/` | Flowers — Shared/Daily/Explore | 花 — 共享/日常/探索 |
| `runtime/` | Runtime — Lifecycle | 运行时 — 生命周期 |
| `observatory/` | Observatory — System Observation | 观测台 — 系统观测 |
| `ui/` | UI — Visualization | 界面 — 可视化 |
| `compost/` | Compost — 66 Nutrients | 堆肥 — 66 营养层 |
| `sandbox/` | Sandbox — Experiments | 沙箱 — 临时实验 |

---

## docs/ — 文档

| 目录 | 英文 | 中文 |
|------|------|------|
| `docs/foundation/` | Founding documents | 立项文档 |
| `docs/architecture/` | Architecture decisions | 架构设计 |
| `docs/decisions/` | ADR — Architecture Decision Records | 架构决策记录 |
| `docs/roadmap/` | Roadmaps | 路线图 |
| `docs/experiments/` | Experiment logs | 实验记录 |
| `docs/migrations/` | Migration guides | 迁移指南 |

---

## root/ — 根：模型能力

模型是智能的唯一来源。只关心模型，不关心工具、workflow、eval。

| 目录 | 英文 | 中文 |
|------|------|------|
| `root/models/router/` | Intent classification model | 意图分类模型 |
| `root/models/planner/` | Deep reasoning model | 深度推理模型 |
| `root/models/coder/` | Code generation model | 代码生成模型 |
| `root/models/vision/` | Image recognition model | 图像识别模型 |
| `root/models/guard/` | Safety guardrail model | 安全审核模型 |
| `root/providers/local/` | Local model backend | 本地模型后端 |
| `root/providers/api/` | Cloud API backend | 云端 API 后端 |
| `root/providers/hybrid/` | Local + API mix | 本地+API 混合 |
| `root/capabilities/` | Model capability definitions | 模型能力定义 |

---

## stem/ — 茎：通用工具

工具是模型接触世界的器官。工具不知道自己在哪朵花里。

| 目录 | 英文 | 中文 |
|------|------|------|
| `stem/registry/` | Tool registration | 工具注册 |
| `stem/tools/weather/` | Weather query | 天气查询 |
| `stem/tools/search/` | Web search | 网络搜索 |
| `stem/tools/translate/` | Translation | 翻译 |
| `stem/tools/filesystem/` | File operations | 文件操作 |
| `stem/tools/code/` | Code execution | 代码执行 |
| `stem/tools/memory/` | Memory operations | 记忆操作 |
| `stem/adapters/` | Tool adapters / wrappers | 工具适配器 |
| `stem/contracts/` | Tool interface definitions | 工具接口定义 |

---

## flowers/ — 花

双花分工：左花 = Operational Flower（日常），右花 = Experimental Flower（探索）。

| 目录 | 解释 | 来源 |
|------|------|------|
| `flowers/daily/` | **左花:** execute / adapt / remember / respond | — |
| `flowers/explore/` | **右花:** 实验、验证、观测、演化 | — |
| `explore/evals/` | Exploration Bench — 能力成长测试 | 66 的 Eval Harness 重构 |
| `explore/shadow/` | Cognitive Shadow — 假设性执行对比 | 66 的 Shadow 轻量化 |
| `explore/drift/` | Preference Drift — 认知漂移观测 | 66 的 Drift Detection |
| `explore/evidence/` | Evidence Graph — 因果可视化 | 66 的 Evidence Graph |
| `explore/model_lab/` | Model Lab — 离线模型对比 | 66 的 Router A/B 限定 |

---

## runtime/ — 运行时：生命周期

不属于根、茎、花。它们共用的运行基础设施。

| 目录 | 英文 | 中文 |
|------|------|------|
| `runtime/engine/` | Core engine | 核心引擎 |
| `runtime/orchestration/` | Lifecycle orchestration | 生命周期编排 |
| `runtime/pipeline/` | Pipeline state machine | 管线状态机 |
| `runtime/execution/` | Tool execution | 工具执行 |
| `runtime/sessions/` | Session management | 会话管理 |
| `runtime/state/` | State management | 状态管理 |
| `runtime/policy/` | Model & routing policy | 模型与路由策略 |

---

## observatory/ — 观测台：系统观测

不是 eval。不是测试。是系统观测。

| 目录 | 英文 | 中文 |
|------|------|------|
| `observatory/graph/` | Evidence / knowledge graph | 证据/知识图谱 |
| `observatory/evidence/` | Evidence collection | 证据采集 |
| `observatory/metrics/` | Metrics & monitoring | 指标与监控 |
| `observatory/smoke/` | Smoke testing | 烟雾测试 |
| `observatory/drift/` | Routing drift detection | 路由漂移检测 |
| `observatory/freeze/` | Release freeze & gates | 发布冻结与门禁 |
| `observatory/reports/` | Generated reports | 生成报告 |

---

## ui/ — 界面：可视化

让你"看见系统"。

| 目录 | 英文 | 中文 |
|------|------|------|
| `ui/graph/` | Graph visualization | 图谱可视化 |
| `ui/dashboard/` | System dashboard | 系统仪表盘 |
| `ui/notebook/` | Notebook / canvas | 笔记本/画布 |
| `ui/twin-map/` | TwinFlower system map | 双生花系统图 |

---

## compost/ — 堆肥：66 营养层

66 不是被丢弃，是回归到养分来源。不参与主编译。

| 目录 | 英文 | 中文 |
|------|------|------|
| `compost/66-notes/` | Notes from 66 | 66 笔记 |
| `compost/extracted-patterns/` | Patterns extracted from 66 | 从 66 提取的模式 |
| `compost/deprecated/` | Deprecated code | 废弃代码 |
| `compost/experiments/` | Past experiments | 历史实验 |
| `compost/migrations/` | Migration scripts | 迁移脚本 |

---

## sandbox/ — 沙箱：临时实验

临时想法、新模型测试、疯狂想法。永不进入主线。

| 目录 | 英文 | 中文 |
|------|------|------|
| `sandbox/model-tests/` | Model evaluation | 模型评估 |
| `sandbox/routing-tests/` | Routing experiments | 路由实验 |
| `sandbox/weird-ideas/` | Wild ideas | 疯狂想法 |
| `sandbox/temp/` | Temporary files | 临时文件 |
