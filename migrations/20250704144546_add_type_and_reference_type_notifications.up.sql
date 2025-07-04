ALTER TABLE notifications 
MODIFY COLUMN type ENUM(
    'friendRequestReceived',
    'friendRequestAccepted',
    'tripInvitationReceived',
    'tripGenerated',
    'tripGeneratedFailed',
    'tripStartingSoon'
) NOT NULL,
MODIFY COLUMN reference_entity_type ENUM(
    'system',
    'friendInvitation',
    'tripInvitation',
    'tripGeneration',
    'tripReminder'
) NOT NULL;
