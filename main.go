package main

import (
	"net/http"

	"github.com/renriquez/ticketing_system/models"
	"github.com/uadmin/uadmin"
)

func main() {
	DBConfig()
	RegisterModels()
	RegisterInlines()
	Server()
	image()

}

func RegisterModels() {
	uadmin.Register(
		models.Ticket{},
		models.CloseTicket{},
		models.OpenTicket{},
		models.Employee{},
	)

}

func RegisterInlines() {
	uadmin.RegisterInlines(
		uadmin.User{},
		map[string]string{
			"Employee": "UserID",
		},
	)
}

func image() {
	http.Handle("/media/images/ticket/", http.StripPrefix("/media/images/ticket/", http.FileServer(http.Dir("./media/images/ticket/"))))
}

func DBConfig() {
	uadmin.Database = &uadmin.DBSettings{
		Type:     "mysql",
		Name:     "ticketing_system",
		User:     "root",
		Password: "Allen is Great 200%",
		Host:     "localhost",
		Port:     3306,
	}
}

func Server() {
	uadmin.RootURL = "/"
	uadmin.SiteName = "Ticketing System"
	uadmin.Port = 3233
	uadmin.StartServer()
}

// func RegisterSetting() {
// 	uadmin.EmailFrom = "renriquez@integranet.ph"
// 	uadmin.EmailUsername = "renriquez@integranet.ph"
// 	uadmin.EmailPassword = "100320"
// 	uadmin.EmailSMTPServer = "smtp.integranet.ph"
// 	uadmin.EmailSMTPServerPort = 587
// }
