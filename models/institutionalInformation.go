package models

import "encoding/xml"

//XmlFaculty Struct for get name of Faculty
type XmlFaculty struct {
	Collection  xml.Name `xml:"infoInstitucionalColleccion"`
	NameFaculty string   `xml:"infoInstitucional>facultad"`
}

//XmlMatricula Struct for get value of Enrollment
type XmlMatricula struct {
	Collection xml.Name `xml:"matriculaCollection"`
	Value      int      `xml:"matriculas>valor"`
}

//XmlEstado Struct for get state of student
type XmlEstado struct {
	Collection xml.Name `xml:"estadoCollection"`
	State      string   `xml:"estados>estado"`
}

//XmlBasic Struct for get name of student
type XmlBasic struct {
	Collection xml.Name `xml:"datosCollection"`
	Name       string   `xml:"datos>nombre"`
}
