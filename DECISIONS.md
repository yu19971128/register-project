# DECISIONS

> 记录重要技术决策及其原因，避免反复讨论。

## ADR 模板

```markdown
### [YYYY-MM-DD] 决策标题
- **背景**: 
- **决策**: 
- **原因**: 
- **影响文档**: （本决策创建或修改了哪些文档？AI 必须同步更新这些文档）
- **后果**: 
```

---

## 架构决策记录 (ADR)

### [2026-04-29] 计划审查 HIGH 问题修复
- **背景**: 阶段 2 计划审查（plan-eng-review + plan-design-review + plan-devex-review + plan-ceo-review）发现 11 个去重后 HIGH 问题，其中 3 项阻塞下一阶段。
- **决策**: 按用户选择「按建议全部修改」，逐项修复所有 HIGH 问题；MEDIUM/LOW 问题（共 57 项）暂不处理，留待实现阶段逐步修复。
- **原因**: 阻塞项涉及核心业务流程缺陷（改号事务、号源扣减竞态、数据完整性），必须在进入 TDD 实现前消除。
- **影响文档**: 
  - `docs/06-api/order.md` — 改号流程改为单一事务 + 乐观锁
  - `docs/06-api/schedule.md` — 扣减/回滚接口补充并发策略
  - `docs/06-api/patient.md` — 删除接口补充「未完成」定义；鉴权区分 H5 与管理端
  - `docs/05-database/schemas/schedule.md` — 补充 WAL + 乐观锁 + 并发上限
  - `docs/05-database/schemas/patient.md` — 唯一约束从 id_card 改为 id_card_encrypted
  - `docs/05-database/schemas/registration.md` — 明确 orders 表归属 order 模块；添加复合索引
  - `docs/05-database/schemas/order.md` — 明确 orders 表归属
  - `docs/02-modules/schedule/README.md` — 删除「仅支持当天」硬约束，改为业务聚焦当天但技术保留扩展
  - `docs/04-frontend/layout.md` — 增加范围声明，明确仅描述 PC 管理端布局
  - `docs/04-frontend/features/registration.md` — 移除定位城市和日期切换
  - `docs/03-architecture/README.md` — 补充 H5 鉴权机制 + SQLite 并发策略
  - `docs/04-frontend/README.md` — 添加 H5 登录页路由
- **后果**: 所有 HIGH 问题已清零，设计文档与实现约束一致，可进入阶段 3 TDD 实现计划。
