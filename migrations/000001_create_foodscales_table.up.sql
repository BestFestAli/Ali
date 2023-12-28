CREATE TABLE IF NOT EXISTS "FoodScales" (
    id bigserial PRIMARY KEY ,
    model text NOT NULL ,
    version bigserial NOT NULL ,
    year integer NOT NULL ,
    price integer NOT NULL ,
    runtime integer NOT NULL ,
    dimensions numeric [] NOT NULL );
