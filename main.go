package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func openDB() error {
	db, err := sql.Open("sqlite3", "./prisma/dev.db")
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func closeDB() error {
	return DB.Close()
}

func main() {
	openDB()
	defer closeDB()

	mux := http.NewServeMux()
	mux.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./src/"))))

	mux.HandleFunc("/home", Home)
	mux.HandleFunc("/blog/", BlogId)

	mux.HandleFunc("/", GetBlogs)
	mux.HandleFunc("/getBlog/", GetBlog)

	err := http.ListenAndServe(":8000", addCORS(mux))
	if err != nil {
		log.Fatal(err)
	}
}
func addCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		h.ServeHTTP(w, r)
	})
}

const path = "./src"

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl := template.Must(template.ParseFiles(path + "/main.html"))

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func BlogId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl := template.Must(template.ParseFiles(path + "/getBlog.html"))

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

type Blog struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Article string `json:"article"`
}

func GetBlogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := DB.Query(`SELECT * FROM "Blog"`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	blogs := []Blog{}
	for rows.Next() {
		var blog Blog
		err := rows.Scan(&blog.Id, &blog.Title, &blog.Article)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		blogs = append(blogs, blog)
	}

	// Convert the blogs slice to JSON
	jsonData, err := json.Marshal(blogs)
	if err != nil {
		http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set content type and send JSON response
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonData)
}

func GetBlog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/getBlog/"):]
	row := DB.QueryRow(`SELECT * FROM "Blog" WHERE id=$1`, id)

	var blog Blog
	err := row.Scan(&blog.Id, &blog.Title, &blog.Article)
	if err != nil {
		http.Error(w, "Error fetching blog", http.StatusInternalServerError)
		return
	}

	// Convert the blog object to JSON
	jsonData, err := json.Marshal(blog)
	if err != nil {
		http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set content type and send JSON response
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonData)
}
