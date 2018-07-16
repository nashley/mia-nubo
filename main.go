package main

import (
	"fmt"
	"log"
	"strconv"
	"os"
	
	"net/http"

	"github.com/jackc/pgx"
)

const (
	port = 41301
)

var pool *pgx.ConnPool

func get_file_info(w http.ResponseWriter, file_id_string string) (file_id int, file_name string, file_path string, err error) {
	file_id, err = strconv.Atoi(file_id_string)
	if err != nil {
		fmt.Fprintf(w, "Unable to determine file ID: %s\nPlease enter a valid integer", err.Error())
		return
	}

	conn, err := pool.Acquire()
	if err != nil {
		fmt.Fprintf(w, "Internal Server Error: %s", err.Error())
		return
	}
	defer pool.Release(conn)
	
	err = pool.QueryRow("SELECT name, path FROM files WHERE id=$1", file_id).Scan(&file_name, &file_path)
	if err != nil {
		switch err {
			case pgx.ErrNoRows:
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "File not found: %s", err.Error())
			default:
				fmt.Fprintf(w, "Internal Server Error: %s", err.Error())
		}
		return
	}
	
	return
}

func serve_info(w http.ResponseWriter, r *http.Request) {
	file_id, file_name, file_path, err := get_file_info(w, r.URL.Path[6:])
	if err == nil {
		fmt.Fprintf(w, "ID: %d\nName: %s\nPath:%s", file_id, file_name, file_path)
	}
}

func stream_files(w http.ResponseWriter, r *http.Request) {
	_, file_name, file_path, err := get_file_info(w, r.URL.Path[8:])
	if err != nil {
		return
	}
	
	file_content, err := os.Open(file_path)
	if err != nil {
		fmt.Fprintf(w, "File not found: %s", file_path)
		w.WriteHeader(http.StatusNotFound)
	}

	defer file_content.Close()
	
	file_status, err := os.Stat(file_path)
	if err != nil {
		fmt.Fprintf(w, "Error getting file status: %s", err.Error())
	}

	http.ServeContent(w, r, file_name, file_status.ModTime(), file_content)
}

func download_files(w http.ResponseWriter, r *http.Request) {
	_, file_name, file_path, err := get_file_info(w, r.URL.Path[10:])
	if err != nil {
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\"" + file_name + "\"")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
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
	http.HandleFunc("/download/", download_files)
	http.HandleFunc("/stream/", stream_files)
	fmt.Println("Running server...")
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), nil))
}
