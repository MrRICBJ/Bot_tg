CREATE TABLE movies (
    id BIGINT PRIMARY KEY,
    user_id integer NOT NULL,
    movie_id integer NOT NULL,
    genre character varying(255) NOT NULL,
    movie_name character varying(255) NOT NULL,
    "time" timestamp without time zone NOT NULL
);
