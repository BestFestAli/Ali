ALTER TABLE "FoodScales" ADD CONSTRAINT foodscales_runtime_check CHECK (runtime >= 0);
ALTER TABLE "FoodScales" ADD CONSTRAINT foodscales_price_check CHECK (price BETWEEN 0 AND 1000);
ALTER TABLE "FoodScales" ADD CONSTRAINT foodscales_year_check CHECK (year BETWEEN 2000 AND date_part ('year ', now ()));
ALTER TABLE "FoodScales" ADD CONSTRAINT dimensions_length_check CHECK (array_length (dimensions , 1) BETWEEN 1 AND 5);