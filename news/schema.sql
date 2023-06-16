DROP TABLE IF EXISTS gonews;

CREATE TABLE gonews (
                       id SERIAL PRIMARY KEY,
                       title TEXT NOT NULL DEFAULT 'empty',
                       content TEXT NOT NULL DEFAULT 'empty',
                       pubtime BIGINT NOT NULL DEFAULT extract (epoch from now()),
                       link TEXT NOT NULL
);