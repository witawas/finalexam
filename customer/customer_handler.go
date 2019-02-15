package customer

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/witawas/finalexam/database"
)

func validateAuthorized(c *gin.Context) {
	authKey := c.GetHeader("Authorization")
	//fmt.Println("authKey===", authKey)
	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}
	c.Next()

}

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(validateAuthorized)
	r.POST("/customers", createCustomerHandler)
	r.GET("/customers", getAllCustomerHandler)
	r.GET("/customers/:id", getCustomerByIDHandler)

	r.PUT("/customers/:id", updCustomerHandler)
	r.DELETE("/customers/:id", delCustomerHandler)
	return r
}

func createCustomerHandler(c *gin.Context) {
	var cus Customer
	if err := c.ShouldBindJSON(&cus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// err := database.InsertTInsertCustomero(cus.Name, cus.Email, cus.Status)
	// var id int
	// err := row.Scan(&id)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	// 	return
	// }

	// stmt, err := database.GetCustomer(id)
	// if err != nil {
	// 	log.Fatal("Fail prepare stmt : ", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	// 	return
	// }

	// row = stmt.QueryRow(id)
	// if err = row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status); err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"message": "no data found"})
	// 	return
	// }

	id, err := database.InsertCustomer(cus.Name, cus.Email, cus.Status)
	if err.Message != "" {
		c.JSON(http.StatusInternalServerError, err.Message)
		return
	}

	row, err := database.GetCustomerByID(id)
	if err.Message != "" {
		c.JSON(http.StatusInternalServerError, err.Message)
		return
	}
	if errR := row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status); errR != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no data found"})
		return
	}

	c.JSON(http.StatusCreated, cus)
}

func getAllCustomerHandler(c *gin.Context) {
	results := []Customer{}
	// stmt, err := database.GetCustomer(0)
	// if err != nil {
	// 	log.Fatal("fail prepare stmt :", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	// }

	// rows, err := stmt.Query()
	// if err != nil {
	// 	log.Fatal("fail query :", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	// }
	rows, myErr := database.GetAllCustomer()
	if myErr.Message != "" {
		c.JSON(http.StatusInternalServerError, myErr.Message)
		return
	}
	for rows.Next() {
		t := Customer{}
		if err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Status); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "no data found"})
			return
		}
		results = append(results, t)
	}

	c.JSON(http.StatusOK, results)
}

func getCustomerByIDHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	cus := Customer{}
	row, myErr := database.GetCustomerByID(id)
	if myErr.Message != "" {
		c.JSON(http.StatusInternalServerError, myErr.Message)
		return
	}
	if errR := row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status); errR != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no data found"})
		return
	}

	c.JSON(http.StatusOK, cus)
}

func updCustomerHandler(c *gin.Context) {
	var cus Customer
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&cus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	myErr := database.UpdateCustomer(id, cus.Name, cus.Email, cus.Status)
	if myErr.Message != "" {
		c.JSON(http.StatusInternalServerError, myErr.Message)
		return
	}
	// stmt, err := database.UpdateCustomer()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	// }

	// if _, err := stmt.Exec(id, cus.Name, cus.Name, cus.Email, cus.Email, cus.Status, cus.Status); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	// }

	// stmt, err = database.GetCustomer(id)
	// if err != nil {
	// 	log.Fatal("Fail prepare stmt : ", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	// 	return
	// }

	// row := stmt.QueryRow(id)
	// if err = row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status); err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"message": "no data found"})
	// 	return
	// }

	row, myErr := database.GetCustomerByID(id)
	if myErr.Message != "" {
		c.JSON(http.StatusInternalServerError, myErr.Message)
		return
	}
	if errR := row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status); errR != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no data found"})
		return
	}

	c.JSON(http.StatusOK, cus)
}

func delCustomerHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = database.DeleteCustomerByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}

func CreateTable() {

	_, err := database.CreateTable()
	if err != nil {
		log.Fatal("can't create table", err)
	}
}
