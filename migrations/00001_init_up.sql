BEGIN;

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = ON;
SET check_function_bodies = FALSE;
SET client_min_messages = WARNING;
SET search_path = public, extensions;
SET default_tablespace = '';
SET default_with_oids = FALSE;

SET SCHEMA 'public';

-- CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- TABLES --
DROP TABLE IF EXISTS absent_reg;
DROP TABLE IF EXISTS pseudo_chapters;
DROP TABLE IF EXISTS pseudo_regulations;


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