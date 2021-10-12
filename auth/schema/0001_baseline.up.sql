CREATE TABLE access_grant (
    id bigint NOT NULL,
    client_id bigint NOT NULL,
    user_id bigint NOT NULL
);

CREATE SEQUENCE access_grant_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE access_grant_id_seq OWNED BY access_grant.id;

CREATE TABLE client (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    client_id character varying(20) NOT NULL,
    client_secret character varying(64) NOT NULL,
    redirect_uri character varying(500),
    trusted boolean NOT NULL,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE client_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE client_id_seq OWNED BY client.id;

CREATE TABLE scope (
    id bigint NOT NULL,
    client_id bigint NOT NULL,
    "key" character varying(40) NOT NULL,
    "value" character varying(100) NOT NULL
);

CREATE SEQUENCE scope_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE scope_id_seq OWNED BY scope.id;

CREATE TABLE "user" (
    id bigint NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(64),
    password_salt character varying(20),
    reset_password boolean NOT NULL DEFAULT FALSE,
    last_password_reset timestamp with time zone,
    firstname character varying(40),
    lastname character varying(40),
    language character varying(2),
    accept_marketing boolean NOT NULL DEFAULT FALSE,
    active boolean NOT NULL DEFAULT FALSE,
    last_login timestamp with time zone DEFAULT now(),
    registration_code character varying(40),
    registration_step integer NOT NULL DEFAULT 0,
    registration_mails_send integer NOT NULL DEFAULT 0,
    new_email character varying(255),
    new_email_code character varying(40),
    login_attempts integer NOT NULL DEFAULT 0,
    last_login_attempt timestamp with time zone DEFAULT now(),
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

CREATE TABLE "token_blacklist" (
    id bigint NOT NULL,
    token text NOT NULL,
    ttl timestamp with time zone
);

CREATE SEQUENCE token_blacklist_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE token_blacklist_id_seq OWNED BY token_blacklist.id;

CREATE TABLE "email_blacklist" (
   id bigint NOT NULL,
   domain text NOT NULL
);

CREATE SEQUENCE email_blacklist_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE email_blacklist_id_seq OWNED BY email_blacklist.id;

ALTER TABLE ONLY access_grant ALTER COLUMN id SET DEFAULT nextval('access_grant_id_seq'::regclass);

ALTER TABLE ONLY client ALTER COLUMN id SET DEFAULT nextval('client_id_seq'::regclass);

ALTER TABLE ONLY scope ALTER COLUMN id SET DEFAULT nextval('scope_id_seq'::regclass);

ALTER TABLE ONLY "user" ALTER COLUMN id SET DEFAULT nextval('user_id_seq'::regclass);

ALTER TABLE ONLY token_blacklist ALTER COLUMN id SET DEFAULT nextval('token_blacklist_id_seq'::regclass);

ALTER TABLE ONLY email_blacklist ALTER COLUMN id SET DEFAULT nextval('email_blacklist_id_seq'::regclass);

ALTER TABLE ONLY client
    ADD CONSTRAINT client_pkey PRIMARY KEY (id),
    ADD CONSTRAINT client_client_id_unique UNIQUE (client_id),
    ADD CONSTRAINT client_client_secret_unique UNIQUE (client_secret),
    ADD CONSTRAINT client_name_unique UNIQUE (name),
    ADD CONSTRAINT client_redirect_uri_unique UNIQUE (redirect_uri);

ALTER TABLE ONLY "user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id),
    ADD CONSTRAINT user_email_unique UNIQUE (email),
    ADD CONSTRAINT user_new_email_unique UNIQUE (new_email);

ALTER TABLE ONLY access_grant
    ADD CONSTRAINT access_grant_pkey PRIMARY KEY (id),
    ADD CONSTRAINT access_grant_client_fk FOREIGN KEY (client_id) REFERENCES client(id),
    ADD CONSTRAINT access_grant_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY scope
    ADD CONSTRAINT scope_pkey PRIMARY KEY (id),
    ADD CONSTRAINT scope_client_fk FOREIGN KEY (client_id) REFERENCES client(id);

ALTER TABLE ONLY token_blacklist
    ADD CONSTRAINT token_blacklist_pkey PRIMARY KEY (id);

ALTER TABLE ONLY email_blacklist
    ADD CONSTRAINT email_blacklist_pkey PRIMARY KEY (id);

CREATE INDEX access_grant_client_fk_index ON access_grant(client_id);
CREATE INDEX access_grant_user_fk_index ON access_grant(user_id);
CREATE INDEX scope_client_fk_index ON scope(client_id);
CREATE INDEX token_blacklist_token_index ON token_blacklist(token);

CREATE OR REPLACE FUNCTION update_mod_time_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.mod_time = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_client_mod_time BEFORE UPDATE
    ON "client" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_user_mod_time BEFORE UPDATE
    ON "user" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE FUNCTION token_blacklist_cleanup()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM token_blacklist WHERE ttl < NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER token_blacklist_cleanup_trigger
    AFTER INSERT ON token_blacklist EXECUTE PROCEDURE
    token_blacklist_cleanup();
