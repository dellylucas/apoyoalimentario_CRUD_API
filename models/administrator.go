package models

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/utility"
	"fmt"
	"strconv"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//ConfigurationOptions - Model of configuarations of administrator
type ConfigurationOptions struct {
	Mensajeestudiantes string   `json:"mensajeestudiantes" bson:"mensajeestudiantes"`
	Moduloactivo       bool     `json:"moduloactivo" bson:"moduloactivo"`
	Refrigerionocturno []string `json:"refrigerionocturno" bson:"refrigerionocturno"`
	Reminder           string   `json:"reminder" bson:"reminder"`
	Modulomodified     bool     `json:"modulomodified" bson:"modulomodified"`
	Salariominimo      int      `json:"salariominimo" bson:"salariominimo"`
}

//GetInscription - all records for current semester by Sede
func GetInscription(session *mgo.Session, State string, model ReportsType) ([]StudentInformation, error) {

	MainSession := db.Cursor(session, utility.CollectionGeneral)
	i, _ := strconv.Atoi(State)
	defer session.Close()

	var InfoGeneralComplete []StudentInformation

	query := []bson.M{
		{
			"$lookup": bson.M{ // lookup the documents table here
				"from": utility.CollectionEconomic,
				"let":  bson.M{"general_id": "$_id"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{"estadoprograma": i,
							"periodo":  model.Periodo,
							"semestre": model.Semestre,
							"$expr": bson.M{
								"$eq": []string{"$$general_id", "$id"},
							}},
					}},
				"as": utility.CollectionEconomic,
			}}}
	err := MainSession.Pipe(query).All(&InfoGeneralComplete)
	if strings.Compare(model.TSede, "ALL") != 0 {
		InfoGeneralComplete = Getname(InfoGeneralComplete, model.TSede)
	} else {
		var AllStudents []StudentInformation
		for fil := range InfoGeneralComplete {
			if len(InfoGeneralComplete[fil].Informacioneconomica) > 0 {
				AllStudents = append(AllStudents, InfoGeneralComplete[fil])
			}
		}
		InfoGeneralComplete = AllStudents
	}
	if err != nil {
		fmt.Println(err)
	}
	return InfoGeneralComplete, err
}

//GetConfiguration - View message administrator
func GetConfiguration(session *mgo.Session) (ConfigurationOptions, error) {

	BDMessage := db.Cursor(session, utility.CollectionAdministrator)
	defer session.Close()
	var MessageComplete ConfigurationOptions
	err := BDMessage.Find(nil).One(&MessageComplete)
	if err != nil {
		fmt.Println(err)
	}
	return MessageComplete, err
}

//UpdateInformationConfig - Update the information economic of student
func UpdateInformationConfig(session *mgo.Session, newInfo ConfigurationOptions) error {
	BDMessage := db.Cursor(session, utility.CollectionAdministrator)
	defer session.Close()
	err := BDMessage.Update(nil, &newInfo)

	return err
}

/* function bonus */

//Getname - Get name of student for faculty
func Getname(model []StudentInformation, SedeChecker string) []StudentInformation {
	var ModelBasic XmlBasic
	var Getinfo []StudentInformation
	var ModelFacul XmlFaculty

	for fil := range model {
		utility.GetServiceXML(&ModelFacul, utility.FacultyService+model[fil].Codigo)
		facultadEstudiante := strings.Replace(ModelFacul.NameFaculty, "/", "-", -1)
		if strings.Compare(SedeChecker, facultadEstudiante) == 0 && len(model[fil].Informacioneconomica) > 0 {
			utility.GetServiceXML(&ModelBasic, utility.BasicService+model[fil].Codigo)
			model[fil].Nombre = ModelBasic.Name
			Getinfo = append(Getinfo, model[fil])
		}
	}
	return Getinfo
}
