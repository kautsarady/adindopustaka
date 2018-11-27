package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kautsarady/adindopustaka/api"
	"github.com/kautsarady/adindopustaka/model"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	cs := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DBNAME"))

	dao, err := model.Make(cs)
	if err != nil {
		log.Fatal(err)
	}

	controller := api.Make(dao)

	addr := ":" + os.Getenv("PORT")
	log.Fatal(controller.Router.Run(addr))
}
