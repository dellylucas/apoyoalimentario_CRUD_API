package models

//ReportsType Struct for map reports
type ReportsType struct {
	Columnas  []int  `json:"columnas" bson:"columnas"`
	TSede     string `json:"tsede" bson:"tsede"`
	NameSheet string `json:"nameSheet" bson:"nameSheet"`
}

//MappingColumn Struct for map reports
type MappingColumn struct {
	Value  int
	Key    string
	Result interface{}
}

//MakeThing - Construc of model
func MakeThing(Val int, Keys string) MappingColumn {
	return MappingColumn{Val, Keys, ""}
}

//GEtMappingColumn - Get values Metadata for reports
func GEtMappingColumn() []MappingColumn {
	var Global []MappingColumn

	Global = append(Global, MakeThing(1, "Codigo"))
	Global = append(Global, MakeThing(2, "Fechainscripcion"))
	Global = append(Global, MakeThing(3, "Estrato"))
	Global = append(Global, MakeThing(4, "Ingresos"))
	Global = append(Global, MakeThing(5, "SostePropia"))
	Global = append(Global, MakeThing(6, "SosteHogar"))
	Global = append(Global, MakeThing(7, "Nucleofam"))
	Global = append(Global, MakeThing(8, "PersACargo"))
	Global = append(Global, MakeThing(9, "EmpleadArriendo"))
	Global = append(Global, MakeThing(10, "ProvBogota"))
	Global = append(Global, MakeThing(11, "Ciudad"))
	Global = append(Global, MakeThing(12, "PobEspecial"))
	Global = append(Global, MakeThing(13, "Discapacidad"))
	Global = append(Global, MakeThing(14, "PatAlimenticia"))
	Global = append(Global, MakeThing(15, "SerPiloPaga"))
	Global = append(Global, MakeThing(16, "Sisben"))
	Global = append(Global, MakeThing(17, "Periodo"))
	Global = append(Global, MakeThing(18, "Semestre"))
	Global = append(Global, MakeThing(19, "Matricula"))
	Global = append(Global, MakeThing(20, "EstadoProg"))
	Global = append(Global, MakeThing(21, "TipoSubsidio"))
	Global = append(Global, MakeThing(22, "Tipoapoyo"))
	Global = append(Global, MakeThing(23, "Mensaje"))
	Global = append(Global, MakeThing(23, "Telefono"))
	Global = append(Global, MakeThing(24, "Correo"))
	Global = append(Global, MakeThing(25, "Antiguedad"))
	/*columns of services*/
	Global = append(Global, MakeThing(26, "Nombre"))
	Global = append(Global, MakeThing(27, "Localidad"))
	Global = append(Global, MakeThing(28, "Direccion"))
	Global = append(Global, MakeThing(29, "TDocument"))
	Global = append(Global, MakeThing(30, "Document"))
	Global = append(Global, MakeThing(31, "Facultad"))
	Global = append(Global, MakeThing(32, "Proyecto"))
	Global = append(Global, MakeThing(33, "Genero"))
	Global = append(Global, MakeThing(34, "Semestre"))
	Global = append(Global, MakeThing(35, "Promedio"))
	return Global
}
