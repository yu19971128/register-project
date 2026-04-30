# 挂号订单管理 — API 契约

## 对外接口契约

### 接口 1：查询订单列表
- **Method**：GET
- **Path**：`/api/v1/orders`
- **鉴权**：JWT（管理端）/ 访客手机号（H5）
- **输入**：
  | 参数 | 类型 | 必填 | 说明 |
  |------|------|------|------|
  | date | string | 否 | 出诊日期，格式 YYYY-MM-DD |
  | department | string | 否 | 按科室筛选 |
  | doctor_name | string | 否 | 按医生姓名筛选 |
  | status | string | 否 | 按状态筛选：confirmed / cancelled / completed |
  | keyword | string | 否 | 按就诊人姓名/订单号搜索（管理端） |
  | page | int | 否 | 页码，默认 1 |
  | page_size | int | 否 | 每页条数，默认 10，最大 100 |
- **权限**：
  - 移动端：仅返回当前访客手机号关联的订单
  - 管理端：返回全部订单（支持 keyword 搜索）
- **输出**：
  ```json
  {
    "code": 200,
    "data": {
      "total": 100,
      "list": [
        {
          "id": 1,
          "order_no": "GH20260429001",
          "patient_name": "张三",
          "department": "内科",
          "doctor_name": "王医生",
          "date": "2026-04-29",
          "start_time": "09:00",
          "end_time": "10:00",
          "status": "confirmed",
          "created_at": "2026-04-29T09:15:00Z"
        }
      ]
    },
    "message": "ok"
  }
  ```

### 接口 2：查询订单详情
- **Method**：GET
- **Path**：`/api/v1/orders/:id`
- **鉴权**：JWT / 访客手机号
- **权限**：移动端仅限本人的订单；管理端可查看任意
- **输出**：
  ```json
  {
    "code": 200,
    "data": {
      "id": 1,
      "order_no": "GH20260429001",
      "status": "confirmed",
      "schedule": {
        "id": 1,
        "department": "内科",
        "doctor_name": "王医生",
        "date": "2026-04-29",
        "start_time": "09:00",
        "end_time": "10:00",
        "location": "1号楼 2层 内科诊室"
      },
      "patient": {
        "id": 1,
        "name": "张三",
        "gender": "male",
        "age": 32,
        "phone": "138****8888"
      },
      "visitor_phone": "13800138000",
      "created_at": "2026-04-29T09:15:00Z",
      "updated_at": "2026-04-29T09:15:00Z"
    },
    "message": "ok"
  }
  ```
- **错误场景**：
  | 错误码 | 场景 |
  |--------|------|
  | 403 | 无权查看该订单 |
  | 404 | 订单不存在 |

### 接口 3：退号
- **Method**：PUT
- **Path**：`/api/v1/orders/:id/cancel`
- **鉴权**：JWT / 访客手机号
- **权限**：移动端仅限本人的待就诊订单；管理端可退任意待就诊订单
- **处理流程**：
  1. 校验订单状态为 `confirmed`
  2. 校验当前时间在就诊开始前 30 分钟以上（退号时限）
  3. 启动事务：更新订单状态为 `cancelled` → 调用号源回滚接口
  4. 记录退号时间和原因（管理端必填，移动端可选）
- **输入**：
  | 参数 | 类型 | 必填 | 说明 |
  |------|------|------|------|
  | reason | string | 否 | 退号原因 |
- **输出**：
  ```json
  {
    "code": 200,
    "data": {
      "id": 1,
      "order_no": "GH20260429001",
      "status": "cancelled",
      "cancelled_at": "2026-04-29T08:30:00Z",
      "refund_status": "none"
    },
    "message": "ok"
  }
  ```
- **错误场景**：
  | 错误码 | 场景 |
  |--------|------|
  | 400 | 订单状态非待就诊，或已超过退号时限 |
  | 403 | 无权退号 |
  | 404 | 订单不存在 |

### 接口 4：改号
- **Method**：PUT
- **Path**：`/api/v1/orders/:id/change`
- **鉴权**：JWT / 访客手机号
- **权限**：同退号
- **处理流程**（必须在单一数据库事务内完成）：
  1. 开启事务（`BEGIN IMMEDIATE`）
  2. 校验原订单状态为 `confirmed` 且未过退号时限
  3. 对新号源执行乐观锁扣减：`UPDATE schedules SET remaining = remaining - 1 WHERE id = ? AND remaining > 0`
  4. 创建新订单记录
  5. 更新原订单状态为 `cancelled`，回滚旧号源余量（`remaining + 1`）
  6. 提交事务；若任何步骤失败则回滚整个事务
  7. 返回新订单信息
- **事务一致性说明**：改号等价于「原子化的退旧号+挂新号」。整个操作在单一事务内完成，确保不会出现「新订单已创建但旧号源未释放」的中间状态。
- **输入**：
  | 参数 | 类型 | 必填 | 说明 |
  |------|------|------|------|
  | new_schedule_id | int | 是 | 新号源 ID |
- **输出**：同「查询订单详情」，返回新订单数据
- **错误场景**：
  | 错误码 | 场景 |
  |--------|------|
  | 400 | 原订单状态不允许改号，或新号源余量为 0 |
  | 403 | 无权改号 |
  | 404 | 订单或新号源不存在 |

## 事件契约

暂无。

---

## 接口冻结声明

**冻结时间**：2026-04-29
**冻结状态**：✅ 已冻结
**说明**：本模块接口契约已通过设计审查，进入实现阶段后若需变更，必须记录理由到 DECISIONS.md 并同步更新所有引用方文档。
