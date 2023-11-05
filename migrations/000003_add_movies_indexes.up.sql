CREATE INDEX IF NOT EXISTS foodscales_model_idx ON foodscales USING GIN (to_tsvector('simple', model));
CREATE INDEX IF NOT EXISTS foodscales_dimensions_idx ON foodscales USING GIN (dimensions);
