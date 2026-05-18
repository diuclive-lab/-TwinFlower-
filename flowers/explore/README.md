# 探索花 — Experimental Flower

验证线在双生花中不消失，但重新定位：从"免疫系统"变成"研究器官"。

右花不参与 runtime 决策。它观测、实验、记录、报告。

## 域

| 域 | 66 来源 | 新定位 |
|----|---------|--------|
| `evals/` | Eval Harness → | Exploration Bench — 能力成长测试，非发布门禁 |
| `shadow/` | Shadow Router → | Cognitive Shadow — 假设性执行对比，不参与 runtime |
| `drift/` | Drift Detection → | Preference Drift — 认知漂移观测 |
| `evidence/` | Evidence Graph → | 因果可视化 — 为什么系统变成这样 |
| `model_lab/` | Router A/B → | Offline model comparison — 非线上双路执行 |

## 原则

- 右花不阻塞左花。探索实验失败不影响日常使用。
- 右花不产生 gate。观测结果不阻止发布。
- 右花可以沉淀到左花：成熟的能力可以通过 skills 迁移到日常花。
