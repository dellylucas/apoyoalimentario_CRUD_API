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

//GetInscription - Retorna todos los estudiantes de un semestre, año(periodo), estado y/o sede  establecido
//Param session  IN	"sesion de base de datos"
//Param State	IN	"estado en el programa de los estudiantes a consultar"
//Param model	IN	"modelo determina el semestre, año y sede a consultar"
//Param InfoGeneralComplete	OUT  "Devuelve todos los estudiantes a consultar"
//Param err		OUT   "error si es que existe"
func GetInscription(session *mgo.Session, State string, model *ReportsType) (*[]StudentInformation, error) {

	MainSession := db.Cursor(session, utility.CollectionGeneral)
	i, _ := strconv.Atoi(State)

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
		InfoGeneralComplete = Getname(&InfoGeneralComplete, model.TSede)
	} else {
		var AllStudents []StudentInformation
		for _, fil := range InfoGeneralComplete {
			if len(fil.Informacioneconomica) > 0 {
				AllStudents = append(AllStudents, fil)
			}
		}
		InfoGeneralComplete = AllStudents
	}
	if err != nil {
		fmt.Println(err)
	}
	return &InfoGeneralComplete, err
}

//GetConfiguration - Retorna la configuracion del administrador
//Param session		IN  	"sesion de base de datos"
//Param Config		OUT   "Retorna la configuracion del administrador"
//Param err		OUT   "error si es que existe"
func GetConfiguration(session *mgo.Session) (*ConfigurationOptions, error) {

	BDMessage := db.Cursor(session, utility.CollectionAdministrator)
	defer session.Close()
	var Config ConfigurationOptions
	err := BDMessage.Find(nil).One(&Config)
	return &Config, err
}

//UpdateInformationConfig - Actualiza la informacion de la configuarcion el administrador
//Param session		IN  	"sesion de base de datos"
//Param newInfo		IN  	"tiene el modelo de las variables de configuracion a actualizar"
//Param err		OUT   "error si es que existe"
func UpdateInformationConfig(session *mgo.Session, newInfo *ConfigurationOptions) error {
	BDMessage := db.Cursor(session, utility.CollectionAdministrator)
	defer session.Close()
	err := BDMessage.Update(nil, &newInfo)

	return err
}

/* function bonus */

//Getname - optiene la informacion basica de los estudiantes que poseen informacion economica del actual semestre y pertenecen a una sede especifica
//Param model		IN  	"estudiantes optenidos de base de datos"
//Param SedeChecker		IN  	"filtro de sede"
//Param Getinfo		OUT 	"estudiantes post filtro"
func Getname(model *[]StudentInformation, SedeChecker string) []StudentInformation {
	var ModelBasic XmlBasic
	var Getinfo []StudentInformation
	var ModelFacul XmlFaculty

	for _, student := range *model {
		utility.GetServiceXML(&ModelFacul, utility.FacultyService+student.Codigo, nil)
		facultadEstudiante := strings.Replace(ModelFacul.NameFaculty, "/", "-", -1)
		if strings.Compare(SedeChecker, facultadEstudiante) == 0 && len(student.Informacioneconomica) > 0 {
			utility.GetServiceXML(&ModelBasic, utility.BasicService+student.Codigo, nil)
			student.Nombre = ModelBasic.Name
			Getinfo = append(Getinfo, student)
		}
	}
	return Getinfo
}
