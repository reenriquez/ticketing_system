package models

import (
	"math/rand"
	"time"

	"github.com/uadmin/uadmin"
	"gopkg.in/gomail.v2"
)

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Ticket struct {
	uadmin.Model
	CreatedAt    time.Time   `uadmin:"read_only"`
	TicketNumber string      `uadmin:"read_only;hidden"`
	CreatedBy    string      `uadmin:"read_only;hidden;list_exclude"`
	User         uadmin.User `uadmin:"read_only;hidden"`
	UserID       uint
	Department   uadmin.UserGroup `uadmin:"read_only;hidden"`
	DepartmentID uint
	Subject      string `uadmin:"required"`
	Issue        string `uadmin:"html;required " sql:"type:LONGTEXT"`
	Attachment   string `uadmin:"image; upload_to:media/ticket/"`
	Status       bool   `uadmin:"read_only;hidden"`
}

func (t *Ticket) String() string {
	uadmin.Preload(t)
	return t.TicketNumber
}

func (t *Ticket) GetActiveSession() {
	active := uadmin.Session{}
	use := uadmin.User{}
	uadmin.Get(&use, "username = ?", t.CreatedBy)
	uadmin.Get(&active, "user_id = ? ", use.ID)
	uadmin.Trail(uadmin.DEBUG, "ACTIVE: %v", active.UserID)
	t.UserID = active.UserID
	return
}

func (t *Ticket) Save() {

	ticket := Ticket{}
	code := ""
	code = RandStringBytes(10)
	for uadmin.Count(&ticket, "ticket_number = ?", code) > 0 {
		code = RandStringBytes(10)
	}
	t.TicketNumber = code
	opensum := OpenTicket{}
	uadmin.Preload(t)
	open := false

	em := Employee{}
	us := uadmin.User{}

	if t.Status == open {
		t.GetActiveSession()
		uadmin.Get(&us, "id = ?", t.UserID)
		opensum.UserID = t.UserID
		opensum.CreatedBy = t.CreatedBy
		opensum.DepartmentID = us.UserGroupID
		opensum.Subject = t.Subject
		opensum.Issue = t.Issue
		opensum.Attachment = t.Attachment
		opensum.TicketNumber = t.TicketNumber
		opensum.CreatedAt = t.CreatedAt
		opensum.Status = t.Status
		uadmin.Save(t)
		uadmin.Save(&opensum)
		uadmin.Delete(t)

		// * Validation of Department and sending email
		uadmin.Get(&us, "user_group_id = ?", opensum.DepartmentID)
		user := []uadmin.User{}
		uadmin.Filter(&user, "user_group_id = ?  ", opensum.DepartmentID)
		for j := range user {
			uadmin.Preload(&user[j])

			if opensum.DepartmentID == us.UserGroupID {
				m := gomail.NewMessage()
				m.SetHeader("From", us.Email)
				m.SetHeader("To", user[j].Email)
				m.SetAddressHeader("Cc", "jmandal@integranet.ph", "Developers")
				m.SetHeader("Subject", t.Subject)
				m.SetBody("text/html", "We have received your issue. Kindly wait for reply from our support. Thank you!")
				m.Attach("/home/Alex/lolcat.jpg")

				d := gomail.NewPlainDialer("smtp.integranet.ph", 587, user.Email, user.Password)
				if err := d.DialAndSend(m); err != nil {
					panic(err)
				}
			}
		}
	}
}

// * Generator for Ticket Number
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
