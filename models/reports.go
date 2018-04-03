package models

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/utility"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/tealeg/xlsx"
	mgo "gopkg.in/mgo.v2"
)

//ReportsType Struct for map reports
type ReportsType struct {
	Columnas   []int  `json:"columnas" bson:"columnas"`
	TSede      string `json:"tsede" bson:"tsede"`
	NameSheet  string `json:"nameSheet" bson:"nameSheet"`
	TypeReport int    `json:"typeReport" bson:"typeReport"`
	Periodo    int    `json:"periodo" bson:"periodo"`
	Semestre   int    `json:"semestre" bson:"semestre"`
}

//MappingColumn Struct for map reports
type MappingColumn struct {
	ColumnName string
	Value      int
	Key        string
	Result     interface{}
	Score      string
}

//ReportsGeneric - Genera reporte genericos Dinamicos
//Param sa		IN   "estudiantes a generar  reporte"
//Param NameSheet		IN   "nombre del libro excel asignado"
//Param column		IN   "columnas a generar dentro del reporte"
func ReportsGeneric(sa *[]StudentInformation, NameSheet string, column *[]int) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var MapingNow []MappingColumn
	var err error
	file = xlsx.NewFile()
	sheet, err = file.AddSheet(NameSheet)

	Maping := GEtMappingColumn()
	for _, numuno := range *column {
		for _, num := range *Maping {
			if num.Value == numuno {
				MapingNow = append(MapingNow, num)
				break
			}
		}
	}
	var cell *xlsx.Cell
	row = sheet.AddRow()
	for _, numdo := range MapingNow {
		cell = row.AddCell()
		cell.Value = numdo.ColumnName
	}
	// cell document

	for _, fil := range *sa {
		row = sheet.AddRow()
		RescueInformation(&fil)
		for _, numdo := range MapingNow {
			cell = row.AddCell()
			//reflection
			MapingBD(&fil, &numdo)
			if numdo.Result != nil {
				cell.Value = ProcessinData(&numdo)
			}
		}
	}

	err = file.Save("tempfile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

//ReportGeneral - Generar reporte final con puntajes de estudiantes
//Param session		IN   "sesion de base de datos"
//Param students		IN   "estudiantes a generar  reporte"
//Param name		IN   "nombre del libro excel asignado"
func ReportGeneral(session *mgo.Session, students *[]StudentInformation, name string) {
	BDSMLV := db.Cursor(session, utility.CollectionAdministrator)
	var salario ConfigurationOptions
	err := BDSMLV.Find(nil).One(&salario)
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var MapingNow []MappingColumn
	var count int
	file = xlsx.NewFile()
	sheet, err = file.AddSheet(name)
	var column []int
	column = append(column, 2, 24, 25, 32, 28, 29, 1, 31, 30, 3, 19, 4, 5, 6, 7, 8, 9, 10, 12, 13, 14, 35, 20, 21, 27, 26, 23, 22, 33, 34, 15, 16)
	Maping := GEtMappingColumn()
	for _, numuno := range column {
		for _, num := range *Maping {
			if num.Value == numuno {
				MapingNow = append(MapingNow, num)
				break
			}
		}
	}
	var cell *xlsx.Cell
	row = sheet.AddRow()
	for _, numdo := range MapingNow {
		cell = row.AddCell()
		cell.Value = numdo.ColumnName
		if strings.Compare(numdo.Score, "Si") == 0 {
			cell = row.AddCell()
			cell.Value = "PUNTAJE " + numdo.ColumnName
		}
	}
	// cell document

	for _, fil := range *students {
		count = 0
		row = sheet.AddRow()
		RescueInformation(&fil)

		for _, numdo := range MapingNow {
			cell = row.AddCell()
			MapingBD(&fil, &numdo)
			if numdo.Result != nil {
				cell.Value = ProcessinData(&numdo)
				if strings.Compare(numdo.Score, "Si") == 0 {
					cell = row.AddCell()
					localcount := Evaluation(&numdo, salario.Salariominimo)
					count += localcount
					cell.Value = strconv.Itoa(localcount)
				}
			}
			if strings.Compare(numdo.Key, "Total") == 0 {
				cell.Value = strconv.Itoa(count)
			}
		}
	}
	err = file.Save("tempfile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

//OthersReports - Genera reportes ser pilo paga - estudiantes con sisben y total de inscritos por tipo de apoyo
//Param students		IN   "estudiantes a generar  reporte"
func OthersReports(students *[]StudentInformation) {

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var MapingNow []MappingColumn
	file = xlsx.NewFile()

	//sheet 1 -- sisben
	sheet, err := file.AddSheet("Sisben")
	var column []int
	column = append(column, 1, 25, 31, 30)
	Maping := GEtMappingColumn()
	for _, numuno := range column {
		for _, num := range *Maping {
			if num.Value == numuno {
				MapingNow = append(MapingNow, num)
				break
			}
		}
	}
	//name Columns
	var cell *xlsx.Cell
	row = sheet.AddRow()
	for _, numdo := range MapingNow {
		cell = row.AddCell()
		cell.Value = numdo.ColumnName
	}

	// cell document
	for _, fil := range *students {
		if strings.Compare(fil.Informacioneconomica[0].Sisben, "no") == 0 {
			continue
		}
		row = sheet.AddRow()
		RescueInformation(&fil)
		for _, numdo := range MapingNow {
			cell = row.AddCell()
			MapingBD(&fil, &numdo)
			if numdo.Result != nil {
				cell.Value = numdo.Result.(string)
			}
		}
	}
	//sheet 2 -- Ser Pilo Paga
	sheet, err = file.AddSheet("Ser Pilo Paga")

	//name Columns
	var celldos *xlsx.Cell
	row = sheet.AddRow()
	for _, numdo := range MapingNow {
		celldos = row.AddCell()
		celldos.Value = numdo.ColumnName
	}

	// cell document
	for _, fil := range *students {
		if strings.Compare(fil.Informacioneconomica[0].SerPiloPaga, "no") == 0 {
			continue
		}
		row = sheet.AddRow()
		RescueInformation(&fil)
		for _, numdo := range MapingNow {
			celldos = row.AddCell()
			MapingBD(&fil, &numdo)
			if numdo.Result != nil {
				celldos.Value = numdo.Result.(string)
			}
		}
	}
	//sheet 3 -- sisben
	sheet, err = file.AddSheet("Totales")

	//name Columns
	var cellSIN *xlsx.Cell
	var cellA *xlsx.Cell
	var cellB *xlsx.Cell
	var cellCON *xlsx.Cell
	row = sheet.AddRow()
	cellSIN = row.AddCell()
	cellA = row.AddCell()
	cellB = row.AddCell()
	cellCON = row.AddCell()
	cellSIN.Value = "SIN SUBSIDIO"
	cellA.Value = "TIPO A"
	cellB.Value = "TIPO B"
	cellCON.Value = "SUBSIDIO TOTAL"
	var TSIN int //count SIN SUBSIDIO
	var TA int   // count TIPO A
	var TB int   // count TIPO B
	var TT int   // count SUBSIDIO TOTAL
	// cell document
	for _, fil := range *students {
		switch os := fil.Informacioneconomica[0].TipoSubsidio; os {
		case "t":
			TT++
		case "a":
			TA++
		case "b":
			TB++
		case "ss":
			TSIN++
		}
	}
	row = sheet.AddRow()
	cellSIN = row.AddCell()
	cellA = row.AddCell()
	cellB = row.AddCell()
	cellCON = row.AddCell()
	cellSIN.Value = strconv.Itoa(TSIN)
	cellA.Value = strconv.Itoa(TA)
	cellB.Value = strconv.Itoa(TB)
	cellCON.Value = strconv.Itoa(TT)

	err = file.Save("tempfile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

/* bonus functions*/

//MakeThing - construye modelo de columnas sin puntaje
//Param Col		IN   "nombre de columna"
//Param Val		IN   "valor o id"
//Param Keys	IN   "nombre asignado en base de datos"
func MakeThing(Col string, Val int, Keys string) MappingColumn {
	return MappingColumn{Col, Val, Keys, "", ""}
}

//MakeThingD - construye modelo de columnas con puntaje
//Param Col		IN   "nombre de columna"
//Param Val		IN   "valor o id"
//Param Keys	IN   "nombre asignado en base de datos"
//Param Score	IN   "Puntaje asignado a dicha respuesta"
func MakeThingD(Col string, Val int, Keys string, Score string) MappingColumn {
	return MappingColumn{Col, Val, Keys, "", Score}
}

//MapingBD - mapeo de valores por columna dynamicos (identificacion de columnas y valores)
//Param sa		IN   "estudiante"
//Param values		IN   "valores de columnas"
func MapingBD(sa *StudentInformation, values *MappingColumn) {
	values.Result = sa.reflect(values.Key)
}

//reflect - reflexion information general encontrar valor
//Param ret		IN   "nombre de columna"
//Param res		OUT   "retorna el valor de la columna"
func (f *StudentInformation) reflect(ret string) interface{} {
	val := reflect.ValueOf(f).Elem()
	var res interface{}
	res = nil
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		if strings.Compare(typeField.Type.String(), "[]models.Economic") == 0 {
			res = fmt.Sprintf("%v", f.Informacioneconomica[0].reflectEcono(ret))
			if res != "<nil>" {
				break
			}
		} else {
			if strings.Compare(ret, typeField.Name) == 0 {
				res = fmt.Sprintf("%v", valueField.Interface())
				break
			}
		}
	}
	return res
}

//reflectEcono - reflexion information economica encontrar valor
//Param ret		IN   "nombre de columna"
//Param res		OUT   "retorna el valor de la columna"
func (f *Economic) reflectEcono(ret string) interface{} {
	val := reflect.ValueOf(f).Elem()
	var res interface{}
	res = nil
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

//RescueInformation - optiene la informacion basica de un estudiante para reportes
//Param sa		IN   "estudiante a consultar y retornar informacion"
func RescueInformation(sa *StudentInformation) {
	var wg sync.WaitGroup
	wg.Add(3)
	var ModelFacult XmlFaculty
	var ModelBasic XmlBasic
	var ModelAcademic XmlAcademic
	go utility.GetServiceXML(&ModelFacult, utility.FacultyService+sa.Codigo, &wg)
	go utility.GetServiceXML(&ModelBasic, utility.BasicService+sa.Codigo, &wg)
	go utility.GetServiceXML(&ModelAcademic, utility.AcademicService+sa.Codigo, &wg)
	wg.Wait()
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
}

//Evaluation  -evaluacion de reglas de negocio para reportes
//Param maping		IN   "clave de la columna"
//Param salario		IN   "salario minimo configurado por el administrador"
//Param i	OUT   "puntaje"
func Evaluation(maping *MappingColumn, salario int) int {
	i := 0
	switch con := maping.Key; con {
	case "Estrato":
		conv, _ := strconv.Atoi(maping.Result.(string))
		if conv <= 3 {
			i = 10
		}
	case "Matricula":
		conv, _ := strconv.Atoi(maping.Result.(string))
		if conv <= 200000 {
			i = 20
		} else if conv <= 400000 {
			i = 16
		} else if conv <= 600000 {
			i = 12
		} else if conv <= 800000 {
			i = 8
		} else if conv <= 900000 {
			i = 4
		}
	case "Ingresos":
		conv, _ := strconv.Atoi(maping.Result.(string))
		if conv <= salario {
			i = 30
		} else if conv <= salario*2 {
			i = 20
		} else if conv <= salario*3 {
			i = 10
		} else if conv <= salario*4 {
			i = 5
		} else {
			i = 0
		}
	case "SostePropia":
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 5
		}
	case "SosteHogar":
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 5
		}
	case "Nucleofam":
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 4
		}
	case "PersACargo":
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 6
		}
	case "EmpleadArriendo":
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 5
		}
	case "ProvBogota":
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 5
		}
	case "PobEspecial":
		if strings.Compare(maping.Result.(string), "D") == 0 || strings.Compare(maping.Result.(string), "I") == 0 || strings.Compare(maping.Result.(string), "M") == 0 || strings.Compare(maping.Result.(string), "A") == 0 || strings.Compare(maping.Result.(string), "MC") == 0 {
			i = 5
		}
	case "Discapacidad":
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 5
		}
	case "PatAlimenticia":
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 5
		}
	}

	return i
}

//ProcessinData  - mapeo de valores para reportes
//Param data		IN   "clave de la columna"
//Param temp	OUT   "valor a remplzar"
func ProcessinData(data *MappingColumn) string {
	var temp string
	temp = data.Result.(string)

	if strings.Compare(data.Key, "TipoSubsidio") == 0 {
		switch conv := data.Result.(string); conv {
		case "ss":
			temp = "SIN SUBSIDIO"
		case "a":
			temp = "TIPO A"
		case "b":
			temp = "TIPO B"
		case "t":
			temp = "SUBSIDIO TOTAL"
		}
	}
	if strings.Compare(data.Key, "PobEspecial") == 0 {
		switch conv := data.Result.(string); conv {
		case "N":
			temp = "NINGUNA"
		case "D":
			temp = "DESPLAZADO"
		case "I":
			temp = "INDIGENA"
		case "M":
			temp = "MINORIAS ETNICAS"
		case "A":
			temp = "AFRODESCENDIENTE"
		case "MC":
			temp = "MADRE CABEZA HOGAR"
		}
	}
	return temp
}

//GEtMappingColumn - optiene los valores de metadata de las columnas para generar reportes
//Param Global	OUT   "columnas"
func GEtMappingColumn() *[]MappingColumn {
	var Global []MappingColumn

	Global = append(Global, MakeThing("CODIGO", 1, "Codigo"))
	Global = append(Global, MakeThing("FECHA DE INSCRIPCION", 2, "Fechainscripcion"))
	Global = append(Global, MakeThingD("ESTRATO SOCIOECONÓMICO", 3, "Estrato", "Si"))
	Global = append(Global, MakeThingD("INGRESOS PROPIOS O FAMILIARES", 4, "Ingresos", "Si"))
	Global = append(Global, MakeThingD("SE SOSTIENE ECONÓMICAMENTE  A SÍ MISMO", 5, "SostePropia", "Si"))
	Global = append(Global, MakeThingD("SOSTIENE EL HOGAR EN QUE VIVE", 6, "SosteHogar", "Si"))
	Global = append(Global, MakeThingD("VIVE FUERA DE SU NÚCLEO FAMILIAR", 7, "Nucleofam", "Si"))
	Global = append(Global, MakeThingD("TIENE PERSONAS A CARGO", 8, "PersACargo", "Si"))
	Global = append(Global, MakeThingD("VIVE EN CASA DEL EMPLEADOR O PAGA ARRIENDO", 9, "EmpleadArriendo", "Si"))
	Global = append(Global, MakeThingD("PROVIENE DE CIUDADES O MUNICIPIOS DISTINTOS A BOGOTÁ", 10, "ProvBogota", "Si"))
	Global = append(Global, MakeThing("CIUDAD O MUNICIPIO", 11, "Ciudad"))
	Global = append(Global, MakeThingD("ESTÁ CERTIFICADO COMO POBLACIÓN ESPECIAL", 12, "PobEspecial", "Si"))
	Global = append(Global, MakeThingD("DISCAPACIDAD FÍSICA O MENTAL", 13, "Discapacidad", "Si"))
	Global = append(Global, MakeThingD("PATOLOGÍA ASOCIADA CON LA NUTRICIÓN", 14, "PatAlimenticia", "Si"))
	Global = append(Global, MakeThing("SER PILO PAGA", 15, "SerPiloPaga"))
	Global = append(Global, MakeThing("SISBEN", 16, "Sisben"))
	Global = append(Global, MakeThing("AÑO", 17, "Periodo"))
	Global = append(Global, MakeThing("SEMESTRE", 18, "SemestreIns"))
	Global = append(Global, MakeThingD("MATRICULA", 19, "Matricula", "Si"))
	Global = append(Global, MakeThing("TIPO DE SUBSIDIO", 20, "TipoSubsidio"))
	Global = append(Global, MakeThing("TIPO DE APOYO ALIMENTARIO", 21, "Tipoapoyo"))
	Global = append(Global, MakeThing("TELEFONO", 22, "Telefono"))
	Global = append(Global, MakeThing("CORREO", 23, "Correo"))
	Global = append(Global, MakeThing("ANTIGUEDAD PROGRAMA", 24, "Antiguedad"))
	/*columns of services*/
	Global = append(Global, MakeThing("APELLIDOS Y NOMBRES", 25, "Nombre"))
	Global = append(Global, MakeThing("LOCALIDAD", 26, "Localidad"))
	Global = append(Global, MakeThing("DIRECCION", 27, "Direccion"))
	Global = append(Global, MakeThing("TIPO DE DOCUMENTO", 28, "TDocument"))
	Global = append(Global, MakeThing("NUMERO DE DOCUMENTO", 29, "Document"))
	Global = append(Global, MakeThing("FACULTAD", 30, "Facultad"))
	Global = append(Global, MakeThing("PROYECTO CURRICULAR", 31, "Proyecto"))
	Global = append(Global, MakeThing("GENERO", 32, "Genero"))
	Global = append(Global, MakeThing("SEMESTRE", 33, "Semestre"))
	Global = append(Global, MakeThing("PROMEDIO", 34, "Promedio"))
	Global = append(Global, MakeThing("TOTAL PUNTAJE", 35, "Total"))
	return &Global
}
