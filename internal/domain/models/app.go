package models

type App struct {
	ID     int
	Name   string
	Secret string // подпись для jwt-токена
}
