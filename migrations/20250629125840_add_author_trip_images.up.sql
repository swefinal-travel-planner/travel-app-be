ALTER TABLE trip_images ADD COLUMN author INT NOT NULL;
ALTER TABLE trip_images ADD CONSTRAINT fk_trip_images_author FOREIGN KEY (author) REFERENCES users(id) ON DELETE CASCADE;
