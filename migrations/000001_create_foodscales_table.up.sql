CREATE TABLE IF NOT EXISTS foodscales (
    server_id bigserial PRIMARY KEY ,
    created_at timestamp (0) with time zone NOT NULL DEFAULT NOW (),
    model text NOT NULL ,
    special_code bigserial NOT NULL ,
    year integer NOT NULL ,
    price integer NOT NULL ,
    runtime integer NOT NULL ,
    dimensions integer [] NOT NULL ,
    version integer NOT NULL DEFAULT 1
    );
