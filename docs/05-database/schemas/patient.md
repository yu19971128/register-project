# schemas — 就诊人管理

## 表名
`patients`

## 表用途
存储就诊人档案信息，支持移动端患者自助管理和管理端全局管理。

## 字段定义

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | 主键 |
| name | TEXT | NOT NULL | 姓名 |
| id_card | TEXT | NOT NULL | 身份证号（脱敏存储，如 110101********1234） |
| id_card_encrypted | TEXT | NOT NULL, UNIQUE | 身份证号密文（AES 加密），唯一约束基于完整密文 |
| phone | TEXT | NOT NULL | 手机号（脱敏存储，如 138****8888） |
| phone_encrypted | TEXT | NOT NULL | 手机号密文（AES 加密） |
| gender | TEXT | | 性别：male / female / unknown |
| age | INTEGER | | 年龄，0-150 |
| address | TEXT | | 住址，最长 200 字符 |
| visitor_phone | TEXT | NOT NULL, INDEX | H5 端访客手机号（关联就诊人归属） |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | 更新时间 |

## 索引

| 索引名 | 类型 | 字段 | 说明 |
|--------|------|------|------|
| pk_patients | 主键 | id | |
| uk_patients_id_card_encrypted | 唯一 | id_card_encrypted | 身份证号密文唯一（避免脱敏值冲突） |
| idx_patients_visitor_phone | 普通 | visitor_phone | H5 端按访客手机号查询 |
| idx_patients_name | 普通 | name | 管理端按姓名搜索 |
| idx_patients_phone | 普通 | phone | 管理端按手机号搜索 |

## 命名规范检查
- 表名：小写复数 `patients` ✅
- 字段名：小写 + 下划线 ✅
- 主键：`id` ✅
- 索引名：`pk_表名` / `uk_表名_字段名` / `idx_表名_字段名` ✅

## SQLite DDL

```sql
CREATE TABLE IF NOT EXISTS patients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    id_card TEXT NOT NULL,
    id_card_encrypted TEXT NOT NULL UNIQUE,
    phone TEXT NOT NULL,
    phone_encrypted TEXT NOT NULL,
    gender TEXT CHECK(gender IN ('male', 'female', 'unknown')),
    age INTEGER CHECK(age >= 0 AND age <= 150),
    address TEXT,
    visitor_phone TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_patients_visitor_phone ON patients(visitor_phone);
CREATE INDEX IF NOT EXISTS idx_patients_name ON patients(name);
CREATE INDEX IF NOT EXISTS idx_patients_phone ON patients(phone);
```

## 安全说明
- 身份证号、手机号必须加密存储（`id_card_encrypted`、`phone_encrypted`）
- 查询展示时使用脱敏值（`id_card`、`phone`）
- 加密密钥通过环境变量注入，禁止硬编码
