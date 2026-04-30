# 任务卡片：TASK-004-order

## 基础信息
- 任务 ID：TASK-004
- 所属模块：order
- 优先级：P3
- 预估工时：10-12 小时
- 创建时间：2026-04-29

## 需求描述
实现挂号订单的查询、退号、改号和状态流转，覆盖管理端（全部订单操作）和移动端（患者自助查询/退号）。引用 `docs/02-modules/order/README.md` 中的职责描述。

## 验收标准（AC）
- [ ] AC1：管理端可按日期/科室/医生/状态筛选查看全部挂号订单，支持分页
- [ ] AC2：移动端可查看自己的挂号记录，待就诊订单可申请退号
- [ ] AC3：退号时校验时限（就诊前 30 分钟以上），成功后回滚号源余量
- [ ] AC4：改号 = 退旧号 + 挂新号，新号挂号成功后旧号才释放

## 技术约束
- 后端语言/框架：Go + Gin
- 前端技术栈：React 18 + Vite，H5 使用 antd-mobile@5，管理端使用 antd@5，样式库 Tailwind CSS@3
- 数据库：SQLite
- 接口规范：引用 `docs/06-api/order.md`
- 其他约束：退号/改号需保证事务性（订单状态更新 + 号源回滚）；复用 P2 创建的 orders 表并扩展字段

## 测试流程
- [ ] 单元测试：覆盖率 ≥ 80%
- [ ] 接口测试：引用 `docs/06-api/order.md` 中的错误场景
- [ ] 集成测试：模块内各原子任务联调通过
- [ ] 端到端测试：默认必填，若无 E2E 必须记录豁免原因
- [ ] 代码审查：必须通过 `/review`
- [ ] Mock 使用：仅在用户明确确认后允许

## 状态流转
- [x] 规划中
- [ ] 开发中
- [ ] 审查中
- [ ] 测试中
- [ ] 已完成

## 原子任务清单（🔴 实现前必须确认，确认后冻结）

> 每个原子任务 = 一个可独立完成的最小功能点，必须明确引用设计文档的具体章节。
> 认知约束：每个原子任务可被人类在 3-5 分钟内理解
> 技术约束：每个原子任务代码变更不超过 3-5 个文件

| 序号 | 原子任务 ID | 任务名称 | 类型 | 输入设计文档 | 输出代码文件 | 优先级 | 状态 |
|------|------------|---------|------|-------------|-------------|--------|------|
| 1 | atom-001 | 扩展 orders 表字段 | db | `docs/05-database/schemas/order.md` | `backend/migrations/004_alter_orders.sql` | P3 | ⏳ 待实现 |
| 2 | atom-002 | 扩展 OrderRepository 查询与状态更新 | repo | `docs/02-modules/order/README.md` + `docs/05-database/schemas/order.md` | `backend/repo/order_repo.go`（追加方法）+ `order_repo_test.go` | P3 | ⏳ 待实现 |
| 3 | atom-003 | 实现 OrderService 订单业务逻辑 | svc | `docs/02-modules/order/README.md` + `docs/06-api/order.md` | `backend/service/order_service.go` + `order_service_test.go` | P3 | ⏳ 待实现 |
| 4 | atom-004 | 实现 OrderHandler REST API | api | `docs/06-api/order.md` | `backend/handler/order_handler.go` + `order_handler_test.go` | P3 | ⏳ 待实现 |
| 5 | atom-005 | 在路由中注册订单接口 | cfg | `docs/06-api/order.md` + `docs/03-architecture/README.md` | `backend/router/router.go`（追加注册） | P3 | ⏳ 待实现 |
| 6 | atom-006 | 实现 H5 挂号记录列表页 | page | `docs/04-frontend/features/order.md` | `frontend/h5/pages/OrderListPage.tsx` | P3 | ⏳ 待实现 |
| 7 | atom-007 | 实现 H5 订单详情页 | page | `docs/04-frontend/features/order.md` | `frontend/h5/pages/OrderDetailPage.tsx` | P3 | ⏳ 待实现 |
| 8 | atom-008 | 实现管理端订单列表页 | page | `docs/04-frontend/features/order.md` + `docs/04-frontend/layout.md` | `frontend/admin/pages/OrderListPage.tsx` | P3 | ⏳ 待实现 |
| 9 | atom-009 | 实现管理端订单详情页 | page | `docs/04-frontend/features/order.md` + `docs/04-frontend/layout.md` | `frontend/admin/pages/OrderDetailPage.tsx` | P3 | ⏳ 待实现 |
| 10 | atom-010 | 挂号订单管理模块集成测试 | test | 全部设计文档 | `backend/tests/order_integration_test.go` | P3 | ⏳ 待实现 |

**状态图例**：
- ⏳ 待确认 — 清单产出后等待用户确认
- ⏳ 待实现 — 用户确认后冻结
- 🔴 Red 已完成 — 失败测试已写出并提交
- 🟢 Green 已完成 — 最小实现已写出并通过测试
- ✅ 已完成 — 实现确认通过，代码已提交
- ⏸️ 已跳过 — 3 次 FAIL 后标记为 TODO

## 实现计划概要

- BDD 场景：
  1. 给定管理员已登录 PC，当查看当天挂号订单时，则列表展示全部订单且可按状态筛选
  2. 给定患者有已确认订单且就诊时间在 1 小时后，当申请退号时，则订单状态变为已退号且号源余量 +1
  3. 给定患者有已确认订单且就诊时间在 10 分钟后，当申请退号时，则返回 400 错误（已过退号时限）
- Red-Green 配对：
  - atom-001：无测试（纯 DDL）
  - atom-002：Red = `TestOrderRepository_List_UpdateStatus` / Green = `OrderRepository` 扩展方法
  - atom-003：Red = `TestOrderService_Cancel_Change_Atomic` / Green = `OrderService.Cancel` + `Change`（改号在单一事务内完成：扣减新号源 → 创建新订单 → 更新原订单 → 回滚旧号源）
  - atom-004：Red = `TestOrderHandler_List_Get_Cancel_Change` / Green = `OrderHandler` REST 接口
  - atom-005：Red = `TestOrderRoutes_Registered` / Green = `router.go` 追加路由
  - atom-006-009：前端页面通过 Storybook 或手动验证
  - atom-010：集成测试覆盖全链路
- 关键文件路径：
  - `backend/repo/order_repo.go`
  - `backend/service/order_service.go`
  - `backend/handler/order_handler.go`
  - `frontend/h5/pages/OrderListPage.tsx`
  - `frontend/admin/pages/OrderListPage.tsx`
- 验证命令：
  - `go test ./repo/...`
  - `go test ./service/...`
  - `go test ./handler/...`
  - `npm run test:unit`（前端）
- 依赖标注：
  - atom-002 依赖 atom-001（及 P2 的 order_repo.go）
  - atom-003 依赖 atom-002（及 patient/schedule/registration 模块的 Service）
  - atom-004 依赖 atom-003
  - atom-005 依赖 atom-004
  - atom-006-009 依赖 atom-005
  - atom-010 依赖 atom-001 至 atom-009

## 验收证据
- [ ] 端到端测试结果：
- [ ] 豁免说明（若无 E2E）：

## 文档引用索引

| 文档类型 | 文档路径 | 引用章节 |
|---------|---------|---------|
| 模块设计 | `docs/02-modules/order/README.md` | 职责/边界 |
| 接口契约 | `docs/06-api/order.md` | 全部接口 |
| 前端设计 | `docs/04-frontend/features/order.md` | 全部页面 |
| 数据库设计 | `docs/05-database/schemas/order.md` | 扩展字段/DDL |
| 架构设计 | `docs/03-architecture/README.md` | 技术选型/部署 |
