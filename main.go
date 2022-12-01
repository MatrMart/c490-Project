package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	//"strconv"

	"github.com/gorilla/sessions"

	"github.com/go-sql-driver/mysql"
)

var tpl *template.Template

// var logFail *template.Template
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
	Userid      int
	DisplayName string
	Message     string
	Message2    string
	CanPost     string
	Fname       string
	Lname       string
	PicURL      string
	Posts       []feedPost
}

type msg struct {
	txt string
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

	//connect to the database
	db, err = sql.Open("mysql", cfg.FormatDSN())
	checkError(err)
	pingErr := db.Ping()
	checkError(pingErr)
	fmt.Println("Connected!") //output to console for demonstration purposes

	mux.Handle("/login", http.HandlerFunc(login))
	mux.Handle("/loginhandler", http.HandlerFunc(loginHandler))
	mux.Handle("/logout", http.HandlerFunc(logout))
	mux.Handle("/profileView", http.HandlerFunc(profileView))
	mux.Handle("/userFeed", http.HandlerFunc(userFeed))
	mux.Handle("/allFeed", http.HandlerFunc(allFeed))
	mux.Handle("/newPost", http.HandlerFunc(newPost))
	mux.Handle("/addPost", http.HandlerFunc(addPost))
	mux.Handle("/adjFriend", http.HandlerFunc(adjFriend))
	mux.Handle("/newPic", http.HandlerFunc(newPic))
	mux.Handle("/addNewPic", http.HandlerFunc(addNewPic))

	log.Fatal(http.ListenAndServe(":8000", mux))
}

// just a normal everyday err checker
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// calls the login template which then redirects to loginHandler once the email and password forms have been completed
func login(w http.ResponseWriter, r *http.Request) {
	var m msg
	m.txt = ""
	tpl.ExecuteTemplate(w, "login.html", m)

}

// loginHandler uses the information retrieved from the login template to query a relational database
// looking for a tuple containing the user's info.  Two tables are used.  One is SocUser which stores most
// of the user's information, and SecurePass which stores just a password and a userid.  This was done
// to make sure the only time we access a user's password can be limited as other handlers will also
// access the SocUser table.
func loginHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("loginHandler") //output to console for demonstration purposes
	var u int
	var p string
	session, _ := store.Get(r, "Logged-in")
	em := r.FormValue("email")
	pwrd := r.FormValue("password")
	fmt.Println(em, pwrd) //output to console for demonstration purposes

	//grab just the userid from the SocUser table.
	row := db.QueryRow("SELECT userid FROM SocUser WHERE email = ?", em)
	err := row.Scan(&u)

	//row.Scan can return an error in the form of ErrNoRows.  Instead of running that through checkError and
	//potentially getting log.Fatal triggered just because the email didn't match any user, ErrNoRows is checked.
	//This way, if someone mistyped their email address or doesn't have an account, authFail is called which
	//will redirect them back to the login handler with an authentication error message.
	if err == sql.ErrNoRows {
		authFail(w, r)
		return
	}

	fmt.Println(u, pwrd) //output to console for demonstration purposes

	//using the userid from the previous query, table SecurePass is queried to get
	//that individual user's password
	row = db.QueryRow("SELECT pword FROM SecurePass WHERE userid = ?", u)
	err = row.Scan(&p)
	fmt.Println(u, p) //output to console for demonstration purposes.

	//Again checking to ensure that there are results otherwise call authFail
	if err == sql.ErrNoRows {
		authFail(w, r)
		return
	}

	fmt.Println(p, pwrd) //output to console for demonstration purposes.

	//if the password provided matches the account's password,  the session
	//is saved as being logged in. Then they are redirected to their user feed.
	//otherwise, authentication failed so call authFail.
	if p == pwrd {
		session.Values["authenticated"] = true
		session.Values["userID"] = u
		session.Save(r, w)
		fmt.Println(session.Values["authenticated"]) //output to console for demonstration purposes
		userFeed(w, r)
	} else {
		authFail(w, r)
	}
}

// When the user has clicked the logout button, their session is ended and they
// are sent back to the login screen
func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "Logged-in")
	session.Values["authenticated"] = false
	delete(session.Values, "userID")
	session.Save(r, w)

	fmt.Println(session.Values["authenticated"]) //output to console for demonstration purposes
	login(w, r)
}

// since Authentication can fail at several steps in the login process,
// authFail was created to prevent repeated code.  First a message to pass
// to the template is created, then session information is reset to show
// not logged in, then call login again with the attached message.
func authFail(w http.ResponseWriter, r *http.Request) {
	msg := make(map[string]string)
	msg["txt"] = "Authentication failed.  Please try again."

	session, _ := store.Get(r, "Logged-in")
	session.Values["authenticated"] = false
	session.Values["userID"] = ""
	session.Save(r, w)
	tpl.ExecuteTemplate(w, "login.html", msg)

}

//To keep unauthorized users from going straight to each handler without logging
//in first, loggedIn is called to check if a session is valid before executing
//handlers where appropriate.

func loggedIn(w http.ResponseWriter, r *http.Request) bool {
	session, _ := store.Get(r, "Logged-in")
	ok := session.Values["authenticated"].(bool)
	return ok
}

// Profile view is specific to viewing another user's profile.  A struct specific to
// this handler was added for formatting when checking to see if the user and the
// person they are viewing are friends.  If they are not, they are not allowed
// to view the contents of the target profile.  They are given an option to
// add the user as a friend.  If they are already friends, the profile
// is displayed and also an option to remove them as a friend is presented.

func profileView(w http.ResponseWriter, r *http.Request) {
	if !(loggedIn(w, r)) {
		login(w, r)
		return
	}

	//this struct is used to send specific information to the addFriend template
	type friends struct {
		userId      int
		friendId    int
		DisplayName string
		PicURL      string
		Message     string
	}
	var (
		uf         uFeed
		friendtest friends
	)

	session, _ := store.Get(r, "Logged-in")
	usr := session.Values["userID"].(int)

	//this is the information entered into the search box
	srch := r.FormValue("search")

	fmt.Println(srch)

	// this query checks the search term agains fname, lname, and email trying to find a match.
	// if no match was made, then the user is dropped back to the userFeed.  If a match was made,
	// a user struct is populated with information to be formatted to send to one of several other
	// templates further in the handler

	row := db.QueryRow("SELECT userid, fname, lname, displayname, picurl FROM SocUser WHERE fname = ? OR lname = ? OR email = ? ", srch, srch, srch)
	err := row.Scan(&uf.Userid, &uf.Fname, &uf.Lname, &uf.DisplayName, &uf.PicURL)
	if err == sql.ErrNoRows {
		fmt.Println("no user found")
		userFeed(w, r)
		return
	}

	//format messages for the profile template.  Message 2 is link directions for passing values through the template back into r for
	//the handlers that add or remove friends.
	uf.Message = " has been saying:"
	if uf.Userid != usr {
		uf.Message2 = fmt.Sprintf("/adjFriend?action=unfriend&friendid=%v", uf.Userid)
	}

	//if there are no rows, the current user is not friends with the searched user.  friendtest is populated with appropriate
	//information to send to the addFriend template.
	row = db.QueryRow("SELECT * FROM FriendList WHERE userid = ? AND friendid =?", usr, uf.Userid)
	err = row.Scan(&friendtest.userId, &friendtest.friendId)

	if err == sql.ErrNoRows {

		friendtest.DisplayName = uf.DisplayName
		friendtest.PicURL = uf.PicURL
		friendtest.friendId = uf.Userid
		friendtest.Message = fmt.Sprintf("/adjFriend?action=friend&friendid=%v", uf.Userid)
		tpl.ExecuteTemplate(w, "addFriend.html", friendtest)
		return
	}

	//since there was a row from checking for a friend, the rest of the ufeed struct is populated to display the friend's
	// profile information, recent posts, and an option to unfriend is presented.
	rows, err := db.Query("SELECT tstamp, txt FROM post WHERE userid= ? ORDER BY tstamp DESC", uf.Userid)
	checkError(err)
	defer rows.Close()

	//the rows are scanned to add information to the posts we'll display on the feed in the profile template. first fd is populated
	//as a feedPost, then it's appended to the slice of feedPosts in the userFeed struct, uf.
	for rows.Next() {
		var fd feedPost
		err := rows.Scan(&fd.Tstamp, &fd.Txt)
		checkError(err)
		fd.FriendName = uf.DisplayName
		fd.FriendPic = uf.PicURL
		uf.Posts = append(uf.Posts, fd)
	}

	fmt.Println(uf) //output to console for demonstration purposes

	tpl.ExecuteTemplate(w, "profile.html", uf)

}

// When a user first logs in, they are sent to userFeed.  userFeed grabs their pertinent
// information from their tuple in the SocUser table.  It also builds a copy of information/posts
// to display in their feed.  Posts in their feed will be laid out in descending order by date
// from all of their friends posts.
func userFeed(w http.ResponseWriter, r *http.Request) {
	//force login if not logged in
	if !(loggedIn(w, r)) {
		login(w, r)
		return
	}

	var (
		u   user
		qry string
		uf  uFeed
	)

	session, _ := store.Get(r, "Logged-in")
	usr := session.Values["userID"].(int)
	//grab user info
	row := db.QueryRow("SELECT * FROM SocUser WHERE userid = ?", usr)
	err := row.Scan(&u.Userid, &u.Fname, &u.Lname, &u.DisplayName, &u.Email, &u.PicURL)

	//at this point, if there are no rows, something has gone wrong with the database
	//log.Fatal will be a valid option
	checkError(err)
	fmt.Println(u) //output to console for demonstration purposes

	//uf is the user feed struct.  This is the information used in the profile template
	//it will include the users displayname, picture, and slice of friends posts that contain the
	//friends picture, display name, text, and timestamp of the posts in descending order by time.
	uf.DisplayName = u.DisplayName
	uf.PicURL = u.PicURL
	uf.Message = "'s friends are saying:"
	uf.CanPost = "yes"

	//grab the user's friend list
	fRows, err := db.Query("SELECT friendid FROM FriendList WHERE userid = ?", usr)
	checkError(err)
	defer fRows.Close()

	//to query for posts from an unknown number of friends in the list, the query is built
	//by piecing together a string in the format " userid = 1 OR userid = 2..." etc.
	//The loop below continues to add a "userid = ?" until all friends userid's are in
	//the string.
	for fRows.Next() {
		var f int
		err := fRows.Scan(&f)
		checkError(err)
		fmt.Println("friend ID = ", f) //output to console for demonstration purposes
		qry = qry + fmt.Sprintf(" userid = %v OR", f)

	}

	//there will be an extra "OR" at the end of the qry statement at the end.
	//this statement trims it off.
	qry = qry[:len(qry)-2]
	//the full query is now built into stmnt to be passed to db.Query
	stmnt := "SELECT * FROM post WHERE" + qry + " ORDER BY tstamp DESC"
	fmt.Println(stmnt) //output to console for demonstration purposes

	//grab all posts from the friend list and sort descending by  time
	pRows, err := db.Query(stmnt)
	checkError(err)
	defer pRows.Close()

	//store the posts in our user struct to be used later
	for pRows.Next() {
		var p post
		err := pRows.Scan(&p.Postid, &p.Userid, &p.Tstamp, &p.Txt)
		checkError(err)
		u.Posts = append(u.Posts, p)
	}

	//this takes the relevent information from each of the stored user posts, and
	//grabs the display name and picture url from each freind that owned the post.
	//it then stores it in the uFeed struct to be passed on to the profile template
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
	}

	tpl.ExecuteTemplate(w, "profile.html", uf)

}

// allFeed displays all posts in the database sorted by date in descending order.  This shows the user the most recent
// posts first. This handler is for demonstration purposes
func allFeed(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "Logged-in")
	ok := session.Values["authenticated"].(bool)

	fmt.Println(ok)

	if !ok {
		login(w, r)
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

// newPost sends user info the to newPost template.  From there the form info is sent to addPost which adds a new entry to the post table
// in the database
func newPost(w http.ResponseWriter, r *http.Request) {
	if !(loggedIn(w, r)) {
		login(w, r)
		return
	}
	var u user
	session, _ := store.Get(r, "Logged-in")
	usr := session.Values["userID"].(int)

	row := db.QueryRow("SELECT * FROM SocUser WHERE userid = ?", usr)
	err := row.Scan(&u.Userid, &u.Fname, &u.Lname, &u.DisplayName, &u.Email, &u.PicURL)
	checkError(err)

	tpl.ExecuteTemplate(w, "newPost.html", u)
}

// addPost inserts a new post into the post table in the database using info from the user session, the current time,
// and the text provided by the user.  The postid number is automatically generated by the database
func addPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("adding post now") //output to console for demonstration purposes
	session, _ := store.Get(r, "Logged-in")
	usr := session.Values["userID"].(int)
	newtxt := r.FormValue("posttxt")

	//grab the current time and format it for the sql datatype timestamp
	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")
	fmt.Println(newtxt) //output to console for demonstration purposes
	fmt.Println(ts)     //output to console for demonstration purposes

	//a new post is inserted into the table using the
	_, err := db.Exec("INSERT INTO post (userid, tstamp, txt) VALUES (?,?,?)", usr, ts, newtxt)
	checkError(err)
	fmt.Println("post added") //output to console for demonstration purposes
	userFeed(w, r)
}

// adjFriend adds or removes a row from the FriendList table.  When adjFriend is called, the parameters from
// the url instruct it to either add or remove.  Then the user is sent back to their feed.
func adjFriend(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "Logged-in")
	usr := session.Values["userID"].(int)
	valuesMap, err := url.ParseQuery(r.URL.RawQuery)
	checkError(err)
	action := (valuesMap["action"][0])
	trg := (valuesMap["friendid"][0])

	fmt.Println(action, trg)
	if action == "unfriend" {
		stmnt := fmt.Sprintf("DELETE FROM FriendList WHERE userid = %v AND friendid = %s;", usr, trg)
		fmt.Println(stmnt) //output to console for demonstration purposes
		_, err := db.Exec(stmnt)
		checkError(err)

	} else if action == "friend" {
		stmnt := fmt.Sprintf("INSERT INTO FriendList (userid, friendid) VALUES( %v, %s);", usr, trg)
		fmt.Println(stmnt) //output to console for demonstration purposes
		_, err := db.Exec(stmnt)
		checkError(err)

	}

	userFeed(w, r)
}

// newPic first gathers user information from our table then sends it to the NewPic template.
func newPic(w http.ResponseWriter, r *http.Request) {
	if !(loggedIn(w, r)) {
		login(w, r)
		return
	}
	var u user
	session, _ := store.Get(r, "Logged-in")
	usr := session.Values["userID"].(int)

	row := db.QueryRow("SELECT * FROM SocUser WHERE userid = ?", usr)
	err := row.Scan(&u.Userid, &u.Fname, &u.Lname, &u.DisplayName, &u.Email, &u.PicURL)
	checkError(err)

	tpl.ExecuteTemplate(w, "newPic.html", u)
}

// once the newPic template form has been completed, an appropriate sql statement is created from
// the currently logged in user's info, and the picURL string received from the newPic template
func addNewPic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("updating new pic url") //output to console for demonstration purposes
	session, _ := store.Get(r, "Logged-in")
	usr := session.Values["userID"].(int)
	newURL := r.FormValue("picURL")

	stmnt := fmt.Sprintf("UPDATE SocUser SET picurl = %s WHERE userid = '%v'", newURL, usr)

	_, err := db.Exec(stmnt)
	checkError(err)

	userFeed(w, r)

	fmt.Println(stmnt) //output to console for demonstration purposes

}
