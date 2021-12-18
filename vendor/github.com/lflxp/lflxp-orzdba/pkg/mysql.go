package pkg

// https://github.com/go-sql-driver/mysql#examples
import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func MysqlConn(username, password, ip, port, dbname string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, ip, port, dbname))
	if err != nil {
		// log.Fatalf("Open database error: %s\n", err)
		return nil, err
	}
	return db, nil
	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	// log.Fatal(err)
	// 	return nil, err
	// }
	// return nil, nil

	// insert(db)
	// rows, err := db.Query("select id, username from user where id = ?", 1)

	// x := db.Stats()
	// fmt.Println(x)

	// rows, err := db.Query("show variables;")
	// rows, err := db.Query("show global variables;")
	// rows, err := db.Query("show slave status")
	// rows, err := db.Query("show global status")
	// rows, err := db.Query("show status")
	// rows, err := db.Query("show databases")
	// if err != nil {
	// 	log.Println(err)
	// }

	// fmt.Println(rows)
	// var Variable_name, Value string
	// defer rows.Close()
	// for rows.Next() {
	// 	err = rows.Scan(&Variable_name, &Value)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Println(Variable_name, Value)
	// }
}
