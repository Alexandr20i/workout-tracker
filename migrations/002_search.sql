-- Добавляем tsvector колонку для полнотекстового поиска
ALTER TABLE exercises ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- Заполняем вектор из name и description
UPDATE exercises
SET search_vector = to_tsvector('russian', name || ' ' || description);

-- Индекс для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_exercises_search ON exercises USING GIN(search_vector);

-- Триггер — автоматически обновляет вектор при INSERT/UPDATE
CREATE OR REPLACE FUNCTION exercises_search_trigger() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := to_tsvector('russian', NEW.name || ' ' || NEW.description);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER exercises_search_update
    BEFORE INSERT OR UPDATE ON exercises
    FOR EACH ROW EXECUTE FUNCTION exercises_search_trigger();