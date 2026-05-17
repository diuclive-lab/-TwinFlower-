# 双生花 · TwinFlower

> 一个以模型能力为根、通用工具为茎、日常与探索双花共生的智能体系统。
>
> 不是框架。不是聊天框。是**可移植认知系统**。

---

## 快速体验

```bash
# 启动本地模型（需要 Qwen3.6-27B GGUF）
llama-server --model qwen3.6-27b-q4_k_m.gguf --host 127.0.0.1 --port 8090

# 运行
go run . "北京天气"
# → 北京现在的温度是18.0°C。

# 模糊查询（自动澄清）
go run . "查一下杭州"
# → 猜你是想查天气（杭州）→ 18.0°C

# 高模糊（主动问）
go run . "看看苹果怎么样"
# → 你是想问水果、手机还是品牌？
```

---

## 架构

```
providers    → 传输层（本地模型 HTTP 调用）
cognition    → 认知适配层（Dense/MoE 策略 + clarify 阈值）
  skills     → 编排契约（流程 + 约束）
  preferences → 偏好学习（EWMA + 遗忘衰减 + 纠正）
    tools    → 原子能力（weather / filesystem / search / translate）
engine       → 三层决策：直接执行 → soft execute → 正式澄清
```

### 三层决策

| 级别 | 条件 | 行为 |
|------|------|------|
| Level 1 | 高置信度 | 直接执行 |
| Level 2 | 中等模糊 + 偏好匹配 | Soft Execute（可纠正） |
| Level 3 | 高模糊 | 正式 Clarify |

---

## 项目结构

| 目录 | 说明 |
|------|------|
| `root/` | 模型能力（providers / cognition / preferences） |
| `stem/` | 通用工具（weather / filesystem / search / translate） |
| `vascular/` | 技能契约（clarify / tool_selection） |
| `flowers/` | 双花（开发中） |
| `runtime/` | 引擎 + 状态管理 |
| `docs/` | 架构文档 + 认知系统说明 |

详细结构见 [`STRUCTURE.md`](STRUCTURE.md)。

---

## 路线图

| 阶段 | 状态 | 内容 |
|------|------|------|
| Phase 1 | ✅ 完成 | 最小闭环：model → skill → tool → result |
| Phase 2 | ✅ 完成 | 认知层：clarify / preference / calibration |
| Phase 3 | 🔄 开发中 | Flask 全特性（精度/轨迹/仪表盘） |
| Phase 4 | 📋 计划 | 探索花自主能力 |

---

## 设计文档

- [`FOUNDATION.md`](FOUNDATION.md) — 立项宪章、世界观、设计原则
- [`ROADMAP.md`](ROADMAP.md) — 四阶段路线图
- [`CHANGELOG.md`](CHANGELOG.md) — 完整开发日志与架构演化
- [`docs/architecture/cognitive_system.md`](docs/architecture/cognitive_system.md) — 认知系统架构

---

## 起源

双生花继承自 [66 (FangLab)](https://github.com/diuclive-lab/66) 的经验，但不继承它的结构。

66 积累了约 50K 行主线代码和 12K 行验证代码，但面临验证线吞噬主线的根本问题。双生花是重新立项——不是重构，是重新定义系统哲学。

详见 [`CHANGELOG.md`](CHANGELOG.md)。
