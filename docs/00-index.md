# 00-index — 文档导航与需求对照表

> AI 开发前必须阅读本文档，根据「功能点名称」定位关联文档，禁止直接读取无关模块文档。

## 快速定位表

| 功能点/需求 | feature_id | 优先级 | 需求文档 | 模块设计 | 前端设计 | 数据库设计 | API 文档 | 当前状态 | 任务拆解 | 备注 |
|------------|------------|--------|----------|----------|----------|------------|----------|----------|----------|------|
| 当天挂号系统 | F001 | P0 | docs/01-background/README.md | docs/02-modules/README.md | docs/04-frontend/features/patient.md, docs/04-frontend/features/schedule.md, docs/04-frontend/features/registration.md, docs/04-frontend/features/order.md | docs/05-database/schemas/patient.md, docs/05-database/schemas/schedule.md, docs/05-database/schemas/registration.md, docs/05-database/schemas/order.md | docs/06-api/README.md | 🧪 测试中 | TASK-001 ✅, TASK-002 ✅, TASK-003 ✅, TASK-004 ✅ | 含移动端H5+管理端PC；全部四个模块已实现，进入质量检查阶段 |

## 模块目录

- docs/02-modules/patient/README.md — 就诊人管理模块设计
- docs/02-modules/schedule/README.md — 号源管理模块设计
- docs/02-modules/registration/README.md — 当天挂号模块设计
- docs/02-modules/order/README.md — 挂号订单管理模块设计

## 状态图例（🔴 强制统一）

AI 在更新「当前状态」列时，**必须**从以下选项中选择，禁止自创状态词：

- 📝 需求澄清中 — 正在收集需求、确认功能范围
- 🎨 设计中 — 正在进行模块拆分、接口契约、前端/数据库设计
- ✅ 设计已完成 — 该功能点的所有设计文档已产出并确认
- 🔨 开发中 — 正在写代码、TDD 实现循环中
- 🧪 测试中 — 代码已实现，正在 QA/集成测试/联调
- 🏁 已完成 — 代码已合并，QA 通过，已发布或已交付
