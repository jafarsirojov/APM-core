package core

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestLoginManager_QueryError(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = LoginManager("", "", db)
	// errors.Is vs errors.As
	var typedErr *QueryError
	if ok := errors.As(err, &typedErr); !ok {
		t.Errorf("error not maptch QueryError: %v", err)
	}
}

func TestLoginManager_NoSuchLoginForEmptyDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	// Crash Early
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	// shift 2 раза -> sql dialect
	_, err = db.Exec(`
	CREATE TABLE manager (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute query: %v", err)
	}

	result, err := LoginManager("", "", db)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	if result != false {
		t.Error("Login result not false for empty table")
	}
}

func TestLoginManager_LoginOk(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	// shift 2 раза -> sql dialect
	_, err = db.Exec(`
	CREATE TABLE manager (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO manager(id, login, password) VALUES (1, 'vasya', 'secret')`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	result, err := LoginManager("vasya", "secret", db)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	if result != true {
		t.Error("Login result not true for existing account")
	}
}

func TestLoginManager_LoginNotOkForInvalidPassword(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	// shift 2 раза -> sql dialect
	_, err = db.Exec(`
	CREATE TABLE manager (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO manager(id, login, password) VALUES (1, 'vasya', 'secret')`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = LoginManager("vasya", "password", db)
	if !errors.Is(err, ErrInvalidPass) {
		t.Errorf("Not ErrInvalidPass error for invalid pass: %v", err)
	}
}

func TestLoginUsers_QueryError(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = LoginUsers("", "", db)
	// errors.Is vs errors.As
	var typedErr *QueryError
	if ok := errors.As(err, &typedErr); !ok {
		t.Errorf("error not maptch QueryError: %v", err)
	}
}

func TestLoginUsers_NoSuchLoginForEmptyDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	// Crash Early
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	// shift 2 раза -> sql dialect
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
	passportSeries TEXT NOT NULL UNIQUE,
	phoneNumber INTEGER NOT NULL,
	hideShow INTEGER NOT NULL
);`)
	if err != nil {
		t.Errorf("can't execute query: %v", err)
	}

	result, err := LoginUsers("", "", db)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	if result != false {
		t.Error("Login result not false for empty table")
	}
}

func TestLoginUsers_LoginOk(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	// shift 2 раза -> sql dialect
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
	passportSeries TEXT NOT NULL UNIQUE,
	phoneNumber INTEGER NOT NULL,
	hideShow INTEGER NOT NULL
);`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO users( name, login, password, passportSeries, phoneNumber, hideShow) VALUES ('Vasya','vasya', 'secret','A132323',9001,3)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	result, err := LoginUsers("vasya", "secret", db)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	if result != true {
		t.Error("Login result not true for existing account")
	}
}

func TestLoginUsers_LoginNotOkForInvalidPassword(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
	passportSeries TEXT NOT NULL UNIQUE,
	phoneNumber INTEGER NOT NULL,
	hideShow INTEGER NOT NULL
);`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO users( name, login, password, passportSeries, phoneNumber, hideShow) VALUES ('Vasya','vasya', 'secret','A132323',9001,3)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = LoginUsers("vasya", "password", db)
	if !errors.Is(err, ErrInvalidPass) {
		t.Errorf("Not ErrInvalidPass error for invalid pass: %v", err)
	}
}

func TestAddAtm_NoBd(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	err = AddAtm("T1", "rudaki 65", db)
	if err == nil {
		t.Errorf("can't execute add atm: %v", err)
	}

}

func TestAddAtm_HasBd(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS atm
	(
		id      INTEGER PRIMARY KEY AUTOINCREMENT,
		name    TEXT    NOT NULL UNIQUE,
		address TEXT NOT NULL
	);`)

	err = AddAtm("T1", "rudaki 65", db)
	if err != nil {
		t.Errorf("can't execute add atm: %v", err)
	}
}

func TestGetAllAtms_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = GetAllAtms(db)
	if err == nil {
		t.Errorf("can't get all atm: %v", err)
	}
}

func TestGetAllAtms_HasDbError(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS atm
	(
		id      INTEGER PRIMARY KEY AUTOINCREMENT,
		name    TEXT    NOT NULL UNIQUE,
		address TEXT NOT NULL
	);`)

	atms, err := GetAllAtms(db)
	if err != nil {
		t.Errorf("can't get all atm: %v", err)
	}
	if atms != nil {
		t.Errorf("can't get all atm: %v", err)
	}
}

func TestGetAllAtms_HasDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS atm
	(
		id      INTEGER PRIMARY KEY AUTOINCREMENT,
		name    TEXT    NOT NULL UNIQUE,
		address TEXT NOT NULL
	);`)
	if err != nil {
		t.Errorf("can't creat atm to get all atm: %v", err)
	}

	_, err = db.Exec(`INSERT INTO atm(name, address) VALUES ("t1","rudaki 43")`)
	if err != nil {
		t.Errorf("can't get all atm, add atm: %v", err)
	}

	_, err = db.Exec(`INSERT INTO atm(name, address) VALUES ("t2","somoni 77")`)
	if err != nil {
		t.Errorf("can't get all atm, add atm: %v", err)
	}
	atms, err := GetAllAtms(db)
	if err != nil {
		t.Errorf("can't get all atm: %v", err)
	}

	if atms == nil {
		t.Errorf("can't get all atm: %v", err)
	}
}

func TestAddService_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = AddService("Internet", db)
	if err == nil {
		t.Errorf("can't add service Internet: %v", err)
	}
}

func TestAddService_HasDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS services
	(
		id      INTEGER PRIMARY KEY AUTOINCREMENT,
		name    TEXT    NOT NULL,
		balance INTEGER NOT NULL
	);`)

	err = AddService("Internet", db)
	if err != nil {
		t.Errorf("can't add service Internet: %v", err)
	}
}

//-------------
func TestGetAllServices_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = GetAllServices(db)
	if err == nil {
		t.Errorf("can't get all atm: %v", err)
	}
}

func TestGetAllServices_HasDbError(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS services
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL 
);`)

	atms, err := GetAllAtms(db)
	if err == nil {
		t.Errorf("can't get all atm: %v", err)
	}
	if atms != nil {
		t.Errorf("can't get all atm: %v", err)
	}
}

func TestGetAllServices_HasDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS services
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL 
);`)
	if err != nil {
		t.Errorf("can't creat atm to get all atm: %v", err)
	}

	_, err = db.Exec(`INSERT INTO services(name , balance) VALUES("Internet",0)`)
	if err != nil {
		t.Errorf("can't get all services, add atm: %v", err)
	}

	_, err = db.Exec(`INSERT INTO services(name , balance) VALUES("Water",0)`)
	if err != nil {
		t.Errorf("can't get all services, add atm: %v", err)
	}
	services, err := GetAllServices(db)
	if err != nil {
		t.Errorf("can't get all seervices: %v", err)
	}

	if services == nil {
		t.Errorf("can't get all services: %v", err)
	}
}

func TestAddCard_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = AddCard("AlifMobi", 100, 1, db)
	if err == nil {
		t.Errorf("can't add card: %v", err)
	}
}

func TestAddCard_HasDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cards
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   numberCard TEXT NOT NULL,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL CHECK ( balance > 0 ),
   user_id INTEGER REFERENCES users(id)
);`)
	if err != nil {
		t.Errorf("can't add card: %v", err)
	}

	err = AddCard("AlifMobi", 100, 1, db)
	if err != nil {
		t.Errorf("can't add card: %v", err)
	}
}

func TestGetAllCards_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	services, err := GetAllServices(db)

	if err == nil {
		t.Errorf("can't get all cards: %v", err)
	}
	if services != nil {
		t.Errorf("can't get all cards: %v", err)
	}
}

func TestGetAllCards_HasDbError(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cards
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   numberCard TEXT NOT NULL,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL CHECK ( balance > 0 ),
   user_id INTEGER REFERENCES users(id)
);`)

	cards, err := GetAllCards(db)
	if err != nil {
		t.Errorf("can't get all atm: %v", err)
	}
	if cards != nil {
		t.Errorf("can't get all atm: %v", err)
	}
}

func TestGetAllCards_HasDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cards
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   numberCard TEXT NOT NULL,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL CHECK ( balance > 0 ),
   user_id INTEGER REFERENCES users(id)
);`)
	if err != nil {
		t.Errorf("can't creat atm to get all atm: %v", err)
	}

	_, err = db.Exec(`INSERT INTO cards(id,name, balance, user_id, numberCard) VALUES (1,"AlifMobi",200,1,"20216000000000001")`)
	if err != nil {
		t.Errorf("can't get all card, add card: %v", err)
	}

	_, err = db.Exec(`INSERT INTO cards(id,name, balance, user_id, numberCard) VALUES (2,"AlifMobi",400,2,"20216000000000002")`)
	if err != nil {
		t.Errorf("can't get all card, add card: %v", err)
	}
	cards, err := GetAllCards(db)
	if err != nil {
		t.Errorf("can't get all cards: %v", err)
	}

	if cards == nil {
		t.Errorf("can't get all card: %v", err)
	}
}

func TestAddUser_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = AddUser("User1", "user1", "secret", "A242342", 9002, db)
	if err == nil {
		t.Errorf("can't add card: %v", err)
	}
}

func TestAddUser_HasDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
	passportSeries TEXT NOT NULL UNIQUE,
	phoneNumber INTEGER NOT NULL,
	hideShow INTEGER NOT NULL
);`)
	if err != nil {
		t.Errorf("can't create table users, add user: %v", err)
	}

	err = AddUser("User1", "user1", "secret", "A242342", 9002, db)
	if err != nil {
		t.Errorf("can't add user: %v", err)
	}
}

func TestGetAllUsers_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = GetAllUsers(db)
	if err == nil {
		t.Errorf("can't get all users: %v", err)
	}
}

func TestGetAllUsers_HasDbError(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
	passportSeries TEXT NOT NULL UNIQUE,
	phoneNumber INTEGER NOT NULL,
	hideShow INTEGER NOT NULL
);`)
	if err != nil {
		t.Errorf("can't get all users: %v", err)
	}

	users, err := GetAllUsers(db)
	if err != nil {
		t.Errorf("can't get all users: %v", err)
	}
	if users != nil {
		t.Errorf("can't get all users: %v", err)
	}
}

func TestGetAllUsers_HasDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
	passportSeries TEXT NOT NULL UNIQUE,
	phoneNumber INTEGER NOT NULL,
	hideShow INTEGER NOT NULL
);`)
	if err != nil {
		t.Errorf("can't creat table users to get all users: %v", err)
	}

	_, err = db.Exec(`INSERT INTO users( name, login, password, passportSeries, phoneNumber, hideShow) VALUES ('Vasya','vasya', 'secret','A132323',9001,3)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO users( name, login, password, passportSeries, phoneNumber, hideShow) VALUES ('Petya','petya', 'secret','A000009',9004,3)`)
	if err != nil {
		t.Errorf("can't add users to get all users: %v", err)
	}
	cards, err := GetAllUsers(db)
	if err != nil {
		t.Errorf("can't get all users: %v", err)
	}

	if cards == nil {
		t.Errorf("can't get all users: %v", err)
	}
}

func TestGetUserCards_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	cards, err := GetUserCards(db)
	if err == nil {
		t.Errorf("can't get card user: %v", err)
	}
	if cards != nil {
		t.Errorf("can't get card user: %v", err)
	}
}

func TestGetUserCards_HasDbError(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cards
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   numberCard TEXT NOT NULL,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL CHECK ( balance > 0 ),
   user_id INTEGER REFERENCES users(id)
);`)
	if err != nil {
		t.Errorf("can't creat table user, get user cards: %v", err)
	}

	users, err := GetUserCards(db)
	if err != nil {
		t.Errorf("can't get user cards: %v", err)
	}
	if users != nil {
		t.Errorf("can't get user cards: %v", err)
	}
}

func TestGetUserCards_HasDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
	passportSeries TEXT NOT NULL UNIQUE,
	phoneNumber INTEGER NOT NULL,
	hideShow INTEGER NOT NULL
);`)
	if err != nil {
		t.Errorf("can't create table users, add user: %v", err)
	}
	_, err = db.Exec(`INSERT INTO users( id,name, login, password, passportSeries, phoneNumber, hideShow) VALUES (1,'Vasya','vasya', 'secret','A132323',9001,3)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cards
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   numberCard TEXT NOT NULL,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL CHECK ( balance > 0 ),
   user_id INTEGER REFERENCES users(id)
);`)
	if err != nil {
		t.Errorf("can't creat table user, get user cards: %v", err)
	}

	_, err = db.Exec(`INSERT INTO cards(id,name, balance, user_id, numberCard) VALUES (1,"AlifMobi",200,0,"20216000000000001")`)
	if err != nil {
		t.Errorf("can't get user cards, add card: %v", err)
	}

	_, err = db.Exec(`INSERT INTO cards(id,name, balance, user_id, numberCard) VALUES (2,"AlifMobi",400,1,"20216000000000002")`)
	if err != nil {
		t.Errorf("can't get user cards, add card: %v", err)
	}
	cards, err := GetUserCards(db)
	if err != nil {
		t.Errorf("can't get user cards: %v", err)
	}

	if cards == nil {
		t.Errorf("can't get user cards: %v", cards)
	}
}

func TestTransferMoneyForPhoneNumber_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = TransferMoneyForPhoneNumber(9001, db)
	if err == nil {
		t.Errorf("can't search id cards for number phone: %v", err)
	}
}

func TestTransferMoneyForPhoneNumber_HasDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
	passportSeries TEXT NOT NULL UNIQUE,
	phoneNumber INTEGER NOT NULL,
	hideShow INTEGER NOT NULL
);`)
	if err != nil {
		t.Errorf("can't create table users, add user: %v", err)
	}

	err = AddUser("User1", "user1", "secret", "A242342", 9001, db)
	if err != nil {
		t.Errorf("can't add user: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cards
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   numberCard TEXT NOT NULL,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL CHECK ( balance > 0 ),
   user_id INTEGER REFERENCES users(id)
);`)

	err = AddCard("AlifMobi", 100, 1, db)
	if err != nil {
		t.Errorf("can't add card: %v", err)
	}

	err = TransferMoneyForPhoneNumber(9001, db)
	if err != nil {
		t.Errorf("can't search id cards for number phone: %v", err)
	}
}

func TestTransferMoneyCardNumber_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = TransferMoneyCardNumber("20216000000000001", db)
	if err == nil {
		t.Errorf("can't search id cards for number phone: %v", err)
	}
}

func TestTransferMoneyCardNumber_HasDbError(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cards
	(
		id      INTEGER PRIMARY KEY AUTOINCREMENT,
		numberCard TEXT NOT NULL,
		name    TEXT    NOT NULL,
		balance INTEGER NOT NULL CHECK ( balance > 0 ),
		user_id INTEGER NOT NULL 
	);`)

	if err != nil {
		t.Errorf("can't creat table cards for number card: %v", err)
	}
	_, err = db.Exec(`INSERT INTO cards(id,name, balance, user_id, numberCard) VALUES (1,"AlifMobi",200,0,"20216000000000001")`)
	if err != nil {
		t.Errorf("can't get all card, add card: %v", err)
	}

	err = TransferMoneyCardNumber("20216000000000001", db)
	if err != nil {
		t.Errorf("can't search id cards for number card: %v", err)
	}
}

func TestTransferMoney_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = TransferMoney(100, db)
	if err == nil {
		t.Errorf("can't trancfer money: %v", err)
	}

}

//func TestTransferMoney_HasDb(t *testing.T) {
//	db, err := sql.Open("sqlite3", ":memory:")
//	if err != nil {
//		t.Errorf("can't open db: %v", err)
//	}
//	defer func() {
//		if err := db.Close(); err != nil {
//			t.Errorf("can't close db: %v", err)
//		}
//	}()
//
//	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sumTransferUsers
//(
//    id      INTEGER PRIMARY KEY AUTOINCREMENT,
//    balance INTEGER NOT NULL
//);`)
//	if err != nil {
//		t.Errorf("can't trancfer money: %v", err)
//	}
//	_, err = db.Exec(`INSERT INTO sumTransferUsers(id,balance)
//VALUES (1,0)
//       ON CONFLICT DO NOTHING;`)
//	if err != nil {
//		t.Errorf("can't trancfer money: %v", err)
//	}
//
//	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cards
//(
//   id      INTEGER PRIMARY KEY AUTOINCREMENT,
//   numberCard TEXT NOT NULL,
//   name    TEXT    NOT NULL,
//   balance INTEGER NOT NULL CHECK ( balance > 0 ),
//   user_id INTEGER REFERENCES users(id)
//);`)
//	if err != nil {
//		t.Errorf("can't creat atm to get all atm: %v", err)
//	}
//
//	_, err = db.Exec(`INSERT INTO cards(id,name, balance, user_id, numberCard) VALUES (1,"AlifMobi",200,0,"20216000000000001")`)
//	if err != nil {
//		t.Errorf("can't get all card, add card: %v", err)
//	}
//
//	_, err = db.Exec(`INSERT INTO cards(id,name, balance, user_id, numberCard) VALUES (2,"AlifMobi",400,1,"20216000000000002")`)
//	if err != nil {
//		t.Errorf("can't get all card, add card: %v", err)
//	}
//
//	err = TransferMoney(100, db)
//	if err != nil {
//		t.Errorf("can't trancfer money: %v", err)
//	}
//
//}
