package db

import (
	"database/sql"
	"employe/internal/server/domain"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type DBStruct struct {
	internalDB *sql.DB
}

var (
	employees = []domain.Employe{domain.Employe{Name: "Joe Cocker"},
		domain.Employe{Name: "JFK"}}
	Roles        = []domain.Role{domain.Role{RoleID: 1, Description: "Manager"}, domain.Role{RoleID: 2, Description: "Waiter"}}
	employeRoles = []domain.EmployeeRole{domain.EmployeeRole{EmployeeID: 1, RoleID: 1, Enabled: true},
		domain.EmployeeRole{EmployeeID: 1, RoleID: 2, Enabled: true},
		domain.EmployeeRole{EmployeeID: 2, RoleID: 2, Enabled: true}}
)

/*
Generates and fills Inmemory db with employe
*/
func NewDBStruct() (*DBStruct, error) {
	inMemoryDB, err := sql.Open("sqlite3", "file:employe.db?cache=shared&mode=memory&_fk=true")
	if err != nil {
		return nil, err
	}
	db := DBStruct{internalDB: inMemoryDB}
	log.Info("Prefilling DB")
	if err := db.createInMemoryTables(); err != nil {
		return nil, errors.Wrap(err, "Failed to create tables")
	}

	if err := db.insertMockIntoDB(); err != nil {
		return nil, errors.Wrap(err, "Failed to insert employe into DB")
	}
	log.Info("DB filled sussesfully")
	return &db, nil
}

func NewDBStructWithDBProvided(inMemoryDB *sql.DB) *DBStruct {
	return &DBStruct{internalDB: inMemoryDB}
}

/*
Creates in memory tables such as employe and employe role
*/
func (db *DBStruct) createInMemoryTables() error {
	log.Info("Creating employe table")
	if _, err := db.internalDB.Exec(createEmployeTable); err != nil {
		return errors.Wrapf(err, "Failed to create employe table")
	}
	log.Info("Employe table created sussesfully")
	log.Info("Creating employe role table")
	if _, err := db.internalDB.Exec(createEmployeRoleTable); err != nil {
		return errors.Wrapf(err, "Failed to create employe role table")
	}
	log.Info("Employe role table created sussesfully")
	log.Info("Creating attendance table")
	if _, err := db.internalDB.Exec(createAttendanceTable); err != nil {
		return errors.Wrapf(err, "Failed to create attendance table")
	}
	log.Info("Attendance role table created sussesfully")
	return nil
}

func (db *DBStruct) insertMockIntoDB() error {
	log.Println("Inserting mock records")
	tx, err := db.internalDB.Begin()
	if err != nil {
		return errors.Wrap(err, "Failed to start a transaction")
	}
	statement, err := tx.Prepare(InsertIntoEmployeTable)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare stmt")
	}
	for _, employe := range employees {
		_, err := statement.Exec(employe.Name)
		if err != nil {
			return errors.Wrap(err, "Failed to execute a prepared statement")
		}
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "Failed to commit a transaction")
	}
	log.Info("Employes records inserted sussesfully")

	tx, err = db.internalDB.Begin()
	if err != nil {
		return errors.Wrap(err, "Failed to start a transaction")
	}
	statement, err = tx.Prepare(InsertIntoEmployeRoleTable)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare stmt")
	}
	for _, employe := range employeRoles {
		_, err := statement.Exec(employe.EmployeeID,
			employe.RoleID,
			employe.Enabled)
		if err != nil {
			return errors.Wrap(err, "Failed to execute a prepared statement")
		}
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "Failed to commit a transaction")
	}
	log.Info("Employes records inserted sussesfully")

	return nil
}

func (db *DBStruct) BeginTransaction() (*sql.Tx, error) {
	return db.internalDB.Begin()
}

func (db *DBStruct) Query(query string) (*sql.Rows, error) {
	return db.internalDB.Query(query)
}

func (db *DBStruct) Prepare(query string) (*sql.Stmt, error) {
	return db.internalDB.Prepare(query)
}
