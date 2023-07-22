package db

import (
	"database/sql"
	"github.com/Rosya-edwica/position-to-demand/logger"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Host       string
	Port       string
	User       string
	Password   string
	Name       string
	Connection *sql.DB
}

func (d *Database) NewConnection() {
	connection, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.User, d.Password, d.Host, d.Port, d.Name))
	checkErr(err)
	d.Connection = connection
	logger.Log.Printf("Успешно подключились к БД: %s", d.Name)
}

func (d *Database) CloseConnection() {
	d.Connection.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}