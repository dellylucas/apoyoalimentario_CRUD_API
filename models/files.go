package models

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/utility"
	"os"
	"strconv"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//FilesStudents Struct of save history files
type FilesStudents struct {
	ID         bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Name       string        `json:"nombre" bson:"nombre"`
	Size       int64         `json:"longitud" bson:"longitud"`
	Dateinsert time.Time     `json:"fecha" bson:"fecha"`
	Code       string        `json:"codigo" bson:"codigo"`
	Urlfile    string        `json:"url" bson:"url"`
}

//Deletefile - check files not used and delete
func Deletefile(session *mgo.Session, code string, claves *[]string) {

	FileSession := db.Cursor(session, utility.CollectionHistoricFiles)
	fromDate, toDate := utility.GetInitEnd()
	path := utility.FileSavePath + code + "\\" + strconv.Itoa(time.Now().UTC().Year()) + "-" + strconv.Itoa(utility.Semester()) + "\\"
	for _, element := range *claves {
		_ = FileSession.Remove(bson.M{"codigo": code, "nombre": element, "fecha": bson.M{"$gt": fromDate, "$lt": toDate}})
		_ = os.Remove(path + element + ".pdf")
	}
}

//Completefile - Check files complete
func Completefile(session *mgo.Session, code string, clave *[]string) error {

	FileSession := db.Cursor(session, utility.CollectionHistoricFiles)
	fromDate, toDate := utility.GetInitEnd()
	var Infofilepath FilesStudents
	var errP error
	for _, element := range *clave {
		errP = FileSession.Find(bson.M{"codigo": code, "nombre": element, "fecha": bson.M{"$gt": fromDate, "$lt": toDate}}).One(&Infofilepath)
		if errP != nil {
			break
		}
	}
	return errP
}

//GetFiles - get all files by code in current semester
func GetFiles(session *mgo.Session, code string) (*[]FilesStudents, error) {

	FileSession := db.Cursor(session, utility.CollectionHistoricFiles)
	fromDate, toDate := utility.GetInitEnd()
	defer session.Close()
	var Infofilepath []FilesStudents
	errP := FileSession.Find(bson.M{"codigo": code, "fecha": bson.M{"$gt": fromDate, "$lt": toDate}}).All(&Infofilepath)
	return &Infofilepath, errP
}

//Insertfile - Insert file(s) in path
func Insertfile(session *mgo.Session, name string, size int64, code string) { //nombre ,tamañ, fech, autor
	FileSession := db.Cursor(session, utility.CollectionHistoricFiles)
	var Exists FilesStudents
	var Hist FilesStudents
	Hist.Code = code
	Hist.Name = name
	Hist.Size = size
	Hist.Urlfile = utility.ServerPath + code + "/" + strconv.Itoa(time.Now().UTC().Year()) + "-" + strconv.Itoa(utility.Semester()) + "/" + name + ".pdf"
	Hist.Dateinsert = time.Now().UTC()
	fromDate, toDate := utility.GetInitEnd()
	err := FileSession.Find(bson.M{"codigo": code, "nombre": name, "fecha": bson.M{"$gt": fromDate, "$lt": toDate}}).One(&Exists)
	if err != nil {
		FileSession.Insert(Hist)
	} else {
		_ = FileSession.Update(bson.M{"_id": Exists.ID, "codigo": code, "nombre": name, "fecha": bson.M{"$gt": fromDate, "$lt": toDate}}, &Hist)
	}
}
