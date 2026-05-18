# Evidence Graph

66 的 Evidence v3 在这里变成因果可视化系统。

## 追踪的关系

```
clarify_threshold
  → affects → soft_execute_rate
filesystem_skill  
  → improves → recovery_success
preference decay_lambda
  → affects → pattern_lifetime
```

## 可视化的结论

- 为什么系统变成这样？
- 哪个参数变化导致了行为漂移？
- skill 改进是否真的降低了 clarify 率？

## 原则

- 只读分析
- 数据来自 calibration/evals/drift
- 不参与决策回路
