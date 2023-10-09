CREATE OR REPLACE FUNCTION insert_created_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = (now() at time zone 'utc');
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_game_created_at BEFORE INSERT ON game FOR EACH ROW EXECUTE PROCEDURE  insert_created_at_column();
CREATE TRIGGER update_game_quarter_created_at BEFORE INSERT ON game_quarter_score FOR EACH ROW EXECUTE PROCEDURE  insert_created_at_column();


CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = (now() at time zone 'utc');
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_game_updated_at BEFORE UPDATE ON game FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();
CREATE TRIGGER update_game_quarter_updated_at BEFORE UPDATE ON game_quarter_score FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();

