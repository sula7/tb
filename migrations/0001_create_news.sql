-- +goose Up
CREATE TABLE news (
    id bigserial NOT NULL,
    tag_content varchar,
    url varchar(255),
    created_at timestamp(0) DEFAULT now()
);
CREATE UNIQUE INDEX news_id_uindex ON news (id);
ALTER TABLE news ADD CONSTRAINT news_id_pk PRIMARY KEY (id);

-- +goose Down
DROP TABLE IF EXISTS news;
