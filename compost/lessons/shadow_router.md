# Shadow Router 教训

## 66 经历

66 实现了完整的 Function Calling Shadow Router：旁路评估、agreement 检测、schema drift、cutover gate。技术上是正确的，但代价过高。

## 问题

- 需要维护两套路由逻辑（production + shadow）
- shadow agreement 检测本身增加了复杂度
- cutover gate 变成了另一个 gate

## 在双生花中的转化

不是"有没有 shadow"，而是：

```
66: 先 shadow → 再 agreement → 再 cutover
TF: 先 soft execute → 再 preference learning → 再 implicit promotion
```

双生花的三层决策（直接执行 / soft execute / clarify）本质上解决了同一个问题——不确定时怎么办——但不需要 shadow infrastructure。

## 结论

Shadow Router 是对的意图，但用了过重的实现。双生花的 lighter-weight 版本：用 clarify threshold + preference confidence 替代 shadow + agreement。
