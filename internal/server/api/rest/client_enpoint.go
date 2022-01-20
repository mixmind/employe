package rest

import (
	"employe/internal/server/db"
	"employe/internal/server/domain"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type RestProcessor struct {
	dbStruct *db.DBStruct
	Router   *mux.Router
}

/*
Creates router and defines REST API's
*/
func NewServer(dbStruct *db.DBStruct) (*mux.Router, error) {
	log.Info("Launching REST API's")
	rtr := mux.NewRouter()
	restProcessor := RestProcessor{dbStruct: dbStruct}
	rtr.Handle("/GetEmployeeList", domain.WrapREST(restProcessor.employees)).Methods(http.MethodGet)
	rtr.Handle("/GetEmployeeRoles", domain.WrapREST(restProcessor.employeRoles)).Methods(http.MethodGet)
	rtr.Handle("/ClockIn", domain.WrapREST(restProcessor.clockIn)).Methods(http.MethodPost)
	restProcessor.Router = rtr
	return rtr, nil
}
