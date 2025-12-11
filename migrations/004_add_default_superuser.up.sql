-- Insert default superuser for system administration
INSERT INTO users (tg_user_id, name, phone_number, tg_link, role_id)
VALUES (
    999999999,                                  -- Reserved Telegram ID for system admin
    'System Administrator',                     -- Display name
    '+79999999999',                            -- System phone number
    'https://t.me/system_admin',               -- Telegram link
    3                                          -- superuser role_id
)
ON CONFLICT (tg_user_id) DO NOTHING;

COMMENT ON TABLE users IS 'Default superuser: tg_user_id=999999999 (System Administrator)';
