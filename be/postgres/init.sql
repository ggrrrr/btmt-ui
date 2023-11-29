CREATE TABLE IF NOT EXISTS dev_auth (
    email TEXT,
    passwd TEXT,
    "status" TEXT,
    system_roles TEXT[],
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(email)
);