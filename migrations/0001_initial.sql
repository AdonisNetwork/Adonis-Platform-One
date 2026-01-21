CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    type VARCHAR(255),
    payload JSONB,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
