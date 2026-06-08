ALTER TABLE users DROP CONSTRAINT IF EXISTS users_active_char_id_fkey;
DROP TRIGGER IF EXISTS characters_set_updated_at ON characters;
DROP TABLE IF EXISTS characters;
