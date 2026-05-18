# Exploration Bench

66 的 Eval Harness 在这里变成能力成长测试，不是发布门禁。

## 用途

跟踪双生花各能力域的演化轨迹。例如：

```
"查一下北京"
  week1: clarify
  week4: soft execute  
  week8: personalized execute
```

## 测试套件

| 套件 | 检测内容 | 对应 skill |
|------|---------|-----------|
| ambiguity suite | 歧义检测准确率 | search_skill |
| clarify suite | 澄清触发/消退曲线 | clarify |
| filesystem suite | 路径恢复成功率 | filesystem_skill |
| preference drift suite | 偏好漂移检测 | preferences |

## 原则

- 不设 pass/fail 门禁
- 只记录演化轨迹
- 数据供 drift/evidence 域使用
