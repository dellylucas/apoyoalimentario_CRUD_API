package db

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Cursor -  conexion a la coleccion deseada
// Param	session		 	IN	"sesion actual"
// Param	Collection		IN	"coleccion la cual se desea conectar"
// Param  c 			OUT		"retorna la conexion a dicha coleccion"
func Cursor(session *mgo.Session, Collection string) *mgo.Collection {
	mongoDB := beego.AppConfig.String("mongo_db")
	c := session.DB(mongoDB).C(Collection)
	return c
}

//GetSession - conexion de Base de Datos
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
		fmt.Println("error session DB!")
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
