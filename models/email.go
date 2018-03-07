package models

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/utility"
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/gomail.v2"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Email - Model of email to send
type Email struct {
	Name        string `json:"name" bson:"name"`
	Server      string `json:"server" bson:"server"`
	Port        int    `json:"port" bson:"port"`
	SecuritySSL bool   `json:"securitySSL" bson:"securitySSL"`
	SecurityTLS bool   `json:"securityTLS" bson:"securityTLS"`
	EmailCon    string `json:"emailCon" bson:"emailCon"`
	Pass        string `json:"pass" bson:"pass"`
	Subject     string `json:"subject" bson:"subject"`
	Text        string `json:"text" bson:"text"`
}

//BodyEmail - body to send
type BodyEmail struct {
	EBody   string `json:"eBody" bson:"eBody"`
	EtoSend string `json:"etoSend" bson:"etoSend"`
	EName   string `json:"eName" bson:"eName"`
}

//SearchInfor - get configuration of email
func SearchInfor(session *mgo.Session) (Email, error) {
	MainSession := db.Cursor(session, utility.CollectionAdministrator)
	var InfoConfig Email
	errd := MainSession.Find(bson.M{"name": "email"}).One(&InfoConfig)
	if errd != nil {
		panic(errd)
	}
	return InfoConfig, errd
}

//EmailSender - send email
func EmailSender(Bod BodyEmail, session *mgo.Session) error {
	var Info Email
	MainSession := db.Cursor(session, utility.CollectionAdministrator)
	err := MainSession.Find(bson.M{"name": "email"}).One(&Info)
	if err == nil {
		// "smtp.gmail.com", 465, ssl y tls
		d := gomail.NewDialer(Info.Server, Info.Port, Info.EmailCon, Info.Pass)

		if Info.SecuritySSL == true {
			d.SSL = true
		}
		if Info.SecurityTLS == true {
			d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}

		s, err := d.Dial()
		if err == nil {
			m := gomail.NewMessage()
			m.SetHeader("From", Info.EmailCon)
			m.SetAddressHeader("To", Bod.EtoSend, Bod.EName)
			m.SetHeader("Subject", Info.Subject)
			m.SetBody("text/html", fmt.Sprintf(Info.Text+Bod.EBody))
			err = gomail.Send(s, m)
			m.Reset()
		}
	}

	return err
}

//TestConnection - test conn
func TestConnection(Info Email) error {

	d := gomail.NewDialer(Info.Server, Info.Port, Info.EmailCon, Info.Pass)

	if Info.SecuritySSL == true {
		d.SSL = true
	}
	if Info.SecurityTLS == true {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	_, err := d.Dial()

	return err
}
