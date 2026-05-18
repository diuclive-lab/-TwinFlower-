# 双生花 · 路线图

> 2026-05-17 | 四阶段 | 从根到花

---

## Phase 1 — 生根

目标：模型能跑通，能聊天。

**交付：**
- 主模型调用（Qwen3.6-27B，HTTP OpenAI 协议）
- 简单 prompt 模板
- 路由模型调用（DeepSeek-R1-Distill-7B，意图分类）
- `go run .` 能对话

**不做的：**
- 工具
- 内存
- 多轮记忆
- 流式

**验收：**
```bash
echo "北京天气" | go run .
# 输出: 意图=天气查询 置信度=0.95
```

---

## Phase 2 — 长茎

目标：工具能跑通，两花共用。

**交付：**
- 工具注册（map，无框架）
- 核心工具：weather / stock / translate / file
- 工具参数校验
- 模型 → 工具 → 结果的闭环

**验收：**
```bash
echo "北京天气" | go run .
# 输出: 北京 25°C 晴
```

---

## Phase 3 — 左花（日常）

目标：日常对话可用，稳定可靠。

**交付：**
- 状态机 pipeline（控制 → 路由 → 执行 → 交付）
- DeepSeek 辅助路由（置信度 < 0.6 时咨询模型）
- 多轮记忆（最近 N 轮）
- checkpoint / resume
- 失败重试 + 降级

**验收：**
- 连续 10 轮对话不丢失上下文
- 工具调用成功/失败有反馈
- 单轮 < 5s

---

## Phase 4 — 右花（探索）

目标：探索花成为双生花的研究器官，而非验证系统。

**交付：**
- `flowers/explore/evals/` — Exploration Bench（能力成长测试，非门禁）
- `flowers/explore/shadow/` — Cognitive Shadow（假设性执行对比，不参与 runtime）
- `flowers/explore/drift/` — Preference Drift（认知漂移观测）
- `flowers/explore/evidence/` — Evidence Graph（因果可视化）
- `flowers/explore/model_lab/` — Model Lab（离线模型对比）

**不做的：**
- Release Freeze
- Gate Explosion
- Full Shadow Everything
- 线上双路执行

**验收：**
- 右花能独立运行实验
- 实验失败不影响左花日常使用
- 能力可从右花沉淀到左花 Skills
