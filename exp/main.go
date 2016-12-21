package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/EvanLib/me_users/models"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	ug, err := models.NewUserGorm("root:lol626465@/me_schedule?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	//user2 := ug.ByID(2)
	user := ug.ByEmail("nope@nope.com")

	//fmt.Println(user2)
	fmt.Println(user)
}

func getInfo() (name, email string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Name")
	name, _ = reader.ReadString('\n')

	fmt.Println("Email")
	email, _ = reader.ReadString('\n')
	return name, email
}
