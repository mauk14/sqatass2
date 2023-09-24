CREATE TABLE IF NOT EXISTS receipts (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    upload_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    author text NOT NULL,
    description text NOT NULL
)