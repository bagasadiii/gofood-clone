CREATE TABLE IF NOT EXISTS users (
  user_id UUID PRIMARY KEY,
  username VARCHAR(25) NOT NULL UNIQUE,
  email VARCHAR(100) NOT NULL UNIQUE,
  password VARCHAR(255),
  role VARCHAR(50) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  phone VARCHAR(20),
  balance BIGINT DEFAULT 0,
  name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS merchants (
  merchant_id UUID PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  rating REAL,
  address VARCHAR(100),
  category VARCHAR(50),
  user_id UUID NOT NULL UNIQUE,
  owner VARCHAR(25) NOT NULL,
  CONSTRAINT fk_merchant_user FOREIGN KEY(user_id)
    REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS drivers (
  driver_id UUID PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  rating REAL,
  license VARCHAR(10) NOT NULL UNIQUE,
  area VARCHAR(25),
  income INT DEFAULT 0,
  user_id UUID NOT NULL UNIQUE,
  username VARCHAR(25) NOT NULL,
  CONSTRAINT fk_driver_user FOREIGN KEY(user_id)
    REFERENCES users(user_id) ON DELETE CASCADE
);
