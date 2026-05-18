# 双生花开发日志

> 格式：`YYYY-MM-DD | 主题 | 关键决策 / 里程碑`

---

## 2026-05-17

### 立项 — 从 66 到双生花

**背景：** 66 (FangLab) 积累了约 50K 行主线代码和 12K 行验证层代码，但面临一个根本问题：验证线（smoke/freeze/gate/eval）远远大于演化线（workflow/memory/routing）。连续 9 轮迭代全部投入验证层，主线几乎停滞。

**关键决策：** 结束 66 的无限膨胀阶段，进入双生花新立项。不继承 66 的结构，只继承经验。

**核心理念：** 智能体不是"模型 + 工具"，而是"根（模型能力）→ 茎（通用工具）→ 花（日常/探索行为）"的四层结构。后续进一步演化为五层 + 两花。

### 架构演化

| 迭代 | 结构 | 里程碑 |
|------|------|--------|
| 初始 | 三层：模型→工具→花 | 立项宪章 |
| V1 | 四层：根/茎/维管束/花 | 发现 Skill 层是缺失的关键 |
| V2 | 五层 + cognition | 认知适配层（Dense/MoE + clarify 阈值） |
| V3 | 五层 + preferences | Adaptive Clarify + Soft Execute |

### 八个顶层域

```
root/           根 — 模型能力（providers + cognition + preferences）
stem/           茎 — 通用原子工具
vascular/       维管束 — 工具编排契约（skills）
flowers/        花 — 共享/日常/探索
runtime/        运行时 — 引擎/生命周期
observatory/    观测台 — 系统观测（graph/evidence/metrics）
compost/        堆肥 — 66 营养层
docs/           文档 — 架构决策记录
```

### 关键认知发现

1. **Skill ≠ Tool** — 工具是原子能力，Skill 是"怎么用工具的知识"（流程 + 约束）。两者的分离是系统不再乱用工具的关键。
2. **Provider ≠ Cognition** — Provider 解决"怎么调用模型"，Cognition 解决"怎么和模型对话"。前者是传输层，后者是认知适配层。
3. **Clarify 不是错误恢复，是一级路由** — 当模型在多个意图之间摇摆时，让它问而不是猜。但需要 Adaptive Clarify 防止"澄清成瘾"。
4. **偏好需要遗忘机制** — 长期 agent 的最大问题不是不会学，而是不会降权旧经验。Recency decay + correction learning 是必须的。

### 模型阵容

| 角色 | 模型 | 大小 | 说明 |
|------|------|------|------|
| 主力 | Qwen3.6-27B (Dense) | 16GB | 日常任务，中文强 |
| 路由 | DeepSeek-R1-Distill-7B (MoE) | 4.4GB | 意图分类 |
| 护卫 | Llama-3.2-3B | 1.9GB | 安全审核 |
| 备用 | Gemma-4-31B (Dense) | 17GB | 主力替身 |

### 核心闭环

```
Phase 1 最小闭环： "Beijing weather" → 26.2°C（model → skill → tool → result）
Phase 1.5:         filesystem / search / 三个预定义 skill
Phase 2 核心:      clarify 契约 / preference learning / soft execute
```

### 三层决策系统

```
Level 1 (明确)    → 直接执行
Level 2 (偏好)    → soft execute（"猜你是想查天气" + 执行，可纠正）
Level 3 (模糊)    → 正式 clarify（"你是想查天气还是新闻？"）
```

### 关键提交

| Commit | 内容 |
|--------|------|
| `4f6c67c` | Phase 1 最小生根 — 第一个闭环 |
| `6f38ed4` | 认知适配层（Dense/MoE profile） |
| `7dd5c35` | 垂直切片 — filesystem + calibration |
| `8c50fa4` | Clarify Contract v1 |
| `e6b1e47` | Adaptive clarify — preference learning |
| `d8c992d` | Translate baseline — 认知不污染简单任务 |
| `97758dd` | Recency decay — 偏好遗忘机制 |
| `e93b02e` | Soft execute — 三层决策完整闭环 |
| `584055f` | Phase 3A: Filesystem Skill Deepening |
| `ddc693f` | Phase 3B: Search Skill Deepening |

---

## 2026-05-18

### Phase 3B: Search Skill Deepening

search_skill 从直通式升级为五阶段管线。

**交付内容：**
- ambiguity detection：高频歧义词库（16 组）+ disambiguator 检测
- clarify integration：高模糊返回 clarify 问题，无需模型调用
- soft execute：中模糊执行 + hint 提示（"如果不是请纠正"）
- ToolRunner 解耦：通过 tool interface 调用，消除与 stem/tools/search 的代码重复
- 结果成形：结构化 header + 去重 + 截断 + 结果计数

**附随修复：**
- stem/tools/search: 添加 User-Agent 头解决 DDG bot 检测（202→200）
- engine: isFilesystemRequest 移除 "search" 关键词，修复 bypass 抢占

**里程碑含义：**
Phase 3B 完成标志着双生花的 Skill 层从"路由契约"进入"程序化智能"阶段。搜索不再是 query→fetch→format 的直通，而是包含 ambiguity resolution 的完整认知管线。

---

### 从 66 移植基础模块

从 66 (FangLab) 移植了 5 个基础模块到双生花。

**基础设施（3 个）：**
- `internal/envutil` — 类型化环境变量读取（String/Int/Bool/Duration），来自 66 的 `envutil`
- `internal/httputil` — 共享 HTTP 客户端，合理超时/连接池/TLS 配置，来自 66 的 `httputil`
- `internal/retry` — 指数退避 HTTP 重试（支持网络错误/5xx/429 自动重试），来自 66 的 `model/retry`

**业务工具（2 个）：**
- `stem/tools/stock` — Yahoo Finance API 股价查询（无 API Key），Provider 接口 + HTTP Provider，来自 66 的 `hermes/stock`
- `stem/tools/currency` — open.er-api.com 汇率转换（无 API Key），Provider 接口 + HTTP Provider，来自 66 的 `hermes/currency`

**修复：**
- currency 工具修复了 66 中 API 字段名 bug（`conversion_rates` → `rates`，API v6 已变更）

**战略含义：**
这是双生花第一次从 66 吸收代码而非仅吸收经验。证明了两件事：
1. 66 的干净模块可以零摩擦地适配双生花的简化接口
2. 新架构下工具扩展只需：写实现 → 注册到 engine → 模型自动路由

### 当前系统完整管线

```
input
  → 根: cognition (profile + clarify threshold)
  → 维管束: skill intent routing
    → filesystem bypass (intent parser + path recovery + orchestration)
    → search bypass (ambiguity detection → clarify/soft execute → execute → shape)
    → 模型 plan (tool selection with contract constraints)
      → preference resolution (EWMA + decay + correction)
      → clarify (formal clarify with candidates)
  → 茎: tool execution
  → 反馈: calibration logging
  → 根: finalize (natural response formatting)
```

---

## 2026-05-18 (part 2)

### 重构 compost: 66 从"堆肥"重新定义为"演化土壤"

66 在双生花中不再只是代码源，而是四层养分体系：

```
compost/
  nutrient/         矿物质 — 可直接移植的代码 (stock/currency/envutil/httputil/retry)
  lessons/          腐殖质 — 失败经验沉淀 (shadow router/gate overgrowth/intent routing)
  anti_patterns/    毒素标记 — 禁止迁移的模式
  decision_history/ 菌根网络 — 架构决策记录 (P2 freeze/skill≠tool/bypass)
```

### 重构右花: 验证线从"免疫系统"重构成"研究器官"

66 的验证线不搬进双生花, 而是重新解释为探索花的 5 个实验域:

| 域 | 66 来源 | 新定位 | 关键改造 |
|----|---------|--------|---------|
| `explore/evals/` | Eval Harness | Exploration Bench | 去 gate 语义，加 Delta 演化对比 |
| `explore/shadow/` | Shadow Router | Cognitive Shadow | 去 cutover/promotion，纯记录 |
| `explore/drift/` | Drift Detection | Preference Drift | 去 gate 触发，Severity 分级 |
| `explore/evidence/` | Evidence Graph | 因果可视化 | 从验证工具变为探索工具 |
| `explore/model_lab/` | Router A/B | 离线模型对比 | 从线上双路变为 offline research |

**右花原则:**
- 右花不阻塞左花, 右花不产生 gate
- 能力从右花沉淀到左花 Skills (成熟后再迁移)

### 关键提交

| Commit | 内容 |
|--------|------|
| `db8b490` | 重构 compost: 66 → 演化土壤 |
| `cfbfc23` | 重构右花: 验证线 → 研究器官 |
| `6aee09a` | 探索花 evals: Exploration Bench |
| `b397cc9` | 探索花: shadow/drift/evidence/model_lab |
