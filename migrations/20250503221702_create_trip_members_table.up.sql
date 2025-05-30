CREATE TABLE trip_members (
   id INT AUTO_INCREMENT PRIMARY KEY,
   trip_id INT NOT NULL,
   user_id INT NOT NULL,
   role ENUM('administrator', 'staff', 'normal_user') NOT NULL,
   CONSTRAINT fk_trip_trip_member FOREIGN KEY (trip_id) REFERENCES trips(id) ON DELETE CASCADE,
   CONSTRAINT fk_user_trip_memeber FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP NULL DEFAULT NULL
);