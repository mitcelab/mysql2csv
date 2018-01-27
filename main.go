package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"os"
)

func main() {
	query := flag.String("query", "", "the query to run")
	output := flag.String("output", "", "the output file (result.csv)")
	user := flag.String("user", "", "MySQL username")
	pass := flag.String("pass", "", "MySQL password")
	host := flag.String("host", "127.0.0.1", "MySQL host")
	port := flag.Int("port", 3306, "MySQL port")
	dbname := flag.String("dbname", "", "database name")

	ds := DataSource{
		user:   user,
		pass:   pass,
		host:   host,
		port:   port,
		dbname: dbname,
	}
	flag.Parse()
	if err := ds.validate(); err != nil {
		log.Fatal(err)
	}

	db := Connect(ds)

	var writer io.Writer
	if *output != "" {
		fout, err := os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
		writer = bufio.NewWriter(fout)
	} else {
		writer = bufio.NewWriter(os.Stdout)
	}
	QueryToCSV(db, *query, writer)
}

func Connect(ds DataSource) *sql.DB {
	db, err := sql.Open("mysql", ds.dsn())
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func QueryToCSV(db *sql.DB, query string, fout io.Writer) {
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	header, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(fout)
	defer writer.Flush()
	writer.Write(header)

	counter := 0
	for rows.Next() {
		columns := make([]sql.NullString, len(header))
		columnPointers := make([]interface{}, len(header))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			log.Fatal(err)
		}

		values := make([]string, len(header))
		for i, column := range columns {
			values[i] = column.String
		}
		writer.Write(values)

		counter += 1
		if counter%10000 == 0 {
			log.Printf("Wrote %v rows...", counter)
		}
	}
}

type DataSource struct {
	user   *string
	pass   *string
	host   *string
	port   *int
	dbname *string
}

func (ds *DataSource) dsn() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", *ds.user, *ds.pass, *ds.host, *ds.port, *ds.dbname)
}

func (ds *DataSource) validate() error {
	if *ds.port < 0 || *ds.port > 65535 {
		return errors.New("port should be between 0 and 65535")
	}
	if *ds.user == "" || *ds.pass == "" {
		return errors.New("user/pass must be non-empty")
	}
	if *ds.dbname == "" {
		return errors.New("dbname must be non-empty")
	}
	return nil
}
