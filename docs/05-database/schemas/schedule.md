# schemas — 号源管理

## 表名
`schedules`

## 表用途
存储当天号源信息，包括科室、医生、时间段、总号数及实时余量。

## 字段定义

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | 主键 |
| date | TEXT | NOT NULL | 出诊日期，格式 YYYY-MM-DD |
| department | TEXT | NOT NULL | 科室名称 |
| doctor_name | TEXT | NOT NULL | 医生姓名 |
| start_time | TEXT | NOT NULL | 开始时间，格式 HH:MM |
| end_time | TEXT | NOT NULL | 结束时间，格式 HH:MM |
| total_quota | INTEGER | NOT NULL, CHECK(total_quota >= 1) | 总号数 |
| remaining | INTEGER | NOT NULL, CHECK(remaining >= 0) | 剩余号数 |
| status | TEXT | NOT NULL, DEFAULT 'available' | 状态：available / full / stopped |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | 更新时间 |

## 索引

| 索引名 | 类型 | 字段 | 说明 |
|--------|------|------|------|
| pk_schedules | 主键 | id | |
| uk_schedules_date_doctor_time | 唯一 | date, doctor_name, start_time | 同一医生同一时间段唯一 |
| idx_schedules_date | 普通 | date | 按出诊日期查询 |
| idx_schedules_department | 普通 | department | 按科室筛选 |

## 命名规范检查
- 表名：小写复数 `schedules` ✅
- 字段名：小写 + 下划线 ✅
- 主键：`id` ✅
- 索引名：`pk_表名` / `uk_表名_字段名` / `idx_表名_字段名` ✅

## SQLite DDL

```sql
CREATE TABLE IF NOT EXISTS schedules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    department TEXT NOT NULL,
    doctor_name TEXT NOT NULL,
    start_time TEXT NOT NULL,
    end_time TEXT NOT NULL,
    total_quota INTEGER NOT NULL CHECK(total_quota >= 1),
    remaining INTEGER NOT NULL CHECK(remaining >= 0),
    status TEXT NOT NULL DEFAULT 'available' CHECK(status IN ('available', 'full', 'stopped')),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_schedules_date_doctor_time 
    ON schedules(date, doctor_name, start_time);
CREATE INDEX IF NOT EXISTS idx_schedules_date ON schedules(date);
CREATE INDEX IF NOT EXISTS idx_schedules_department ON schedules(department);
```

## 业务约束说明
- `remaining` 不得大于 `total_quota`
- 当 `remaining` = 0 时，`status` 应自动更新为 `full`
- 删除号源前必须检查 `remaining` < `total_quota`（即已有预约）
- 更新 `total_quota` 时，新值不得小于 `total_quota - remaining`（已预约数）

## 并发策略与乐观锁
- **SQLite 模式**：启用 WAL（Write-Ahead Logging）模式以提升并发读性能
- **写操作串行化**：通过单一数据库连接池（或 `sync.Mutex`）串行化所有写操作，避免 SQLite 的写锁竞争
- **乐观锁扣减**：号源扣减使用 `UPDATE schedules SET remaining = remaining - 1 WHERE id = ? AND remaining > 0`，通过检查影响行数判断是否扣减成功
- **并发上限预期**：单机 SQLite 约支持 50-100 TPS 的写操作，满足中小型诊所日常挂号需求
