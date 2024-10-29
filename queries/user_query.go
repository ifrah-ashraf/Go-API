package queries

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Artist struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  string `json:"sex"`
}

type UserData struct {
	id         int
	name       string
	age        int
	sex        string
	created_at time.Time
}

// to use it in the other packages name the function starting letter in uppercase

func InsertQuery(db *sql.DB, a Artist) (int64, error) {
	iquery := `INSERT INTO Artist (name , age , sex) values ($1 , $2 ,$3) RETURNING *`

	res, err := db.Exec(iquery, a.Name, a.Age, a.Sex)
	if err != nil {
		return 0, fmt.Errorf("error while inserting: %v", err)

	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error in row insertion operation: %v", err)
	}
	fmt.Print("Artist inserted successfully\n")

	return rowsAffected, nil

}

func GetUsers(db *sql.DB) chan Artist {

	artistChannel := make(chan Artist)

	go func() {
		defer close(artistChannel)

		gquery := `SELECT name , age , sex FROM  Artist `
		rows, err := db.Query(gquery)
		if err != nil {
			fmt.Printf("Error while fetching user : %v", err)
		}

		defer rows.Close()
		for rows.Next() {
			var a Artist
			err := rows.Scan(&a.Name, &a.Age, &a.Sex)
			if err != nil {
				log.Printf("Error scanning row: %v", err)
				return
			}

			artistChannel <- a

			//fmt.Println("\n", name, age, sex)
		}
	}()
	return artistChannel
}

// Don't use unnecessary channel and go routine
func DelUser(db *sql.DB, n string) (UserData, error) {
	value := n
	var u UserData

	gquery := `SELECT * FROM Artist WHERE name = $1`
	err := db.QueryRow(gquery, value).Scan(&u.id, &u.name, &u.age, &u.sex, &u.created_at)
	if err == sql.ErrNoRows {
		return UserData{}, fmt.Errorf("no user found with the specified name: %v", err)
	} else if err != nil {
		return UserData{}, fmt.Errorf("error while fetching user details: %v", err)
	}

	dquery := `DELETE FROM Artist WHERE name = $1`
	res, err := db.Exec(dquery, value)
	if err != nil {
		return UserData{}, fmt.Errorf("error while deleting user: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return UserData{}, fmt.Errorf("error while checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return UserData{}, fmt.Errorf("no user deleted")
	}

	fmt.Println("User deleted successfully.")
	return u, nil
}

//only use defer rows.close() for multiple Query (coming from query db.query("select * ...")) to close the cursor over multiple rows
