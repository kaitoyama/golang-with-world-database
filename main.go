package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type City struct {
	ID          int    `json:"id,omitempty"  db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

func main() {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}

	fmt.Println("Connected!")
	var city City
	cities := []City{}

	query_city := "SELECT * FROM city WHERE Name=?"
	query_cities := "SELECT * FROM city WHERE CountryCode=?"
	countryCity := "INSERT INTO city (CountryCode, Name) VALUES (?, ?)"
	delete_city := "DELETE FROM city WHERE Name=?"

	switch os.Args[1] {
	case "city":
		input_city := os.Args[2]
		if err := db.Get(&city, query_city, input_city); errors.Is(err, sql.ErrNoRows) {
			log.Printf("no such city Name = %s", input_city)
		} else if err != nil {
			log.Fatalf("DB Error: %s", err)
		}

		fmt.Printf(input_city+"の人口は%d人です\n", city.Population)

	case "cities":
		input_country := os.Args[2]
		if err := db.Select(&cities, query_cities, input_country); errors.Is(err, sql.ErrNoRows) {
			log.Printf("no such country Name = %s", input_country)
		} else if err != nil {
			log.Fatalf("DB Error: %s", err)
		}

		fmt.Println("日本の都市一覧")
		for _, city := range cities {
			fmt.Printf("都市名: %s, 人口: %d人\n", city.Name, city.Population)

		}
	case "add":
		input_city_country := os.Args[2]
		input_city_name := os.Args[3]
		db.Exec(countryCity, input_city_country, input_city_name)

	case "delete":
		input_city_name := os.Args[2]
		if err := db.Select(&city, delete_city, input_city_name); errors.Is(err, sql.ErrNoRows) {
			log.Printf("no such country Name = %s", input_city_name)
		}
	}

}
