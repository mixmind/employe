package domain

type (
	RestResponse struct {
		ResponseMessage interface{} `json:"responseMessage,omitempty"`
		ResponseError   interface{} `json:"responseError,omitempty"`
	}

	Employe struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	EmployeeRole struct {
		EmployeeID int  `json:"employeeID"`
		RoleID     int  `json:"roleID"`
		Enabled    bool `json:"enabled"`
	}
	Attendance struct {
		EmployeeID int    `json:"employeeID"`
		RoleID     int    `json:"roleID"`
		ActionTime string `json:"actionTime"`
	}
	Role struct {
		RoleID      int    `json:"roleID"`
		Description string `json:"description"`
	}
)
