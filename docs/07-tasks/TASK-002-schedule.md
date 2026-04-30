# 任务卡片：TASK-002-schedule

## 基础信息
- 任务 ID：TASK-002
- 所属模块：schedule
- 优先级：P1
- 预估工时：8-10 小时
- 创建时间：2026-04-29

## 需求描述
实现当天号源的配置管理与余量控制，覆盖管理端的号源增删改查和模块间的余量扣减/回滚接口。引用 `docs/02-modules/schedule/README.md` 中的职责描述。

## 验收标准（AC）
- [ ] AC1：管理端可完成号源的添加、编辑、删除，按日期查看号源列表
- [ ] AC2：挂号时原子性扣减号源余量，余量为 0 时自动标记为已满
- [ ] AC3：退号/改号时回滚号源余量，回滚后不超过总号数
- [ ] AC4：已有预约的号源禁止删除，更新时总号数不得小于已预约数

## 技术约束
- 后端语言/框架：Go + Gin
- 前端技术栈：React 18 + Vite，管理端使用 antd@5，样式库 Tailwind CSS@3
- 数据库：SQLite
- 接口规范：引用 `docs/06-api/schedule.md`
- 其他约束：余量扣减需保证原子性（SQLite 事务）；扣减/回滚接口为模块间内部调用

## 测试流程
- [ ] 单元测试：覆盖率 ≥ 80%
- [ ] 接口测试：引用 `docs/06-api/schedule.md` 中的错误场景
- [ ] 集成测试：模块内各原子任务联调通过
- [ ] 端到端测试：默认必填，若无 E2E 必须记录豁免原因
- [ ] 代码审查：必须通过 `/review`
- [ ] Mock 使用：仅在用户明确确认后允许

## 状态流转
- [x] 规划中
- [x] 开发中
- [x] 审查中
- [x] 测试中
- [x] 已完成

## 原子任务清单（🔴 实现前必须确认，确认后冻结）

> 每个原子任务 = 一个可独立完成的最小功能点，必须明确引用设计文档的具体章节。
> 认知约束：每个原子任务可被人类在 3-5 分钟内理解
> 技术约束：每个原子任务代码变更不超过 3-5 个文件

| 序号 | 原子任务 ID | 任务名称 | 类型 | 输入设计文档 | 输出代码文件 | 优先级 | 状态 |
|------|------------|---------|------|-------------|-------------|--------|------|
| 1 | atom-001 | 创建 schedules 表迁移与模型 | db | `docs/05-database/schemas/schedule.md` | `backend/migrations/002_create_schedules.sql` + `backend/models/schedule.go` | P1 | ✅ 已完成 |
| 2 | atom-002 | 实现 ScheduleRepository CRUD | repo | `docs/02-modules/schedule/README.md` + `docs/05-database/schemas/schedule.md` | `backend/repo/schedule_repo.go` + `schedule_repo_test.go` | P1 | ✅ 已完成 |
| 3 | atom-003 | 实现 ScheduleService 含余量扣减回滚 | svc | `docs/02-modules/schedule/README.md` + `docs/06-api/schedule.md` | `backend/service/schedule_service.go` + `schedule_service_test.go` | P1 | ✅ 已完成 |
| 4 | atom-004 | 实现 ScheduleHandler REST API | api | `docs/06-api/schedule.md` | `backend/handler/schedule_handler.go` + `schedule_handler_test.go` | P1 | ✅ 已完成 |
| 5 | atom-005 | 在路由中注册号源接口 | cfg | `docs/06-api/schedule.md` + `docs/03-architecture/README.md` | `backend/router/router.go`（追加注册） | P1 | ✅ 已完成 |
| 6 | atom-006 | 实现管理端号源列表页 | page | `docs/04-frontend/features/schedule.md` + `docs/04-frontend/layout.md` | `frontend/admin/pages/ScheduleListPage.tsx` | P1 | ✅ 已完成 |
| 7 | atom-007 | 实现管理端添加/编辑号源页 | page | `docs/04-frontend/features/schedule.md` + `docs/04-frontend/layout.md` | `frontend/admin/pages/ScheduleEditPage.tsx` | P1 | ✅ 已完成 |
| 8 | atom-008 | 号源模块集成测试 | test | 全部设计文档 | `backend/tests/schedule_integration_test.go` | P1 | ✅ 已完成 |

**状态图例**：
- ⏳ 待确认 — 清单产出后等待用户确认
- ⏳ 待实现 — 用户确认后冻结
- 🔴 Red 已完成 — 失败测试已写出并提交
- 🟢 Green 已完成 — 最小实现已写出并通过测试
- ✅ 已完成 — 实现确认通过，代码已提交
- ⏸️ 已跳过 — 3 次 FAIL 后标记为 TODO

## 实现计划概要

- BDD 场景：
  1. 给定管理员配置内科王医生 09:00-10:00 总号数 20，当保存后，则号源列表展示该记录且余量为 20
  2. 给定号源余量为 1，当患者提交挂号时，则余量扣减为 0 且状态变为已满
  3. 给定号源已有 5 人预约，当管理员编辑总号数为 4 时，则返回 400 错误
- Red-Green 配对：
  - atom-001：无测试（纯 DDL）
  - atom-002：Red = `TestScheduleRepository_Create_Get_List_Update_Delete` / Green = `ScheduleRepository` CRUD 方法
  - atom-003：Red = `TestScheduleService_Deduct_And_Rollback_With_OptimisticLock` / Green = `ScheduleService.Deduct` + `Rollback`（乐观锁：UPDATE schedules SET remaining = remaining - 1 WHERE id = ? AND remaining > 0）
  - atom-004：Red = `TestScheduleHandler_Create_List_Deduct` / Green = `ScheduleHandler` REST 接口
  - atom-005：Red = `TestScheduleRoutes_Registered` / Green = `router.go` 追加路由
  - atom-006-007：前端页面通过 Storybook 或手动验证
  - atom-008：集成测试覆盖全链路
- 关键文件路径：
  - `backend/repo/schedule_repo.go`
  - `backend/service/schedule_service.go`
  - `backend/handler/schedule_handler.go`
  - `frontend/admin/pages/ScheduleListPage.tsx`
- 验证命令：
  - `go test ./repo/...`
  - `go test ./service/...`
  - `go test ./handler/...`
  - `npm run test:unit`（前端）
- 依赖标注：
  - atom-002 依赖 atom-001
  - atom-003 依赖 atom-002
  - atom-004 依赖 atom-003
  - atom-005 依赖 atom-004
  - atom-006-007 依赖 atom-005
  - atom-008 依赖 atom-001 至 atom-007

## 验收证据
- [ ] 端到端测试结果：
- [ ] 豁免说明（若无 E2E）：

## 文档引用索引

| 文档类型 | 文档路径 | 引用章节 |
|---------|---------|---------|
| 模块设计 | `docs/02-modules/schedule/README.md` | 职责/边界 |
| 接口契约 | `docs/06-api/schedule.md` | 全部接口 |
| 前端设计 | `docs/04-frontend/features/schedule.md` | 全部页面 |
| 数据库设计 | `docs/05-database/schemas/schedule.md` | 表结构/DDL |
| 架构设计 | `docs/03-architecture/README.md` | 技术选型/部署 |
