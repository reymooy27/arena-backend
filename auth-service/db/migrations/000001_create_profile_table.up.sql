CREATE TABLE IF NOT EXISTS profiles(
  id INT PRIMARY KEY,
  user_id INT NOT NULL UNIQUE,
  email TEXT,
  profile_image TEXT
);
