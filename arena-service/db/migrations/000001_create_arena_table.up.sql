CREATE TABLE IF NOT EXISTS "arenas"(
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP DEFAULT NOW(),
  name TEXT,
  description TEXT
);
