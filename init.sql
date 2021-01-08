CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE users
(
    Nickname citext PRIMARY KEY,
    FullName text NOT NULL,
    About    text,
    Email    citext UNIQUE
);

CREATE UNLOGGED TABLE forum
(
    Slug    citext PRIMARY KEY,
    "user"  citext REFERENCES "users" (Nickname),
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
    Author   citext REFERENCES "users" (nickname),
    Created  timestamp with time zone default now(),
    Forum    citext,
    isEdited BOOLEAN                  DEFAULT FALSE,
    Message  text NOT NULL,
    Parent   BIGINT                   DEFAULT 0,
    Thread   INT,
    Path     BIGINT[]                 DEFAULT ARRAY []::INTEGER[],
    FOREIGN KEY (forum) REFERENCES "forum" (slug),
    FOREIGN KEY (thread) REFERENCES "thread" (id)
--     FOREIGN KEY (parent) REFERENCES "post" (id)
);

CREATE UNLOGGED TABLE votes
(
    id     BIGSERIAL PRIMARY KEY,
    Author citext REFERENCES "users" (nickname),
    Voice INT NOT NULL,
    Thread INT,
    FOREIGN KEY (thread) REFERENCES "thread" (id)
);


CREATE OR REPLACE FUNCTION insertVotes() RETURNS TRIGGER AS
$update_vote$
BEGIN
    UPDATE thread SET votes=(votes+NEW.voice) WHERE id=NEW.thread;
    return NEW;
end
$update_vote$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION updateVotes() RETURNS TRIGGER AS
$update_vote$
BEGIN
    UPDATE thread SET votes=(votes+NEW.Voice*2) WHERE id=NEW.Thread;
    return NEW;
end
$update_vote$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION updateCountOfThreads() RETURNS TRIGGER AS
$update_forum$
BEGIN
    UPDATE forum SET Threads=(Threads+1) WHERE LOWER(slug)=LOWER(NEW.forum);
    return NEW;
end
$update_forum$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION updatePath() RETURNS TRIGGER AS
$update_path$
DECLARE
    parentPath         BIGINT[];
    first_parent_thread INT;
BEGIN
    IF (NEW.parent IS 0) THEN
        NEW.path := array_append(new.path, new.id);
    ELSE
        SELECT path FROM post WHERE id = new.parent INTO parentPath;
        SELECT thread FROM post WHERE id = parentPath[1] INTO first_parent_thread;
        IF NOT FOUND OR first_parent_thread != NEW.thread THEN
            RAISE EXCEPTION 'parent is from different thread' USING ERRCODE = '00409';
        end if;

        NEW.path := NEW.path || parentPath || new.id;
    end if;
    UPDATE forum SET Posts=Posts + 1 WHERE lower(forum.slug) = lower(new.forum);
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