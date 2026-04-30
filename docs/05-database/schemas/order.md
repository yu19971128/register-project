# schemas — 挂号订单管理

## 表名
`orders`

## 表用途
`orders` 表由「挂号订单管理」模块拥有并维护 schema。虽然订单在「当天挂号」流程中创建，但表结构定义、索引设计和扩展字段均由 order 模块负责。「当天挂号」模块通过调用 order 模块的服务接口创建订单记录。

## 核心字段（已由当天挂号模块创建）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | 主键 |
| order_no | TEXT | NOT NULL, UNIQUE | 订单号 |
| schedule_id | INTEGER | NOT NULL | 号源 ID |
| patient_id | INTEGER | NOT NULL | 就诊人 ID |
| visitor_phone | TEXT | NOT NULL | 访客手机号 |
| status | TEXT | NOT NULL, DEFAULT 'confirmed' | confirmed / cancelled / completed |
| created_at | DATETIME | NOT NULL | 创建时间 |
| updated_at | DATETIME | NOT NULL | 更新时间 |

## 扩展字段（由本模块添加）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| cancel_reason | TEXT | | 退号原因 |
| cancelled_at | DATETIME | | 退号时间 |
| completed_at | DATETIME | | 就诊完成时间 |
| operated_by | TEXT | | 操作人标识（管理端用户名 / 移动端 visitor_phone） |

## SQLite DDL

```sql
-- 扩展 orders 表字段
ALTER TABLE orders ADD COLUMN cancel_reason TEXT;
ALTER TABLE orders ADD COLUMN cancelled_at DATETIME;
ALTER TABLE orders ADD COLUMN completed_at DATETIME;
ALTER TABLE orders ADD COLUMN operated_by TEXT;
```

## 业务约束说明
- `cancelled_at` 仅在 `status` = 'cancelled' 时有值
- `completed_at` 仅在 `status` = 'completed' 时有值
- `operated_by` 记录执行退号/改号/完成操作的操作人
- 退号时限校验由应用层实现（就诊开始前 30 分钟），数据库层不做限制

## 索引复用
沿用当天挂号模块已创建的索引：
- `uk_orders_order_no` — 订单号唯一
- `idx_orders_visitor_phone` — 按访客查询
- `idx_orders_schedule_id` — 按号源查询
- `idx_orders_status` — 按状态筛选
