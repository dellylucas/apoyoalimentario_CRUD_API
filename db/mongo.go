package db

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Cursor - sd
func Cursor(session *mgo.Session, Collection string) *mgo.Collection {
	mongoDB := beego.AppConfig.String("mongo_db")
	c := session.DB(mongoDB).C(Collection)
	return c
}

//GetSession - dfs
func GetSession() (*mgo.Session, error) {

	mongoHost := beego.AppConfig.String("mongo_host")
	mongoUser := beego.AppConfig.String("mongo_user")
	mongoPassword := beego.AppConfig.String("mongo_pass")
	mongoDatabase := beego.AppConfig.String("mongo_db")

	info := &mgo.DialInfo{
		Addrs:    []string{mongoHost},
		Timeout:  60 * time.Second,
		Database: mongoDatabase,
		Username: mongoUser,
		Password: mongoPassword,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		fmt.Println("Helo this is an error!")
		panic(err)
	} else {
		session.SetMode(mgo.Monotonic, true)
	}

	return session, err
}

//GetAll - records from Collection
func GetAll(session *mgo.Session, collection string) []bson.M {
	c := Cursor(session, collection)
	defer session.Close()
	var records []bson.M
	err := c.Find(bson.M{}).All(&records)
	if err != nil {
		fmt.Println(err)
	}

	return records
}
