package models

import (
	"time"

	"github.com/uadmin/uadmin"
)

type CloseTicket struct {
	uadmin.Model
	CreatedAt    time.Time   `uadmin:"read_only"`
	TicketNumber string      `uadmin:"read_only"`
	User         uadmin.User `uadmin:"read_only"`
	UserID       uint
	Department   uadmin.UserGroup
	DepartmentID uint
	Subject      string `uadmin:"read_only"`
	Issue        string `uadmin:"html;read_only" sql:"type:LONGTEXT"`
	Attachment   string `uadmin:"image"`
	Remarks      string `uadmin:"html;read_only" sql:"type:LONGTEXT"`
	Solution     string `uadmin:"file"`
	ClosedBy     uadmin.User
	ClosedByID   uint
	Status       bool `uadmin:"read_only"`
}

func (c *CloseTicket) String() string {
	uadmin.Preload(c)
	return c.TicketNumber
}
