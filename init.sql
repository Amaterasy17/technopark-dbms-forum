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
    Message  text   NOT NULL,
    Parent   BIGINT                   DEFAULT 0,
    Thread   INT,
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