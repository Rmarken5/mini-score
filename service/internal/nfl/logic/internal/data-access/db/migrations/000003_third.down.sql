DROP TRIGGER IF EXISTS update_game_updated_at ON game;
DROP TRIGGER IF EXISTS update_game_quarter_updated_at ON game_quarter_score;

DROP TRIGGER IF EXISTS insert_game_created_at ON game;
DROP TRIGGER IF EXISTS insert_game_quarter_created_at ON game_quarter_score;

DROP FUNCTION IF EXISTS insert_created_at_column();
DROP FUNCTION IF EXISTS update_updated_at_column();
