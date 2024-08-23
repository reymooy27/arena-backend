CREATE TABLE IF NOT EXISTS payments (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP DEFAULT NOW(),
  user_id INT,
  total_amount INT,
  payment_method TEXT,
  booking_id INT
);
