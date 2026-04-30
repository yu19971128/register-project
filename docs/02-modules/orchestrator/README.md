# 模块：orchestrator — AI 自循环编排引擎

## 职责
负责在 GSD 文档约束下，主动驱动 AI 完成从需求澄清到发布的完整开发流程。

## 边界
- **不做代码生成**：代码生成交给 `subagent-driven-development` 和 `test-driven-development` skill
- **不替代文档**：文档的创建和修改由 orchestrator 触发，但关键节点需经用户确认
- **不替代用户决策**：🔴 级操作（如 git push、架构大改）必须经用户确认

## 依赖模块
- `docs/00-index.md`、`docs/01-background/`、`docs/02-modules/`、`docs/03-architecture/`
- `docs/04-frontend/`（条件）、`docs/05-database/`（条件）、`docs/06-api/`（条件）
- `docs/07-tasks/`、`AGENTS.md`

## 对外暴露
### Orchestrator 复合 skill 接口
| 命令 | 职责 |
|------|------|
| `/stack:init` | 初始化项目后主动发起需求澄清 |
| `/stack:full-dev` | 端到端完整开发 |
| `/stack:design` | 仅设计阶段 |
| `/stack:impl` | 仅实现阶段 |
| `/stack:bugfix` | Bug 修复（S1/S2/S3 分级） |
| `/stack:iterate` | 小改动快速路径 |
