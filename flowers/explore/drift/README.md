# Preference Drift

双生花已有 preference learning + EWMA + calibration，自然需要漂移检测。

## 检测什么

| 模式 | 示例 | 含义 |
|------|------|------|
| Behavior shift | 查城市 → weather → 查城市 → news | 用户偏好变了 |
| Threshold drift | clarify 率从 18% 升到 35% | 认知模型不再匹配 |
| Skill degradation | filesystem recovery 成功率下降 | skill 需要更新 |

## 输出

```
pattern drift detected:
  old: 查一下{city} → weather (confidence 0.82)
  new: 查一下{city} → news  (confidence 0.64)
  gap: 22%
```

## 原则

- 不触发 gate
- 只通知观测台
- 数据供左花 preference 层参考
