package main

import (
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Config struct {
	Output output
	Database database
}

type database struct {
	Server string
	Port string
	Database string
	User string
	Password string
}

type output struct {
	Directory string
	Format string
}


func main(){
	currentTime := time.Now()
	var conf Config
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", conf)

	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.Database.User, conf.Database.Password, conf.Database.Server, conf.Database.Port, conf.Database.Database)

	db, err := sql.Open("mysql", connString)

	if err != nil {
		fmt.Println(err.Error())
	}
	//var first_name string
	defer db.Close()

	router := mux.NewRouter()

	http.ListenAndServe(":8000", router)


	//err = db.QueryRow("SELECT first_name from persons").Scan(&first_name)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(first_name)
	fmt.Println("ShortYear : ", currentTime.Format("06-Jan-02"))
}