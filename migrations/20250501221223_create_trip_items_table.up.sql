CREATE TABLE trip_items (
   id INT AUTO_INCREMENT PRIMARY KEY,
   trip_id INT NOT NULL,
   place_id VARCHAR(255) UNIQUE NOT NULL,
   trip_day INT NOT NULL,
   order_in_day INT NOT NULL,
   time_in_date ENUM('morning', 'afternoon', 'evening', 'night') NOT NULL,
   CONSTRAINT fk_trip_trip_item FOREIGN KEY (trip_id) REFERENCES trips(id) ON DELETE CASCADE,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP NULL DEFAULT NULL
);