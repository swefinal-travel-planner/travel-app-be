ALTER TABLE trip_images DROP COLUMN place_id;
ALTER TABLE trip_images ADD COLUMN trip_item_id INT NULL;
ALTER TABLE trip_images ADD CONSTRAINT fk_trip_images_trip_item FOREIGN KEY (trip_item_id) REFERENCES trip_items(id) ON DELETE CASCADE; 