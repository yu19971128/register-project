# 06-api

## 接口协议规范
- 协议：REST
- 数据格式：JSON
- 字符编码：UTF-8

## 鉴权方式
- 方式：JWT（管理端）+ 访客手机号（H5）
- Token 有效期：JWT 24 小时；访客手机号通过 session/localStorage 关联
- 刷新机制：JWT 支持 refresh token（可选）

## 跨模块核心 endpoint 列表

| 模块 | Method | Path | 说明 | 文档 |
|------|--------|------|------|------|
| 就诊人管理 | POST | /api/v1/patients | 创建就诊人 | docs/06-api/patient.md |
| 就诊人管理 | GET | /api/v1/patients | 查询就诊人列表 | docs/06-api/patient.md |
| 就诊人管理 | GET | /api/v1/patients/:id | 查询就诊人详情 | docs/06-api/patient.md |
| 就诊人管理 | PUT | /api/v1/patients/:id | 更新就诊人 | docs/06-api/patient.md |
| 就诊人管理 | DELETE | /api/v1/patients/:id | 删除就诊人 | docs/06-api/patient.md |
| 号源管理 | POST | /api/v1/schedules | 创建号源 | docs/06-api/schedule.md |
| 号源管理 | GET | /api/v1/schedules | 查询号源列表 | docs/06-api/schedule.md |
| 号源管理 | GET | /api/v1/schedules/:id | 查询号源详情 | docs/06-api/schedule.md |
| 号源管理 | PUT | /api/v1/schedules/:id | 更新号源 | docs/06-api/schedule.md |
| 号源管理 | DELETE | /api/v1/schedules/:id | 删除号源 | docs/06-api/schedule.md |
| 当天挂号 | POST | /api/v1/registrations | 提交挂号 | docs/06-api/registration.md |
| 当天挂号 | GET | /api/v1/registrations/ticket/:id | 查询挂号凭证 | docs/06-api/registration.md |
| 挂号订单管理 | GET | /api/v1/orders | 查询订单列表 | docs/06-api/order.md |
| 挂号订单管理 | GET | /api/v1/orders/:id | 查询订单详情 | docs/06-api/order.md |
| 挂号订单管理 | PUT | /api/v1/orders/:id/cancel | 退号 | docs/06-api/order.md |
| 挂号订单管理 | PUT | /api/v1/orders/:id/change | 改号 | docs/06-api/order.md |

## 模块 API 设计映射

| 模块名称 | 接口契约路径 | 状态 |
|----------|-------------|------|
| 就诊人管理 | docs/06-api/patient.md | ✅ 已冻结 |
| 号源管理 | docs/06-api/schedule.md | ✅ 已冻结 |
| 当天挂号 | docs/06-api/registration.md | ✅ 已冻结 |
| 挂号订单管理 | docs/06-api/order.md | ✅ 已冻结 |

## 错误码统一
| 错误码 | 含义 | 场景 |
|--------|------|------|
| 200 | 成功 | |
| 400 | 请求参数错误 | 格式校验失败、业务规则冲突（如退号时限已过） |
| 401 | 未授权 | Token 缺失或过期 |
| 403 | 禁止访问 | 无权访问该资源 |
| 404 | 资源不存在 | 就诊人/号源/订单不存在 |
| 409 | 资源冲突 | 身份证号已存在 |
| 429 | 请求过于频繁 | 限流触发 |
| 500 | 服务器内部错误 | |

## 版本策略
- URL 版本控制：/api/v1/...
- 向后兼容策略：v1 阶段不做破坏性变更，如需变更升级 v2

## 限流策略
- 单 IP 限制：挂号接口 10 次/分钟
- 单用户限制：同一访客手机号创建就诊人 20 次/小时
