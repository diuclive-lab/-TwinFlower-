# 菌根网络 — 架构决策记录

66 的隐藏知识散落在 commits、evals、notes、freeze logs、smoke cases 中。
这些记录保留决策语境，防止双生花重蹈覆辙。

## 决策记录格式

```json
{
  "decision": "<决策内容>",
  "why": "<为什么这样做>",
  "origin": "<66 的教训来源>",
  "context": "<决策时的约束条件>",
  "alternatives": ["<被否定的方案>"],
  "consequences": "<已知影响>"
}
```

## 已记录的决策

| 决策 | 原因 | 来源 |
|------|------|------|
| P2 封箱冻结 | 继续优化 cognition 会进入 meta-cognition overgrowth | 66 gate overgrowth |
| Skill ≠ Tool | 工具是原子能力，Skill 是怎么用工具的知识 | 66 routing complexity |
| Bypass 模式 | 确定性操作不走模型，减少 latency 和 cost | 66 over-verification |
| Translate 绕过 cognition | 避免简单任务被 preference 污染 | 66 routing complexity |
| 不建立独立验证层 | 验证层不应成为开发主导 | 66 gate explosion |
