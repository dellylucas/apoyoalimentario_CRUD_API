package models

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/utility"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Sede Struct for ver...
type Sede struct {
	ID            bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Name          string        `json:"name" bson:"name"`
	Nombre        string        `json:"nombre" bson:"nombre"`
	Verificadores []string      `json:"verificadores" bson:"verificadores"`
}

//GetVerifier - Optiene todos las asociaciones de todos los verificadores
//Param session		IN   "sesion de base de datos"
//Param SedeVerif		OUT   "Optiene la asiciacion"
//Param errd		OUT   "error si es que existe"
func GetVerifier(session *mgo.Session) (*[]Sede, error) {
	MainSession := db.Cursor(session, utility.CollectionAdministrator)

	var SedeVerif []Sede

	err := MainSession.Find(bson.M{"name": "verificadores"}).All(&SedeVerif)

	return &SedeVerif, err
}

//UpdateVerifier - Actualiza las asociaciones de todos los verificadores
//Param session		IN   "sesion de base de datos"
//Param newInfo		IN   "asociaciones verificadoes y sedes"
//Param err		OUT   "error si es que existe"
func UpdateVerifier(session *mgo.Session, newInfo *[]Sede) error {
	MainSession := db.Cursor(session, utility.CollectionAdministrator)
	var err error
	err = nil
	var OldInfo Sede
	MainSession.RemoveAll(bson.M{"name": "verificadores"})
	for _, element := range *newInfo {
		err = MainSession.Find(bson.M{"name": "verificadores", "nombre": element.Nombre}).One(&OldInfo)
		if err != nil {
			element.Name = "verificadores"
			MainSession.Insert(element)
		} else {
			err = MainSession.Update(bson.M{"_id": OldInfo.ID}, &element)
		}
	}

	return err
}

//GetSede - optiene la(s) sede(s) que esta asociado un verificador
//Param session		IN   "sesion de base de datos"
//Param name		IN   "nombre del verificador a consultar"
//Param sedes		OUT   "Sedes a las cuales esta asociado el verificador"
//Param err		OUT   "error si es que existe"
func GetSede(session *mgo.Session, name string) (*[]string, error) {
	MainSession := db.Cursor(session, utility.CollectionAdministrator)

	var SedeVerif []Sede
	var sedes []string
	err := MainSession.Find(bson.M{"name": "verificadores"}).All(&SedeVerif)
	for _, element := range SedeVerif {
		for _, verifi := range element.Verificadores {
			if strings.Compare(name, verifi) == 0 {
				sedes = append(sedes, element.Nombre)
				break
			}
		}
	}
	return &sedes, err
}
