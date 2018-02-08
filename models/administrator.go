package models

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/utility"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"

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
							"$expr": bson.M{
								"$eq": []string{"$$general_id", "$id"},
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
			utility.GetServiceXML(&ModelBasic, utility.BasicService+model[fil].Codigo)
			model[fil].Nombre = ModelBasic.Name
			PruebaGetinfo = append(PruebaGetinfo, model[fil])
		}
	}
	return PruebaGetinfo
}

//Reports - Generate Reports
func Reports(sa []StudentInformation, SedeChecker string) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var cellnombre *xlsx.Cell
	var celldos *xlsx.Cell
	var cellfecha *xlsx.Cell
	var celltipoApoyo *xlsx.Cell
	var celltipoSubs *xlsx.Cell
	var cellObserv *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Prueba1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	cell = row.AddCell()
	cellnombre = row.AddCell()
	celldos = row.AddCell()
	cellfecha = row.AddCell()
	celltipoApoyo = row.AddCell()
	celltipoSubs = row.AddCell()
	cellObserv = row.AddCell()
	cell.Value = "codigo"
	cellnombre.Value = "Nombre"
	celldos.Value = "Ciudad"
	cellfecha.Value = "Fecha de inscripcion"
	celltipoApoyo.Value = "Tipo de Apoyo"
	celltipoSubs.Value = "Subsidio"
	cellObserv.Value = "Observaciones"
	for fil := range sa {
		row = sheet.AddRow()
		cell = row.AddCell()
		cellnombre = row.AddCell()
		celldos = row.AddCell()
		cellfecha = row.AddCell()
		celltipoApoyo = row.AddCell()
		celltipoSubs = row.AddCell()
		cellObserv = row.AddCell()
		cell.Value = sa[fil].Codigo
		cellnombre.Value = sa[fil].Nombre
		celldos.Value = sa[fil].Informacioneconomica[0].Ciudad
		cellfecha.SetDate(sa[fil].Fechainscripcion)
		celltipoApoyo.Value = sa[fil].Informacioneconomica[0].Tipoapoyo
		celltipoSubs.Value = sa[fil].Informacioneconomica[0].TipoSubsidio
		cellObserv.Value = sa[fil].Informacioneconomica[0].Mensaje
	}

	err = file.Save(SedeChecker + ".xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}
