-- UP: 创建 patients 表
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

-- DOWN: 回滚
-- DROP INDEX IF EXISTS idx_patients_phone;
-- DROP INDEX IF EXISTS idx_patients_name;
-- DROP INDEX IF EXISTS idx_patients_visitor_phone;
-- DROP TABLE IF EXISTS patients;
