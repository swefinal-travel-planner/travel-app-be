ALTER TABLE trip_members
  MODIFY COLUMN role ENUM('administrator', 'member', 'staff') NOT NULL;