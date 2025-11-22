CREATE TABLE profiles (
                          id SERIAL PRIMARY KEY,
                          user_id INTEGER UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          username VARCHAR(30) UNIQUE NOT NULL,
                          display_name VARCHAR(100),
                          bio TEXT,
                          website VARCHAR(255),
                          created_at TIMESTAMP DEFAULT NOW(),
                          updated_at TIMESTAMP DEFAULT NOW()
);


CREATE INDEX idx_profiles_user_id ON profiles(user_id);
CREATE INDEX idx_profiles_username ON profiles(username);
CREATE INDEX idx_profiles_created_at ON profiles(created_at DESC);
