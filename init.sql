CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE "users"
(
    About    text,
    Email    citext UNIQUE,
    FullName text NOT NULL,
    Nickname citext PRIMARY KEY
);

CREATE UNLOGGED TABLE forum
(
    "user"  citext,
    Posts   BIGINT DEFAULT 0,
    Slug    citext PRIMARY KEY,
    Threads INT    DEFAULT 0,
    title   text,
    FOREIGN KEY ("user") REFERENCES "users" (nickname)
);

CREATE UNLOGGED TABLE thread
(
    author  citext,
    created timestamp with time zone default now(),
    forum   citext,
    id      SERIAL PRIMARY KEY,
    message text NOT NULL,
    slug    citext UNIQUE,
    title   text not null,
    votes   INT                      default 0,
    FOREIGN KEY (author) REFERENCES "users" (nickname),
    FOREIGN KEY (forum) REFERENCES "forum" (slug)
);

CREATE UNLOGGED TABLE post
(
    author   citext NOT NULL,
    created  timestamp with time zone default now(),
    forum    citext,
    id       BIGSERIAL PRIMARY KEY,
    isEdited BOOLEAN                  DEFAULT FALSE,
    message  text   NOT NULL,
    parent   BIGINT                   DEFAULT 0,
    thread   INT,
    FOREIGN KEY (author) REFERENCES "users" (nickname),
    FOREIGN KEY (forum) REFERENCES "forum" (slug),
    FOREIGN KEY (thread) REFERENCES "thread" (id),
    FOREIGN KEY (parent) REFERENCES "post" (id)
);