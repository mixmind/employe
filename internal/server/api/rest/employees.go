package rest

import (
	"employe/internal/server/cmds"
	"employe/internal/server/domain"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (restPr *RestProcessor) employees(writer http.ResponseWriter, request *http.Request) {
	employeProcessor := cmds.NewEmployeProcessor(restPr.dbStruct)
	responseCode := http.StatusOK
	var responseMessage interface{}
	var err error
	switch request.Method {
	case http.MethodGet:
		responseMessage, err = employeProcessor.GetEmployeesFromDB()
	default:
		responseCode = http.StatusBadRequest
		responseMessage = "This method is not allowed"
	}

	if _, err := domain.WriteResponse(writer, responseCode, responseMessage, err); err != nil {
		log.Error(errors.Wrap(err, "Error occurred during writing response"))
	}

}

func (restPr *RestProcessor) employeRoles(writer http.ResponseWriter, request *http.Request) {
	employeProcessor := cmds.NewEmployeProcessor(restPr.dbStruct)
	responseCode := http.StatusOK
	var responseMessage interface{}
	var err error
	var employeeID int
	skipProcessing := false
	urlValues := request.URL.Query()
	if len(urlValues) != 0 {
		if val, isExists := urlValues["employeeId"]; isExists && len(val) == 1 {
			employeeID, err = strconv.Atoi(val[0])
			if err != nil {
				skipProcessing = true
				responseCode = http.StatusBadRequest
			}
		}
	}
	if !skipProcessing {
		switch request.Method {
		case http.MethodGet:
			responseMessage, err = employeProcessor.GetEmployeRoleByID(employeeID)
		default:
			responseCode = http.StatusBadRequest
			responseMessage = "This method is not allowed"
		}
	}

	if _, err := domain.WriteResponse(writer, responseCode, responseMessage, err); err != nil {
		log.Error(errors.Wrap(err, "Error occurred during writing response"))
	}
}

func (restPr *RestProcessor) clockIn(writer http.ResponseWriter, request *http.Request) {
	employeProcessor := cmds.NewEmployeProcessor(restPr.dbStruct)
	responseCode := http.StatusOK
	var responseMessage interface{}
	var err error
	var role domain.Attendance
	skipProcessing := false
	err = parseBodyToObj(request, &role)
	if err != nil {
		log.Error(err)
		responseCode = http.StatusBadRequest
		skipProcessing = true
	}
	if !skipProcessing {
		switch request.Method {
		case http.MethodPost:
			result, err := employeProcessor.GetEmployeRoleByID(role.EmployeeID)
			if err != nil {
				responseCode = http.StatusBadRequest
			} else {
				isMatch := false
				for _, row := range result {
					if row.RoleID == role.RoleID {
						isMatch = true
						break
					}
				}
				if !isMatch || len(result) == 0 {
					log.Errorf("Failed to find such employee %d with provided role %d", role.EmployeeID, role.RoleID)
					responseCode = http.StatusBadRequest
					if _, err := domain.WriteResponse(writer, responseCode, responseMessage, err); err != nil {
						log.Error(errors.Wrap(err, "Error occurred during writing response"))
					}
					return
				}

				employeProcessor.InsertAttendanceInDB(role)
			}
		default:
			responseCode = http.StatusBadRequest
			responseMessage = "This method is not allowed"
		}

	}

	if _, err := domain.WriteResponse(writer, responseCode, responseMessage, err); err != nil {
		log.Error(errors.Wrap(err, "Error occurred during writing response"))
	}
}

func parseBodyToObj(request *http.Request, obj interface{}) error {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return errors.Wrap(err, "Incorrect body in request")
	}
	log.Debugf("Request body [%s]", string(requestBody))
	err = json.Unmarshal(requestBody, obj)
	if err != nil {
		return errors.Wrap(err, "Incorrect format of body")
	}
	return nil
}
