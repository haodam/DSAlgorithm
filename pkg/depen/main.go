package main

import "fmt"

type Database interface {
	GetUser(id int) string
}

type PostgresDB struct {
}

func (p *PostgresDB) GetUser(id int) string {
	return fmt.Sprintf("User %d", id)
}

type MockDB struct{}

func (m *MockDB) GetUser(id int) string {
	return "MockDB User "
}

type UserService struct {
	db Database
}

func NewUserService(db Database) *UserService {
	return &UserService{db: db}
}

func (u *UserService) GetUserName(id int) string {
	return u.db.GetUser(id)
}

func main() {
	// Inject real database
	realDB := &PostgresDB{}
	userService := NewUserService(realDB)
	fmt.Println(userService.GetUserName(1))

	// Inject mock database for testing
	mockDB := &MockDB{}
	userService = NewUserService(mockDB)
	fmt.Println(userService.GetUserName(2))
}

//// Database interacts with the DB
//type Database struct{}
//
//func (d *Database) GetUser(id int) string {
//	return fmt.Sprintf("User %d", id)
//}
//
//// UserService depends on Database
//
//type UserService struct {
//	db Database
//}
//
//func (us *UserService) GetUserName(id int) string {
//	return us.db.GetUser(id)
//}

//func main() {
//	service := UserService{} // Can`t replace Database easily
//	fmt.Println(service.GetUserName(1))
//
//}
