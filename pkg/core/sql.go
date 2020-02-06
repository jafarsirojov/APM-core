package core

const managerDDL = `
CREATE TABLE IF NOT EXISTS managers
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
);`

const usersDDL = `
CREATE TABLE IF NOT EXISTS managers
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
);`


