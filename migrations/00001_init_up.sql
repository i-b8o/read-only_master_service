-- TABLES --
DROP TABLE IF EXISTS absent_reg;
DROP TABLE IF EXISTS pseudo_chapters;
DROP TABLE IF EXISTS pseudo_regulations;
DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS speech;


CREATE TABLE pseudo_regulations (
    r_id integer,
    pseudo TEXT NOT NULL CHECK (pseudo != '')
);

CREATE TABLE pseudo_chapters (
    c_id integer,
    pseudo TEXT NOT NULL CHECK (pseudo != '')
);

CREATE TABLE absent_reg (
    id SERIAL PRIMARY KEY,
    pseudo TEXT NOT NULL CHECK (pseudo != ''),
    done BOOLEAN NOT NULL DEFAULT false,
    paragraph_id integer  
);

CREATE TABLE links (
    id INT NOT NULL UNIQUE,
    paragraph_num INT NOT NULL CHECK (paragraph_num >= 0),
    c_id integer REFERENCES chapters,
    r_id integer REFERENCES regulations
);

CREATE TABLE speech (
    id SERIAL PRIMARY KEY,
    order_num INT NOT NULL CHECK (order_num >= 0),
    content TEXT,
    paragraph_id INT NOT NULL CHECK (paragraph_id >= 0)
);