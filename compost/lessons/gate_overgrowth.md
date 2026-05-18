# Gate Overgrowth 教训

## 66 经历

66 从 R4 到 R14 连续 9 轮迭代全部投入验证层：smoke / freeze / gate / eval。最终 ~50K 行主线 + ~12K 行验证代码。

## 问题

验证系统过度生长，导致：

- **免疫系统定义身体**：gate 决定什么能发布，而不是产品能力决定
- **主线几乎停滞**：9 轮没有任何 workflow/memory/routing 的实质改进
- **验证层成为目的本身**：为了通过 gate 而写 gate

## 在双生花中的转化

```
66:  写代码 → 加 gate → 修 gate → 再加 gate
TF:  写代码 → 验证核心行为 → 冻结迭代
```

- Phase 2 (cognition) 完成后立即封箱，不继续优化
- 不建立独立的验证层（没有 smoke/freeze/gate/eval）
- 用 calibration logging 替代 gate

## 结论

验证不是目的，是手段。当验证层开始主导开发节奏时，就是需要 freeze 的信号。
