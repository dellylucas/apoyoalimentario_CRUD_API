package models

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/utility"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//ConfigurationOptions - Model of configuarations of administrator
type ConfigurationOptions struct {
	Mensajeestudiantes string   `json:"mensajeestudiantes" bson:"mensajeestudiantes"`
	Moduloactivo       bool     `json:"moduloactivo" bson:"moduloactivo"`
	Refrigerionocturno []string `json:"refrigerionocturno" bson:"refrigerionocturno"`
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
	if strings.Compare(SedeChecker, "ALL") != 0 {
		InfoGeneralComplete = Getname(InfoGeneralComplete, SedeChecker)
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

//ReportsAll - Generate GENERIC Reports dynamic
func ReportsAll(sa []StudentInformation, NameSheet string, column []int) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var Maping []MappingColumn
	var MapingNow []MappingColumn
	var err error
	file = xlsx.NewFile()
	sheet, err = file.AddSheet(NameSheet)

	Maping = GEtMappingColumn()
	for numuno := range column {
		for num := range Maping {
			if Maping[num].Value == column[numuno] {
				MapingNow = append(MapingNow, Maping[num])
				break
			}
		}
	}
	var cell *xlsx.Cell
	row = sheet.AddRow()
	for numdo := range MapingNow {
		cell = row.AddCell()
		cell.Value = MapingNow[numdo].Key
	}
	// cell document

	for fil := range sa {
		row = sheet.AddRow()
		sa[fil] = RescueInformation(sa[fil])
		for numdo := range MapingNow {
			cell = row.AddCell()
			MapingNow[numdo] = MapingBD(sa[fil], MapingNow[numdo])
			if MapingNow[numdo].Result != nil {
				cell.Value = MapingNow[numdo].Result.(string)
			}
		}
	}

	err = file.Save("tempfile" + ".xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

//ReportSPADIES - Generate Reports all students "SPADIES"
func ReportSPADIES(sa []StudentInformation) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	// cell document
	var cNumber *xlsx.Cell
	var cName *xlsx.Cell
	var cTypeDoc *xlsx.Cell
	var cDocument *xlsx.Cell
	var cCode *xlsx.Cell
	var cNameProgram *xlsx.Cell
	//var service
	var ModelBasic XmlBasic
	var ModelFacul XmlFaculty
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("SPADIES " + strconv.Itoa(time.Now().UTC().Year()) + " - " + strconv.Itoa(utility.Semester()))
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	cNumber = row.AddCell()
	cName = row.AddCell()
	cTypeDoc = row.AddCell()
	cDocument = row.AddCell()
	cCode = row.AddCell()
	cNameProgram = row.AddCell()
	//name column
	cNumber.HMerge = 50000
	cNumber.Value = "N"
	cName.Value = "NOMBRE Y APELLIDOS"
	cTypeDoc.Value = "Tipo Documento"
	cDocument.Value = "Documento"
	cCode.Value = "Codigo"
	cNameProgram.Value = "Nombre Programa"
	for fil := range sa {
		utility.GetServiceXML(&ModelFacul, utility.FacultyService+sa[fil].Codigo)
		utility.GetServiceXML(&ModelBasic, utility.BasicService+sa[fil].Codigo)
		//add row
		row = sheet.AddRow()
		cNumber = row.AddCell()
		cName = row.AddCell()
		cTypeDoc = row.AddCell()
		cDocument = row.AddCell()
		cCode = row.AddCell()
		cNameProgram = row.AddCell()
		//add value
		cNumber.Value = strconv.Itoa(fil + 1)
		cName.Value = ModelBasic.Name
		cTypeDoc.Value = ModelBasic.TypeDoc
		cDocument.Value = ModelBasic.Document
		cCode.Value = sa[fil].Codigo
		cNameProgram.Value = ModelFacul.Proyecto
	}

	err = file.Save("SPADIES" + strconv.Itoa(time.Now().UTC().Year()) + "-" + strconv.Itoa(utility.Semester()) + ".xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

/* function bonus */

//Getname - Get name of student for faculty
func Getname(model []StudentInformation, SedeChecker string) []StudentInformation {
	var ModelBasic XmlBasic
	var Getinfo []StudentInformation
	var ModelFacul XmlFaculty

	for fil := range model {
		utility.GetServiceXML(&ModelFacul, utility.FacultyService+model[fil].Codigo)
		str := strings.Replace(ModelFacul.NameFaculty, "/", "-", -1)
		if strings.Compare(SedeChecker, str) == 0 && len(model[fil].Informacioneconomica) > 0 {
			utility.GetServiceXML(&ModelBasic, utility.BasicService+model[fil].Codigo)
			model[fil].Nombre = ModelBasic.Name
			Getinfo = append(Getinfo, model[fil])
		}
	}
	return Getinfo
}
func (f *StudentInformation) reflect(ret string) interface{} {
	val := reflect.ValueOf(f).Elem()
	var res interface{}
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		if strings.Compare(typeField.Type.String(), "[]models.Economic") == 0 {
			res = fmt.Sprintf("%v", f.Informacioneconomica[0].reflectEcono(ret))
		} else {
			if strings.Compare(ret, typeField.Name) == 0 {
				res = fmt.Sprintf("%v", valueField.Interface())
				break
			}
		}
	}
	return res
}

func (f *Economic) reflectEcono(ret string) interface{} {
	val := reflect.ValueOf(f).Elem()
	var res interface{}
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		if strings.Compare(ret, typeField.Name) == 0 {
			res = valueField.Interface()
			break
		}
	}
	return res
}

//MapingBD - map BD to Collumn Dynamic
func MapingBD(sa StudentInformation, values MappingColumn) MappingColumn {
	values.Result = sa.reflect(values.Key)
	return values
}

//RescueInformation - rescue info student
func RescueInformation(sa StudentInformation) StudentInformation {
	var ModelFacult XmlFaculty
	var ModelBasic XmlBasic
	var ModelAcademic XmlAcademic
	utility.GetServiceXML(&ModelFacult, utility.FacultyService+sa.Codigo)
	utility.GetServiceXML(&ModelBasic, utility.BasicService+sa.Codigo)
	utility.GetServiceXML(&ModelAcademic, utility.AcademicService+sa.Codigo)
	sa.Nombre = ModelBasic.Name
	sa.Localidad = ModelBasic.Localidad
	sa.Direccion = ModelBasic.Direccion
	sa.Genero = ModelBasic.Genero
	sa.TDocument = ModelBasic.TypeDoc
	sa.Document = ModelBasic.Document
	sa.Facultad = ModelFacult.NameFaculty
	sa.Proyecto = ModelFacult.Proyecto
	sa.Semestre = ModelAcademic.Semestre
	sa.Promedio = ModelAcademic.Promedio
	return sa
}
