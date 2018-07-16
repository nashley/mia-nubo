package main

import (
	"fmt"
	"log"
	"strconv"
	"net/http"

	"github.com/jackc/pgx"
)

const (
	port = 41301
)

var pool *pgx.ConnPool

func get_file_info(file_id int) (file_name string, file_path string, err error) {
	conn, err := pool.Acquire()
	if err != nil {
		return
	}
	defer pool.Release(conn)
	
	err = pool.QueryRow("SELECT name, path FROM files WHERE id=$1", file_id).Scan(&file_name, &file_path)
	if err != nil {
		return
	}
	
	return
}

func serve_info(w http.ResponseWriter, r *http.Request) {
	file_id, err := strconv.Atoi(r.URL.Path[6:])
	if err != nil {
		fmt.Fprintf(w, "Unable to determine file ID: %s\nPlease enter a valid integer", err.Error())
		return
	}
	
	file_name, file_path, err := get_file_info(file_id)
	if err != nil {
		fmt.Fprintf(w, "DB Error: %s", err.Error())
		return
	}
	
	fmt.Fprintf(w, "ID: %d\nName: %s\nPath:%s", file_id, file_name, file_path)
}

func serve_files(w http.ResponseWriter, r *http.Request) {
	file_id, err := strconv.Atoi(r.URL.Path[7:])
	if err != nil {
		fmt.Fprintf(w, "Unable to determine file ID: %s\nPlease enter a valid integer", err.Error())
		return
	}

	_, file_path, err := get_file_info(file_id)
	if err != nil {
		fmt.Fprintf(w, "DB Error: %s", err.Error())
		return
	}

	http.ServeFile(w, r, file_path)
}

func main() {
	var err error

	sqlConnConfig := pgx.ConnConfig{
		User:              "mianubo",
		Password:          "mianubo",
		Host:              "localhost",
		Port:              5432,
		Database:          "mianubo",
		TLSConfig:         nil,
		UseFallbackTLS:    false,
		FallbackTLSConfig: nil,
	}

	pool, err = pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: sqlConnConfig})
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	http.HandleFunc("/info/", serve_info)
	http.HandleFunc("/files/", serve_files)
	fmt.Println("Running server...")
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), nil))
}
