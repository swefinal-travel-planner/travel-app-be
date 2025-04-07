ALTER TABLE invitation_friends
ADD COLUMN status ENUM('accepted', 'rejected', 'blocked', 'pending') NOT NULL DEFAULT 'pending';