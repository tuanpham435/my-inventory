package main

import "fmt"

func main() {
	app := new(App)
	err := app.Initialise(DbUser, DbPassword, DbName)
	if err != nil {
		fmt.Println(err)
		return
	}
	app.Run("localhost:8080")
}
