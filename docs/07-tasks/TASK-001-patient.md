# 任务卡片：TASK-001-patient

## 基础信息
- 任务 ID：TASK-001
- 所属模块：patient
- 优先级：P1
- 预估工时：10-12 小时
- 创建时间：2026-04-29

## 需求描述
实现就诊人档案的全生命周期管理，覆盖移动端（H5）患者自助管理和管理端（PC）全局管理。引用 `docs/02-modules/patient/README.md` 中的职责描述。

## 验收标准（AC）
- [ ] AC1：移动端可完成就诊人的添加、编辑、删除，列表展示脱敏信息
- [ ] AC2：管理端可查看全部就诊人列表，支持按姓名/手机号/身份证号搜索
- [ ] AC3：身份证号/手机号加密存储，查询时脱敏展示
- [ ] AC4：存在未完成挂号的就诊人禁止删除，返回明确错误提示

## 技术约束
- 后端语言/框架：Go + Gin
- 前端技术栈：React 18 + Vite，H5 使用 antd-mobile@5，管理端使用 antd@5，样式库 Tailwind CSS@3
- 数据库：SQLite
- 接口规范：引用 `docs/06-api/patient.md`
- 其他约束：敏感数据 AES 加密，加密密钥环境变量注入；JWT 鉴权（管理端）+ 访客手机号关联（H5）

## 测试流程
- [ ] 单元测试：覆盖率 ≥ 80%
- [ ] 接口测试：引用 `docs/06-api/patient.md` 中的错误场景
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
| 1 | atom-001 | 创建 patients 表迁移与模型 | db | `docs/05-database/schemas/patient.md` | `backend/migrations/001_create_patients.sql` + `backend/models/patient.go` | P1 | ⏳ 待实现 |
| 2 | atom-002 | 实现 PatientRepository CRUD | repo | `docs/02-modules/patient/README.md` + `docs/05-database/schemas/patient.md` | `backend/repo/patient_repo.go` + `patient_repo_test.go` | P1 | ⏳ 待实现 |
| 3 | atom-003 | 实现 PatientService 业务逻辑 | svc | `docs/02-modules/patient/README.md` + `docs/06-api/patient.md` | `backend/service/patient_service.go` + `patient_service_test.go` | P1 | ⏳ 待实现 |
| 4 | atom-004 | 实现 PatientHandler REST API | api | `docs/06-api/patient.md` | `backend/handler/patient_handler.go` + `patient_handler_test.go` | P1 | ⏳ 待实现 |
| 5 | atom-005 | 配置路由与 JWT 中间件 | cfg | `docs/06-api/patient.md` + `docs/03-architecture/README.md` | `backend/router/router.go` + `backend/middleware/jwt.go` | P1 | ⏳ 待实现 |
| 6 | atom-006 | 实现 H5 就诊人列表页 | page | `docs/04-frontend/features/patient.md` | `frontend/h5/pages/PatientListPage.tsx` | P1 | ⏳ 待实现 |
| 7 | atom-007 | 实现 H5 添加/编辑就诊人页 | page | `docs/04-frontend/features/patient.md` | `frontend/h5/pages/PatientEditPage.tsx` | P1 | ⏳ 待实现 |
| 8 | atom-008 | 实现管理端就诊人列表页 | page | `docs/04-frontend/features/patient.md` + `docs/04-frontend/layout.md` | `frontend/admin/pages/PatientListPage.tsx` | P1 | ⏳ 待实现 |
| 9 | atom-009 | 实现管理端就诊人详情页 | page | `docs/04-frontend/features/patient.md` + `docs/04-frontend/layout.md` | `frontend/admin/pages/PatientDetailPage.tsx` | P1 | ⏳ 待实现 |
| 10 | atom-010 | 就诊人模块集成测试 | test | 全部设计文档 | `backend/tests/patient_integration_test.go` | P1 | ⏳ 待实现 |

**状态图例**：
- ⏳ 待确认 — 清单产出后等待用户确认
- ⏳ 待实现 — 用户确认后冻结
- 🔴 Red 已完成 — 失败测试已写出并提交
- 🟢 Green 已完成 — 最小实现已写出并通过测试
- ✅ 已完成 — 实现确认通过，代码已提交
- ⏸️ 已跳过 — 3 次 FAIL 后标记为 TODO

## 实现计划概要

- BDD 场景：
  1. 给定患者已登录 H5，当添加姓名/身份证/手机号合法的就诊人时，则成功创建并展示在列表中
  2. 给定管理员已登录 PC，当搜索身份证号时，则返回匹配的就诊人（脱敏展示）
  3. 给定就诊人存在未完成的挂号订单，当执行删除时，则返回 400 错误
- Red-Green 配对：
  - atom-001：无测试（纯 DDL）
  - atom-002：Red = `TestPatientRepository_Create_Get_Update_Delete` / Green = `PatientRepository` CRUD 方法
  - atom-003：Red = `TestPatientService_Create_With_Encryption` / Green = `PatientService` 业务方法
  - atom-004：Red = `TestPatientHandler_Create_List_Get_Update_Delete` / Green = `PatientHandler` REST 接口
  - atom-005：Red = `TestJWTAuth_Middleware` / Green = `JWTMiddleware` + `Router`
  - atom-006-009：前端页面通过 Storybook 或手动验证
  - atom-010：集成测试覆盖全链路
- 关键文件路径：
  - `backend/repo/patient_repo.go`
  - `backend/service/patient_service.go`
  - `backend/handler/patient_handler.go`
  - `frontend/h5/pages/PatientListPage.tsx`
  - `frontend/admin/pages/PatientListPage.tsx`
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
  - atom-006-009 依赖 atom-005
  - atom-010 依赖 atom-001 至 atom-009

## 验收证据
- [ ] 端到端测试结果：
- [ ] 豁免说明（若无 E2E）：

## 文档引用索引

| 文档类型 | 文档路径 | 引用章节 |
|---------|---------|---------|
| 模块设计 | `docs/02-modules/patient/README.md` | 职责/边界 |
| 接口契约 | `docs/06-api/patient.md` | 全部接口 |
| 前端设计 | `docs/04-frontend/features/patient.md` | 全部页面 |
| 数据库设计 | `docs/05-database/schemas/patient.md` | 表结构/DDL |
| 架构设计 | `docs/03-architecture/README.md` | 技术选型/部署 |
