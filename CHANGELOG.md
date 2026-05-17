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
