package cmds

import (
	"employe/internal/server/db"
	"employe/internal/server/domain"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

type EmployeProcessor struct {
	dbStruct *db.DBStruct
}

func NewEmployeProcessor(dbStruct *db.DBStruct) *EmployeProcessor {
	return &EmployeProcessor{dbStruct: dbStruct}
}

func (employePr *EmployeProcessor) InsertAttendanceInDB(att domain.Attendance) (int64, error) {
	tx, err := employePr.dbStruct.BeginTransaction()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to start a transaction")
	}
	statement, err := tx.Prepare(db.InsertIntoAttendanceTable)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to prepare stmt")
	}
	res, err := statement.Exec(att.EmployeeID,
		att.RoleID)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to execute a prepared statement")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to execute a extract last id")
	}
	if err := tx.Commit(); err != nil {
		return 0, errors.Wrap(err, "Failed to commit a transaction")
	}

	return id, nil
}

/*
Get employees from DB
*/
func (employePr *EmployeProcessor) GetEmployeesFromDB() ([]domain.Employe, error) {
	rows, err := employePr.dbStruct.Query(db.SelectEmployees)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute a sql query")
	}

	var result []domain.Employe
	for rows.Next() {
		var receivedRow domain.Employe
		err = rows.Scan(&receivedRow.ID, &receivedRow.Name)
		if err != nil {
			log.Error(err)
			continue
		}

		result = append(result, receivedRow)
	}
	return result, nil
}

func (employePr *EmployeProcessor) GetEmployeRoleByID(id int) ([]domain.Role, error) {
	rows, err := employePr.dbStruct.Query(fmt.Sprintf(db.SelectEmployeRole, id))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute a sql query")
	}

	var result []domain.Role
	for rows.Next() {
		var receivedRow domain.Role
		err = rows.Scan(&receivedRow.RoleID)
		if err != nil {
			log.Error(err)
			continue
		}

		result = append(result, receivedRow)
	}
	for index, row := range result {
		for _, role := range db.Roles {
			if row.RoleID == role.RoleID {
				result[index].Description = role.Description
				break
			}
		}
	}
	return result, nil
}

func (employePr *EmployeProcessor) ClockIn(employee domain.EmployeeRole) error {
	tx, err := employePr.dbStruct.BeginTransaction()
	if err != nil {
		return errors.Wrap(err, "Failed to start a transaction")
	}
	statement, err := tx.Prepare(db.InsertIntoAttendanceTable)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare stmt")
	}

	_, err = statement.Exec(employee.EmployeeID, employee.RoleID)
	if err != nil {
		return errors.Wrap(err, "Failed to execute a prepared statement")
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "Failed to commit a transaction")
	}
	log.Info("Employes records inserted sussesfully")
	return nil
}
