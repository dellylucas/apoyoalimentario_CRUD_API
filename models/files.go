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

//Deletefile - chequea archivos que estudiante elimina y borra (historial y servidor)
//Param session		IN   "sesion de base de datos"
//Param code		IN   "Codigo del estudiante a consultar"
//Param claves	IN   "lista de nombres de archivos que no deben de estar en servidor ni historico"
func Deletefile(session *mgo.Session, code string, claves *[]string) {

	FileSession := db.Cursor(session, utility.CollectionHistoricFiles)
	fromDate, toDate := utility.GetInitEnd()
	path := utility.FileSavePath + code + "\\" + strconv.Itoa(time.Now().UTC().Year()) + "-" + strconv.Itoa(utility.Semester()) + "\\"
	for _, element := range *claves {
		_ = FileSession.Remove(bson.M{"codigo": code, "nombre": element, "fecha": bson.M{"$gt": fromDate, "$lt": toDate}})
		_ = os.Remove(path + element + ".pdf")
	}
}

//Completefile - verifica que todos los documentos esten para realizar inscripcion
//Param session		IN   "sesion de base de datos"
//Param code		IN   "Codigo del estudiante a consultar"
//Param clave	IN   "lista de nombres de archivos que si deben de estar en servidor e historico"
//Param errP		OUT   "error si es que existe"
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

//GetFiles - optiene todos los arcivos del historial para un estudiante en el semestre actual
//Param session		IN   "sesion de base de datos"
//Param code		IN   "Codigo del estudiante a consultar"
//Param Infofilepath	OUT   "Archivos que tiene el estudiante subidos al servidor"
//Param errP		OUT   "error si es que existe"
func GetFiles(session *mgo.Session, code string) (*[]FilesStudents, error) {

	FileSession := db.Cursor(session, utility.CollectionHistoricFiles)
	fromDate, toDate := utility.GetInitEnd()
	defer session.Close()
	var Infofilepath []FilesStudents
	errP := FileSession.Find(bson.M{"codigo": code, "fecha": bson.M{"$gt": fromDate, "$lt": toDate}}).All(&Infofilepath)
	return &Infofilepath, errP
}

//Insertfile - Inserta archivos en el historico
//Param session		IN   "sesion de base de datos"
//Param name		IN   "nombre del archivo"
//Param size	IN   "tamaño dl archivo"
//Param code		IN   "codigo del estudiante"
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
