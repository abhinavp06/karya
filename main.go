package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	ID          int    `json:"id"`
	ParentID    *int   `json:"parent_id"`
	Description string `json:"description"`
}

var db *sql.DB

// TODO: check golang naming conventions (functions, variables, etc.)
func health_check(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Do the task!")
}

// TODO: add validations and pagination
func get_items(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, parent_id, description FROM items ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID, &it.ParentID, &it.Description); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// TODO: add 3 levels validation
func create_item(w http.ResponseWriter, r *http.Request) {
	var it Item
	if err := json.NewDecoder(r.Body).Decode(&it); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO items (parent_id, description) VALUES (?, ?)", it.ParentID, it.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	it.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(it)
}

// TODO: split into net.Listen + http.Serve.. print only if listen works
func main() {
	addr := ":3010"

	fmt.Println("Starting server on", addr)

	http.HandleFunc("/", health_check)
	http.HandleFunc("GET /items", get_items)
	http.HandleFunc("POST /items", create_item)

	dbFile := "karya.db"

	var err error
	db, err = sql.Open("sqlite3", "file:"+dbFile+"?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL&_foreign_keys=on")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	schema := `
CREATE TABLE IF NOT EXISTS items (
  id INTEGER PRIMARY KEY,
  parent_id INTEGER REFERENCES items(id),
  description TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_items_parent ON items(parent_id);`

	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
