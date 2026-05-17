# Phase 1：最小生根 (Rooting)

> 目标：跑通 4 个测试用例。
>
> 规则：只搬能形成闭环的代码。

---

## 搬运清单

### ① root/providers/ — 模型抽象

```
providers/
    provider.go      (接口定义)
    local.go          (本地模型 HTTP 调用)
    api.go            (API 模型调用，待定)
```

### ② runtime/engine.go — 最小循环

```
input → model.Plan → skill → tool → model.Finalize → result
```

无 shadow、无 gate、无 freeze。

### ③ stem/tools/ — 10 个原子工具

| 工具 | 来源 |
|------|------|
| weather | 66 直接搬 |
| translate | 66 直接搬 |
| stock | 66 直接搬 |
| filesystem | 66 只搬 read/list/search |
| search | 66 直接搬 |
| math | 66 直接搬 |
| browser | 保留，skill contract 约束 |
| code | 66 只搬 symbol_search + ast_slice |
| time | 轻量实现 |
| system | 轻量实现 |

### ④ vascular/skills/ — 3 个技能

| Skill | 约束 |
|-------|------|
| `tool_selection` | 意图 → 工具映射，约束禁止 |
| `business_query` | 天气/翻译/股价，确定性路由 |
| `filesystem_compare` | 只读，禁止 browser |

---

## Definition of Done

```bash
echo "北京天气" | go run .       # → 北京 25°C 晴
echo "翻译 hello" | go run .     # → 你好
echo "比较 A.go B.go" | go run . # → diff 输出，没有 browser
echo "列出文件" | go run .       # → 当前目录文件列表
```
