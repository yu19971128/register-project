# 04-frontend

> 本目录存放前端设计的 **Markdown 规范文档**，包含页面/组件的布局、字体、样式约束。
> **注意**：此处不生成实际的 HTML 页面，仅输出设计规范，后续开发按此规范实现。

## 设计文档规范

### 全局布局框架设计（管理端/后台系统必做）

若项目包含 PC 管理端（后台/仪表盘/运营后台），**必须**在根目录下产出 `layout.md`：
- **文件**：`docs/04-frontend/layout.md`
- **内容**：全局布局 ASCII 图、导航菜单结构表、路由结构表、布局组件清单、布局规格、交互状态、面包屑规则、权限控制点（共 8 项强制章节，详见 `stack:design` skill）
- **约束**：布局模块文档中**不得**包含任何业务逻辑（如具体表单字段、数据表格列定义），仅定义视觉框架和导航入口

### 业务页面前端设计

每个页面或独立组件应有一个对应的 Markdown 文件，存放于 `features/` 子目录，内容至少包含：
- **页面/组件名称**
- **布局结构**（**必须用 ASCII 图示逐区域展示主内容区内的布局**）
- **首行声明**：「本页面基于全局布局框架 `docs/04-frontend/layout.md` 渲染，不独立设计导航菜单、全局路由、面包屑。」（管理端项目必填）
- **字体规范**（字号、字重、行高、字体族）
- **颜色与样式约束**（主色、辅助色、边框、圆角、阴影）
- **交互状态说明**（hover、active、disabled、loading）
- **响应式断点要求**（如有）

## 页面路由清单
| 页面 | 端 | 路由 | 权限 | 前端设计文档 |
|------|-----|------|------|-------------|
| H5 登录/手机号验证 | H5 | /h5/login | 首次验证后免登录 | docs/03-architecture/README.md（安全设计章节） |
| 就诊人列表 | H5 | /h5/patients | 访客手机号关联 | docs/04-frontend/features/patient.md |
| 添加/编辑就诊人 | H5 | /h5/patients/edit | 访客手机号关联 | docs/04-frontend/features/patient.md |
| 就诊人管理列表 | PC | /patients | JWT 管理端 | docs/04-frontend/features/patient.md |
| 就诊人详情 | PC | /patients/:id | JWT 管理端 | docs/04-frontend/features/patient.md |
| 号源配置列表 | PC | /schedules | JWT 管理端 | docs/04-frontend/features/schedule.md |
| 号源配置详情/编辑 | PC | /schedules/:id | JWT 管理端 | docs/04-frontend/features/schedule.md |
| 当天挂号首页 | H5 | /h5/register | 访客手机号关联 | docs/04-frontend/features/registration.md |
| 挂号确认 | H5 | /h5/register/confirm | 访客手机号关联 | docs/04-frontend/features/registration.md |
| 挂号凭证 | H5 | /h5/register/ticket | 访客手机号关联 | docs/04-frontend/features/registration.md |
| 挂号记录列表 | H5 | /h5/orders | 访客手机号关联 | docs/04-frontend/features/order.md |
| 订单详情（H5） | H5 | /h5/orders/:id | 访客手机号关联 | docs/04-frontend/features/order.md |
| 挂号订单列表 | PC | /orders | JWT 管理端 | docs/04-frontend/features/order.md |
| 订单详情（PC） | PC | /orders/:id | JWT 管理端 | docs/04-frontend/features/order.md |

## 模块前端设计映射

| 模块名称 | 页面/组件名 | 前端设计文档路径 |
|----------|------------|-----------------|
| 就诊人管理 | 就诊人列表(H5)、就诊人管理(PC) | docs/04-frontend/features/patient.md |
| 号源管理 | 号源配置列表/编辑(PC) | docs/04-frontend/features/schedule.md |
| 当天挂号 | 挂号首页、号源选择、挂号确认、挂号凭证(H5) | docs/04-frontend/features/registration.md |
| 挂号订单管理 | 挂号记录/订单列表、订单详情(H5+PC) | docs/04-frontend/features/order.md |

## 组件规范
- 组件命名：PascalCase
- 目录结构：按页面/功能分目录，公共组件放 `components/`
- 状态管理方案：React Context + useReducer（轻量场景）/ Zustand（复杂场景）

## UI 设计约束
- 设计系统：Ant Design
- 响应式断点：
  - ≥1440px：管理端大屏适配
  - ≥1024px：管理端标准屏
  - ≥768px：平板适配
  - <768px：H5 移动端
- 支持浏览器：Chrome ≥ 90, Safari ≥ 14, Edge ≥ 90, 微信内置浏览器

## 技术栈
- **组件库**: antd-mobile@5（H5 端）+ antd@5（管理端 PC）
- **样式库**: Tailwind CSS@3
- **图标库**: @ant-design/icons
- **构建工具**: Vite
- **状态管理**: React Context + useReducer（轻量场景）/ Zustand（复杂场景）
