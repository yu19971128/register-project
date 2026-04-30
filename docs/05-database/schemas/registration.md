# schemas — 当天挂号

## 表名
`orders`

## 表用途
存储挂号订单信息。`orders` 表由「挂号订单管理」模块拥有并维护 schema；「当天挂号」模块通过调用 order 模块的服务接口来创建订单，不直接写表。

## 字段定义

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | 主键 |
| order_no | TEXT | NOT NULL, UNIQUE | 订单号，格式 GH + YYYYMMDD + 4位序号 |
| schedule_id | INTEGER | NOT NULL | 号源 ID，外键关联 schedules |
| patient_id | INTEGER | NOT NULL | 就诊人 ID，外键关联 patients |
| visitor_phone | TEXT | NOT NULL | 访客手机号（H5 端标识） |
| status | TEXT | NOT NULL, DEFAULT 'confirmed' | 状态：confirmed / cancelled / completed |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | 更新时间 |

## 索引

| 索引名 | 类型 | 字段 | 说明 |
|--------|------|------|------|
| pk_orders | 主键 | id | |
| uk_orders_order_no | 唯一 | order_no | 订单号唯一 |
| idx_orders_visitor_phone | 普通 | visitor_phone | 按访客手机号查询 |
| idx_orders_schedule_id | 普通 | schedule_id | 按号源查询 |
| idx_orders_status | 普通 | status | 按状态筛选 |
| idx_orders_status_created_at | 普通 | status, created_at | 管理端按状态+时间筛选 |
| idx_orders_query | 普通 | visitor_phone, status, created_at | 复合索引覆盖常见查询 |

## 命名规范检查
- 表名：小写复数 `orders` ✅
- 字段名：小写 + 下划线 ✅
- 主键：`id` ✅
- 索引名：`pk_表名` / `uk_表名_字段名` / `idx_表名_字段名` ✅

## SQLite DDL

```sql
CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_no TEXT NOT NULL UNIQUE,
    schedule_id INTEGER NOT NULL,
    patient_id INTEGER NOT NULL,
    visitor_phone TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'confirmed' CHECK(status IN ('confirmed', 'cancelled', 'completed')),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (schedule_id) REFERENCES schedules(id),
    FOREIGN KEY (patient_id) REFERENCES patients(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_orders_order_no ON orders(order_no);
CREATE INDEX IF NOT EXISTS idx_orders_visitor_phone ON orders(visitor_phone);
CREATE INDEX IF NOT EXISTS idx_orders_schedule_id ON orders(schedule_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_status_created_at ON orders(status, created_at);
CREATE INDEX IF NOT EXISTS idx_orders_query ON orders(visitor_phone, status, created_at);
```

## 业务约束说明
- `order_no` 必须全局唯一，建议由应用层生成（非数据库自增）
- `status` 初始值为 `confirmed`，后续由「挂号订单管理」模块更新为 `cancelled` 或 `completed`
- 外键约束确保号源和就诊人数据一致性
- 删除号源或就诊人前需检查是否有关联订单

## 扩展预留
以下字段由「挂号订单管理」模块后续通过 ALTER TABLE 或文档更新添加：
- `cancel_reason` TEXT — 退号原因
- `cancelled_at` DATETIME — 退号时间
- `completed_at` DATETIME — 就诊完成时间
