# Phase 2：意图认知 (Intent)

> 核心原则：**降低理解难度，而不是提升模型能力。**
>
> 当模型在多个意图之间摇摆时，不是让它猜，而是让它问。

---

## 核心发现

Phase 1 跑通了一个闭环，但它暴露了一个根本问题：

**中文自然语言的意图本来就不是确定的。**

```
"北京怎么样"
  → weather? travel? stock? filesystem?
  → 没有模型能"真正理解"，因为输入本身信息不足
```

解决方案不是更大的模型，而是：

```
降低 decision surface + 主动澄清
```

---

## 四个方向

### 1. Clarify Skill（一级路由）

把澄清从"错误恢复"升级为**正式 skill**。

```yaml
name: clarify
trigger: 置信度 < 0.55
output: 追问问题
```

当模型在多个意图之间摇摆时（最大分差 < threshold），主动问用户，而不是猜。

案例：

```
用户: "北京怎么样"
router: weather=0.42 travel=0.38 stock=0.17
→ "你是想问北京天气、旅游，还是其他内容？"
```

### 2. Intent Confidence + Evidence

每次路由输出完整的证据链，而不是只输出选中结果。

```go
type IntentEvident struct {
    Route      string
    Score      float64
    Evidence   []string   // 命中的关键词
    Missing    []string   // 缺失的关键词
    Confidence float64
}
```

### 3. Failure Learning

记录每次路由的输入、预测、实际结果、成功/失败。

积累自己的意图数据。

### 4. Chinese Intent Corpus

建立自己的中文意图数据集 `datasets/intents_zh/`。

```yaml
- input: "北京怎么样"
  expected: clarify
- input: "把hello翻成日语"
  expected: translate
- input: "看看当前目录"
  expected: filesystem
```

1000 条后，本地模型的理解力会开始明显提升。

---

## 和 Phase 1 的关系

Phase 1 闭环证明了"模型 → skill → tool → 结果"的链路通了。

Phase 2 在这个链路上加三层：
- **before routing:** 置信度评估 + 证据收集
- **at ambiguity:** Clarify Skill（问用户而不是猜）
- **after execution:** 失败记录（为 future improvement 积累数据）
