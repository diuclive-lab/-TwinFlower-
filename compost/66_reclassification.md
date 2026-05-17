# 66 模块重新分类

> 2026-05-17 | 从"搬不搬"到"放哪里"
>
> 四层归属：stem/tools → vascular/skills → runtime → compost

---

## 一、stem/tools/ — 原子能力

只做一件事。无流程，无约束。

| 66 模块 | 新位置 | 说明 |
|---------|--------|------|
| weather_api | `stem/tools/weather/` | 查天气，输入城市输出温度 |
| translate_api | `stem/tools/translate/` | 翻译，输入文本输出译文 |
| stock_api | `stem/tools/stock/` | 查股价 |
| currency_converter | `stem/tools/currency/` | 汇率转换 |
| filesystem_list | `stem/tools/filesystem/` | 列出目录 |
| filesystem_read | `stem/tools/filesystem/` | 读取文件 |
| filesystem_search | `stem/tools/filesystem/` | 搜索文件 |
| web_search | `stem/tools/search/` | 网络搜索 |
| web_extract | `stem/tools/search/` | 网页内容提取 |
| math_calculate | `stem/tools/math/` | 数学计算 |
| math_stats | `stem/tools/math/` | 统计计算 |
| browser_extract | `stem/tools/browser/` | 浏览器内容提取 |
| code_ast_slice | `stem/tools/code/` | 代码 AST 切片 |
| code_symbol_search | `stem/tools/code/` | 代码符号搜索 |
| code_lsp_lookup | `stem/tools/code/` | LSP 查询 |
| mcp_tools | `stem/tools/mcp/` | MCP 工具（保留接口） |

**注意：** tool 层不包含调用策略。`weather_api` 不知道什么时候该用它，也不知道失败了怎么办。

---

## 二、vascular/skills/ — 工具编排契约

核心发现层。Skill = 流程 + 约束。

| 66 模块 | 新位置 | 说明 |
|---------|--------|------|
| 业务路由 (weather/translate/stock/currency) | `vascular/skills/business_query/` | 确定性的工具选择 + 参数提取 + fallback |
| semantic fast path | `vascular/skills/light_conversation/` | 轻量闲聊直接回复，不走工具 |
| code_review_skill | `vascular/skills/code_review/` | 理解代码 → 定位 → 对比 → 建议 |
| filesystem_compare_skill | `vascular/skills/file_compare/` | 读文件 → 对比 → 输出差异 |
| repo_analysis_skill | `vascular/skills/repo_analysis/` | 扫描 → 理解结构 → 生成报告 |
| tool_selection_skill | `vascular/skills/tool_selection/` | 根据意图选择工具（替代 EvidenceRouter） |
| write_refusal_skill | `vascular/skills/safe_write/` | 写入检查 → 拒绝/允许 |
| clarify_skill | `vascular/skills/clarify/` | 信息不足时主动询问 |
| evidence_router (Top-N 部分) | `vascular/skills/tool_selection/` | 保留候选排序逻辑，去掉关键词硬编码 |
| EWMA keyword weights | `vascular/skills/tool_selection/` | 作为 skill 的学习数据，不作为路由本体 |

**Skill 的固定结构：**
```yaml
name: <skill_name>
goal: <任务目标>
allowed_tools: [工具列表]
forbidden_tools: [禁止工具]
steps: [执行步骤]
constraints: [额外约束]
success: [成功标准]
fallback: [降级策略]
```

---

## 三、runtime/ — 生命周期

不承载业务逻辑，只承载"怎么跑"。

| 66 模块 | 新位置 | 说明 |
|---------|--------|------|
| pipeline 状态机 | `runtime/pipeline/` | Run() 状态循环（纯逻辑，去门禁） |
| retry executor | `runtime/execution/` | retryWithBackoff |
| tool.Run 执行 | `runtime/execution/` | 工具执行主循环 |
| timeout / cancellation | `runtime/execution/` | Budget.TimeoutSeconds |
| session 管理 | `runtime/sessions/` | 会话生命周期 |
| state 管理 | `runtime/state/` | Context 结构体 |
| engine 入口 | `runtime/engine/` | New() + Run() |
| checkpoint | `runtime/pipeline/` | 步骤持久化（Phase 2 再启用） |

**runtime 不包含：**
- 路由策略（那是 skills）
- 工具实现（那是 stem）
- 门禁/验证（那是 observatory，Phase 3+）

---

## 四、compost/ — 留在 66

不带走，不删除。

| 模块 | 原因 |
|------|------|
| shadow router | 验证层产物 |
| router-ab-report | 验证层产物 |
| function_call_shadow | 验证层产物 |
| schema drift | 膨胀产物 |
| all eval scenarios | 验证层 |
| agent-eval cmd | 验证层 |
| project-health | 验证层 |
| doctor | 验证层 |
| freeze gate | 验证层 |
| sidecar diagnostics | 66 特定实现 |
| trace bundle | 66 特定实现 |
| trajectory recorder | 66 特定实现 |
| knowledge base | 过于特定，未来按 skill 重做 |
| LlamaGuardrail | 独立审核模型，当前不需要 |
| ETL engine | E4B 专用，当前不需要 |
| vision model | 当前不需要 |
| audio model | 当前不需要 |
| MCP 全量 | 保留接口定义，不搬实现 |

---

## 五、迁移顺序

```
Phase 1: stem/tools (10个核心工具) → runtime/engine 最小循环
Phase 2: vascular/skills (3个核心技能: business_query, tool_selection, clarify)
Phase 3: flowers/daily 日常对话可用
Phase 4: flowers/explore 探索能力
```

每阶段只动对应层的模块。不跨层。
