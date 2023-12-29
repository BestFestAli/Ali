CREATE TABLE IF NOT EXISTS "Users" (
    id bigserial PRIMARY KEY ,
    created_at timestamp (0) with time zone NOT NULL DEFAULT NOW (),
    name text NOT NULL ,
    email text UNIQUE NOT NULL ,
    password_hash bytea NOT NULL ,
    activated bool NOT NULL ,
    version integer NOT NULL DEFAULT 1);

INSERT INTO "Users" (id, created_at, name, email, password_hash, activated, version)
VALUES (1, '2023-12-29 11:18:49+02','ali', 'am429602@gmail.com', '\123', TRUE, 1)
