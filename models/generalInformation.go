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
	ID                   bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Codigo               string        `json:"codigo" bson:"codigo"`
	Fechainscripcion     time.Time     `json:"ultimafechainscripcion"  bson:"ultimafechainscripcion"`
	Nombre               string        `json:",omitempty" bson:",omitempty"`
	Informacioneconomica []Economic    `json:",omitempty" bson:",omitempty"`
	Localidad            string        `json:",omitempty" bson:",omitempty"`
	Direccion            string        `json:",omitempty" bson:",omitempty"`
	TDocument            string        `json:",omitempty" bson:",omitempty"`
	Document             string        `json:",omitempty" bson:",omitempty"`
	Facultad             string        `json:",omitempty" bson:",omitempty"`
	Proyecto             string        `json:",omitempty" bson:",omitempty"`
	Genero               string        `json:",omitempty" bson:",omitempty"`
	Semestre             string        `json:",omitempty" bson:",omitempty"`
	Promedio             string        `json:",omitempty" bson:",omitempty"`
}

//GetStatus - get status current of student
func GetStatus(session *mgo.Session, code string) (state int) {

	var StateUniversity XmlEstado
	utility.GetServiceXML(&StateUniversity, utility.StateService+code)
	ValidateAdministator := db.Cursor(session, utility.CollectionAdministrator)
	var ModuleActive ConfigurationOptions
	err := ValidateAdministator.Find(nil).One(&ModuleActive)
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
			InfoEconomic = TemplatenewEcon(InfoEconomic, InfoGeneral.ID, code)
			EconomicSession.Insert(InfoEconomic)
			err = nil
		}
		err = EconomicSession.Find(bson.M{"id": InfoGeneral.ID, "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}).One(&InfoEconomic)
		if err != nil {
			InfoEconomic = TemplatenewEcon(InfoEconomic, InfoGeneral.ID, code)
			EconomicSession.Insert(InfoEconomic)
			err = EconomicSession.Find(bson.M{"id": InfoGeneral.ID, "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}).One(&InfoEconomic)
		}
		if InfoEconomic.EstadoProg == 0 || InfoEconomic.EstadoProg == 4 {
			var FacultadName XmlFaculty
			count := 0
			utility.GetServiceXML(&FacultadName, utility.FacultyService+code)
			for _, element := range ModuleActive.Refrigerionocturno {
				if element == FacultadName.NameFaculty {
					count = 1
				}
			}
			if count == 1 {
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
	} else if strings.Compare(StateUniversity.State, "ACTIVO") == 0 && ModuleActive.Modulomodified == true {
		MainSession := db.Cursor(session, utility.CollectionGeneral)
		EconomicSession := db.Cursor(session, utility.CollectionEconomic)
		var InfoGeneral StudentInformation
		var InfoEconomic Economic
		err = MainSession.Find(bson.M{"codigo": code}).One(&InfoGeneral)
		if err != nil {
			state = 0 //OTHER ERROR EXIT
		} else {
			err = EconomicSession.Find(bson.M{"id": InfoGeneral.ID, "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}).One(&InfoEconomic)
			if err != nil {
				state = 0 //OTHER ERROR EXIT
			} else {
				if InfoEconomic.EstadoProg == 4 {
					state = 3 // ONLY ModificatioN IN
				} else {
					state = -1 //INSCRIPTION EXIT
				}
			}

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

	err := EconomicSession.Find(bson.M{"id": InfoGeneralU.ID, "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}).One(&InfoEcoOldU)
	UpdateS := LastState(InfoEcoOldU)
	err = EconomicSession.Update(bson.M{"id": InfoGeneralU.ID, "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}, &UpdateS)
	if err != nil {
		panic(errd)
	}
	return err
}

/* functions Bonus*/

//TemplatenewP - create new template for New students
func TemplatenewP(j StudentInformation, cod string) StudentInformation {

	j.Codigo = cod
	return j
}

//TemplatenewEcon - create new template for the economic information of students
func TemplatenewEcon(j Economic, id bson.ObjectId, cod string) Economic {
	var v XmlMatricula
	utility.GetServiceXML(&v, utility.EnrollmentService+cod)

	j.EstadoProg = 0
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

	return old
}

//LastState - Update Information economic empty
func LastState(old Economic) Economic {
	old.EstadoProg = 1
	if strings.Compare(old.Ciudad, "") == 0 {
		old.Ciudad = "bogota"
	}
	if strings.Compare(old.Tipoapoyo, "") == 0 {
		old.Tipoapoyo = "A"
	}
	return old
}
