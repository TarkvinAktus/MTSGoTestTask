package main

import (
	"database/sql"
	"testing"
	"time"
)

var testLogin = "testLogin"

func TestCreateUser(t *testing.T) {

	result := new(string)
	userArgs := new(Args)

	var user User
	var err error

	//DB
	connStr := "user=postgres password=root dbname=users sslmode=disable"
	DBConn, err = sql.Open("postgres", connStr)
	if err != nil {
		t.Error("DB connection error")
	}
	defer DBConn.Close()

	//args for user.create
	userArgs.ID = 1
	userArgs.Login = testLogin

	err = user.Create(userArgs, result)
	if err != nil {
		t.Error("User create error")
	}

	//Remove from DB
	_, err = DBConn.Exec("DELETE FROM users WHERE login = $1", testLogin)
	if err != nil {
		t.Error("User delete error")
	}
}

func TestGetUser(t *testing.T) {

	resultUser := new(User)
	userArgs := new(Args)

	var user User
	var err error

	//DB
	connStr := "user=postgres password=root dbname=users sslmode=disable"
	DBConn, err = sql.Open("postgres", connStr)
	if err != nil {
		t.Error("DB connection error")
	}
	defer DBConn.Close()

	//Insert testLogin
	_, err = DBConn.Exec("INSERT INTO users (login,registration_date) VALUES ($1,$2)", testLogin, time.Now().Unix())
	if err != nil {
		t.Error("DB create error")
	}

	row := DBConn.QueryRow("SELECT * FROM users WHERE login = $1", testLogin)

	err = row.Scan(&user.UUID, &user.Login, &user.RegistrationDate)
	if err != nil {
		t.Error("User get id error")
	}

	//Get testLogin
	userArgs.ID = user.UUID

	err = user.Get(userArgs, resultUser)
	if err != nil || resultUser.Login != testLogin {
		t.Error("User get error")
	}

	//Remove from DB
	_, err = DBConn.Exec("DELETE FROM users WHERE login = $1", testLogin)
	if err != nil {
		t.Error("User delete error")
	}
}

func TestUpdateUser(t *testing.T) {

	resultUser := new(User)
	result := new(string)
	userArgs := new(Args)

	var user User
	var err error

	//DB
	connStr := "user=postgres password=root dbname=users sslmode=disable"
	DBConn, err = sql.Open("postgres", connStr)
	if err != nil {
		t.Error("DB connection error")
	}
	defer DBConn.Close()

	//Insert testLogin
	_, err = DBConn.Exec("INSERT INTO users (login,registration_date) VALUES ($1,$2)", testLogin, time.Now().Unix())
	if err != nil {
		t.Error("DB create error")
	}

	row := DBConn.QueryRow("SELECT * FROM users WHERE login = $1", testLogin)

	err = row.Scan(&user.UUID, &user.Login, &user.RegistrationDate)
	if err != nil {
		t.Error("User get id error")
	}

	//Update testLogin to Login test
	userArgs.ID = user.UUID
	userArgs.Login = "Login test"

	err = user.Update(userArgs, result)
	if err != nil {
		t.Error("User update error")
	}

	//Get Login test != Test Login
	err = user.Get(userArgs, resultUser)
	if err != nil || resultUser.Login != "Login test" {
		t.Error("User get error")
	}

	//Remove from DB
	_, err = DBConn.Exec("DELETE FROM users WHERE login='Login test'")
	if err != nil {
		t.Error("User delete error")
	}
}
