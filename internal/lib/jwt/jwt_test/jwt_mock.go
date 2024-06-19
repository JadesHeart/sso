package main

type MockUser struct {
	ID    int
	Email string
}

type MockApp struct {
	ID     int
	Secret string
}

func NewMockUser(id int, email string) MockUser {
	return MockUser{ID: id, Email: email}
}

func NewMockApp(id int, secret string) MockApp {
	return MockApp{ID: id, Secret: secret}
}
