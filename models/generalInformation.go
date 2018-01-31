package models

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/utility"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//StudentInformation struct of information fgeneral of student
type StudentInformation struct {
	ID               bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Codigo           string        `json:"codigo" bson:"codigo"`
	Fechainscripcion time.Time     `json:"ultimafechainscripcion"  bson:"ultimafechainscripcion"`
	EstadoProg       int           `json:"estadoprograma" bson:"estadoprograma"`
	Nombre           string        `json:",omitempty" bson:",omitempty"`
}

//GetStatus - get status current of student
func GetStatus(session *mgo.Session, code string) (state int) {

	var StateUniversity XmlEstado
	utility.GetServiceXML(&StateUniversity, utility.StateService+code)
	ValidateAdministator := db.Cursor(session, utility.CollectionAdministrator)
	var ModuleActive ConfigurationOptions
	err := ValidateAdministator.Find(nil).Select(bson.M{"moduloactivo": 1}).One(&ModuleActive)
	if strings.Compare(StateUniversity.State, "ACTIVO") == 0 && ModuleActive.Moduloactivo == true {
		MainSession := db.Cursor(session, utility.CollectionGeneral)
		EconomicSession := db.Cursor(session, utility.CollectionEconomic)
		var InfoGeneral StudentInformation
		var InfoEconomic Economic
		err = MainSession.Find(bson.M{"codigo": code}).One(&InfoGeneral)
		if err != nil {
			InfoGeneral = TemplatenewP(InfoGeneral, code)
			MainSession.Insert(InfoGeneral)
			_ = MainSession.Find(bson.M{"codigo": code}).One(&InfoGeneral)
			InfoEconomic = TemplatenewEcon(InfoEconomic, InfoGeneral.ID.Hex(), code)
			EconomicSession.Insert(InfoEconomic)
			err = nil
		}
		err = EconomicSession.Find(bson.M{"id": InfoGeneral.ID.Hex(), "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}).One(&InfoEconomic)
		if err != nil {
			InfoEconomic = TemplatenewEcon(InfoEconomic, InfoGeneral.ID.Hex(), code)
			EconomicSession.Insert(InfoEconomic)
			err = EconomicSession.Find(bson.M{"id": InfoGeneral.ID.Hex(), "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}).One(&InfoEconomic)
		}
		if InfoGeneral.EstadoProg == 0 {
			var FacultadName XmlFaculty
			utility.GetServiceXML(&FacultadName, utility.FacultyService+code)
			aa, _ := ValidateAdministator.Find(bson.M{"refrigerionocturno": bson.M{"$in": []string{FacultadName.NameFaculty}}}).Count()
			if aa == 1 {
				state = 2 //ACCES OK almuerzo y refrigerio
			} else {
				state = 1 //ACCES OK almuerzo
			}
		} else {
			state = -1 //INSCRIPTION EXIT
		}
		if err != nil {
			state = -2 //OTHER ERROR EXIT
		}
	} else {
		state = 0 //INACTIVE USER OR MODULE off --> EXIT
	}
	return state
}

//UpdateState - update state in schedule of student
func UpdateState(session *mgo.Session, cod string) error {
	MainSession := db.Cursor(session, utility.CollectionGeneral)
	EconomicSession := db.Cursor(session, utility.CollectionEconomic)
	var InfoGeneralU StudentInformation
	var InfoEcoOldU Economic
	errd := MainSession.Find(bson.M{"codigo": cod}).One(&InfoGeneralU)
	UpdateDate := LastDate(InfoGeneralU)
	errd = MainSession.Update(bson.M{"codigo": cod}, &UpdateDate)

	err := EconomicSession.Find(bson.M{"id": InfoGeneralU.ID.Hex(), "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}).One(&InfoEcoOldU)
	UpdateS := LastState(InfoEcoOldU)
	err = EconomicSession.Update(bson.M{"id": InfoGeneralU.ID.Hex(), "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}, &UpdateS)
	if err != nil {
		panic(errd)
	}
	return err
}

/* functions Bonus*/

//TemplatenewP - create new template for New students
func TemplatenewP(j StudentInformation, cod string) StudentInformation {

	j.Codigo = cod
	j.EstadoProg = 0
	return j
}

//TemplatenewEcon - create new template for the economic information of students
func TemplatenewEcon(j Economic, id string, cod string) Economic {
	var v XmlMatricula
	utility.GetServiceXML(&v, utility.EnrollmentService+cod)

	j.ID = bson.NewObjectId()
	j.Idc = id
	j.Periodo = time.Now().UTC().Year()
	j.Semestre = utility.Semester()
	j.Matricula = v.Value
	j.TipoSubsidio = "na"
	return j
}

//LastDate - Update date of inscription
func LastDate(old StudentInformation) StudentInformation {

	old.Fechainscripcion = time.Now().UTC()
	old.EstadoProg = 1
	return old
}

//LastState - Update Information economic empty
func LastState(old Economic) Economic {

	if strings.Compare(old.Ciudad, "") == 0 {
		old.Ciudad = "bogota"
	}
	if strings.Compare(old.Tipoapoyo, "") == 0 {
		old.Tipoapoyo = "A"
	}
	return old
}
