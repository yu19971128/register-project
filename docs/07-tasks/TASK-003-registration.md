# 任务卡片：TASK-003-registration

## 基础信息
- 任务 ID：TASK-003
- 所属模块：registration
- 优先级：P2
- 预估工时：10-12 小时
- 创建时间：2026-04-29

## 需求描述
实现移动端当天挂号的完整流程，包括号源浏览、就诊人选择、提交挂号和生成挂号凭证。引用 `docs/02-modules/registration/README.md` 中的职责描述。

## 验收标准（AC）
- [ ] AC1：移动端可浏览当天可用号源，按科室筛选，已满号源不可点击
- [ ] AC2：选择号源后进入确认页，可选择已有就诊人或添加新就诊人
- [ ] AC3：提交挂号后原子性扣减号源余量并创建订单，返回唯一挂号凭证
- [ ] AC4：挂号凭证页展示二维码、订单号、就诊信息和温馨提示

## 技术约束
- 后端语言/框架：Go + Gin
- 前端技术栈：React 18 + Vite，H5 使用 antd-mobile@5，样式库 Tailwind CSS@3
- 数据库：SQLite
- 接口规范：引用 `docs/06-api/registration.md`
- 其他约束：提交挂号必须保证事务性（号源扣减 + 订单创建）；订单号格式 GH + YYYYMMDD + 4位序号

## 测试流程
- [ ] 单元测试：覆盖率 ≥ 80%
- [ ] 接口测试：引用 `docs/06-api/registration.md` 中的错误场景
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
| 1 | atom-001 | 创建 orders 表迁移与模型 | db | `docs/05-database/schemas/registration.md` | `backend/migrations/003_create_orders.sql` + `backend/models/order.go` | P2 | ✅ 已完成 |
| 2 | atom-002 | 实现 OrderRepository | repo | `docs/02-modules/registration/README.md` + `docs/05-database/schemas/registration.md` | `backend/repo/order_repo.go` + `order_repo_test.go` | P2 | ✅ 已完成 |
| 3 | atom-003 | 实现 RegistrationService 挂号业务逻辑 | svc | `docs/02-modules/registration/README.md` + `docs/06-api/registration.md` | `backend/service/registration_service.go` + `registration_service_test.go` | P2 | ✅ 已完成 |
| 4 | atom-004 | 实现 RegistrationHandler REST API | api | `docs/06-api/registration.md` | `backend/handler/registration_handler.go` + `registration_handler_test.go` | P2 | ✅ 已完成 |
| 5 | atom-005 | 在路由中注册挂号接口 | cfg | `docs/06-api/registration.md` + `docs/03-architecture/README.md` | `backend/router/router.go`（追加注册） | P2 | ✅ 已完成 |
| 6 | atom-006 | 实现 H5 号源浏览页 | page | `docs/04-frontend/features/registration.md` | `frontend/h5/pages/RegisterPage.tsx` | P2 | ✅ 已完成 |
| 7 | atom-007 | 实现 H5 挂号确认页 | page | `docs/04-frontend/features/registration.md` | `frontend/h5/pages/RegisterConfirmPage.tsx` | P2 | ✅ 已完成 |
| 8 | atom-008 | 实现 H5 挂号凭证页 | page | `docs/04-frontend/features/registration.md` | `frontend/h5/pages/RegisterTicketPage.tsx` | P2 | ✅ 已完成 |
| 9 | atom-009 | 当天挂号模块集成测试 | test | 全部设计文档 | `backend/tests/registration_integration_test.go` | P2 | ✅ 已完成 |

**状态图例**：
- ⏳ 待确认 — 清单产出后等待用户确认
- ⏳ 待实现 — 用户确认后冻结
- 🔴 Red 已完成 — 失败测试已写出并提交
- 🟢 Green 已完成 — 最小实现已写出并通过测试
- ✅ 已完成 — 实现确认通过，代码已提交
- ⏸️ 已跳过 — 3 次 FAIL 后标记为 TODO

## 实现计划概要

- BDD 场景：
  1. 给定患者已添加就诊人张三，当浏览内科号源并选择王医生 09:00-10:00 后，则进入确认页并展示张三为可选就诊人
  2. 给定号源余量为 1，当患者提交挂号时，则成功创建订单 GH20260429001 且号源余量变为 0
  3. 给定号源余量已为 0，当患者点击立即挂号时，则按钮置灰不可点击
- Red-Green 配对：
  - atom-001：无测试（纯 DDL）
  - atom-002：Red = `TestOrderRepository_Create_Get` / Green = `OrderRepository` 方法
  - atom-003：Red = `TestRegistrationService_Submit_With_OrderService` / Green = `RegistrationService.SubmitRegistration`（调用 order 模块 OrderService.CreateOrder，不直接写 orders 表）
  - atom-004：Red = `TestRegistrationHandler_Submit_GetTicket` / Green = `RegistrationHandler` REST 接口
  - atom-005：Red = `TestRegistrationRoutes_Registered` / Green = `router.go` 追加路由
  - atom-006-008：前端页面通过 Storybook 或手动验证
  - atom-009：集成测试覆盖全链路
- 关键文件路径：
  - `backend/service/registration_service.go`
  - `backend/handler/registration_handler.go`
  - `frontend/h5/pages/RegisterPage.tsx`
  - `frontend/h5/pages/RegisterTicketPage.tsx`
- 验证命令：
  - `go test ./repo/...`
  - `go test ./service/...`
  - `go test ./handler/...`
  - `npm run test:unit`（前端）
- 依赖标注：
  - atom-002 依赖 atom-001
  - atom-003 依赖 atom-002（及 patient/schedule 模块的 Repo/Service）
  - atom-004 依赖 atom-003
  - atom-005 依赖 atom-004
  - atom-006-008 依赖 atom-005
  - atom-009 依赖 atom-001 至 atom-008

## 验收证据
- [ ] 端到端测试结果：
- [ ] 豁免说明（若无 E2E）：

## 文档引用索引

| 文档类型 | 文档路径 | 引用章节 |
|---------|---------|---------|
| 模块设计 | `docs/02-modules/registration/README.md` | 职责/边界 |
| 接口契约 | `docs/06-api/registration.md` | 全部接口 |
| 前端设计 | `docs/04-frontend/features/registration.md` | 全部页面 |
| 数据库设计 | `docs/05-database/schemas/registration.md` | 表结构/DDL |
| 架构设计 | `docs/03-architecture/README.md` | 技术选型/部署 |
