
CREATE TABLE notifications (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    type ENUM('friendRequestReceived', 'friendRequestAccepted') NOT NULL,
    is_seen BOOLEAN NOT NULL DEFAULT false,
    trigger_entity_type ENUM('user', 'system') NOT NULL,
    trigger_entity_avatar VARCHAR(255) NULL,
    trigger_entity_name VARCHAR(255),
    trigger_entity_id INTEGER,
    reference_entity_type ENUM('system', 'friendInvitation') NOT NULL,
    reference_entity_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id)
);