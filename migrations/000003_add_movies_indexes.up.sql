CREATE INDEX IF NOT EXISTS foodscales_model_idx ON "FoodScales" USING GIN (to_tsvector('simple', model));
CREATE INDEX IF NOT EXISTS foodscales_dimensions_idx ON "FoodScales" USING GIN (dimensions);
