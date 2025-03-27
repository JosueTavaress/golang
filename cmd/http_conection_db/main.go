package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

func setupDatabase() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/movies_db")
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

type Actor struct {
	ID        int
	CreatedAt interface{}
	UpdatedAt interface{}
	FirstName string
	LastName  string
	Rating    float64
}

// / ------------
func main() {
	r := chi.NewRouter()
	db := setupDatabase()
	defer db.Close()

	actors, err := getActors(db)
	if err != nil {
		log.Fatal(err)
	}


	// query select actors
	r.Get("/actors", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		bytes, err := json.Marshal(actors)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(bytes))
	})

	http.ListenAndServe(":8080", r)
}

func getActors(db *sql.DB) ([]Actor, error) {
	rows, err := db.Query("SELECT id, created_at, updated_at, first_name, last_name, rating FROM actors;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []Actor

	for rows.Next() {
		var actor Actor
		if err := rows.Scan(&actor.ID, &actor.CreatedAt, &actor.UpdatedAt, &actor.FirstName, &actor.LastName, &actor.Rating); err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actors, nil
}
