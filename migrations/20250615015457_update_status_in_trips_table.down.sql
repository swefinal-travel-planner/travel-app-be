ALTER TABLE trips 
MODIFY COLUMN status ENUM(
  'not_started',
  'in_progress',
  'completed',
  'cancel',
) NOT NULL DEFAULT 'not_started';
