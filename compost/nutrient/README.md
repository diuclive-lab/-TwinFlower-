# 矿物质 — 可直接移植的模块

这些模块成熟、稳定、低耦合、功能清晰。验证过的代码，直接进入双生花对应层。

## 已移植

| 模块 | 来源 (66) | 目标 | 状态 |
|------|-----------|------|------|
| weather | `tools/hermes/weather.go` | `stem/tools/weather/` | ✅ Phase 1 |
| translate | `tools/hermes/translate.go` | `stem/tools/translate/` | ✅ Phase 1 |
| filesystem | `tools/filesystem/` | `stem/tools/filesystem/` | ✅ Phase 1 |
| search | `tools/hermes/web.go` | `stem/tools/search/` | ✅ Phase 1 |
| stock | `tools/hermes/stock.go` | `stem/tools/stock/` | ✅ 2026-05-18 |
| currency | `tools/hermes/currency.go` | `stem/tools/currency/` | ✅ 2026-05-18 |
| envutil | `envutil/` | `internal/envutil/` | ✅ 2026-05-18 |
| httputil | `httputil/` | `internal/httputil/` | ✅ 2026-05-18 |
| retry | `model/retry.go` | `internal/retry/` | ✅ 2026-05-18 |

## 待评估

| 模块 | 来源 (66) | 说明 | 优先级 |
|------|-----------|------|--------|
| tool registry | `tools/registry.go` | 结构化注册 + Lookup | 中 |
| execution helpers | `execution/execution_helpers.go` | args/results 处理 | 中 |
| radix tree | `tools/radixtree.go` | 轻量前缀匹配 | 低 |
| session cache | `memory/session_cache.go` | Store/Lookup/PurgeStale | 低 |
