-- Drop the foreign key constraint first
ALTER TABLE trip_images DROP FOREIGN KEY fk_trip_images_author;

-- Then drop the column
ALTER TABLE trip_images DROP COLUMN author;
