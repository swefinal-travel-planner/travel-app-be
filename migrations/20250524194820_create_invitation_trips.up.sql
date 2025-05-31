CREATE TABLE invitation_trips (
    id INT AUTO_INCREMENT PRIMARY KEY,
    trip_id INT NOT NULL,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    status ENUM('accepted', 'rejected', 'pending') NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    CONSTRAINT fk_trip_invitation FOREIGN KEY (trip_id) REFERENCES trips(id) ON DELETE CASCADE,
    CONSTRAINT fk_sender_invitation FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_receiver_invitation FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
);