package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type MyError struct {
	Message string `json:"message`
	Code    int    `json:"code"`
}

func (e *MyError) Error() string {
	return e.Message // code
}

var db *sql.DB

func Conn() *sql.DB {
	if db != nil {
		return db
	}
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Fail connect db")
	}
	//defer db.Close()
	return db
}

func InsertCustomer(name, email, status string) (int, MyError) {

	row := Conn().QueryRow("insert into customers (name,email,status) values ($1,$2,$3) returning id", name, email, status)
	var id int
	err := row.Scan(&id)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return 0, MyError{Message: err.Error(), Code: 111}
	}

	fmt.Println("row insert return", id)
	// stmt, err := GetCustomer(id)
	// if err != nil {
	// 	log.Fatal("Fail prepare stmt : ", err)
	// 	//c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	// 	return MyError{Message: err.Error(), Code: 222}
	// }

	// row = stmt.QueryRow(id)
	// if err = row.Scan(id, name, email, status); err != nil {
	// 	//c.JSON(http.StatusNotFound, gin.H{"message": "no data found"})
	// 	return MyError{Message: err.Error(), Code: 222}
	// }

	return id, MyError{}
}

func UpdateCustomer() (*sql.Stmt, error) {
	return Conn().Prepare(`update customers set 
						name = CASE WHEN $2 = '' THEN name ELSE $3 END 
						, email = CASE WHEN $4 = '' THEN email ELSE $5 END 
						,status = CASE WHEN $6 = '' THEN status ELSE $7 END
						where id = $1`)
}

func GetCustomer(id int) (*sql.Stmt, error) {

	if id != 0 {
		return Conn().Prepare("select id,name,email,status from customers where id = $1")
	} else {
		return Conn().Prepare("select id,name,email,status from customers")
	}
}

func GetCustomerByID(id int) (*sql.Row, MyError) {
	stmt, err := Conn().Prepare("select id,name,email,status from customers where id = $1")
	if err != nil {
		//log.Fatal("Fail prepare stmt : ", err)
		//c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		//return
		return nil, MyError{Message: err.Error(), Code: 111}
	}

	//cus := Customer{}
	row := stmt.QueryRow(id)
	//if err = row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status); err != nil {
	// c.JSON(http.StatusNotFound, gin.H{"message": "no data found"})
	// return
	//	return nil, MyError{Message: err.Error(), Code: 222}
	//}
	if row == nil {
		return nil, MyError{Message: "No data found!", Code: 111}
	}

	return row, MyError{}
}

func GetAllCustomer() (*sql.Rows, MyError) {
	stmt, err := Conn().Prepare("select id,name,email,status from customers ")
	if err != nil {
		//log.Fatal("Fail prepare stmt : ", err)
		//c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		//return
		return nil, MyError{Message: err.Error(), Code: 111}
	}

	//cus := Customer{}
	rows, err := stmt.Query()
	if err != nil {
		//log.Fatal("Fail prepare stmt : ", err)
		//c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		//return
		return nil, MyError{Message: err.Error(), Code: 222}
	}
	//if err = row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status); err != nil {
	// c.JSON(http.StatusNotFound, gin.H{"message": "no data found"})
	// return
	//	return nil, MyError{Message: err.Error(), Code: 222}
	//}
	if rows == nil {
		return nil, MyError{Message: "No data found!", Code: 111}
	}

	return rows, MyError{}
}

func CreateTable() (sql.Result, error) {
	ctb := `
				CREATE TABLE IF NOT EXISTS customers (
				id SERIAL PRIMARY KEY, name TEXT,
				email TEXT,status TEXT
				);
	`
	return Conn().Exec(ctb)

}

func DeleteCustomerByID(id int) error {
	stmt, err := Conn().Prepare("delete from customers where id = $1")
	if err != nil {
		log.Fatal("Fail prepare stm")
		//c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return err
	}

	res, err := stmt.Exec(id)
	//fmt.Println("TESTTTtttt", res.RowsAffected)
	if err != nil {
		log.Fatal("Fail delete")
		//c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Fail delete")
		//c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return err
	}
	if count != 1 {
		log.Println("No record delete or Many record delete!")
		return &MyError{"No record delete or Many record delete!", http.StatusInternalServerError}
	}

	return err
}
