# Cognitive Shadow

66 的 Shadow Router 思想很强，但之前长太大。在右花中保持轻量化：只做认知层的假设性对比。

## 影子实验

| 主行为 | 影子 | 问题 |
|--------|------|------|
| soft execute | 如果用 Gemma？ | 模型替换会改变 soft execute 率吗？ |
| clarify threshold=0.45 | 如果 threshold=0.55？ | clarify 率变化多少？ |
| preference active | 如果不用 preference？ | 没有偏好的基线性能？ |

## 记录格式

```json
{
  "live": "clarify",
  "shadow": "soft_execute", 
  "agreement": false,
  "timestamp": "2026-05-18T22:00:00Z"
}
```

## 原则

- 绝不参与 runtime 决策
- 只读观测，不写回 engine
- 结果供 drift 和 evidence 域分析
