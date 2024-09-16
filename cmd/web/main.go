package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore" // New import
	"github.com/alexedwards/scs/v2"         // New import
	"github.com/go-playground/form"
	_ "github.com/go-sql-driver/mysql"
	"snippetbox.mohit.net/internal/models"
)

/* Defining an application struct to hold the wide dependency for the
web application. */

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	//the value of flag will be stored in addr variable at runtime
	addr := flag.String("addr", ":4000", "HTTP network address")

	dsn := flag.String("dsn", "web:hue@/snippetbox?parseTime=true", "MySQL data source name")

	/* //need to do flag.Pars() and should always use before if not used
	//then default value will be ":4000" always */
	flag.Parse()

	/* creating infolog for information messages
	//flags are joined using bitwise or operator which is |
	*/
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	/* Creating a logger for writting error messages
	using log.Lshortfile flag to include the relevant file name and line number
	*/
	errorLog := log.New(os.Stderr, "ERROR\t", log.
		Ldate|log.Ltime|log.Lshortfile)

	//calling openDb() function for our databae
	db, err := openDB(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}
	//we call defer db.close() so connection pool closes before it main function exits
	defer db.Close()

	templateCahe, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	/* using scs.New() to initialize a new session manager.
	   Then we configure it to use our MySQl database as the session store nd set a lifetime of 12 hrs
	*/
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCahe,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
	/* mux := http.NewServeMux()
	// creating a file server it will help to get realive directory from path
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//for matching path we stric the "/static" prefix before the reques reaches the file server
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/vew", app.snippetView)
	mux.HandleFunc(
		"/snippet/create",
		app.snippetCreate,)*/

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	/* initialize a new http server struct. we set the addr and handler fields
	so that the same networkd add. and routes as before and set the error log field so now it uses custom errorLog loggin in
	*/
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
