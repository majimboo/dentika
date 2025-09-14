-- Update existing admin user to have super_admin role
UPDATE users
SET role = 'super_admin', clinic_id = 1
WHERE username = 'admin';

-- Verify the update
SELECT id, username, role, clinic_id FROM users WHERE username = 'admin';