ALTER TABLE trips 
MODIFY COLUMN status ENUM(
  'not_started',
  'in_progress',
  'completed',
  'cancel',
  'ai_generating',
  'failed'
) NOT NULL DEFAULT 'not_started';
