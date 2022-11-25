package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

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
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	checkError(err)
	pingErr := db.Ping()
	checkError(pingErr)
	fmt.Println("Connected!")

	//feed, err := allFeed()
	//checkError(err)
	//fmt.Printf("Posts found in descending date format\n %v\n", feed)
	allFeed()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// allFeed displays all posts in the database sorted by date in descending order.  This shows the user the most recent
// posts first.
func allFeed() {

	// A post slice to hold data from returned rows.
	var posts []post

	rows, err := db.Query("SELECT * FROM post ORDER BY tstamp DESC")
	checkError(err)

	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
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
		row := db.QueryRow("SELECT displayname FROM SocUser WHERE userid = ?", v.Userid)
		if err := row.Scan(&u); err != nil {
			if err == sql.ErrNoRows {
				fmt.Errorf("%d: no user", v.Userid)
			}
			fmt.Errorf("user %d: %v", v.Userid, err)
		}
		fmt.Println(u)
		fmt.Println(v.Tstamp, " ", v.Txt)
	}
	//return posts, err
}
