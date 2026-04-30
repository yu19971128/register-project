# PROJECT

## 项目概述
为医疗机构提供一套轻量化的当天挂号系统，覆盖患者移动端自助挂号和后台管理端的订单/就诊人管理。

## 核心目标
- 患者可通过 H5 自助完成当天挂号
- 管理员可通过 PC 管理端查看订单、管理就诊人、配置号源
- 轻量 SQLite 单机部署，适合中小型诊所

## 技术栈
- 语言: Go
- 框架: Gin (后端) / React 18 + Vite (前端)
- 构建工具: Vite
- 数据库: SQLite
- 部署: Docker + Nginx

## 非目标（明确不做什么）
- 不支持在线支付
- 不支持预约未来日期（仅限当天）
- 不做实名认证（仅做格式校验）
- 不做细粒度 RBAC（仅区分登录/未登录）

## 设计文档清单（开发前必须全部完成）
- [x] `docs/00-index.md` — 文档导航与需求-文档对照表
- [x] `docs/01-background/` — 项目背景、用户场景、核心痛点
- [x] `docs/02-modules/` — 模块拆分（每个模块含 overview + requirements + interfaces）
- [x] `docs/03-architecture/` — 技术架构、部署图、核心数据流
- [x] `docs/04-frontend/layout.md` — 全局布局框架设计（管理端/后台系统必做，含导航菜单、路由结构、面包屑、权限控制点）
- [x] `docs/04-frontend/features/` — 业务页面前端设计（按页面/组件拆分，Markdown 格式，非 HTML）
- [x] `docs/05-database/` — 数据库规范（按表/域拆分）
- [x] `docs/06-api/` — 接口规范（按模块/功能点拆分）
- [x] `docs/07-tasks/` — Task Card 清单（每个任务独立文件）

## 当前状态
当前阶段：设计已完成，待进入计划审查
