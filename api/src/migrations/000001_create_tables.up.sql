CREATE TABLE IF NOT EXISTS feedback(
   review_id VARCHAR (255) PRIMARY KEY NOT NULL,
   rating INT NOT NULL,
   pros TEXT,
   cons TEXT,
   original_is_genuine INT,
   user_is_genuine INT,
   created_by VARCHAR(255),
   created_at TIMESTAMP 
);