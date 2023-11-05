CREATE TABLE IF NOT EXISTS foodscales (
    id bigserial PRIMARY KEY ,
    model text NOT NULL ,
    specialcode bigserial NOT NULL ,
    year integer NOT NULL ,
    price integer NOT NULL ,
    runtime integer NOT NULL ,
    dimensions integer [] NOT NULL ,);
