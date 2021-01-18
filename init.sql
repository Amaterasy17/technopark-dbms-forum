CREATE EXTENSION IF NOT EXISTS citext;

-- ALTER SYSTEM SET
--     checkpoint_completion_target = '0.9';
-- ALTER SYSTEM SET
--     wal_buffers = '6912kB';
-- ALTER SYSTEM SET
--     default_statistics_target = '100';
-- ALTER SYSTEM SET
--     random_page_cost = '1.1';
-- ALTER SYSTEM SET
--     effective_io_concurrency = '200';


CREATE UNLOGGED TABLE users
(
    Nickname citext Primary Key,
    FullName text NOT NULL,
    About    text,
    Email    citext UNIQUE
);

CREATE UNLOGGED TABLE forum
(
    Slug    citext PRIMARY KEY,
    "user"  citext,
    Title   text NOT NULL,
    Posts   BIGINT DEFAULT 0,
    Threads INT    DEFAULT 0
);

CREATE UNLOGGED TABLE thread
(
    id      SERIAL PRIMARY KEY,
    Title   text not null,
    Author  citext REFERENCES "users" (Nickname),
    Created timestamp with time zone default now(),
    Forum   citext REFERENCES "forum" (slug),
    Message text NOT NULL,
    slug    citext UNIQUE,
    Votes   INT default 0
);

CREATE UNLOGGED TABLE post
(
    id       BIGSERIAL PRIMARY KEY,
    Author citext,
    Created  timestamp with time zone default now(),
    Forum    citext,
    isEdited BOOLEAN                  DEFAULT FALSE,
    Message  text NOT NULL,
    Parent   BIGINT                   DEFAULT 0,
    Thread   INT,
    Path     BIGINT[]                 DEFAULT ARRAY []::INTEGER[],
--     FOREIGN KEY (forum) REFERENCES "forum" (slug),
    FOREIGN KEY (thread) REFERENCES "thread" (id),
    FOREIGN KEY (author) REFERENCES "users"  (nickname)
);

CREATE UNLOGGED TABLE votes
(
    id     BIGSERIAL PRIMARY KEY,
    Author citext REFERENCES "users" (nickname),
    Voice INT NOT NULL,
    Thread INT,
    FOREIGN KEY (thread) REFERENCES "thread" (id),
    UNIQUE (Author, Thread)
);

CREATE UNLOGGED TABLE users_forum
(
    nickname citext NOT NULL,
    fullname TEXT NOT NULL,
    about    TEXT,
    email    CITEXT,
    slug     citext NOT NULL,
    FOREIGN KEY (nickname) REFERENCES "users" (nickname),
    FOREIGN KEY (slug) REFERENCES "forum" (slug),
    UNIQUE (nickname, slug)
);


CREATE OR REPLACE FUNCTION insertVotes() RETURNS TRIGGER AS
$update_vote$
BEGIN
    UPDATE thread SET votes=(votes+NEW.voice) WHERE id=NEW.thread;
    return NEW;
end
$update_vote$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION updatePostUserForum() RETURNS TRIGGER AS
$update_forum_post$
DECLARE
    m_fullname CITEXT;
    m_about    CITEXT;
    m_email CITEXT;
BEGIN
    SELECT fullname, about, email FROM users WHERE nickname = NEW.author INTO m_fullname, m_about, m_email;
    INSERT INTO users_forum (nickname, fullname, about, email, Slug)
     VALUES (New.Author,m_fullname, m_about, m_email, NEW.forum) on conflict do nothing;
    return NEW;
end
$update_forum_post$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION updateThreadUserForum() RETURNS TRIGGER AS
$update_forum_thread$
DECLARE
    author_nick citext;
    m_fullname CITEXT;
    m_about    CITEXT;
    m_email CITEXT;
BEGIN
    SELECT Nickname, fullname, about, email FROM users WHERE Nickname = new.Author INTO author_nick, m_fullname, m_about, m_email;
    INSERT INTO users_forum (nickname, fullname, about, email, Slug)
     VALUES (author_nick,m_fullname, m_about, m_email, NEW.forum) on conflict do nothing;
    return NEW;
end
$update_forum_thread$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION updateVotes() RETURNS TRIGGER AS
$update_vote$
BEGIN
    IF OLD.Voice <> NEW.Voice THEN
        UPDATE thread SET votes=(votes+NEW.Voice*2) WHERE id=NEW.Thread;
    END IF;
    return NEW;
end
$update_vote$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION updateCountOfThreads() RETURNS TRIGGER AS
$update_forum$
BEGIN
    UPDATE forum SET Threads=(Threads+1) WHERE slug=NEW.forum;
    return NEW;
end
$update_forum$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION updatePath() RETURNS TRIGGER AS
$update_path$
DECLARE
    parentPath         BIGINT[];
    first_parent_thread INT;
BEGIN
    IF (NEW.parent IS NULL) THEN
        NEW.path := array_append(new.path, new.id);
    ELSE
        SELECT path FROM post WHERE id = new.parent INTO parentPath;
        SELECT thread FROM post WHERE id = parentPath[1] INTO first_parent_thread;
        IF NOT FOUND OR first_parent_thread <> NEW.thread THEN
            RAISE EXCEPTION 'parent is from different thread' USING ERRCODE = '00409';
        end if;

        NEW.path := NEW.path || parentPath || new.id;
    end if;
    UPDATE forum SET Posts=Posts + 1 WHERE forum.slug = new.forum;
    RETURN new;
end
$update_path$ LANGUAGE plpgsql;



CREATE TRIGGER addThreadInForum
    BEFORE INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE updateCountOfThreads();


CREATE TRIGGER add_voice
    BEFORE INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE insertVotes();

CREATE TRIGGER edit_voice
    BEFORE UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE updateVotes();

CREATE TRIGGER update_path_trigger
    BEFORE INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE updatePath();

CREATE TRIGGER post_insert_user_forum
    AFTER INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE updatePostUserForum();

CREATE TRIGGER thread_insert_user_forum
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE updateThreadUserForum();



CREATE INDEX IF NOT EXISTS user_nickname ON users using hash (nickname);
CREATE INDEX IF NOT EXISTS user_email ON users using hash (email);
CREATE INDEX IF NOT EXISTS forum_slug ON forum using hash (slug);
CREATE UNIQUE INDEX IF NOT EXISTS  forum_users_unique on users_forum (slug, nickname);
-- CLUSTER users_forum USING forum_users_unique;
CREATE INDEX IF NOT EXISTS  thr_slug ON thread using hash (slug);
CREATE INDEX IF NOT EXISTS  thr_date ON thread (created);
CREATE INDEX IF NOT EXISTS  thr_forum ON thread using hash (forum);
CREATE INDEX IF NOT EXISTS  thr_forum_date ON thread (forum, created);
CREATE INDEX IF NOT EXISTS post_id_path on post (id, (path[1]));
CREATE INDEX IF NOT EXISTS post_thread_id_path1_parent on post (thread, id, (path[1]), parent);
CREATE INDEX IF NOT EXISTS post_thread_path_id on post (thread, path, id);
CREATE INDEX IF NOT EXISTS post_path1 on post ((path[1]));
CREATE INDEX IF NOT EXISTS post_thread_id on post (thread, id);
CREATE INDEX IF NOT EXISTS post_thr_id ON post (thread);
CREATE UNIQUE INDEX IF NOT EXISTS  vote_unique on votes (Author, Thread);
CREATE INDEX IF NOT EXISTS post_path1_path_id_desc ON post ((path[1]) ASC, path, id);

-- CREATE INDEX if not exists user_nickname ON users using hash (nickname);
-- CREATE INDEX if not exists user_email ON users using hash (email);
-- CREATE INDEX if not exists forum_slug ON forum using hash (slug);
--
-- create unique index if not exists forum_users_unique on users_forum (slug, nickname);
-- -- cluster users_forum using forum_users_unique;
--
-- CREATE INDEX if not exists thr_slug ON thread using hash (slug);
-- CREATE INDEX if not exists thr_date ON thread (created);
-- CREATE INDEX if not exists thr_forum ON thread using hash (forum);
-- CREATE INDEX if not exists thr_forum_date ON thread (forum, created);
--
-- create index if not exists post_id_path on post (id, (path[1]));
-- create index if not exists post_thread_id_path1_parent on post (thread, id, (path[1]), parent);
-- create index if not exists post_thread_path_id on post (thread, path, id);
-- -- create index if not exists post_thread_parent_id on post (id, thread, parent) WHERE parent is Null;
--
-- -- create index if not exists post_thread_path_id_desc on post (thread, id Desc, path DESC);
-- -- create index if not exists post_thread_path_id_asc on post (thread, id ASC, path ASC);
--
-- create index if not exists post_path1 on post ((path[1]));
-- create index if not exists post_thread_id on post (thread, id);
-- CREATE INDEX if not exists post_thr_id ON post (thread);
--
-- create unique index if not exists vote_unique on votes (Author, Thread);

-- CREATE INDEX IF NOT EXISTS post_path1_path_id_desc ON post ((path[1]) DESC, path, id);
-- CREATE INDEX IF NOT EXISTS post_path1_path_id_asc ON post ((path[1]) ASC, path, id);




ANALYZE post;
ANALYZE users_forum;
ANALYZE thread;



