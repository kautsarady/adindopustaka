package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/kautsarady/adindopustaka/api"
	"github.com/kautsarady/adindopustaka/model"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	cs := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DBNAME"))

	dao, err := model.Make(cs)
	if err != nil {
		log.Fatalf("Failed connecting to database: %v", err)
	}
	defer dao.DB.Close()

	controller := api.Make(dao)

	addr := ":" + os.Getenv("PORT")
	if err := controller.Router.Run(addr); err != nil {
		log.Fatalf("Failed running server: %v", err)
	}
}
