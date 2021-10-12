CREATE TABLE "language" (
    id bigint NOT NULL UNIQUE,
    organization_id bigint NOT NULL,
    name character varying(40) NOT NULL,
    code character varying(2) NOT NULL,
    "default" boolean NOT NULL DEFAULT FALSE,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE language_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE language_id_seq OWNED BY "language".id;

CREATE TABLE "user" (
    id bigint NOT NULL UNIQUE,
    email character varying(255) NOT NULL,
    firstname character varying(40) NOT NULL,
    lastname character varying(40) NOT NULL,
    "language" character varying(2),
    info text,
    picture character varying(255),
    accept_marketing boolean NOT NULL DEFAULT FALSE,
    color_mode integer NOT NULL DEFAULT 0,
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

CREATE TABLE article (
    id bigint NOT NULL UNIQUE,
    organization_id bigint NOT NULL,
    views integer NOT NULL DEFAULT 0,
    claps integer NOT NULL DEFAULT 0,
    wip integer,
    read_everyone boolean NOT NULL,
    write_everyone boolean NOT NULL,
    private boolean NOT NULL,
    client_access boolean NOT NULL,
    archived character varying(100),
    published timestamp with time zone,
    pinned boolean NOT NULL DEFAULT FALSE,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE article_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE article_id_seq OWNED BY article.id;

CREATE TABLE article_content (
    id bigint NOT NULL UNIQUE,
    article_id bigint NOT NULL,
    language_id bigint NOT NULL,
    user_id bigint NOT NULL,
    title character varying(100) NOT NULL,
    content text NOT NULL,
    version integer NOT NULL,
    "commit" character varying(100),
    wip boolean NOT NULL DEFAULT FALSE,
    content_tsvector tsvector NOT NULL,
    title_tsvector tsvector NOT NULL,
    reading_time integer NOT NULL DEFAULT 0,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE article_content_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE article_content_id_seq OWNED BY article_content.id;

CREATE TABLE article_access (
    id bigint NOT NULL UNIQUE,
    user_id bigint,
    user_group_id bigint,
    article_id bigint NOT NULL,
    "write" boolean NOT NULL DEFAULT FALSE,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE article_access_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE article_access_id_seq OWNED BY article_access.id;

CREATE TABLE article_content_author (
    id bigint NOT NULL UNIQUE,
    article_content_id bigint,
    user_id bigint,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE article_content_author_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE article_content_author_id_seq OWNED BY article_content_author.id;

CREATE TABLE organization (
    id bigint NOT NULL UNIQUE,
    owner_user_id bigint NOT NULL,
    name character varying(60) NOT NULL,
    name_normalized character varying(20) NOT NULL,
    picture character varying(255),
    expert boolean NOT NULL DEFAULT FALSE,
    max_storage_gb integer NOT NULL,
    create_group_admin boolean NOT NULL DEFAULT FALSE,
    create_group_mod boolean NOT NULL DEFAULT FALSE,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE organization_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE organization_id_seq OWNED BY organization.id;

CREATE TABLE organization_member (
    id bigint NOT NULL UNIQUE,
    organization_id bigint NOT NULL,
    user_id bigint NOT NULL,
    language_id bigint NOT NULL,
    username character varying(20) NOT NULL,
    phone character varying(30),
    mobile character varying(30),
    info character varying(100),
    is_moderator boolean NOT NULL DEFAULT FALSE,
    is_admin boolean NOT NULL DEFAULT FALSE,
    send_notifications_interval integer NOT NULL DEFAULT 7,
    desktop_notifications boolean NOT NULL DEFAULT FALSE,
    next_notification_mail timestamp NOT NULL DEFAULT now(),
    recommendation_mail boolean NOT NULL DEFAULT TRUE,
    read_only boolean NOT NULL DEFAULT FALSE,
    active boolean NOT NULL DEFAULT TRUE,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE organization_member_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE organization_member_id_seq OWNED BY organization_member.id;

CREATE TABLE tag (
    id bigint NOT NULL UNIQUE,
    organization_id bigint NOT NULL,
    name character varying(40) NOT NULL,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE tag_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE tag_id_seq OWNED BY tag.id;

CREATE TABLE user_group (
    id bigint NOT NULL UNIQUE,
    organization_id bigint NOT NULL,
    name character varying(40) NOT NULL,
    info character varying(100),
    immutable boolean NOT NULL DEFAULT FALSE,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE user_group_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE user_group_id_seq OWNED BY user_group.id;

CREATE TABLE user_group_member (
    id bigint NOT NULL UNIQUE,
    user_group_id bigint NOT NULL,
    user_id bigint NOT NULL,
    is_moderator boolean NOT NULL DEFAULT FALSE,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE user_group_member_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE user_group_member_id_seq OWNED BY user_group_member.id;

CREATE TABLE feed (
    id bigint NOT NULL UNIQUE,
    organization_id bigint NOT NULL,
    triggered_by_user_id bigint NOT NULL,
    public boolean NOT NULL,
    reason character varying(40) NOT NULL,
    room_id character varying(40),
    deleted_name character varying(100),
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE feed_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE feed_id_seq OWNED BY feed.id;

CREATE TABLE feed_ref (
    id bigint NOT NULL UNIQUE,
    feed_id bigint NOT NULL,
    user_id bigint,
    user_group_id bigint,
    article_id bigint,
    article_content_id bigint,
    article_list_id bigint,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE feed_ref_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE feed_ref_id_seq OWNED BY feed_ref.id;

CREATE TABLE feed_access (
    id bigint NOT NULL UNIQUE,
    user_id bigint NOT NULL,
    feed_id bigint NOT NULL,
    notification boolean NOT NULL,
    read boolean NOT NULL DEFAULT FALSE,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE feed_access_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE feed_access_id_seq OWNED BY feed_access.id;

CREATE TABLE observed_object (
    id bigint NOT NULL UNIQUE,
    user_id bigint NOT NULL,
    article_id bigint,
    article_list_id bigint,
    user_group_id bigint,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE observed_object_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE observed_object_id_seq OWNED BY observed_object.id;

CREATE TABLE article_list (
    id bigint NOT NULL UNIQUE,
    organization_id bigint NOT NULL,
    public boolean NOT NULL,
    pinned boolean NOT NULL DEFAULT FALSE,
    client_access boolean NOT NULL,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE article_list_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE article_list_id_seq OWNED BY article_list.id;

CREATE TABLE article_list_entry (
    id bigint NOT NULL UNIQUE,
    article_list_id bigint NOT NULL,
    article_id bigint NOT NULL,
    position integer NOT NULL,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE article_list_entry_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE article_list_entry_id_seq OWNED BY article_list_entry.id;

CREATE TABLE article_list_member (
    id bigint NOT NULL UNIQUE,
    article_list_id bigint NOT NULL,
    user_id bigint,
    user_group_id bigint,
    is_moderator boolean NOT NULL,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE article_list_member_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE article_list_member_id_seq OWNED BY article_list_member.id;

CREATE TABLE article_list_name (
    id bigint NOT NULL UNIQUE,
    article_list_id bigint NOT NULL,
    language_id bigint NOT NULL,
    name character varying(40) NOT NULL,
    info text,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE article_list_name_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE article_list_name_id_seq OWNED BY article_list_name.id;

CREATE TABLE invitation (
    id bigint NOT NULL UNIQUE,
    organization_id bigint NOT NULL,
    user_id bigint,
    email character varying(255) NOT NULL,
    code character varying(32) NOT NULL,
    read_only boolean NOT NULL,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE invitation_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE invitation_id_seq OWNED BY invitation.id;

CREATE TABLE domain_blacklist (
    id bigint NOT NULL UNIQUE,
    name character varying(20) NOT NULL,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE domain_blacklist_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE domain_blacklist_id_seq OWNED BY domain_blacklist.id;

CREATE TABLE article_tag (
    id bigint NOT NULL UNIQUE,
    article_id bigint NOT NULL,
    tag_id bigint NOT NULL,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE article_tag_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE article_tag_id_seq OWNED BY article_tag.id;

CREATE TABLE file (
    id bigint NOT NULL,
    organization_id bigint,
    user_id bigint NOT NULL,
    article_id bigint,
    room_id character varying(255),
    language_id bigint,
    original_name character varying(255) NOT NULL,
    unique_name character varying(30) NOT NULL,
    "path" character varying(4096) NOT NULL,
    type character varying(20) NOT NULL,
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

CREATE TABLE bookmark (
    id bigint NOT NULL,
    organization_id bigint NOT NULL,
    user_id bigint NOT NULL,
    article_id bigint,
    article_list_id bigint,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE bookmark_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE bookmark_id_seq OWNED BY bookmark.id;

CREATE TABLE support_ticket (
    id bigint NOT NULL,
    organization_id bigint NOT NULL,
    user_id bigint NOT NULL,
    "type" character varying(40),
    "subject" character varying(100),
    "message" text,
    status character varying(10),
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE support_ticket_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE support_ticket_id_seq OWNED BY support_ticket.id;

CREATE TABLE client (
    id bigint NOT NULL,
    organization_id bigint NOT NULL,
    name character varying(40) NOT NULL,
    client_id character varying(20) NOT NULL,
    client_secret character varying(64) NOT NULL,
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

CREATE TABLE client_scope (
    id bigint NOT NULL,
    client_id bigint NOT NULL,
    name character varying(40) NOT NULL,
    read boolean NOT NULL DEFAULT FALSE,
    write boolean NOT NULL DEFAULT FALSE,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE client_scope_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE client_scope_id_seq OWNED BY client_scope.id;

CREATE TABLE newsletter (
    id bigint NOT NULL,
    email character varying(255) NOT NULL,
    list character varying(40),
    confirmed boolean NOT NULL DEFAULT FALSE,
    code character varying(20) NOT NULL,
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

CREATE TABLE article_visit (
     id bigint NOT NULL UNIQUE,
     article_id bigint NOT NULL,
     user_id bigint NOT NULL,
     def_time timestamp with time zone DEFAULT now(),
     mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE article_visit_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE article_visit_id_seq OWNED BY article_visit.id;

CREATE TABLE article_claps (
    id bigint NOT NULL UNIQUE,
    article_id bigint NOT NULL,
    user_id bigint NOT NULL,
    claps integer NOT NULL,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE article_claps_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE article_claps_id_seq OWNED BY article_claps.id;

CREATE TABLE article_recommendation (
    id bigint NOT NULL UNIQUE,
    article_id bigint NOT NULL,
    user_id bigint NOT NULL,
    recommended_to bigint NOT NULL,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE article_recommendation_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE article_recommendation_id_seq OWNED BY article_recommendation.id;

ALTER TABLE ONLY "language" ALTER COLUMN id SET DEFAULT nextval('language_id_seq'::regclass);

ALTER TABLE ONLY article ALTER COLUMN id SET DEFAULT nextval('article_id_seq'::regclass);

ALTER TABLE ONLY article_content ALTER COLUMN id SET DEFAULT nextval('article_content_id_seq'::regclass);

ALTER TABLE ONLY article_access ALTER COLUMN id SET DEFAULT nextval('article_access_id_seq'::regclass);

ALTER TABLE ONLY article_content_author ALTER COLUMN id SET DEFAULT nextval('article_content_author_id_seq'::regclass);

ALTER TABLE ONLY organization ALTER COLUMN id SET DEFAULT nextval('organization_id_seq'::regclass);

ALTER TABLE ONLY organization_member ALTER COLUMN id SET DEFAULT nextval('organization_member_id_seq'::regclass);

ALTER TABLE ONLY tag ALTER COLUMN id SET DEFAULT nextval('tag_id_seq'::regclass);

ALTER TABLE ONLY user_group ALTER COLUMN id SET DEFAULT nextval('user_group_id_seq'::regclass);

ALTER TABLE ONLY user_group_member ALTER COLUMN id SET DEFAULT nextval('user_group_member_id_seq'::regclass);

ALTER TABLE ONLY feed ALTER COLUMN id SET DEFAULT nextval('feed_id_seq'::regclass);

ALTER TABLE ONLY feed_ref ALTER COLUMN id SET DEFAULT nextval('feed_ref_id_seq'::regclass);

ALTER TABLE ONLY feed_access ALTER COLUMN id SET DEFAULT nextval('feed_access_id_seq'::regclass);

ALTER TABLE ONLY observed_object ALTER COLUMN id SET DEFAULT nextval('observed_object_id_seq'::regclass);

ALTER TABLE ONLY article_list ALTER COLUMN id SET DEFAULT nextval('article_list_id_seq'::regclass);

ALTER TABLE ONLY article_list_entry ALTER COLUMN id SET DEFAULT nextval('article_list_entry_id_seq'::regclass);

ALTER TABLE ONLY article_list_member ALTER COLUMN id SET DEFAULT nextval('article_list_member_id_seq'::regclass);

ALTER TABLE ONLY article_list_name ALTER COLUMN id SET DEFAULT nextval('article_list_name_id_seq'::regclass);

ALTER TABLE ONLY invitation ALTER COLUMN id SET DEFAULT nextval('invitation_id_seq'::regclass);

ALTER TABLE ONLY domain_blacklist ALTER COLUMN id SET DEFAULT nextval('domain_blacklist_id_seq'::regclass);

ALTER TABLE ONLY article_tag ALTER COLUMN id SET DEFAULT nextval('article_tag_id_seq'::regclass);

ALTER TABLE ONLY file ALTER COLUMN id SET DEFAULT nextval('file_id_seq'::regclass);

ALTER TABLE ONLY bookmark ALTER COLUMN id SET DEFAULT nextval('bookmark_id_seq'::regclass);

ALTER TABLE ONLY support_ticket ALTER COLUMN id SET DEFAULT nextval('support_ticket_id_seq'::regclass);

ALTER TABLE ONLY client ALTER COLUMN id SET DEFAULT nextval('client_id_seq'::regclass);

ALTER TABLE ONLY client_scope ALTER COLUMN id SET DEFAULT nextval('client_scope_id_seq'::regclass);

ALTER TABLE ONLY newsletter ALTER COLUMN id SET DEFAULT nextval('newsletter_id_seq'::regclass);

ALTER TABLE ONLY article_visit ALTER COLUMN id SET DEFAULT nextval('article_visit_id_seq'::regclass);

ALTER TABLE ONLY article_claps ALTER COLUMN id SET DEFAULT nextval('article_claps_id_seq'::regclass);

ALTER TABLE ONLY article_recommendation ALTER COLUMN id SET DEFAULT nextval('article_recommendation_id_seq'::regclass);

ALTER TABLE ONLY "language"
    ADD CONSTRAINT language_pkey PRIMARY KEY (id),
    ADD CONSTRAINT language_organization_fk FOREIGN KEY (organization_id) REFERENCES organization(id);

ALTER TABLE ONLY "user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id),
    ADD CONSTRAINT user_email_unique UNIQUE (email);

ALTER TABLE ONLY organization
    ADD CONSTRAINT organization_pkey PRIMARY KEY (id),
    ADD CONSTRAINT organization_name_normalized_unique UNIQUE (name_normalized),
    ADD CONSTRAINT organization_owner_fk FOREIGN KEY (owner_user_id) REFERENCES "user"(id);

ALTER TABLE ONLY organization_member
    ADD CONSTRAINT organization_member_pkey PRIMARY KEY (id),
    ADD CONSTRAINT organization_member_organization_fk FOREIGN KEY (organization_id) REFERENCES organization(id),
    ADD CONSTRAINT organization_member_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id),
    ADD CONSTRAINT organization_member_language_fk FOREIGN KEY (language_id) REFERENCES language(id);

ALTER TABLE ONLY article
    ADD CONSTRAINT article_pkey PRIMARY KEY (id),
    ADD CONSTRAINT article_organization_fk FOREIGN KEY (organization_id) REFERENCES organization(id);

ALTER TABLE ONLY article_content
    ADD CONSTRAINT article_content_pkey PRIMARY KEY (id),
    ADD CONSTRAINT article_content_article_fk FOREIGN KEY (article_id) REFERENCES article(id),
    ADD CONSTRAINT article_content_language_fk FOREIGN KEY (language_id) REFERENCES language(id),
    ADD CONSTRAINT article_content_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY article_access
    ADD CONSTRAINT article_access_pkey PRIMARY KEY (id),
    ADD CONSTRAINT article_access_article_fk FOREIGN KEY (article_id) REFERENCES article(id),
    ADD CONSTRAINT article_access_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id),
    ADD CONSTRAINT article_access_user_group_fk FOREIGN KEY (user_group_id) REFERENCES user_group(id);

ALTER TABLE ONLY article_content_author
    ADD CONSTRAINT article_content_author_pkey PRIMARY KEY (id),
    ADD CONSTRAINT article_content_author_article_content_fk FOREIGN KEY (article_content_id) REFERENCES article_content(id),
    ADD CONSTRAINT article_content_author_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY tag
    ADD CONSTRAINT tag_pkey PRIMARY KEY (id),
    ADD CONSTRAINT tag_organization_fk FOREIGN KEY (organization_id) REFERENCES organization(id);

ALTER TABLE ONLY user_group
    ADD CONSTRAINT user_group_pkey PRIMARY KEY (id),
    ADD CONSTRAINT user_group_organization_fk FOREIGN KEY (organization_id) REFERENCES organization(id);

ALTER TABLE ONLY user_group_member
    ADD CONSTRAINT user_group_member_pkey PRIMARY KEY (id),
    ADD CONSTRAINT user_group_member_user_group_fk FOREIGN KEY (user_group_id) REFERENCES user_group(id),
    ADD CONSTRAINT user_group_member_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY feed
    ADD CONSTRAINT feed_pkey PRIMARY KEY (id),
    ADD CONSTRAINT feed_organization_fk FOREIGN KEY (organization_id) REFERENCES organization(id),
    ADD CONSTRAINT feed_triggered_by_user_fk FOREIGN KEY (triggered_by_user_id) REFERENCES "user"(id);

ALTER TABLE ONLY feed_ref
    ADD CONSTRAINT feed_ref_pkey PRIMARY KEY (id),
    ADD CONSTRAINT feed_ref_feed_fk FOREIGN KEY (feed_id) REFERENCES "feed"(id) ON DELETE CASCADE,
    ADD CONSTRAINT feed_ref_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id),
    ADD CONSTRAINT feed_ref_user_group_fk FOREIGN KEY (user_group_id) REFERENCES user_group(id),
    ADD CONSTRAINT feed_ref_article_fk FOREIGN KEY (article_id) REFERENCES article(id),
    ADD CONSTRAINT feed_ref_article_content_fk FOREIGN KEY (article_content_id) REFERENCES article_content(id),
    ADD CONSTRAINT feed_ref_article_list_fk FOREIGN KEY (article_list_id) REFERENCES article_list(id);

ALTER TABLE ONLY feed_access
    ADD CONSTRAINT feed_access_pkey PRIMARY KEY (id),
    ADD CONSTRAINT feed_access_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id),
    ADD CONSTRAINT feed_access_feed_fk FOREIGN KEY (feed_id) REFERENCES feed(id) ON DELETE CASCADE;

ALTER TABLE ONLY observed_object
    ADD CONSTRAINT observed_object_pkey PRIMARY KEY (id),
    ADD CONSTRAINT observed_object_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id),
    ADD CONSTRAINT observed_object_article_fk FOREIGN KEY (article_id) REFERENCES "article"(id),
    ADD CONSTRAINT observed_object_article_list_fk FOREIGN KEY (article_list_id) REFERENCES "article_list"(id),
    ADD CONSTRAINT observed_object_user_group_fk FOREIGN KEY (user_group_id) REFERENCES "user_group"(id);

ALTER TABLE ONLY article_list
    ADD CONSTRAINT article_list_pkey PRIMARY KEY (id),
    ADD CONSTRAINT article_list_organization_fk FOREIGN KEY (organization_id) REFERENCES organization(id);

ALTER TABLE ONLY article_list_entry
    ADD CONSTRAINT article_list_entry_pkey PRIMARY KEY (id),
    ADD CONSTRAINT article_list_entry_article_list_fk FOREIGN KEY (article_list_id) REFERENCES article_list(id),
    ADD CONSTRAINT article_list_entry_article_fk FOREIGN KEY (article_id) REFERENCES article(id);

ALTER TABLE ONLY article_list_member
    ADD CONSTRAINT article_list_member_pkey PRIMARY KEY (id),
    ADD CONSTRAINT article_list_member_article_list_fk FOREIGN KEY (article_list_id) REFERENCES article_list(id),
    ADD CONSTRAINT article_list_member_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY article_list_name
    ADD CONSTRAINT article_list_name_pkey PRIMARY KEY (id),
    ADD CONSTRAINT article_list_name_article_list_fk FOREIGN KEY (article_list_id) REFERENCES article_list(id),
    ADD CONSTRAINT article_list_name_language_id_fk FOREIGN KEY (language_id) REFERENCES language(id);

ALTER TABLE ONLY invitation
    ADD CONSTRAINT invitation_pkey PRIMARY KEY (id),
    ADD CONSTRAINT invitation_organization_fk FOREIGN KEY (organization_id) REFERENCES organization(id),
    ADD CONSTRAINT invitation_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY domain_blacklist
    ADD CONSTRAINT domain_blacklist_pkey PRIMARY KEY (id),
    ADD CONSTRAINT domain_blacklist_name_unique unique ("name");

ALTER TABLE ONLY article_tag
    ADD CONSTRAINT article_tag_pkey PRIMARY KEY (id),
    ADD CONSTRAINT article_tag_article_fk FOREIGN KEY (article_id) REFERENCES article(id),
    ADD CONSTRAINT article_tag_tag_fk FOREIGN KEY (tag_id) REFERENCES tag(id);

ALTER TABLE ONLY file
    ADD CONSTRAINT file_pkey PRIMARY KEY (id),
    ADD CONSTRAINT file_organization_fk FOREIGN KEY (organization_id) REFERENCES organization(id),
    ADD CONSTRAINT file_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id),
    ADD CONSTRAINT file_article_fk FOREIGN KEY (article_id) references article(id),
    ADD CONSTRAINT file_language_fk FOREIGN KEY (language_id) references language(id);

ALTER TABLE ONLY bookmark
    ADD CONSTRAINT bookmark_pkey PRIMARY KEY (id),
    ADD CONSTRAINT bookmark_organization_fk FOREIGN KEY (organization_id) REFERENCES organization(id),
    ADD CONSTRAINT bookmark_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id),
    ADD CONSTRAINT bookmark_article_fk FOREIGN KEY (article_id) references article(id),
    ADD CONSTRAINT bookmark_article_list_fk FOREIGN KEY (article_list_id) references article_list(id);

ALTER TABLE ONLY support_ticket
    ADD CONSTRAINT support_ticket_pkey PRIMARY KEY (id),
    ADD CONSTRAINT support_ticket_organization_fk FOREIGN KEY (organization_id) REFERENCES organization(id),
    ADD CONSTRAINT support_ticket_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY client
    ADD CONSTRAINT client_pkey PRIMARY KEY (id),
    ADD CONSTRAINT client_organization_fk FOREIGN KEY (organization_id) REFERENCES organization(id),
    ADD CONSTRAINT client_client_id_unique UNIQUE (client_id);

ALTER TABLE ONLY client_scope
    ADD CONSTRAINT client_scope_pkey PRIMARY KEY (id),
    ADD CONSTRAINT client_scope_client_fk FOREIGN KEY (client_id) REFERENCES client(id);

ALTER TABLE ONLY newsletter
    ADD CONSTRAINT newsletter_pkey PRIMARY KEY (id),
    ADD CONSTRAINT newsletter_email_list_unique UNIQUE (email, "list"),
    ADD CONSTRAINT newsletter_code_unique UNIQUE (code);

ALTER TABLE ONLY article_visit
    ADD CONSTRAINT article_visit_pkey PRIMARY KEY (id),
    ADD CONSTRAINT article_visit_article_fk FOREIGN KEY (article_id) REFERENCES article(id),
    ADD CONSTRAINT article_visit_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY article_claps
    ADD CONSTRAINT article_claps_pkey PRIMARY KEY (id),
    ADD CONSTRAINT article_claps_article_fk FOREIGN KEY (article_id) REFERENCES article(id),
    ADD CONSTRAINT article_claps_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id);

ALTER TABLE ONLY article_recommendation
    ADD CONSTRAINT article_recommendation_pkey PRIMARY KEY (id),
    ADD CONSTRAINT article_recommendation_article_fk FOREIGN KEY (article_id) REFERENCES article(id),
    ADD CONSTRAINT article_recommendation_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id),
    ADD CONSTRAINT article_recommendation_recommended_to_fk FOREIGN KEY (recommended_to) REFERENCES "user"(id);

CREATE INDEX language_organization_fk_index ON language(organization_id);
CREATE INDEX organization_owner_fk_index ON organization(owner_user_id);
CREATE INDEX organization_member_organization_fk_index ON organization_member(organization_id);
CREATE INDEX organization_member_user_fk_index ON organization_member(user_id);
CREATE INDEX organization_member_language_fk_index ON organization_member(language_id);
CREATE INDEX article_organization_fk_index ON article(organization_id);
CREATE INDEX article_content_article_fk_index ON article_content(article_id);
CREATE INDEX article_content_language_fk_index ON article_content(language_id);
CREATE INDEX article_content_user_fk_index ON article_content(user_id);
CREATE INDEX article_access_article_fk_index ON article_access(article_id);
CREATE INDEX article_access_user_fk_index ON article_access(user_id);
CREATE INDEX article_access_user_group_fk_index ON article_access(user_group_id);
CREATE INDEX article_content_author_article_content_fk_index ON article_content_author(article_content_id);
CREATE INDEX article_content_author_user_fk_index ON article_content_author(user_id);
CREATE INDEX tag_organization_fk_index ON tag(organization_id);
CREATE INDEX user_group_organization_fk_index ON user_group(organization_id);
CREATE INDEX user_group_member_user_group_fk_index ON user_group_member(user_group_id);
CREATE INDEX user_group_member_user_fk_index ON user_group_member(user_id);
CREATE INDEX feed_organization_fk_index ON feed(organization_id);
CREATE INDEX feed_triggered_by_user_fk_index ON feed(triggered_by_user_id);
CREATE INDEX feed_ref_feed_fk_index ON feed_ref(feed_id);
CREATE INDEX feed_ref_user_fk_index ON feed_ref(user_id);
CREATE INDEX feed_ref_user_group_fk_index ON feed_ref(user_group_id);
CREATE INDEX feed_ref_article_fk_index ON feed_ref(article_id);
CREATE INDEX feed_ref_article_content_fk_index ON feed_ref(article_content_id);
CREATE INDEX feed_ref_article_list_fk_index ON feed_ref(article_list_id);
CREATE INDEX feed_access_user_fk_index ON feed_access(user_id);
CREATE INDEX feed_access_feed_fk_index ON feed_access(feed_id);
CREATE INDEX observed_object_user_fk_index ON observed_object(user_id);
CREATE INDEX observed_object_article_fk_index ON observed_object(article_id);
CREATE INDEX observed_object_article_list_fk_index ON observed_object(article_list_id);
CREATE INDEX observed_object_user_group_fk_index ON observed_object(user_group_id);
CREATE INDEX article_list_organization_fk_index ON article_list(organization_id);
CREATE INDEX article_list_entry_article_list_fk_index ON article_list_entry(article_list_id);
CREATE INDEX article_list_entry_article_fk_index ON article_list_entry(article_id);
CREATE INDEX article_list_member_article_list_fk_index ON article_list_member(article_list_id);
CREATE INDEX article_list_member_user_fk_index ON article_list_member(user_id);
CREATE INDEX article_list_name_article_list_fk_index ON article_list_name(article_list_id);
CREATE INDEX article_list_name_language_id_fk_index ON article_list_name(language_id);
CREATE INDEX invitation_organization_fk_index ON invitation(organization_id);
CREATE INDEX invitation_user_fk_index ON invitation(user_id);
CREATE INDEX article_tag_article_fk_index ON article_tag(article_id);
CREATE INDEX article_tag_tag_fk_index ON article_tag(tag_id);
CREATE INDEX file_organization_fk_index ON file(organization_id);
CREATE INDEX file_user_fk_index ON file(user_id);
CREATE INDEX file_article_fk_index ON file(article_id);
CREATE INDEX bookmark_organization_fk_index ON bookmark(organization_id);
CREATE INDEX bookmark_user_fk_index ON bookmark(user_id);
CREATE INDEX bookmark_article_fk_index ON bookmark(article_id);
CREATE INDEX bookmark_article_list_fk_index ON bookmark(article_list_id);
CREATE INDEX support_ticket_organization_fk_index ON support_ticket(organization_id);
CREATE INDEX support_ticket_user_fk_index ON support_ticket(user_id);
CREATE INDEX client_organization_fk_index ON client(organization_id);
CREATE INDEX client_scope_client_fk_index ON client_scope(client_id);
CREATE INDEX article_visit_article_fk_index ON article_visit(article_id);
CREATE INDEX article_visit_user_fk_index ON article_visit(user_id);
CREATE INDEX article_claps_article_fk_index ON article_claps(article_id);
CREATE INDEX article_claps_user_fk_index ON article_claps(user_id);
CREATE INDEX article_recommendation_article_fk_index ON article_recommendation(article_id);
CREATE INDEX article_recommendation_user_fk_index ON article_recommendation(user_id);
CREATE INDEX article_recommendation_recommended_to_fk_index ON article_recommendation(recommended_to);

CREATE INDEX article_content_content_trgm_gin ON article_content USING GIN(content_tsvector);
CREATE INDEX article_content_title_trgm_gin ON article_content USING GIN(title_tsvector);

CREATE OR REPLACE FUNCTION update_mod_time_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.mod_time = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_language_mod_time BEFORE UPDATE
    ON "language" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_user_mod_time BEFORE UPDATE
    ON "user" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_article_mod_time BEFORE UPDATE
    ON "article" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_article_content_mod_time BEFORE UPDATE
    ON "article_content" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_article_access_mod_time BEFORE UPDATE
    ON "article_access" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_article_content_author_mod_time BEFORE UPDATE
    ON "article_content_author" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_organization_mod_time BEFORE UPDATE
    ON "organization" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_organization_member_mod_time BEFORE UPDATE
    ON "organization_member" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_tag_mod_time BEFORE UPDATE
    ON "tag" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_tag_mod_time BEFORE UPDATE
    ON "user_group" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_tag_mod_time BEFORE UPDATE
    ON "user_group_member" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_feed_mod_time BEFORE UPDATE
    ON "feed" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_feed_ref_mod_time BEFORE UPDATE
    ON "feed_ref" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_feed_access_mod_time BEFORE UPDATE
    ON "feed_access" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_observed_object_mod_time BEFORE UPDATE
    ON "observed_object" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_article_list_mod_time BEFORE UPDATE
    ON "article_list" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_article_list_entry_mod_time BEFORE UPDATE
    ON "article_list_entry" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_article_list_member_mod_time BEFORE UPDATE
    ON "article_list_member" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_article_list_name_mod_time BEFORE UPDATE
    ON "article_list_name" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_invitation_mod_time BEFORE UPDATE
    ON "invitation" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_domain_blacklist_mod_time BEFORE UPDATE
    ON "domain_blacklist" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_article_tag_mod_time BEFORE UPDATE
    ON "article_tag" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_file_mod_time BEFORE UPDATE
    ON "file" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_bookmark_mod_time BEFORE UPDATE
    ON "bookmark" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_support_ticket_mod_time BEFORE UPDATE
    ON "support_ticket" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_client_mod_time BEFORE UPDATE
    ON "client" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_client_scope_mod_time BEFORE UPDATE
    ON "client_scope" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_newsletter_mod_time BEFORE UPDATE
    ON "newsletter" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_article_visit_mod_time BEFORE UPDATE
    ON "article_visit" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_article_claps_mod_time BEFORE UPDATE
    ON "article_claps" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();

CREATE TRIGGER update_article_recommendation_mod_time BEFORE UPDATE
    ON "article_recommendation" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();
