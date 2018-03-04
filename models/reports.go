package models

import (
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
}

//MappingColumn Struct for map reports
type MappingColumn struct {
	Value  int
	Key    string
	Result interface{}
	Score  string
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

	err = file.Save("tempfile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

//ReportGeneral - Generate Reports students
func ReportGeneral(students []StudentInformation, sede string) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var Maping []MappingColumn
	var MapingNow []MappingColumn
	var count int
	var err error
	file = xlsx.NewFile()
	sheet, err = file.AddSheet(sede)
	var column []int
	column = append(column, 2, 26, 27, 34, 30, 31, 1, 33, 32, 3, 19, 4, 5, 6, 7, 8, 9, 10, 12, 13, 14, 37, 21, 29, 28, 25, 24, 35, 36)
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
		if strings.Compare(MapingNow[numdo].Score, "Si") == 0 {
			cell = row.AddCell()
			cell.Value = "Puntaje"
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
				cell.Value = MapingNow[numdo].Result.(string)
				if strings.Compare(MapingNow[numdo].Score, "Si") == 0 {
					localcount := 0
					cell = row.AddCell()
					localcount = Evaluation(MapingNow[numdo])
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
}

/* bonus functions*/

//MakeThing - Construc of model
func MakeThing(Val int, Keys string) MappingColumn {
	return MappingColumn{Val, Keys, "", ""}
}

//MakeThingD - Construc of model
func MakeThingD(Val int, Keys string, Score string) MappingColumn {
	return MappingColumn{Val, Keys, "", Score}
}

//GEtMappingColumn - Get values Metadata for reports
func GEtMappingColumn() []MappingColumn {
	var Global []MappingColumn

	Global = append(Global, MakeThing(1, "Codigo"))
	Global = append(Global, MakeThing(2, "Fechainscripcion"))
	Global = append(Global, MakeThingD(3, "Estrato", "Si"))
	Global = append(Global, MakeThingD(4, "Ingresos", "Si"))
	Global = append(Global, MakeThingD(5, "SostePropia", "Si"))
	Global = append(Global, MakeThingD(6, "SosteHogar", "Si"))
	Global = append(Global, MakeThingD(7, "Nucleofam", "Si"))
	Global = append(Global, MakeThingD(8, "PersACargo", "Si"))
	Global = append(Global, MakeThingD(9, "EmpleadArriendo", "Si"))
	Global = append(Global, MakeThingD(10, "ProvBogota", "Si"))
	Global = append(Global, MakeThing(11, "Ciudad"))
	Global = append(Global, MakeThingD(12, "PobEspecial", "Si"))
	Global = append(Global, MakeThingD(13, "Discapacidad", "Si"))
	Global = append(Global, MakeThingD(14, "PatAlimenticia", "Si"))
	Global = append(Global, MakeThing(15, "SerPiloPaga"))
	Global = append(Global, MakeThing(16, "Sisben"))
	Global = append(Global, MakeThing(17, "Periodo"))
	Global = append(Global, MakeThing(18, "Semestre"))
	Global = append(Global, MakeThingD(19, "Matricula", "Si"))
	Global = append(Global, MakeThing(20, "EstadoProg"))
	Global = append(Global, MakeThing(21, "TipoSubsidio"))
	Global = append(Global, MakeThing(22, "Tipoapoyo"))
	Global = append(Global, MakeThing(23, "Mensaje"))
	Global = append(Global, MakeThing(24, "Telefono"))
	Global = append(Global, MakeThing(25, "Correo"))
	Global = append(Global, MakeThing(26, "Antiguedad"))
	/*columns of services*/
	Global = append(Global, MakeThing(27, "Nombre"))
	Global = append(Global, MakeThing(28, "Localidad"))
	Global = append(Global, MakeThing(29, "Direccion"))
	Global = append(Global, MakeThing(30, "TDocument"))
	Global = append(Global, MakeThing(31, "Document"))
	Global = append(Global, MakeThing(32, "Facultad"))
	Global = append(Global, MakeThing(33, "Proyecto"))
	Global = append(Global, MakeThing(34, "Genero"))
	Global = append(Global, MakeThing(35, "Semestre"))
	Global = append(Global, MakeThing(36, "Promedio"))
	Global = append(Global, MakeThing(37, "Total"))
	return Global
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

//Evaluation  - evaluate bussines rules
func Evaluation(maping MappingColumn) int {
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

		switch conv, _ := strconv.Atoi(maping.Result.(string)); conv {
		case 1:
			i = 30
		case 2:
			i = 20
		case 3:
			i = 10
		case 4:
			i = 5
		default:
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
