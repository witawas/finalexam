package main

import (
	"github.com/witawas/finalexam/customer"
)

func main() {
	customer.CreateTable()

	r := customer.NewRouter()

	r.Run(":2019")
}
