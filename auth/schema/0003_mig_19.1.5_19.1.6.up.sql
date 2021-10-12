BEGIN;

CREATE TABLE "login" (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE login_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE login_id_seq OWNED BY "login".id;

ALTER TABLE ONLY "login" ALTER COLUMN id SET DEFAULT nextval('login_id_seq'::regclass);

ALTER TABLE ONLY "login"
    ADD CONSTRAINT login_pkey PRIMARY KEY (id),
    ADD CONSTRAINT login_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id);

CREATE INDEX login_user_fk_index ON "login"(user_id);

CREATE TRIGGER update_login_mod_time BEFORE UPDATE
    ON "login" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

COMMIT;
