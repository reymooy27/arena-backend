CREATE TABLE IF NOT EXISTS "users"(
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP DEFAULT NOW(),
  username TEXT,
  password TEXT
);

