# 双生花开发日志

> 格式：`YYYY-MM-DD | 主题 | 关键决策 / 里程碑`

---

## 2026-05-17

**立项 · Phase 1 最小闭环跑通**

- 项目初始化：FOUNDATION.md / ROADMAP.md / STRUCTURE.md
- 关键技术决策：摒弃"Dense/MoE"二分，转向**二维认知适配矩阵**（架构层 × 行为层）
- 新增 `root/cognition/` 层 — 模型认知适配（prompt 策略、clarify 阈值、tool_bias）
- 新增 `docs/architecture/cognitive_system.md` — 认知系统架构说明
- Phase 1 第一个闭环：`"Beijing weather" → 26.2°C`（model → skill → tool → result）

**关键 commit：**
- `4f6c67c` — Phase 1 最小生根
- `6f38ed4` — 认知适配层

**当前模型：** Qwen3.6-27B（本地，port 8090）
