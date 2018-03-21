package models

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/utility"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
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

//ReportsGeneric - Generate GENERIC Reports dynamic
func ReportsGeneric(sa []StudentInformation, NameSheet string, column []int) {
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
		cell.Value = MapingNow[numdo].ColumnName
	}
	// cell document

	for fil := range sa {
		row = sheet.AddRow()
		sa[fil] = RescueInformation(sa[fil])
		for numdo := range MapingNow {
			cell = row.AddCell()
			MapingNow[numdo] = MapingBD(sa[fil], MapingNow[numdo])
			if MapingNow[numdo].Result != nil {
				var temp string
				temp = ProcessinData(MapingNow[numdo])
				cell.Value = temp
			}
		}
	}

	err = file.Save("tempfile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

//ReportGeneral - Generate Reports students
func ReportGeneral(students []StudentInformation, name string) {
	session, _ := db.GetSession()
	BDSMLV := db.Cursor(session, utility.CollectionAdministrator)
	var salario ConfigurationOptions
	err := BDSMLV.Find(nil).One(&salario)
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var Maping []MappingColumn
	var MapingNow []MappingColumn
	var count int
	file = xlsx.NewFile()
	sheet, err = file.AddSheet(name)
	var column []int
	column = append(column, 2, 24, 25, 32, 28, 29, 1, 31, 30, 3, 19, 4, 5, 6, 7, 8, 9, 10, 12, 13, 14, 35, 20, 21, 27, 26, 23, 22, 33, 34, 15, 16)
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
		cell.Value = MapingNow[numdo].ColumnName
		if strings.Compare(MapingNow[numdo].Score, "Si") == 0 {
			cell = row.AddCell()
			cell.Value = "PUNTAJE " + MapingNow[numdo].ColumnName
		}
	}
	// cell document

	for fil := range students {
		count = 0
		row = sheet.AddRow()
		students[fil] = RescueInformation(students[fil])

		for numdo := range MapingNow {
			cell = row.AddCell()
			MapingNow[numdo] = MapingBD(students[fil], MapingNow[numdo])
			if MapingNow[numdo].Result != nil {
				var temp string
				temp = ProcessinData(MapingNow[numdo])
				cell.Value = temp
				if strings.Compare(MapingNow[numdo].Score, "Si") == 0 {
					localcount := 0
					cell = row.AddCell()
					localcount = Evaluation(MapingNow[numdo], salario.Salariominimo)
					count += localcount
					cell.Value = strconv.Itoa(localcount)
				}

			}
			if strings.Compare(MapingNow[numdo].Key, "Total") == 0 {
				cell.Value = strconv.Itoa(count)
			}
		}
	}

	err = file.Save("tempfile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
	defer session.Close()
}

//OthersReports - Generate Reports students
func OthersReports(students []StudentInformation) {

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var Maping []MappingColumn
	var MapingNow []MappingColumn
	file = xlsx.NewFile()

	//sheet 1 -- sisben
	sheet, err := file.AddSheet("Sisben")
	var column []int
	column = append(column, 1, 25, 31, 30)
	Maping = GEtMappingColumn()
	for numuno := range column {
		for num := range Maping {
			if Maping[num].Value == column[numuno] {
				MapingNow = append(MapingNow, Maping[num])
				break
			}
		}
	}
	//name Columns
	var cell *xlsx.Cell
	row = sheet.AddRow()
	for numdo := range MapingNow {
		cell = row.AddCell()
		cell.Value = MapingNow[numdo].ColumnName
	}

	// cell document
	for fil := range students {
		if strings.Compare(students[fil].Informacioneconomica[0].Sisben, "no") == 0 {
			continue
		}
		row = sheet.AddRow()
		students[fil] = RescueInformation(students[fil])
		for numdo := range MapingNow {
			cell = row.AddCell()
			MapingNow[numdo] = MapingBD(students[fil], MapingNow[numdo])
			if MapingNow[numdo].Result != nil {
				cell.Value = MapingNow[numdo].Result.(string)
			}
		}
	}
	//sheet 2 -- Ser Pilo Paga
	sheet, err = file.AddSheet("Ser Pilo Paga")

	//name Columns
	var celldos *xlsx.Cell
	row = sheet.AddRow()
	for numdo := range MapingNow {
		celldos = row.AddCell()
		celldos.Value = MapingNow[numdo].ColumnName
	}

	// cell document
	for fil := range students {
		if strings.Compare(students[fil].Informacioneconomica[0].SerPiloPaga, "no") == 0 {
			continue
		}
		row = sheet.AddRow()
		students[fil] = RescueInformation(students[fil])
		for numdo := range MapingNow {
			celldos = row.AddCell()
			MapingNow[numdo] = MapingBD(students[fil], MapingNow[numdo])
			if MapingNow[numdo].Result != nil {
				celldos.Value = MapingNow[numdo].Result.(string)
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
	var TSIN int
	var TA int
	var TB int
	var TT int
	// cell document
	for fil := range students {
		if strings.Compare(students[fil].Informacioneconomica[0].TipoSubsidio, "t") == 0 {
			TT++
		} else if strings.Compare(students[fil].Informacioneconomica[0].TipoSubsidio, "a") == 0 {
			TA++
		} else if strings.Compare(students[fil].Informacioneconomica[0].TipoSubsidio, "b") == 0 {
			TB++
		} else if strings.Compare(students[fil].Informacioneconomica[0].TipoSubsidio, "ss") == 0 {
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

//MakeThing - Construc of model
func MakeThing(Col string, Val int, Keys string) MappingColumn {
	return MappingColumn{Col, Val, Keys, "", ""}
}

//MakeThingD - Construc of model
func MakeThingD(Col string, Val int, Keys string, Score string) MappingColumn {
	return MappingColumn{Col, Val, Keys, "", Score}
}

//GEtMappingColumn - Get values Metadata for reports
func GEtMappingColumn() []MappingColumn {
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
	return Global
}
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

//Evaluation  - evaluate bussines rules
func Evaluation(maping MappingColumn, salario int) int {
	i := 0
	if strings.Compare(maping.Key, "Estrato") == 0 {
		conv, _ := strconv.Atoi(maping.Result.(string))
		if conv <= 3 {
			i = 10
		}
	}
	if strings.Compare(maping.Key, "Matricula") == 0 {
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
	}
	if strings.Compare(maping.Key, "Ingresos") == 0 {
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

	}
	if strings.Compare(maping.Key, "SostePropia") == 0 {
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 5
		}
	}
	if strings.Compare(maping.Key, "SosteHogar") == 0 {
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 5
		}
	}
	if strings.Compare(maping.Key, "Nucleofam") == 0 {
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 4
		}
	}
	if strings.Compare(maping.Key, "PersACargo") == 0 {
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 6
		}
	}
	if strings.Compare(maping.Key, "EmpleadArriendo") == 0 {
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 5
		}
	}
	if strings.Compare(maping.Key, "ProvBogota") == 0 {
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 5
		}
	}
	if strings.Compare(maping.Key, "PobEspecial") == 0 {
		if strings.Compare(maping.Result.(string), "D") == 0 || strings.Compare(maping.Result.(string), "I") == 0 || strings.Compare(maping.Result.(string), "M") == 0 || strings.Compare(maping.Result.(string), "A") == 0 || strings.Compare(maping.Result.(string), "MC") == 0 {
			i = 5
		}
	}
	if strings.Compare(maping.Key, "Discapacidad") == 0 {
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 5
		}
	}
	if strings.Compare(maping.Key, "PatAlimenticia") == 0 {
		if strings.Compare(maping.Result.(string), "si") == 0 {
			i = 5
		}
	}
	return i
}

//ProcessinData  - evaluate bussines rules
func ProcessinData(data MappingColumn) string {
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
