# 毒素 — 禁止迁移的模式

这些是 66 的重金属沉淀。不是养分，是必须标红的反模式。

## ☣️ 毒素清单

| 毒素 | 症状 | 禁止原因 |
|------|------|----------|
| Gate Explosion | 每个功能对应一个 gate | 验证层吞噬主线 |
| Shadow Everything | 所有变更都走 shadow → agreement → cutover | 引入不必要的间接层 |
| Over-Verification | 测试行数 > 业务行数 | 收益递减 |
| Monolithic Dispatcher | 2500+ 行的 semanticdispatch | 耦合度过高，难以拆分 |
| Parallel Model Overload | 同时维护 4+ 模型角色 | 组合爆炸 |
| Premature Router Tuning | 在路由上做 A/B 测试 | 优化的是路由本身，不是用户体验 |

## 判断标准

在双生花中，如果遇到类似冲动，问三个问题：

1. 这个直接提高"模型→工具→结果"的成功率吗？
2. 没有它系统会崩溃吗？
3. 它增加的系统复杂度是否小于它解决的问题？

只要有一个"不"，就是潜在的毒素。
