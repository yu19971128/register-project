# 当天挂号 — API 契约

## 对外接口契约

### 接口 1：提交挂号
- **Method**：POST
- **Path**：`/api/v1/registrations`
- **鉴权**：访客手机号（H5 免登录场景）
- **输入**：
  | 参数 | 类型 | 必填 | 说明 |
  |------|------|------|------|
  | schedule_id | int | 是 | 号源 ID |
  | patient_id | int | 是 | 就诊人 ID |
  | visitor_phone | string | 是 | 访客手机号，用于关联就诊人 |
- **处理流程**：
  1. 校验号源是否存在且余量 > 0
  2. 校验就诊人是否属于当前访客
  3. 启动事务：扣减号源余量 → 创建挂号订单
  4. 生成唯一订单号（格式：GH + YYYYMMDD + 4 位序号）
  5. 返回挂号凭证
- **输出**：
  ```json
  {
    "code": 200,
    "data": {
      "order_no": "GH20260429001",
      "schedule": {
        "id": 1,
        "department": "内科",
        "doctor_name": "王医生",
        "date": "2026-04-29",
        "start_time": "09:00",
        "end_time": "10:00"
      },
      "patient": {
        "id": 1,
        "name": "张三",
        "gender": "male",
        "age": 32
      },
      "status": "confirmed",
      "created_at": "2026-04-29T09:15:00Z",
      "ticket_url": "/h5/register/ticket?order_no=GH20260429001"
    },
    "message": "ok"
  }
  ```
- **错误场景**：
  | 错误码 | 场景 |
  |--------|------|
  | 400 | 号源 ID 或就诊人 ID 无效 |
  | 403 | 就诊人不属于当前访客 |
  | 409 | 号源余量已为 0 |
  | 429 | 同一访客同一号源重复提交 |

### 接口 2：查询挂号凭证
- **Method**：GET
- **Path**：`/api/v1/registrations/ticket/:order_no`
- **鉴权**：访客手机号
- **输入**：
  | 参数 | 类型 | 必填 | 说明 |
  |------|------|------|------|
  | order_no | string | 是 | 订单号，路径参数 |
- **输出**：
  ```json
  {
    "code": 200,
    "data": {
      "order_no": "GH20260429001",
      "qrcode_data": "GH20260429001",
      "department": "内科",
      "doctor_name": "王医生",
      "date": "2026-04-29",
      "start_time": "09:00",
      "end_time": "10:00",
      "patient_name": "张三",
      "patient_gender": "male",
      "patient_age": 32,
      "location": "1号楼 2层 内科诊室",
      "status": "confirmed",
      "notice": [
        "请提前 15 分钟到院取号",
        "凭此二维码或订单号就诊",
        "如需退号请提前 30 分钟"
      ]
    },
    "message": "ok"
  }
  ```
- **错误场景**：
  | 错误码 | 场景 |
  |--------|------|
  | 403 | 无权查看该凭证（非本人订单） |
  | 404 | 订单不存在 |

## 事件契约

### 事件：挂号成功（供挂号订单管理模块消费）
- **事件名**：registration.created
- **生产者**：当天挂号模块
- **消费者**：挂号订单管理模块
- **Payload**：
  ```json
  {
    "order_no": "GH20260429001",
    "schedule_id": 1,
    "patient_id": 1,
    "visitor_phone": "13800138000",
    "status": "confirmed",
    "created_at": "2026-04-29T09:15:00Z"
  }
  ```
- **说明**：当前阶段为简化实现，事件通过直接写入 orders 表替代消息队列。后续若引入 RabbitMQ/Kafka，可平滑迁移。

---

## 接口冻结声明

**冻结时间**：2026-04-29
**冻结状态**：✅ 已冻结
**说明**：本模块接口契约已通过设计审查，进入实现阶段后若需变更，必须记录理由到 DECISIONS.md 并同步更新所有引用方文档。
