# 营养审计 — 66 → 双生花

> 2026-05-17
>
> 判断标准：这个东西是否直接提高"大模型→工具→结果"的成功率？

---

## ✅ ROOT-READY — 直接吸收

这些模块成熟、稳定、直接提升核心循环成功率。

| 模块 | 来源 (66) | 理由 |
|------|-----------|------|
| weather tool | `internal/tools/hermes/weather.go` | 核心工具，直接搬 |
| translate tool | `internal/tools/hermes/translate.go` | 核心工具，直接搬 |
| stock tool | `internal/tools/hermes/stock.go` | 核心工具，直接搬 |
| currency tool | `internal/tools/hermes/currency.go` | 核心工具，直接搬 |
| filesystem tools | `internal/tools/filesystem/` | 核心工具，只拿 read/list/search |
| web search | `internal/tools/hermes/web.go` | 核心工具 |
| tool registry | `internal/tools/registry.go` | Register/Lookup/Schema，必需品 |
| tool schema | `internal/tools/tool.go` | Tool 接口定义 |
| execution engine | `internal/execution/execution.go` | tool.Run，最成熟的部分 |
| execution helpers | `internal/execution/execution_helpers.go` | args/results 处理 |
| retry | `internal/orchestrator/retry.go` | retryWithBackoff + fallbackMatrix |
| timeout | `internal/execution/descriptor.go` | Budget.TimeoutSeconds |
| parameter normalize | `internal/execution/execution_helpers.go` | args 格式化 |
| session cache | `internal/memory/session_cache.go` | Store/Lookup/PurgeStale |
| context pipeline | `internal/orchestrator/pipeline.go` | 状态机驱动（纯逻辑，无门禁） |
| routing primitives | `internal/tools/evidence_router.go` | RouteTopN（只拿接口，不拿关键词规则） |
| radix tree | `internal/tools/radixtree.go` | 轻量前缀匹配 |
| EWMA weights | `internal/tools/feedback_store.go` | 关键词权重更新（有用，不依赖 gate） |

---

## 🔶 STEM-LATER — 后续再长

这些有价值，但当前生根阶段不需要。等茎和花长出来再搬。

| 模块 | 理由 | 预计时机 |
|------|------|---------|
| workflow executor | 20K 行代码，当前阶段不需要多步 workflow | Phase 3 之后 |
| workflow contracts | 与 workflow executor 绑定 | Phase 3 之后 |
| task progress | 与 workflow 绑定 | Phase 3 之后 |
| planner model | 主模型自己做规划，不需要独立 planner | Phase 4 |
| long memory | 记忆写入回路刚通，还需验证 | Phase 3 |
| memory bridge | 跨 session 记忆，当前不需要 | Phase 3 |
| evidence graph | 可视化有用但不是核心循环 | Phase 4 |
| code model | 主模型自己写代码，不需要独立代码模型 | 不确定 |
| conclusion engine | 额外抽象层，当前不需要 | 不确定 |
| model pool | 多模型策略，当前单模型够用 | 不确定 |

---

## ❌ COMPOST — 留在 66

这些是验证层膨胀的产物，当前阶段不应带入。未来可能需要其中部分能力，但要以双生花的方式重新长出来，而不是移植。

| 模块 | 原因 |
|------|------|
| shadow router | 验证层产物，当前不需要旁路评估 |
| freeze gate | 发布门禁，当前不需要 |
| schema drift | 工具结构漂移检测，膨胀产物 |
| eval harness | 整个 `cmd/agent-eval/`，太重 |
| eval scenarios | `eval/scenarios/` 下的 YAML 文件 |
| project-health | 项目健康面板，当前不需要 |
| router-ab-report | A/B 比较报告，当前不需要 |
| doctor | 系统诊断工具 |
| LlamaGuardrail | 审核模型，当前不需要额外安全层 |
| guardrails package | 大部分守卫逻辑（保留 identity 和 human_approval 的基础版本） |
| sidecar diagnostics | 侧边车诊断 |
| trace bundle | 轨迹追踪 |
| knowledge base | 知识库工具 |
| MCP 相关 | MCP server/client，当前不需要 |
| vision model | 图像识别，当前不需要 |
| audio model | 语音，当前不需要 |
| ETL engine | Gemma-4-E4B，当前不需要 |

---

## ROOT-READY 清单统计

| 类别 | 数量 |
|------|------|
| 工具 | 8 |
| 基础设施 | 9 |
| 路由 | 3 |
| **总计** | **20** |

这 20 个模块构成双生花 v0 的根。每个都是直接提高"模型→工具→结果"成功率的东西。没有一个是验证层、门禁层、或观测层。
