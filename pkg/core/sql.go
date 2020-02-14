package core

const managerDDL = `
CREATE TABLE IF NOT EXISTS manager
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL
);`

const sumTransferUsersDDL = `
CREATE TABLE IF NOT EXISTS sumTransferUsers
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    balance INTEGER NOT NULL
);`

const usersDDL = `
CREATE TABLE IF NOT EXISTS users
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
	passportSeries TEXT NOT NULL UNIQUE,
	phoneNumber INTEGER NOT NULL,
	hideShow INTEGER NOT NULL
);`

const operationsLoggingDDL = `
CREATE TABLE IF NOT EXISTS operationsLogging
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   time TEXT NOT NULL,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL CHECK ( balance > 0 ),
   user_id INTEGER REFERENCES users(id)
);`

const atmDDL = `
CREATE TABLE IF NOT EXISTS atm
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL UNIQUE,
    address TEXT NOT NULL
);`

const cardsDDL = `
CREATE TABLE IF NOT EXISTS cards
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   numberCard TEXT NOT NULL,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL CHECK ( balance > 0 ),
   user_id INTEGER REFERENCES users(id)
);`

const servicesDDL = `
CREATE TABLE IF NOT EXISTS services
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL 
);`

const managerInitialData = `INSERT INTO manager(name, login, password)
VALUES ('IBank', 'admin', 'boss')
       ON CONFLICT DO NOTHING;`
const sumTransferUsersDDLInitialData = `INSERT INTO sumTransferUsers(id,balance)
VALUES (1,0)
       ON CONFLICT DO NOTHING;`

const loginManagerSQL = `SELECT login, password FROM manager WHERE login = ?`
const loginUsersSQL = `SELECT id, login, password, hideShow FROM users WHERE login = ?`

const selectBalanceSumTransferUsers  = `SELECT balance FROM sumTransferUsers`
const selectIdUserPhoneNumberSQL = `SELECT id FROM users WHERE phoneNumber = ?`
const selectIdCardForTransferPhoneNumberSQL = `SELECT id FROM cards WHERE user_id= ?`
const selectIdCardForTransferCountNumberSQL = `SELECT id FROM cards WHERE numberCard= ?`

const selectIdUserLoginNumberSQL = `SELECT id FROM users WHERE login = ?`

const getAllAtmsSQL = `SELECT id, name, address FROM atm;`
const getAllServicesSQL = `SELECT id, name, balance FROM services;`
const getAllCardsSQL = `SELECT id, name, balance, user_id, numberCard FROM cards;`
const getAllUsersSQL = `SELECT id, name, passportSeries, phoneNumber FROM users;`
const getUserCardsSQL = `SELECT id, name, balance, numberCard FROM cards WHERE user_id = ?`
const getHideUserSQL = `SELECT id, name, passportSeries, phoneNumber FROM users WHERE hideShow = ?`


const insertAtmSQL = `INSERT INTO atm(name, address) VALUES ( :name, :address);`
const insertServiceSQL = `INSERT INTO services(name , balance) VALUES( :name, :balance);`
const insertCardSQL = `INSERT INTO cards(name, balance, user_id, numberCard) VALUES ( :name, :balance, :user_id, :numberCard);`
const insertUserSQL = `INSERT INTO users(name, login, password, passportSeries, phoneNumber, hideShow) VALUES (:name , :login, :password, :passportSeries, :phoneNumber, :hideShow);`


const updateBalanceToCardSenderSQL = `UPDATE cards SET balance=? WHERE user_id = ?`
const updateBalanceToCardRecipientSQL = `UPDATE cards SET balance=? WHERE id = ?`
const updateBalanceSumTransferUsersSQL  = `UPDATE sumTransferUsers SET balance = ?`


const selectBalanceToCardSenderSQL = `SELECT balance FROM cards WHERE user_id = ?`
const selectBalanceToCardRecipientSQL = `SELECT balance FROM cards WHERE id = ?`

const selectDescIdFromCardSQL = `SELECT id FROM cards ORDER BY id DESC LIMIT 1;`

const selectBalanceOnServiceSQL = `SELECT balance FROM services WHERE name = ?`
const updateBalanceServiceSQL = `UPDATE services SET balance=? WHERE name = ?`

const updateHideShowUser = `UPDATE users SET hideShow = ? WHERE id = ?`

const searchUserForPhoneNumberSQL = `SELECT id, name, passportSeries, phoneNumber FROM users WHERE phoneNumber = ?`

const staticCountUserSQL  = `SELECT count(id) FROM users`
const staticSumBalanceUsersSQL  = `SELECT sum(balance) FROM cards`
const staticBalanceOfServicesSQL  = `SELECT sum(balance) FROM services`
const staticBalanceOfServiceSQL  = `SELECT name, balance FROM services`

