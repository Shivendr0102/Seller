package database

import (
	"fmt"
	"time"
)

func (db Database) AddURL(url string, encoded_url string) (string, error) {

	// SQL QUERIES used in this function
	const (
		getallCountqueryStatement = `SELECT COUNT(*) FROM [URL_Table];`
		InsertqueryStatement      = `INSERT INTO [URL_Table] (ID, Url, ShortenedURL, dateTime ) VALUES (@ID, @URL, @SHort_URL, @DateTime);`
	)

	// var storeURLshortener model.URLMap
	var cnt int
	cntStmt, err := db.SqlDb.Prepare(getallCountqueryStatement)
	if err != nil {
		return "", nil
	}
	err = cntStmt.QueryRow().Scan(&cnt)
	if err != nil {
		return "Unsuccessfull", nil
	}

	stmt, err := db.SqlDb.Prepare(InsertqueryStatement)
	if err != nil {
		return "Unsuccessfull", nil
	}
	defer stmt.Close()

	// InCase the Base62 Encode Url comes similar to some old one , Then to make it permanent UNIQUE
	// we will be incrementing the string size with counter ID Digit ( which will always be unique)
	unique_counter := fmt.Sprintf("%d", cnt+1)
	stmt.Exec(cnt+1, url, encoded_url+"_"+unique_counter, time.Now())

	return "Successfull", nil

}
