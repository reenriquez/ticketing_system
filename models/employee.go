package models

import "github.com/uadmin/uadmin"

type Employee struct {
	uadmin.Model
	User          uadmin.User
	UserID        uint
	EmailPassword string `uadmin:"password;list_exclude"`
	ContactNumber string `uadmin:"pattern:^[0-9]*$;pattern_msg:Your input must be a number."`
}

func (e *Employee) String() string {
	uadmin.Preload(e)
	return e.User.Username
}

func (e *Employee) GetActiveSession() {
	active := uadmin.Session{}
	uadmin.Get(&active, "active = ? ", true)
	uadmin.Trail(uadmin.DEBUG, "ACTIVE: %v", active)
	e.UserID = active.UserID
}
