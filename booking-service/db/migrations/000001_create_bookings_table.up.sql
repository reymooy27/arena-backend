CREATE TABLE IF NOT EXISTS bookings(
  id SERIAL PRIMARY KEY,
  arena_id INT, 
  user_id INT,
  created_at TIMESTAMP DEFAULT NOW(),
  booking_slots VARCHAR(255)
);
