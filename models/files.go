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
func Deletefile(session *mgo.Session, code string, claves []string) string {

	FileSession := db.Cursor(session, utility.CollectionHistoricFiles)
	fromDate, toDate := utility.GetInitEnd()
	count := 0
	path := utility.FileSavePath + code + "\\" + strconv.Itoa(time.Now().UTC().Year()) + "-" + strconv.Itoa(utility.Semester()) + "\\"
	var errGlobal string
	for count < len(claves) {

		errP := FileSession.Remove(bson.M{"codigo": code, "nombre": claves[count], "fecha": bson.M{"$gt": fromDate, "$lt": toDate}})
		err := os.Remove(path + claves[count] + ".pdf")
		if errP != nil {
			errGlobal = errGlobal + "/ " + errP.Error() + claves[count] + ".pdf"
		}
		if err != nil {
			errGlobal = errGlobal + "/ " + err.Error() + claves[count] + ".pdf"
		}
		count++
	}
	return errGlobal
}

//Completefile - Check files complete
func Completefile(session *mgo.Session, code string, clave []string) (int, error) {

	FileSession := db.Cursor(session, utility.CollectionHistoricFiles)
	fromDate, toDate := utility.GetInitEnd()
	var Infofilepath FilesStudents
	count := 0
	var result int
	var errP error
	for count < len(clave) {

		errP = FileSession.Find(bson.M{"codigo": code, "nombre": clave[count], "fecha": bson.M{"$gt": fromDate, "$lt": toDate}}).One(&Infofilepath)

		if errP != nil {
			count = len(clave)
			result = 0
		} else {
			count++
			result = 1
		}
	}
	return result, errP
}

//GetFiles - get all files by code in current semester
func GetFiles(session *mgo.Session, code string) ([]FilesStudents, error) {
	//	path := utils.PathRootSaveFile + codigo + "/" + strconv.Itoa() + "-" + strconv.Itoa(utils.Semestre()) + "/"
	FileSession := db.Cursor(session, utility.CollectionHistoricFiles)
	fromDate, toDate := utility.GetInitEnd()
	defer session.Close()
	var Infofilepath []FilesStudents
	errP := FileSession.Find(bson.M{"codigo": code, "fecha": bson.M{"$gt": fromDate, "$lt": toDate}}).All(&Infofilepath)
	return Infofilepath, errP
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
