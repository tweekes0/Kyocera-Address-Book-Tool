package db

import "fmt"

type Entry struct {
	ID       int64
	Name     string
	Username string
	Email    string
}

func (e Entry) Display() {
	fmt.Printf("ID: %d\n", e.ID)
	fmt.Printf("Name: %v\n", e.Name)
	fmt.Printf("Username: %v\n", e.Username)
	fmt.Printf("Email: %v\n", e.Email)

}
