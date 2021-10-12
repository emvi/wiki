BEGIN;

CREATE TABLE newsletter (
    id bigint NOT NULL,
    subject character varying(500) NOT NULL,
    content text NOT NULL,
    scheduled timestamp with time zone NOT NULL,
    status varchar(7) NOT NULL,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE newsletter_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE newsletter_id_seq OWNED BY newsletter.id;
ALTER TABLE ONLY newsletter ALTER COLUMN id SET DEFAULT nextval('newsletter_id_seq'::regclass);

ALTER TABLE ONLY newsletter
    ADD CONSTRAINT newsletter_pkey PRIMARY KEY (id);

CREATE TRIGGER update_newsletter_mod_time BEFORE UPDATE
    ON "newsletter" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TABLE "file" (
    id bigint NOT NULL,
    filename character varying(500) NOT NULL UNIQUE,
    original_filename character varying(500) NOT NULL,
    mime_type character varying(100) NOT NULL,
    size integer NOT NULL,
    md5 character varying(32) NOT NULL,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE file_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE file_id_seq OWNED BY file.id;
ALTER TABLE ONLY file ALTER COLUMN id SET DEFAULT nextval('file_id_seq'::regclass);

ALTER TABLE ONLY "file"
    ADD CONSTRAINT file_pkey PRIMARY KEY (id);

CREATE TRIGGER update_file_mod_time BEFORE UPDATE
    ON "file" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

COMMIT;
