package models

import (
	"time"

	"github.com/uadmin/uadmin"
)

type OpenTicket struct {
	uadmin.Model
	CreatedAt    time.Time   `uadmin:"read_only"`
	TicketNumber string      `uadmin:"read_only"`
	User         uadmin.User `uadmin:"read_only;hidden"`
	UserID       uint
	CreatedBy    string           `uadmin:"read_only;hidden;list_exclude"`
	Department   uadmin.UserGroup `uadmin:"read_only"`
	DepartmentID uint
	Subject      string      `uadmin:"required"`
	Issue        string      `uadmin:"read_only;hidden" sql:"type:LONGTEXT"`
	Attachment   string      `uadmin:"image;hidden"`
	Remarks      string      `uadmin:"html;required" sql:"type:LONGTEXT"`
	Solution     string      `uadmin:"file"`
	ClosedBy     uadmin.User `uadmin:"read_only;hidden"`
	ClosedByID   uint
	Status       bool `uadmin:"required"`
}

func (o *OpenTicket) String() string {
	uadmin.Preload(o)
	return o.TicketNumber
}

func (o *OpenTicket) GetActiveSession() {
	active := uadmin.Session{}
	uadmin.Get(&active, "active = ? ", true)
	uadmin.Trail(uadmin.DEBUG, "ACTIVE: %v", active)
	o.ClosedByID = active.UserID
}

// * Saving of Added ticket
func (o *OpenTicket) Save() {
	closesum := CloseTicket{}
	uadmin.Preload(o)

	em_o := Employee{}
	us_o := uadmin.User{}

	uadmin.Trail(uadmin.DEBUG, "email %v", us_o.Email)

	if o.Status {
		o.GetActiveSession()
		uadmin.Get(&us_o, "id =?", o.ClosedByID)
		closesum.UserID = o.UserID
		closesum.DepartmentID = us_o.UserGroupID
		closesum.Subject = o.Subject
		closesum.Issue = o.Issue
		closesum.Attachment = o.Attachment
		closesum.TicketNumber = o.TicketNumber
		closesum.CreatedAt = o.CreatedAt
		closesum.Remarks = o.Remarks
		closesum.Solution = o.Solution
		closesum.ClosedByID = o.ClosedByID
		closesum.Status = o.Status
		uadmin.Save(&closesum)
		uadmin.Delete(o)

		// * Validation of Department and sending email
		uadmin.Get(&us_o, "user_group_id = ?", closesum.DepartmentID)
		if closesum.DepartmentID == us_o.UserGroupID {
			uadmin.Get(&em_o, "user_id = ?", o.ClosedByID)
			uadmin.EmailFrom = us_o.Email
			uadmin.EmailUsername = us_o.Email
			uadmin.EmailPassword = em_o.EmailPassword
			uadmin.EmailSMTPServer = "smtp.integranet.ph"
			uadmin.EmailSMTPServerPort = 587
		}

		user := []uadmin.User{}
		uadmin.Filter(&user, "user_group_id = ?  ", closesum.DepartmentID)
		for j := range user {
			uadmin.Preload(&user[j])
			body := o.Remarks + o.Solution

			uadmin.SendEmail([]string{user[j].Email}, []string{}, []string{}, o.Subject, body)
		}
	} else {
	}
}
