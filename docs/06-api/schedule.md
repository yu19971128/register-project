# 号源管理 — API 契约

## 对外接口契约

### 接口 1：创建号源
- **Method**：POST
- **Path**：`/api/v1/schedules`
- **鉴权**：JWT（管理端）
- **输入**：
  | 参数 | 类型 | 必填 | 说明 |
  |------|------|------|------|
  | date | string | 是 | 出诊日期，格式 YYYY-MM-DD |
  | department | string | 是 | 科室名称，≤20 字符 |
  | doctor_name | string | 是 | 医生姓名，≤20 字符 |
  | start_time | string | 是 | 开始时间，格式 HH:MM |
  | end_time | string | 是 | 结束时间，格式 HH:MM |
  | total_quota | int | 是 | 总号数，≥1 |
- **输出**：
  ```json
  {
    "code": 200,
    "data": {
      "id": 1,
      "date": "2026-04-29",
      "department": "内科",
      "doctor_name": "王医生",
      "start_time": "09:00",
      "end_time": "10:00",
      "total_quota": 20,
      "remaining": 20,
      "status": "available"
    },
    "message": "ok"
  }
  ```
- **错误场景**：
  | 错误码 | 场景 |
  |--------|------|
  | 400 | 时间段格式错误或开始时间 ≥ 结束时间 |
  | 409 | 同一医生同一时间段号源已存在 |

### 接口 2：查询号源列表
- **Method**：GET
- **Path**：`/api/v1/schedules`
- **鉴权**：无（供移动端浏览）/ JWT（管理端）
- **输入**：
  | 参数 | 类型 | 必填 | 说明 |
  |------|------|------|------|
  | date | string | 否 | 出诊日期，默认当天，格式 YYYY-MM-DD |
  | department | string | 否 | 按科室筛选 |
  | page | int | 否 | 页码，默认 1 |
  | page_size | int | 否 | 每页条数，默认 20 |
- **输出**：
  ```json
  {
    "code": 200,
    "data": {
      "total": 50,
      "list": [
        {
          "id": 1,
          "date": "2026-04-29",
          "department": "内科",
          "doctor_name": "王医生",
          "start_time": "09:00",
          "end_time": "10:00",
          "total_quota": 20,
          "remaining": 15,
          "status": "available"
        }
      ]
    },
    "message": "ok"
  }
  ```

### 接口 3：查询号源详情
- **Method**：GET
- **Path**：`/api/v1/schedules/:id`
- **鉴权**：无 / JWT
- **输出**：同「创建号源」data 结构
- **错误场景**：
  | 错误码 | 场景 |
  |--------|------|
  | 404 | 号源不存在 |

### 接口 4：更新号源
- **Method**：PUT
- **Path**：`/api/v1/schedules/:id`
- **鉴权**：JWT（管理端）
- **输入**：同「创建号源」，全部字段可选
- **约束**：若已有患者预约（total_quota - remaining > 0），禁止缩小 total_quota 至已预约数以下
- **错误场景**：
  | 错误码 | 场景 |
  |--------|------|
  | 400 | 修改后总号数小于已预约数 |
  | 404 | 号源不存在 |

### 接口 5：删除号源
- **Method**：DELETE
- **Path**：`/api/v1/schedules/:id`
- **鉴权**：JWT（管理端）
- **约束**：若已有患者预约（remaining < total_quota），禁止删除
- **错误场景**：
  | 错误码 | 场景 |
  |--------|------|
  | 400 | 该号源已有预约，禁止删除 |
  | 404 | 号源不存在 |

### 接口 6：扣减号源余量（模块间调用）
- **Method**：POST
- **Path**：`/api/v1/schedules/:id/deduct`
- **鉴权**：内部调用（由当天挂号模块触发，不对外暴露给前端）
- **说明**：挂号时原子性扣减余量，remaining - 1。采用乐观锁模式防止并发超卖：
  ```sql
  UPDATE schedules
  SET remaining = remaining - 1
  WHERE id = ? AND remaining > 0;
  ```
  若影响行数为 0，则表示余量已为 0，返回 409。
- **并发策略**：SQLite 启用 WAL 模式；写操作通过连接池串行化；单实例并发上限约 100 TPS。
- **输出**：
  ```json
  {
    "code": 200,
    "data": {
      "id": 1,
      "remaining": 14,
      "deducted": true
    },
    "message": "ok"
  }
  ```
- **错误场景**：
  | 错误码 | 场景 |
  |--------|------|
  | 409 | 号源余量已为 0，扣减失败 |

### 接口 7：回滚号源余量（模块间调用）
- **Method**：POST
- **Path**：`/api/v1/schedules/:id/rollback`
- **鉴权**：内部调用（由挂号订单管理模块触发）
- **说明**：退号/改号时回滚余量，remaining + 1（不超过 total_quota）。需在订单状态更新同一事务内执行，确保原子性。
- **输出**：
  ```json
  {
    "code": 200,
    "data": {
      "id": 1,
      "remaining": 15,
      "rolled_back": true
    },
    "message": "ok"
  }
  ```

## 事件契约

暂无。

---

## 接口冻结声明

**冻结时间**：2026-04-29
**冻结状态**：✅ 已冻结
**说明**：本模块接口契约已通过设计审查，进入实现阶段后若需变更，必须记录理由到 DECISIONS.md 并同步更新所有引用方文档。
