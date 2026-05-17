# 认知系统 — Cognitive Architecture

> 2026-05-17 | 从 prompt 模板到认知适配矩阵

---

## 核心转变

66：`model → route → tool → validate`

TwinFlower：`model → cognition → skill → tool → learning`

---

## 二维适配矩阵

模型适配不是一维的 Dense/MoE 分类，而是两维的矩阵：

```
                     Architecture Layer
                Dense               MoE
Behavior     ┌─────────────────────────────
  Qwen       │  explicit + aggressive      gated + conservative
  Gemma      │  explicit + cautious         gated + cautious
  DeepSeek   │  (n/a)                       gated + conservative
  Llama      │  explicit + balanced         (n/a)
```

### 第一维：Architecture Layer（共性层）

由模型架构决定，同架构的模型共享特征。

| 特征 | Dense | MoE |
|------|-------|-----|
| chain_style | explicit | gated |
| tool_bias | balanced | conservative |
| clarify_threshold | 0.45 | 0.65 |
| prompt_granularity | fine | coarse |

### 第二维：Behavior Profile（个性层）

同一个架构内不同模型的行为差异。

| 特征 | Qwen (Dense) | Gemma (Dense) | DeepSeek (MoE) |
|------|-------------|--------------|----------------|
| overcommit | high | medium | medium |
| tool_aggressiveness | aggressive | cautious | conservative |
| clarify_preference | low | high | medium |
| ambiguity_tolerance | low | medium | low |
| reasoning_strength | medium | medium | high |

---

## 合并策略

运行时动态合并两维：

```
profile = architecture.Dense + behavior.Qwen
→ prompt_style: explicit
→ tool_bias: aggressive
→ clarify_threshold: 0.50
```

换 behavior.Gemma：

```
profile = architecture.Dense + behavior.Gemma
→ clarify_threshold: 0.72
→ tool_bias: cautious
```

换 behavior.DeepSeek（MoE）：

```
profile = architecture.MoE + behavior.DeepSeek
→ chain_style: gated
→ tool_bias: conservative
→ clarify_threshold: 0.70
```

---

## 在线校准 (Cognitive Calibration)

不是静态 profile，而是运行中自动更新。

### 校准表

```
model       weather  filesystem  code  clarify_rate
qwen27b     92%      71%         94%   18%
gemma31b    88%      86%         73%   41%
deepseek    95%      81%         89%   24%
```

### 自适应规则

- 某个模型在某类意图上频繁误判 → 自动提高该类意图的 clarify_threshold
- 某个模型在某类工具上稳定 → 自动降低 clarify_threshold，直接执行
- 校准数据来自 failure learning 的积累

### 目录结构

```
root/cognition/
├── profile.go            (原 Profile 定义)
├── registry.go           (Profile registry)
├── architectures/
│   ├── dense.go
│   └── moe.go
├── behaviors/
│   ├── qwen.go
│   ├── gemma.go
│   ├── deepseek.go
│   └── llama.go
└── calibration/
    ├── table.go          (校准表)
    ├── adapter.go        (自适应策略)
    └── data.go           (校准数据)
```

---

## 和 Phase 2 的关系

Phase 2 的四个方向（clarify skill / intent evidence / failure learning / corpus）与 cognition 层天然耦合：

- **clarify skill** → 由 cognition 层的 clarify_threshold 触发
- **intent evidence** → 由 cognition 层的 PromptSet 定义输出格式
- **failure learning** → 供给 calibration 层的数据源
- **corpus** → calibration 层的验证集
