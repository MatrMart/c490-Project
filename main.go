package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

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
	PicURL      string
	Posts       []post
}

type feedPost struct {
	FriendPic  string
	FriendName string
	Tstamp     string
	Txt        string
}
type uFeed struct {
	DisplayName string
	Fname       string
	Lname       string
	PicURL      string
	Posts       []feedPost
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
	mux.Handle("/profileView", http.HandlerFunc(profileView))
	mux.Handle("/userFeed", http.HandlerFunc(userFeed))
	mux.Handle("/allFeed", http.HandlerFunc(allFeed))
	mux.Handle("/newPost", http.HandlerFunc(newPost))
	mux.Handle("/addPost", http.HandlerFunc(addPost))
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	//session, _ := store.Get(r, "Logged-in")
	tpl.ExecuteTemplate(w, "login.html", "")

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("check1")
	var u int
	var p string
	session, _ := store.Get(r, "Logged-in")
	em := r.FormValue("email")
	pwrd := r.FormValue("password")
	fmt.Println(em, pwrd)

	row := db.QueryRow("SELECT userid FROM SocUser WHERE email = ?", em)
	err := row.Scan(&u)
	if err != nil {
		fmt.Println("did not find user in database")
		authFail(w, r)
	}
	fmt.Println(u, pwrd)

	row = db.QueryRow("SELECT pword FROM SecurePass WHERE userid = ?", u)
	err = row.Scan(&p)
	fmt.Println(u, p)
	if err != nil {
		fmt.Println("did not find user in database")
		authFail(w, r)
	}

	fmt.Println(p, pwrd)
	//if err != nil {
	//tpl.ExecuteTemplate(w, "login.html", "Wrong email or password")
	checkError(err)
	//}

	if p == pwrd {

		session.Values["authenticated"] = true
		session.Values["userID"] = u
		session.Save(r, w)
		fmt.Println(session.Values["authenticated"])
	} else {
		//msg := "Authentication failed.  Please try again"
		authFail(w, r)
		//tpl.ExecuteTemplate(w, "login.html", msg)
	}

	//fmt.Fprintln(w, "Profile view \n\n")
	//profileView(w, r)
	//fmt.Fprintln(w, "\n\nUserFeed\n\n")
	userFeed(w, r)
	//fmt.Fprintln(w, "\n\nAll Feed\n\n")
	//allFeed(w, r)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "Logged-in")
	session.Values["authenticated"] = false
	delete(session.Values, "userID")
	session.Save(r, w)

	fmt.Println(session.Values["authenticated"])
	login(w, r)
}

func authFail(w http.ResponseWriter, r *http.Request) {
	fail := template.Must(template.New("failed").Parse("templates/login.html"))
	session, _ := store.Get(r, "Logged-in")
	session.Values["authenticated"] = false
	session.Values["userID"] = ""
	session.Save(r, w)
	fmt.Println("did it fail here")
	fail.Execute(w, "Authentication failed. Please try again.")
	fmt.Println("or here")
}

func loggedIn(w http.ResponseWriter, r *http.Request) bool {
	session, _ := store.Get(r, "Logged-in")
	ok := session.Values["authenticated"].(bool)
	return ok
}

func profileView(w http.ResponseWriter, r *http.Request) {
	var u user

	session, _ := store.Get(r, "Logged-in")
	usr := session.Values["userID"].(int)

	fmt.Println(usr)
	u.Userid = usr

	row := db.QueryRow("SELECT * FROM SocUser WHERE userid = ?", usr)
	err := row.Scan(&u.Fname, &u.Lname, &u.DisplayName, &u.Email, &u.PicURL)
	checkError(err)

	rows, err := db.Query("SELECT * FROM post WHERE userid= ? ORDER BY tstamp DESC", u.Userid)
	checkError(err)
	for rows.Next() {
		var p post
		err := rows.Scan(&p.Postid, &p.Userid, &p.Tstamp, &p.Txt)
		checkError(err)
		u.Posts = append(u.Posts, p)

	}

	fmt.Println(u)
	fmt.Fprintln(w, u.DisplayName)
	for _, v := range u.Posts {
		fmt.Fprintln(w, v.Tstamp, "\n", v.Txt, "\n")
	}

}

func userFeed(w http.ResponseWriter, r *http.Request) {
	if !(loggedIn(w, r)) {
		login(w, r)
	}

	var (
		u   user
		qry string
		uf  uFeed
	)

	session, _ := store.Get(r, "Logged-in")
	usr := session.Values["userID"].(int)

	row := db.QueryRow("SELECT * FROM SocUser WHERE userid = ?", usr)
	err := row.Scan(&u.Userid, &u.Fname, &u.Lname, &u.DisplayName, &u.Email, &u.PicURL)
	checkError(err)
	fmt.Println(u)
	uf.DisplayName = u.DisplayName
	uf.PicURL = u.PicURL

	fRows, err := db.Query("SELECT friendid FROM FriendList WHERE userid = ?", usr)
	checkError(err)
	defer fRows.Close()

	for fRows.Next() {
		var f int
		err := fRows.Scan(&f)
		checkError(err)
		fmt.Println("friend ID = ", f)
		qry = qry + fmt.Sprintf(" userid = %v OR", f)

	}

	qry = qry[:len(qry)-2]
	stmnt := "SELECT * FROM post WHERE" + qry + " ORDER BY tstamp DESC"
	fmt.Println(stmnt)

	pRows, err := db.Query(stmnt)
	checkError(err)
	defer pRows.Close()
	for pRows.Next() {
		var p post
		err := pRows.Scan(&p.Postid, &p.Userid, &p.Tstamp, &p.Txt)
		checkError(err)
		u.Posts = append(u.Posts, p)
	}

	for _, v := range u.Posts {
		var u, p string
		var fd feedPost
		row := db.QueryRow("SELECT displayname,picurl FROM SocUser where userid = ?", v.Userid)
		err := row.Scan(&u, &p)
		checkError(err)

		fd.FriendName = u
		fd.FriendPic = p
		fd.Tstamp = v.Tstamp
		fd.Txt = v.Txt
		uf.Posts = append(uf.Posts, fd)

		//fmt.Fprintln(w, fd.FriendName, "  ", fd.FriendPic, "\n")
		//fmt.Fprintln(w, v.Tstamp, " ", v.Txt, "\n")

	}

	tpl.ExecuteTemplate(w, "profile.html", uf)

}

// allFeed displays all posts in the database sorted by date in descending order.  This shows the user the most recent
// posts first.
func allFeed(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "Logged-in")
	ok := session.Values["authenticated"].(bool)

	fmt.Println(ok)

	if !ok {
		login(w, r)
	}

	// A post slice to hold data from returned rows.
	var posts []post

	rows, err := db.Query("SELECT * FROM post ORDER BY tstamp DESC")
	checkError(err)

	defer rows.Close()
	// Loop through rows, using Scan to assign Post data to p struct fields.
	for rows.Next() {
		var p post
		err := rows.Scan(&p.Postid, &p.Userid, &p.Tstamp, &p.Txt)
		checkError(err)
		posts = append(posts, p)
	}

	for _, v := range posts {

		var u string
		//get the Display name from our SocUser table based on the owner of the post being displayed
		row := db.QueryRow("SELECT displayname FROM SocUser WHERE userid = ?", v.Userid)
		err := row.Scan(&u)
		checkError(err)
		fmt.Fprintln(w, u)
		fmt.Fprintln(w, v.Tstamp, " ", v.Txt, "\n")
	}

}

func newPost(w http.ResponseWriter, r *http.Request) {
	if !(loggedIn(w, r)) {
		login(w, r)
	}
	tpl.ExecuteTemplate(w, "newPost.html", nil)
}

func addPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("made it this far")
	session, _ := store.Get(r, "Logged-in")
	usr := session.Values["userID"].(int)
	newtxt := r.FormValue("posttxt")
	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")
	fmt.Println(newtxt)
	fmt.Println(ts)

	_, err := db.Exec("INSERT INTO post (userid, tstamp, txt) VALUES (?,?,?)", usr, ts, newtxt)
	checkError(err)

	userFeed(w, r)
}
