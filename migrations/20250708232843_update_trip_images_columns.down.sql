ALTER TABLE trip_images DROP FOREIGN KEY fk_trip_images_trip_item;
ALTER TABLE trip_images DROP COLUMN trip_item_id;
ALTER TABLE trip_images ADD COLUMN place_id VARCHAR(255) NULL; 