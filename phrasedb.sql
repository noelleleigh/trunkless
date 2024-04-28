DROP TABLE phrases;
DROP TABLE sources;
DROP TABLE CORPORA;

CREATE TABLE IF NOT EXISTS corpora (
  id   char(7) PRIMARY KEY,
  name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS sources (
  id       char(7) PRIMARY KEY,
  corpusid char(7) NOT NULL,
  name     TEXT UNIQUE,

  FOREIGN KEY (corpusid) REFERENCES corpora(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS phrases (
  id       SERIAL PRIMARY KEY,
  sourceid char(7) NOT NULL,
  phrase   TEXT,

  FOREIGN KEY (sourceid) REFERENCES sources(id) ON DELETE CASCADE
);

