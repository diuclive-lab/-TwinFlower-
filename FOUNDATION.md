# 双生花

> 一个以模型能力为根、通用工具为茎、日常与探索双花共生的智能体系统。
>
> 立项于 2026-05-17。继承 66 (FangLab) 的经验，但不继承它的结构。

---

## 世界观

智能体不是一个聊天框。

它由四层构成：

```
      左花（日常）              右花（探索）
     workflows（多步任务）    workflows（多步任务）
              \                  /
               \                /
           ─── 维管束：Skills ───
          工具编排契约 / 约束规则
                    │
                    │
           ─── 茎：Tools ─────
              原子能力
                    │
                    │
             根：Model
             思考能力
                    │
                    │
          66（堆肥 / 营养）
```

---

## 四层定义

### 第一层：根 — Model（模型能力）

模型是智能的唯一来源。

- 理解
- 推理
- 规划
- 生成
- 决策

没有模型，不存在智能体。

### 第二层：茎 — Tools（原子工具）

工具是模型接触世界的器官。只做一件事，不思考，不规划。

```go
weather(city)    → temperature
translate(text)  → translated
stock(symbol)    → price
```

工具不关心谁在调用它。日常花和探索花用同一套工具。

### 第三层：维管束 — Skills（工具编排契约）

Skill 不是工具。Skill 是"如何使用工具的知识"。

```
weather_skill:
  工具: weather_api
  流程: 提取城市 → 调用 → 格式化
  约束: 不能同时查多个城市

code_search_skill:
  工具: grep, filesystem_read, symbol_search
  流程: 搜索 → 定位 → 读取
  约束: 只读，不改文件
```

Skill = 流程 + 约束。这是防止模型乱用工具的关键。

### 第四层：花 — Workflows（多步任务）

Workflow 组合多个 Skill 完成复杂目标。

```
repository_analysis_workflow:
  步骤:
    1. code_search_skill (理解结构)
    2. file_compare_skill (对比变化)
    3. proposal_skill (生成建议)
```

日常花的工作流偏向稳定、可预测。探索花的工作流偏向实验、长思考、自我恢复。

---

## 能力迁移机制

```
探索花实验成功
    → 右花积累数据
    → 稳定后写入 Skills
    → 左花使用 Skills
```

新能力先在探索花实验，成熟后沉淀为技能，日常花直接使用。

---

## 营养审计 (66 → 双生花)

详见 [`compost/nutrient_audit.md`](compost/nutrient_audit.md)

三类标记：
- **ROOT-READY** — 直接吸收 (20 模块)
- **STEM-LATER** — 后续再长 (11 模块)
- **COMPOST** — 留在 66 (18 模块)

---

## 认知系统

模型不是"理解"语言，是需要适配的大脑。

详见 [`docs/architecture/cognitive_system.md`](docs/architecture/cognitive_system.md)
- 认知适配矩阵（架构层 × 行为层）
- 在线校准（运行中自动调整 clarify 阈值）
- 与 Phase 2（clarify / intent evidence / failure learning）的关系

---

## 设计原则

1. **模型是根。** 所有能力围绕模型，而不是围绕框架。
2. **工具是死的，技能是活的。** 工具不包含使用知识，技能包含。
3. **Skill = 流程 + 约束。** 防止模型乱用工具。
4. **双花共享茎和维管束。** 工具和技能不分日常还是探索。
5. **新能力先在探索花生根。** 成熟后沉淀为技能，右花到左花。
6. **66 是养分，不是束缚。** 吸收经验，不复制结构。
7. **宪章优先于代码。** 先想清楚，再动手。

---

## 模型体系

| 角色 | 模型 | 内存 |
|------|------|------|
| 主力 | Qwen3.6-27B | ~16GB |
| 路由 | DeepSeek-R1-Distill-7B | ~4.4GB |
| 护卫 | Llama-3.2-3B | ~1.9GB |
| 备用 | Gemma-4-31B | ~17GB |
