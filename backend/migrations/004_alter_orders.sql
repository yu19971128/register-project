-- 扩展 orders 表字段（挂号订单管理模块）
ALTER TABLE orders ADD COLUMN cancel_reason TEXT;
ALTER TABLE orders ADD COLUMN cancelled_at DATETIME;
ALTER TABLE orders ADD COLUMN completed_at DATETIME;
ALTER TABLE orders ADD COLUMN operated_by TEXT;
