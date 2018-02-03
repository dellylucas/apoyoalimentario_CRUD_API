package models

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/utility"
	"fmt"
	"strconv"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//ConfigurationOptions - Model of configuarations of administrator
type ConfigurationOptions struct {
	Mensajeestudiantes string   `json:"mensajeestudiantes,omitempty" bson:"mensajeestudiantes,omitempty"`
	Moduloactivo       bool     `json:"moduloactivo,omitempty" bson:"moduloactivo,omitempty"`
	Refrigerionocturno []string `json:"refrigerionocturno,omitempty" bson:"refrigerionocturno,omitempty"`
}

//TypeRol - Model of type rol of user
type TypeRol struct {
	Rol     string `json:"rol" bson:"rol"`
	Usuario int    `json:"usuario" bson:"usuario"`
	Sede    string `json:"sede" bson:"sede"`
}

//GetTypeUser - Get type user Administrator or checker
func GetTypeUser(session *mgo.Session, user string) (TypeRol, error) {

	Principal := db.Cursor(session, utility.CollectionAdministrator)
	defer session.Close()
	i, _ := strconv.Atoi(user)
	var InfoUser TypeRol

	err := Principal.Find(bson.M{"usuario": i}).One(&InfoUser)

	if err != nil {
		fmt.Println(err)
	}
	return InfoUser, err
}

//GetInscription - all records for current semester by Sede
func GetInscription(session *mgo.Session, State string, SedeChecker string) ([]StudentInformation, error) {

	MainSession := db.Cursor(session, utility.CollectionGeneral)
	i, _ := strconv.Atoi(State)
	defer session.Close()

	var InfoGeneralComplete []StudentInformation

	query := []bson.M{
		{
			"$lookup": bson.M{ // lookup the documents table here
				"from": "informacioneconomica",
				"let":  bson.M{"general_id": "$_id"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{"estadoprograma": i,
							"periodo":  time.Now().UTC().Year(),
							"semestre": utility.Semester(),
							"$expr": bson.M{"$and": []bson.M{
								{"$eq": []string{"$$general_id", "$id"}},
							},
							}},
					}},
				"as": "informacioneconomica",
			}}}
	err := MainSession.Pipe(query).All(&InfoGeneralComplete)
	InfoGeneralComplete = Getname(InfoGeneralComplete, SedeChecker)
	if err != nil {
		fmt.Println(err)
	}
	return InfoGeneralComplete, err
}

//GetMessage - View message administrator
func GetMessage(session *mgo.Session) (string, error) {

	BDMessage := db.Cursor(session, utility.CollectionAdministrator)
	defer session.Close()
	var MessageComplete ConfigurationOptions
	err := BDMessage.Find(nil).Select(bson.M{"mensajeestudiantes": 1}).One(&MessageComplete)
	if err != nil {
		fmt.Println(err)
	}
	return MessageComplete.Mensajeestudiantes, err
}

/* function bonus */

//Getname - Get name of student
func Getname(model []StudentInformation, SedeChecker string) []StudentInformation {
	var ModelBasic XmlBasic
	var PruebaGetinfo []StudentInformation
	var ModelFacul XmlFaculty

	for fil := range model {
		utility.GetServiceXML(&ModelFacul, utility.FacultyService+model[fil].Codigo)
		str := strings.Replace(ModelFacul.NameFaculty, "/", "-", -1)
		if strings.Compare(SedeChecker, str) == 0 && len(model[fil].Informacioneconomica) > 0 {
			s := len(model[fil].Informacioneconomica)
			fmt.Println(s)
			utility.GetServiceXML(&ModelBasic, utility.BasicService+model[fil].Codigo)
			model[fil].Nombre = ModelBasic.Name
			PruebaGetinfo = append(PruebaGetinfo, model[fil])
		}
	}
	return PruebaGetinfo
}
