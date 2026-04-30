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
