# Model Lab

66 的 Router A/B 在这里变成离线模型对比环境。

## 实验类型

| 实验 | 方法 | 用途 |
|------|------|------|
| Model swap | Qwen vs Gemma 同一输入 | 模型行为差异 |
| Threshold sweep | 0.3-0.8 clarify_threshold | 最优阈值 |
| Preference toggle | 开/关 preference | 偏好层净收益 |

## 原则

- 离线对比，不线上双路执行
- 结果沉淀到 decision_history/
- 成熟结论可建议左花调整参数
