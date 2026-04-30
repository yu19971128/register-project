# AGENTS.md

## 项目上下文
启动任何任务前，必须先读取最小恢复上下文:
1. `PROJECT.md` — 项目目标、技术栈、设计文档完成状态。
2. `docs/00-index.md` — 功能点到需求、模块、架构、接口、数据库、Task Card 的索引。
3. `M001-ROADMAP.md` — 当前阶段、进行中任务、阻塞项和最近完成项。
4. `DECISIONS.md` — 已确认的架构决策、例外和偏离记录。
5. `KNOWLEDGE.md` — 团队规范、技术约束和复用知识。
6. `docs/07-tasks/README.md` — Task Card 总览；若 `M001-ROADMAP.md` 或 `docs/00-index.md` 指向进行中的 Task Card，必须继续读取对应 `docs/07-tasks/TASK-*.md`。

恢复中断会话时，必须先根据 `M001-ROADMAP.md` 和 `docs/00-index.md` 判断上次断点，再只读取当前功能点相关文档；禁止为恢复上下文一次性读取整个 `docs/` 目录。

## 语言约束（🔴 强制）

AI 的所有回复、文档产出、代码注释、提交信息必须使用 **中文**。
- 禁止在正文中混入英文单词（专有名词、技术术语、代码片段除外）。
- 若用户明确使用其他语言提问，可先用该语言简短回应，但进入工作流后必须切回中文。
- 代码中的变量名、函数名保持英文，但注释必须用中文。

## 可用工作流

Claude Code 调用格式为 `/<skill-name>`。例如技能 `stack:full-dev` 应调用 `/stack:full-dev`，技能 `office-hours` 应调用 `/office-hours`。

### GSD 上下文管理
- 所有重要的决策、知识、状态变更必须同步写入对应的 Markdown 文件
- 每个任务完成后，更新 M001-ROADMAP.md 的进度

### gstack 决策层
- 产品设计诊断: `/office-hours`
- 产品审查: `/plan-ceo-review`
- 架构审查: `/plan-eng-review`
- 代码审查: `/review`
- QA 验证: `/qa`

### Superpowers 执行层
- 需求澄清: `/brainstorming`
- 制定计划: `/writing-plans`
- 执行计划: `/executing-plans`
- TDD 开发: `/test-driven-development`
- 子代理开发: `/subagent-driven-development` (仅限 Code CLI)
- 代码审查请求: `/requesting-code-review`
- 完成分支: `/finishing-a-development-branch`

### Orchestrator 复合 skill 层（新增）
- 初始化并主动澄清: `/stack:init`
- 端到端完整开发: `/stack:full-dev`
- 仅设计阶段: `/stack:design`
- 仅做模块拆分: `/stack:module-split`
- 仅实现阶段: `/stack:impl`
- Bug 修复: `/stack:bugfix`
- 小改动快速路径: `/stack:iterate`

## 工具适配规则

不同 AI 环境拥有不同的子代理/并行执行机制。AI 应根据当前运行环境，**自动映射**到对应的工具：

| 能力 | Kimi Code CLI | Claude Code / Windsurf | 通用原则 |
|------|---------------|------------------------|----------|
| 子代理调用 | `Agent` 工具 | `Task` 工具或内置 subagent | 优先使用当前环境原生支持的子代理机制 |
| 只读探索 | `explore` Agent | `Read-only` subagent | 禁止让有写权限的代理执行纯探索任务 |
| 规划代理 | `plan` Agent | `plan` subagent | 复杂任务先拆规划再执行 |
| 并行执行 | `run_in_background: true` | 环境原生并行指令 | 独立子任务应尽可能并行以提升效率 |

**核心原则**：
- 有 `Agent` 工具 → 用 `Agent`
- 有 `Task` 工具 → 用 `Task`
- 两者皆无 → 在主线程顺序执行，但须明确告知用户「当前环境不支持子代理，执行时间会延长」

## 强制开发门禁（Hard Gate）

### 规则 0：AI 必须先读 `docs/00-index.md`
启动任何开发任务前，AI **必须**先读取 `docs/00-index.md`，根据用户提到的功能点/需求名称，定位并读取相关文档。
- **禁止一次性读取整个 `docs/` 目录或无关模块的文档。**
- 若 `00-index.md` 中未找到该功能点，需引导用户补充索引或拆分 Task Card。

### 规则 1：禁止无文档开发
用户提出任何开发需求时，AI **必须**先检查 `PROJECT.md` 中的「设计文档清单」，并读取 `docs/00-index.md` 确认该功能点的关联文档链路是否完整。
- 如果存在未勾选的文档，或功能点缺少必要的关联文档，**禁止写任何代码**。
- 必须引导用户：「以下设计文档尚未完成，请先完成规划：xxx」，并调用 `/brainstorming` 或 `/writing-plans` 补全。

### 规则 2：禁止无 Task Card 的开发
任何进入开发阶段的任务，必须能在 `docs/07-tasks/` 目录下找到对应的 Task Card 文件（如 `TASK-001-xxx.md`）。
- 如果任务没有 Task Card，**禁止写任何代码**。
- 必须引导用户：「请为该任务创建 Task Card，明确验收标准和测试流程。」

### 规则 3：任务执行必须遵循闭环
每个 Task Card 的开发必须按以下顺序执行，**不能跳过**：
1. `/writing-plans` — 基于 Task Card 输出实施计划
2. `/test-driven-development` 或 `/executing-plans` — 开发实现
3. `/review` — 代码审查（必须引用 Task Card 的技术约束）
4. `/qa` — 按 Task Card 的测试流程执行验证
5. 更新 `M001-ROADMAP.md` 和 Task Card 状态为「已完成」

### 规则 4：模块拆分优先
如果用户需求涉及多个模块，AI **必须**：
- 先按模块拆分任务
- 为每个模块创建独立的 Task Card
- 明确模块间的接口契约（写入 `docs/06-api/{module}.md`），并在 `docs/00-index.md` 中登记

### 规则 5：跳过审查的强制警告
若用户明确要求跳过代码审查或测试步骤，AI 必须：
- 发出强烈警告
- 将用户的跳过请求及原因记录到 `DECISIONS.md`
- 在 `docs/07-tasks/` 的对应 Task Card 中标记「审查/测试被跳过」

### 规则 6：Orchestrator 优先原则
当用户提出开发需求时，AI 应**优先判断是否有合适的 Orchestrator 复合 skill**可用：
- 新功能 / 大重构 → 调用 `/stack:full-dev`
- 仅做模块建模 / 明确模块边界 → 调用 `/stack:module-split`
- 仅做设计规划 → 调用 `/stack:design`
- 读取已有计划做实现 → 调用 `/stack:impl`
- Bug 修复 → 调用 `/stack:bugfix`
- 小改动 / 配置调整 → 调用 `/stack:iterate`
- 若用户需求与已有 Orchestrator 工作流匹配，应先执行 Orchestrator 的 `Router` 分级，而非直接调用单一基础 skill。

### 规则 7：偏离检测与校正
AI 在执行 Orchestrator 工作流的任何阶段时，若发现**代码实现、文件路径、接口签名与 `docs/` 中的设计文档不一致**，必须立即暂停：
1. **优先修正代码**以匹配文档；
2. 若文档确实需要更新，记录变更理由到 `DECISIONS.md` 并同步更新相关文档；
3. **严禁在不记录的情况下「以代码为准」覆盖文档。**

### 规则 8：前端工作必须遵循声明的组件库
前端开发必须遵循 `docs/04-frontend/` 和 `workflow/frontend-rules.yaml` 中声明的组件/样式/图标库。若缺失，启用默认兜底规则并禁止手写基础组件。

### 规则 9：Mock 数据需确认
Mock 数据或 stub 服务需要用户明确确认，并记录到 Task Card 和 `DECISIONS.md`。

### 规则 10：状态/索引变更优先使用脚本
状态或索引变更必须优先使用 `scripts/doc-sync-v2.sh`、`scripts/update-task-status.sh` 和 `scripts/phase-exit-gate.sh`；`scripts/doc-sync.sh` 仅作为旧版兼容回退。

### 规则 11：交付需要 E2E 证据或用户批准的豁免
交付必须有端到端测试证据，或用户书面批准的豁免记录。
