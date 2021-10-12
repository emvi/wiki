CREATE TABLE "user" (
    id bigint NOT NULL UNIQUE,
    email character varying(255) NOT NULL,
    firstname character varying(40) NOT NULL,
    lastname character varying(40) NOT NULL,
    password character varying(64),
    password_salt character varying(20),
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE user_id_seq OWNED BY "user".id;

ALTER TABLE ONLY "user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id),
    ADD CONSTRAINT user_email_unique unique (email);

CREATE OR REPLACE FUNCTION update_mod_time_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.mod_time = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_mod_time BEFORE UPDATE
    ON "user" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();
