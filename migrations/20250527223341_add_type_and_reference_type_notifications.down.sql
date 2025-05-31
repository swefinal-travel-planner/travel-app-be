ALTER TABLE notifications 
MODIFY COLUMN type ENUM('friendRequestReceived', 'friendRequestAccepted') NOT NULL,
MODIFY COLUMN reference_entity_type ENUM('system', 'friendInvitation') NOT NULL;
