CREATE TABLE IF NOT EXISTS Activities (
    id VARCHAR(255) NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    done_at TIMESTAMP,
    duration_in_minutes INT CHECK (duration_in_minutes >= 1),
    calories_burned DECIMAL(10,2),
    activity_type VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE
);