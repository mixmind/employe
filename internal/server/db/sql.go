package db

var (
	createEmployeTable = `CREATE TABLE employe(id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
					name text);`
	createEmployeRoleTable = `CREATE TABLE employe_role(employee_id INTEGER NOT NULL,
					role_id INTEGER NOT NULL,
					enabled boolean,
					FOREIGN KEY(employee_id) REFERENCES employe(id) ON DELETE RESTRICT
					);`
	createAttendanceTable = `CREATE TABLE attendance(employee_id INTEGER NOT NULL,
						role_id INTEGER NOT NULL,
						action_time TIMESTAMP
								DEFAULT CURRENT_TIMESTAMP,
						FOREIGN KEY(employee_id) REFERENCES employe(id) ON DELETE RESTRICT
						);`
	InsertIntoEmployeTable     = `INSERT INTO employe(name) VALUES (?)`
	InsertIntoEmployeRoleTable = `INSERT INTO employe_role(employee_id,
													role_id,
													enabled) VALUES (?,?,?)`
	InsertIntoAttendanceTable = `INSERT INTO attendance(employee_id,
														role_id) VALUES (?,?)`
	SelectEmployees = `SELECT id,
					name FROM employe`
	SelectEmployeRoleForAttendance = `SELECT employee_id,
							role_id,
							enabled
						FROM employe_role er
						JOIN employe e ON er.employe_role=e.id
						WHERE role_id=? and employee_id=?`
	SelectEmployeRole = `SELECT
						role_id
					FROM employe_role 
					WHERE employee_id=%d`
)
