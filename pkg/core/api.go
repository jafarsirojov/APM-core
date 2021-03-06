package core

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

var ErrInvalidPass = errors.New("invalid password")

var onlineUserID int

var idCardForTransferRecipient int64

const tempNumberCard = 20216000000000000

type QueryError struct {
	Query string
	Err   error
}

type DbError struct {
	Err error
}

type DbTxError struct {
	Err         error
	RollbackErr error
}

type Atm struct {
	Id      int64
	Name    string
	Address string
}

type Service struct {
	Id      int64
	Name    string
	Balance int64
}

type Card struct {
	Id         int64
	Name       string
	Balance    int64
	User_id    int64
	NumberCard int
}

type User struct {
	Id             int64
	Name           string
	Login string
	Password string
	PassportSeries string
	NumberPhone    int
	HideShow       int
}

type UserHide struct {
	Id             int64
	Name           string
	PassportSeries string
	NumberPhone    int
}

type UserShow struct {
	Id             int64
	Name           string
	PassportSeries string
	NumberPhone    int
}

type OperationsLogging struct {
	Id              int64
	Name            string
	Time            string
	RecipientSender string
	Balance         int
	User_id         int
}

func (receiver *QueryError) Unwrap() error {
	return receiver.Err
}

func (receiver *QueryError) Error() string {
	return fmt.Sprintf("can't execute query %s: %s", loginManagerSQL, receiver.Err.Error())
}

func queryError(query string, err error) *QueryError {
	return &QueryError{Query: query, Err: err}
}

func (receiver *DbError) Error() string {
	return fmt.Sprintf("can't handle db operation: %v", receiver.Err.Error())
}

func (receiver *DbError) Unwrap() error {
	return receiver.Err
}

func dbError(err error) *DbError {
	return &DbError{Err: err}
}

func Init(db *sql.DB) (err error) {
	ddls := []string{managerDDL, usersDDL, cardsDDL, atmDDL, servicesDDL, sumTransferUsersDDL, operationsLoggingDDL}
	for _, ddl := range ddls {
		_, err = db.Exec(ddl)
		if err != nil {
			return err
		}
	}

	initialData := []string{managerInitialData, sumTransferUsersDDLInitialData}
	for _, datum := range initialData {
		_, err = db.Exec(datum)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoginManager(login, password string, db *sql.DB) (bool, error) {
	var dbLogin, dbPassword string

	err := db.QueryRow(
		loginManagerSQL,
		login).Scan(&dbLogin, &dbPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, queryError(loginManagerSQL, err)
	}

	if dbPassword != password {
		return false, ErrInvalidPass
	}

	return true, nil
}

func LoginUsers(login, password string, db *sql.DB) (bool, error) {
	var dbLogin, dbPassword string
	var dbHideShow int

	err := db.QueryRow(
		loginUsersSQL,
		login).Scan(&onlineUserID, &dbLogin, &dbPassword, &dbHideShow)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, queryError(loginUsersSQL, err)
	}
	if dbHideShow == 4 {
		fmt.Println("У вас нет доступа!!!\n Вы заблокированы менеджером!!! ")
		return false, nil
	}

	if dbPassword != password {
		return false, ErrInvalidPass
	}

	return true, nil
}

func AddAtm(atmName string, atmAddress string, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		insertAtmSQL,

		sql.Named("name", atmName),
		sql.Named("address", atmAddress),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetAllAtms(db *sql.DB) (atms []Atm, err error) {
	rows, err := db.Query(getAllAtmsSQL)
	if err != nil {
		return nil, queryError(getAllAtmsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			atms, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		atm := Atm{}
		err = rows.Scan(&atm.Id, &atm.Name, &atm.Address)
		if err != nil {
			return nil, dbError(err)
		}
		atms = append(atms, atm)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return atms, nil
}

func AddService(serviceName string, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		insertServiceSQL,

		sql.Named("name", serviceName),
		sql.Named("balance", 0),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetAllServices(db *sql.DB) (services []Service, err error) {
	rows, err := db.Query(getAllServicesSQL)
	if err != nil {
		return nil, queryError(getAllServicesSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			services, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		service := Service{}
		err = rows.Scan(&service.Id, &service.Name, &service.Balance)
		if err != nil {
			return nil, dbError(err)
		}
		services = append(services, service)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return services, nil
}

func AddCard(cardName string, cardBalance int64, cardUser_id int64, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	var cardNumberCard string
	selectDescIdFromCard := 0
	var numberCard int
	_ = tx.QueryRow(selectDescIdFromCardSQL).Scan(&selectDescIdFromCard)
	numberCard = tempNumberCard + selectDescIdFromCard + 1
	cardNumberCard = strconv.Itoa(numberCard)

	_, err = tx.Exec(
		insertCardSQL,

		sql.Named("name", cardName),
		sql.Named("balance", cardBalance),
		sql.Named("user_id", cardUser_id),
		sql.Named("numberCard", cardNumberCard),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetAllCards(db *sql.DB) (cards []Card, err error) {
	rows, err := db.Query(getAllCardsSQL)
	if err != nil {
		return nil, queryError(getAllCardsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			cards, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		card := Card{}
		err = rows.Scan(&card.Id, &card.Name, &card.Balance, &card.User_id, &card.NumberCard)
		if err != nil {
			return nil, dbError(err)
		}
		cards = append(cards, card)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return cards, nil
}

func AddUser(userName string, userLogin string, userPassword string, userPassportSeries string, userPhoneNumber int, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	userHideShow := 3

	_, err = tx.Exec(
		insertUserSQL,

		sql.Named("name", userName),
		sql.Named("login", userLogin),
		sql.Named("password", userPassword),
		sql.Named("passportSeries", userPassportSeries),
		sql.Named("phoneNumber", userPhoneNumber),
		sql.Named("hideShow", userHideShow),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetAllUsers(db *sql.DB) (users []User, err error) {
	rows, err := db.Query(getAllUsersSQL)
	if err != nil {
		return nil, queryError(getAllUsersSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			users, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.Name, &user.PassportSeries, &user.NumberPhone)
		if err != nil {
			return nil, dbError(err)
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return users, nil
}

func GetUserCards(db *sql.DB) (cards []Card, err error) {
	rows, err := db.Query(getUserCardsSQL, onlineUserID)
	if err != nil {
		return nil, queryError(getUserCardsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			cards, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		card := Card{}
		err = rows.Scan(&card.Id, &card.Name, &card.Balance, &card.NumberCard)
		if err != nil {
			return nil, dbError(err)
		}
		cards = append(cards, card)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}
	return cards, err
}

func TransferMoneyForPhoneNumber(phoneNumber int, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	var userIdRecipient int
	err = tx.QueryRow(selectIdUserPhoneNumberSQL, phoneNumber).Scan(&userIdRecipient)
	if err != nil {
		fmt.Println("Клиент с такой номера не зарегистрирован!!!")
		return err
	}

	err = tx.QueryRow(selectIdCardForTransferPhoneNumberSQL, userIdRecipient).Scan(&idCardForTransferRecipient)
	if err != nil {
		fmt.Println("У клиент нет счёта!!!")
		return err
	}

	return nil
}

func TransferMoneyCardNumber(countNumber string, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = tx.QueryRow(selectIdCardForTransferCountNumberSQL, countNumber).Scan(&idCardForTransferRecipient)
	if err != nil {
		fmt.Println("Введен неверный номер счёта!!!")
		return err
	}
	return err
}

func TransferMoney(currency int, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var currencySenderLast int
	err = tx.QueryRow(selectBalanceToCardSenderSQL, onlineUserID).Scan(&currencySenderLast)
	if err != nil {
		return err
	}
	if currencySenderLast < currency {
		fmt.Println("У вас нет таких денег в счету!!!")
		return nil
	}

	var currencySenderFirst int
	currencySenderFirst = currencySenderLast - currency

	_, err = tx.Exec(
		updateBalanceToCardSenderSQL, currencySenderFirst, onlineUserID,
	)
	if err != nil {
		return err
	}

	var currencyRecipientLast int
	err = tx.QueryRow(selectBalanceToCardRecipientSQL, idCardForTransferRecipient).Scan(&currencyRecipientLast)
	var currencyRecipientFirst int
	currencyRecipientFirst = currencyRecipientLast + currency

	_, err = tx.Exec(
		updateBalanceToCardRecipientSQL, currencyRecipientFirst, idCardForTransferRecipient,
	)
	if err != nil {
		return err
	}
	var sumTransferUsers int
	err = tx.QueryRow(selectBalanceSumTransferUsers).Scan(&sumTransferUsers)
	if err != nil {
		return err
	}

	var numberCard string
	err = tx.QueryRow(selectNumberCardToIdCardSQL, idCardForTransferRecipient).Scan(&numberCard)
	if err != nil {
		log.Fatalf("can't operationLogging, select number card for id card: %s", err)
	}
	t := time.Now().String()
	_, err = tx.Exec(insertOperationsLoggingSQL, "translatedToSend", t, numberCard, -currency, onlineUserID)
	if err != nil {
		log.Fatalf("can't operationLogging %s", err)
	}
	//------------
	var idUserGet int
	err = tx.QueryRow(selectUser_idWhereIdCardSQL, idCardForTransferRecipient).Scan(&idUserGet)
	if err != nil {
		log.Fatalf("can't operationLogging, select number card for id card: %s", err)
	}
	err = tx.QueryRow(selectNumberCardFromUser_idCardSQL, onlineUserID).Scan(&numberCard)
	if err != nil {
		log.Fatalf("can't operationLogging, select number card for id card: %s", err)
	}
	_, err = tx.Exec(insertOperationsLoggingSQL, "translatedToGet", t, numberCard, currency, idUserGet)
	if err != nil {
		log.Fatalf("can't operationLogging %s", err)
	}

	sumTransferUsers = sumTransferUsers + currency
	_, err = tx.Exec(updateBalanceSumTransferUsersSQL, sumTransferUsers)
	if err != nil {
		return err
	}
	return nil
}

func TransferServices(currency int, name string, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	var currencyUser int
	err = tx.QueryRow(selectBalanceToCardSenderSQL, onlineUserID).Scan(&currencyUser)
	if err != nil {
		return err
	}
	currencyUser = currencyUser - currency
	_, err = tx.Exec(
		updateBalanceToCardSenderSQL, currencyUser, onlineUserID,
	)
	if err != nil {
		return err
	}

	var currencyService int
	err = tx.QueryRow(selectBalanceOnServiceSQL, name).Scan(&currencyService)
	if err != nil {
		return err
	}
	currencyService = currencyService + currency
	_, err = tx.Exec(
		updateBalanceServiceSQL, currencyService, name,
	)
	t := time.Now().String()

	_,err = tx.Exec(insertOperationsLoggingSQL,"payToService",t,name,-currency,onlineUserID)
	if err != nil {
		return err
	}
	return nil
}

func UserHideManager(userId int, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		updateHideShowUser, 4, userId,
	)

	if err != nil {
		return err
	}

	return nil
}

func UserShowManager(userId int, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		updateHideShowUser, 3, userId,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetHideUsers(db *sql.DB) (users []UserHide, err error) {
	rows, err := db.Query(getHideUserSQL, 4)
	if err != nil {
		return nil, nil
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			users, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		user := UserHide{}
		err = rows.Scan(&user.Id, &user.Name, &user.PassportSeries, &user.NumberPhone)
		if err != nil {
			return nil, dbError(err)
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}
	return users, nil
}

func GetShowUsers(db *sql.DB) (users []UserShow, err error) {
	rows, err := db.Query(getHideUserSQL, 3)
	if err != nil {
		return nil, nil
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			users, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		user := UserShow{}
		err = rows.Scan(&user.Id, &user.Name, &user.PassportSeries, &user.NumberPhone)
		if err != nil {
			return nil, dbError(err)
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}
	return users, nil
}

func SearchUserByPhoneNumber(phoneNumber int, db *sql.DB) (users []User, err error) {
	rows, err := db.Query(
		searchUserForPhoneNumberSQL, phoneNumber,
	)
	if err != nil {
		return nil, queryError(getAllUsersSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			users, err = nil, dbError(innerErr)
		}
	}()
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.Name, &user.PassportSeries, &user.NumberPhone)
		if err != nil {
			return nil, dbError(err)
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return users, nil
}

func StaticCountUsers(db *sql.DB) (count int) {
	db.QueryRow(staticCountUserSQL).Scan(&count)
	return count
}

func StaticSumBalanceUsers(db *sql.DB) (sum int) {
	db.QueryRow(staticSumBalanceUsersSQL).Scan(&sum)
	return sum
}

func StaticBalanceOfServices(db *sql.DB) (sum int) {
	db.QueryRow(staticBalanceOfServicesSQL).Scan(&sum)
	return sum
}

func StaticBalanceSumTransfer(db *sql.DB) (sum int) {
	db.QueryRow(selectBalanceSumTransferUsers).Scan(&sum)
	return sum
}

func ViewOperationsLogging(db *sql.DB) (opLogs []OperationsLogging, err error) {
	rows, err := db.Query(getOperationsLoggingUserSQL, onlineUserID)
	if err != nil {
		return nil, queryError(getOperationsLoggingUserSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			opLogs, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		opLog := OperationsLogging{}
		err = rows.Scan(&opLog.Id, &opLog.Name, &opLog.Time, &opLog.RecipientSender, &opLog.Balance)
		if err != nil {
			return nil, dbError(err)
		}
		opLogs = append(opLogs, opLog)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}
	return opLogs, err
}

func ViewOperationsLoggingToSearch(idUser int, db *sql.DB) (opLogs []OperationsLogging, err error) {
	rows, err := db.Query(getOperationsLoggingUserSQL, idUser)
	if err != nil {
		return nil, queryError(getOperationsLoggingUserSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			opLogs, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		opLog := OperationsLogging{}
		err = rows.Scan(&opLog.Id, &opLog.Name, &opLog.Time, &opLog.RecipientSender, &opLog.Balance)
		if err != nil {
			return nil, dbError(err)
		}
		opLogs = append(opLogs, opLog)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}
	return opLogs, err
}

func ViewAllOperationsLogging(db *sql.DB) (opLogs []OperationsLogging, err error) {
	rows, err := db.Query(getAllOperationsLoggingUserSQL)
	if err != nil {
		return nil, queryError(getAllOperationsLoggingUserSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			opLogs, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		opLog := OperationsLogging{}
		err = rows.Scan(&opLog.Id, &opLog.Name, &opLog.Time, &opLog.RecipientSender, &opLog.Balance)
		if err != nil {
			return nil, dbError(err)
		}
		opLogs = append(opLogs, opLog)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}
	return opLogs, err
}


//--------------------------------

func ExportClientsToJSON(db *sql.DB) error {
	return ExportToFile(db, getAllUsersSQL, "clients.json",
		mapRowToClient, json.Marshal, mapInterfaceSliceToClients)
}
func ExportAtmsToJSON(db *sql.DB) error {
	return ExportToFile(db, getAllAtmsSQL, "atms.json",
		mapRowToAtm, json.Marshal,
		mapInterfaceSliceToAtms)
}

//XML

func ExportClientsToXML(db *sql.DB) error {
	return ExportToFile(db, getAllUsersSQL, "clients.xml",
		mapRowToClient, xml.Marshal, mapInterfaceSliceToClients)
}
func ExportAtmsToXML(db *sql.DB) error {
	return ExportToFile(db, getAllAtmsSQL, "atms.xml",
		mapRowToAtm, xml.Marshal,
		mapInterfaceSliceToAtms)
}

func mapRowToClient(rows *sql.Rows) (interface{}, error) {
	user := User{}
	err := rows.Scan(&user.Id, &user.Login, &user.Password, &user.Name, &user.NumberPhone, &user.HideShow)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func mapRowToAtm(rows *sql.Rows) (interface{}, error) {
	atm := Atm{}
	err := rows.Scan(&atm.Id,&atm.Name, &atm.Address)
	if err != nil {
		return nil, err
	}
	return atm, nil
}
type ClientsExport struct {
	Users []User
}
func mapInterfaceSliceToClients(ifaces []interface{}) interface{} {
	clients := make([]User, len(ifaces))
	for i := range ifaces {
		clients[i] = ifaces[i].(User)
	}
	clientsExport := ClientsExport{Users: clients}
	return clientsExport
}
func mapInterfaceSliceToAtms(ifaces []interface{}) interface{} {
	atms := make([]Atm, len(ifaces))
	for i := range ifaces {
		atms[i] = ifaces[i].(Atm)
	}
	atmsExport := AtmsExport{Atms: atms}
	return atmsExport
}
func ImportClientsFromJSON(db *sql.DB) error {
	return ImportFromFile(
		db,
		"clients.json",
		func(data []byte) ([]interface{}, error) {
			return mapBytesToClients(data, json.Unmarshal)
		},
		insertClientToDB,
	)
}
func ImportAtmsFromJSON(db *sql.DB) error {
	return ImportFromFile(
		db,
		"atms.json",
		func(data []byte) ([]interface{}, error) {
			return mapBytesToAtms(data, json.Unmarshal)
		},
		insertAtmToDB,
	)
}
func ImportClientsFromXML(db *sql.DB) error {
	return ImportFromFile(
		db,
		"clients.xml",
		func(data []byte) ([]interface{}, error) {
			return mapBytesToClients(data, xml.Unmarshal)
		},
		insertClientToDB,
	)
}
func ImportAtmsFromXML(db *sql.DB) error {
	return ImportFromFile(
		db,
		"atms.xml",
		func(data []byte) ([]interface{}, error) {
			return mapBytesToAtms(data, xml.Unmarshal)
		},
		insertAtmToDB,
	)
}
func mapBytesToClients(data []byte, unmarshal func([]byte, interface{}) error,) ([]interface{}, error) {
	clientsExport := ClientsExport{}
	err := unmarshal(data, &clientsExport)
	if err != nil {
		return nil, err
	}
	ifaces := make([]interface{}, len(clientsExport.Users))
	for index := range ifaces {
		ifaces[index] = clientsExport.Users[index]
	}
	return ifaces, nil
}
func insertClientToDB(iface interface{}, db *sql.DB) error {
	client := iface.(User)
	_, err := db.Exec(
		insertUserSQL,
		sql.Named("name", client.Name),
		sql.Named("login", client.Login),
		sql.Named("password", client.Password),
		sql.Named("passportSeries", client.PassportSeries),
		sql.Named("phoneNumber", client.NumberPhone),
		sql.Named("hideShow", client.HideShow),
	)
	if err != nil {
		return err
	}
	return nil
}

type AtmsExport struct {
	Atms []Atm
}

func mapBytesToAtms(data []byte, unmarshal func([]byte, interface{}) error,
) ([]interface{}, error) {
	atmsExport := AtmsExport{}
	err := unmarshal(data, &atmsExport)
	if err != nil {
		return nil, err
	}
	ifaces := make([]interface{}, len(atmsExport.Atms))
	for index := range ifaces {
		ifaces[index] = atmsExport.Atms[index]
	}
	return ifaces, nil
}
func insertAtmToDB(iface interface{}, db *sql.DB) error {
	atm := iface.(Atm)
	_, err := db.Exec(
		insertAtmSQL,
		sql.Named("id", atm.Id),
		sql.Named("name", atm.Name),
		sql.Named("address", atm.Address),
	)
	if err != nil {
		return err
	}
	return nil
}

type MapperRowTo func(rows *sql.Rows) (interface{}, error)
type MapperInterfaceSliceTo func([]interface{}) interface{}
type Marshaller func(interface{}) ([]byte, error)

func ExportToFile(
	db *sql.DB,
	getDataFromDbSQL string,
	filename string,
	mapRow MapperRowTo,
	marshal Marshaller,
	mapDataSlice MapperInterfaceSliceTo) error {

	rows, err := db.Query(getDataFromDbSQL)
	if err != nil {
		return err
	}
	defer func() {
		err = rows.Close()
	}()
	var dataSlice []interface{}
	for rows.Next() {
		dataElement, err := mapRow(rows)
		if err != nil {
			return err
		}
		dataSlice = append(dataSlice, dataElement)
	}
	exportData := mapDataSlice(dataSlice)
	data, err := marshal(exportData)
	err = ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		return err
	}
	return nil
}


type MapperBytesTo func([]byte) ([]interface{}, error)

func ImportFromFile(
	db *sql.DB,
	filename string,
	mapBytes MapperBytesTo,
	insertToDB func(interface{}, *sql.DB) error,
) error {
	itemsData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	sliceData, err := mapBytes(itemsData)

	for _, datum := range sliceData {
		err = insertToDB(datum, db)
		if err != nil {
			return err
		}
	}
	return nil
}