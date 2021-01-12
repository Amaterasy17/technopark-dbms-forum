CREATE EXTENSION IF NOT EXISTS citext;

ALTER SYSTEM SET
    checkpoint_completion_target = '0.9';
ALTER SYSTEM SET
    wal_buffers = '6912kB';
ALTER SYSTEM SET
    default_statistics_target = '100';
ALTER SYSTEM SET
    random_page_cost = '1.1';
ALTER SYSTEM SET
    effective_io_concurrency = '200';


CREATE UNLOGGED TABLE users
(
    Id       Serial PRIMARY KEY,
    Nickname citext UNIQUE,
    FullName text NOT NULL,
    About    text,
    Email    citext UNIQUE
);

CREATE UNLOGGED TABLE forum
(
    Slug    citext PRIMARY KEY,
    "user"  int,
    Title   text NOT NULL,
    Posts   BIGINT DEFAULT 0,
    Threads INT    DEFAULT 0
);

CREATE UNLOGGED TABLE thread
(
    id      SERIAL PRIMARY KEY,
    Title   text not null,
    Author  int,
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
    FOREIGN KEY (thread) REFERENCES "thread" (id)

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
BEGIN
    INSERT INTO users_forum (nickname, Slug) VALUES (New.Author, NEW.forum) on conflict do nothing;
    return NEW;
end
$update_forum_post$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION updateThreadUserForum() RETURNS TRIGGER AS
$update_forum_thread$
DECLARE
    author_nick citext;
BEGIN
    SELECT Nickname FROM users WHERE id = new.Author INTO author_nick;
    INSERT INTO users_forum (nickname, Slug) VALUES (author_nick, NEW.forum) on conflict do nothing;
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





Create index thread_slug_hash_index ON thread using hash (slug);
CREATE INDEX thread_id_forum_index ON thread (id, forum);
CREATE INDEX thread_created_index ON thread (created);
CREATE INDEX thread_forum_lower_index ON thread (forum);


CREATE INDEX forum_index ON forum (Slug);


CREATE INDEX votes_index ON votes (Author, Thread);

CREATE INDEX users_forum_forum_index ON users_forum ((users_forum.Slug));
CREATE INDEX thread_created_index ON thread (Created);

CREATE INDEX post_first_parent_thread_index ON post (thread, (post.path[1]) );
CREATE INDEX post_first_parent_id_index ON post (id, (post.path[1]));
CREATE INDEX post_first_parent_index ON post ((post.path[1]));
CREATE INDEX post_path_index ON post ((post.path));
CREATE INDEX post_thread_id_index ON post (thread, id);
CREATE INDEX post_path_id_index ON post (id, (post.path));
CREATE INDEX post_thread_path_id_index ON post (thread, (post.parent), id);



CREATE INDEX users_nickname_index ON users ((users.Nickname));
CREATE INDEX users_email_index ON users (Email);


CREATE INDEX users_forum_forum_user_index ON users_forum ((users_forum.slug), nickname);
CREATE INDEX users_forum_user_index ON users_forum (nickname);
CREATE INDEX users_forum_forum_index ON users_forum ((users_forum.Slug));

ANALYZE post;
ANALYZE users_forum;
ANALYZE thread;



