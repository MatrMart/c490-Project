package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"

	"github.com/go-sql-driver/mysql"
)

var tpl *template.Template
var db *sql.DB
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type post struct {
	Postid int
	Userid int
	Tstamp string
	Txt    string
}

type user struct {
	Userid      int
	Pword       string
	Fname       string
	Lname       string
	DisplayName string
	Email       string
}

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:                 "DBUser", //os.Getenv("DBUSER"),
		Passwd:               "DBPass", //os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "SocMed",
		AllowNativePasswords: true,
	}
	mux := http.NewServeMux()
	tpl, _ = template.ParseGlob("templates/*.html")

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	checkError(err)
	pingErr := db.Ping()
	checkError(pingErr)
	fmt.Println("Connected!")

	mux.Handle("/login", http.HandlerFunc(login))
	mux.Handle("/loginhandler", http.HandlerFunc(loginHandler))
	mux.Handle("/logout", http.HandlerFunc(logout))
	mux.Handle("/allFeed", http.HandlerFunc(allFeed))
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	//session, _ := store.Get(r, "Logged-in")
	tpl.ExecuteTemplate(w, "login.html", nil)

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("check1")
	var u, p string
	session, _ := store.Get(r, "Logged-in")
	usr := r.FormValue("email")
	pwrd := r.FormValue("password")
	fmt.Println(usr, pwrd)
	row := db.QueryRow("SELECT email,pword FROM SocUser WHERE email = ?", usr)

	err := row.Scan(&u, &p)
	fmt.Println(u, p)
	if err != nil {
		tpl.ExecuteTemplate(w, "login.html", "Wrong email or password")
	}

	if p == pwrd {

		session.Values["authenticated"] = true
		session.Values["userID"] = usr
		session.Save(r, w)
		fmt.Println(session.Values["authenticated"])
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "Logged-in")
	//row := db.QueryRow("SELECT displayname FROM SocUser WHERE userid = ?", v.Userid)
	session.Values["authenticated"] = false
	delete(session.Values, "userID")
	session.Save(r, w)

	fmt.Println(session.Values["authenticated"])
}

// allFeed displays all posts in the database sorted by date in descending order.  This shows the user the most recent
// posts first.
func allFeed(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "Logged-in")

	ok := session.Values["authenticated"].(bool)

	fmt.Println(ok)

	if !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// A post slice to hold data from returned rows.
	var posts []post

	rows, err := db.Query("SELECT * FROM post ORDER BY tstamp DESC")
	checkError(err)

	defer rows.Close()
	// Loop through rows, using Scan to assign Post data to p struct fields.
	for rows.Next() {
		var p post
		if err := rows.Scan(&p.Postid, &p.Userid, &p.Tstamp, &p.Txt); err != nil {
			//return nil, fmt.Errorf("posts %v", err)
			fmt.Errorf("posts %v", err)
		}
		posts = append(posts, p)
	}

	for _, v := range posts {

		var u string
		//get the Display name from our SocUser table based on the owner of the post being displayed
		row := db.QueryRow("SELECT displayname FROM SocUser WHERE userid = ?", v.Userid)
		if err := row.Scan(&u); err != nil {
			if err == sql.ErrNoRows {
				fmt.Errorf("%d: no user", v.Userid)
			}
			fmt.Errorf("user %d: %v", v.Userid, err)
		}
		fmt.Fprintln(w, u)
		fmt.Fprintln(w, v.Tstamp, " ", v.Txt, "\n")
	}
	//return posts, err
}
