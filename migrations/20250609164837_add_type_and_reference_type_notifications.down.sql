ALTER TABLE notifications 
MODIFY COLUMN type ENUM('friendRequestReceived', 'friendRequestAccepted', 'tripInvitationReceived') NOT NULL,
MODIFY COLUMN reference_entity_type ENUM('system', 'friendInvitation', 'tripInvitation') NOT NULL;
