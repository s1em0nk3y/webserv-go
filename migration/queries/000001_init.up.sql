CREATE TABLE IF NOT EXISTS Users (
   id serial PRIMARY KEY,
   nickname character varying(255),
   passhash character varying(300)
);