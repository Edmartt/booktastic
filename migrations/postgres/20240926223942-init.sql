
-- +migrate Up

CREATE TABLE IF NOT EXISTS books(
	id SERIAL PRIMARY KEY,
	uuid character varying(36),
	isbn character varying(30),
	title TEXT,
	pages character varying(30),
	current_page character varying(30),
	author TEXT,
	year character varying(30),
	status TEXT CHECK (status IN ('read', 'to_be_read', 'reading', 'pending_purchase'))
	
);

-- +migrate Down

DROP TABLE books;
