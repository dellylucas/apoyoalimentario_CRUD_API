package models

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/utility"
	"strconv"
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

//GetStatus - get status current of student in system
func GetStatus(session *mgo.Session, code string) (state int) {

	var StateUniversity XmlEstado
	var ModuleActive ConfigurationOptions

	/*obtiene el estado de un estudiante ACTIVO O INACTIVO*/
	utility.GetServiceXML(&StateUniversity, utility.StateService+code)

	/*Obtiene el estado del modulo configurado por el Administrador*/
	ValidateAdministator := db.Cursor(session, utility.CollectionAdministrator)
	err := ValidateAdministator.Find(nil).One(&ModuleActive)

	if strings.Compare(StateUniversity.State, "ACTIVO") == 0 && ModuleActive.Moduloactivo == true {

		var InfoGeneral StudentInformation
		var InfoEconomic Economic
		MainSession := db.Cursor(session, utility.CollectionGeneral)
		EconomicSession := db.Cursor(session, utility.CollectionEconomic)
		/*Encuentra la informacion general de un estudiante en la BD*/
		err = MainSession.Find(bson.M{"codigo": code}).One(&InfoGeneral)

		/*Si no existe la crea con una plantilla por defecto*/
		if err != nil {
			InfoGeneral.Codigo = code
			InfoGeneral.ID = bson.NewObjectId()
			MainSession.Insert(InfoGeneral)
			InfoEconomic = TemplatenewEcon(InfoEconomic, InfoGeneral.ID, code)
			EconomicSession.Insert(InfoEconomic)
			err = nil
		}
		/*Encuentra la informacion economica de un estudiante en la BD*/
		err = EconomicSession.Find(bson.M{"id": InfoGeneral.ID, "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}).One(&InfoEconomic)

		/*Si no existe la crea con una plantilla por defecto*/
		if err != nil {
			InfoEconomic = TemplatenewEcon(InfoEconomic, InfoGeneral.ID, code)
			EconomicSession.Insert(InfoEconomic)
		}

		/*Ingreso para realizar la inscripcion
		0 el estudiante es nuevo
		4 puede modificar despues de una revision de un verificador
		*/
		if InfoEconomic.EstadoProg == 0 || InfoEconomic.EstadoProg == 4 {
			var FacultadName XmlFaculty
			count := 0
			/*Obtiene la facultad del estudiante*/
			utility.GetServiceXML(&FacultadName, utility.FacultyService+code)
			sedeEstudiante := strings.Replace(FacultadName.NameFaculty, "/", "-", -1)

			/*Iteracion de la configuracion de las sedes las cuales tienen refrigerio  nocturno configuradas por el administrador*/
			for _, sederefrigerio := range ModuleActive.Refrigerionocturno {
				if strings.Compare(sederefrigerio, sedeEstudiante) == 0 {
					count = 1
				}
			}
			//Estudiante puede escoger entre refrigerio o almuerzo
			if count == 1 {
				state = 2
			} else { //Estudiante solo puede escoger almuerzo
				state = 1
			}
		} else { //Estudiante solo puede consultar su informacion no modificar
			state = -1
		}
		//si es estudiante activo y esta habilitado el modulo de modificacion
	} else if strings.Compare(StateUniversity.State, "ACTIVO") == 0 && ModuleActive.Modulomodified == true {
		var InfoGeneral StudentInformation
		var InfoEconomic Economic
		MainSession := db.Cursor(session, utility.CollectionGeneral)
		EconomicSession := db.Cursor(session, utility.CollectionEconomic)

		err = MainSession.Find(bson.M{"codigo": code}).One(&InfoGeneral)
		if err != nil { //Estudiante no esta en BD y fuera de fechas de inscripcion
			state = 0
		} else {
			err = EconomicSession.Find(bson.M{"id": InfoGeneral.ID, "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}).One(&InfoEconomic)
			if err != nil { //Estudiante no esta en BD y fuera de fechas de inscripcion
				state = 0
			} else {
				//Si el estudiante fue calificado por un verificador y debe realizar modificaciones es estado 4
				if InfoEconomic.EstadoProg == 4 {
					var FacultadName XmlFaculty
					count := 0
					/*Obtiene la facultad del estudiante*/
					utility.GetServiceXML(&FacultadName, utility.FacultyService+code)
					sedeEstudiante := strings.Replace(FacultadName.NameFaculty, "/", "-", -1)

					/*Iteracion de la configuracion de las sedes las cuales tienen refrigerio  nocturno configuradas por el administrador*/
					for _, sederefrigerio := range ModuleActive.Refrigerionocturno {
						if strings.Compare(sederefrigerio, sedeEstudiante) == 0 {
							count = 1
						}
					}
					//Estudiante puede escoger entre refrigerio o almuerzo
					if count == 1 {
						state = 2
					} else { //Estudiante solo puede escoger almuerzo
						state = 1
					}
				} else { //Solo lectura de la inscripcion ya hecha
					state = -1
				}
			}
		}
	} else { //si el estudiante esta Inactivo y/o modulo de inscripcion y modificacion desactivados
		state = 0
	}
	return state
}

//UpdateState - update state in schedule of student
func UpdateState(session *mgo.Session, cod string) error {
	MainSession := db.Cursor(session, utility.CollectionGeneral)
	EconomicSession := db.Cursor(session, utility.CollectionEconomic)
	BDSMLV := db.Cursor(session, utility.CollectionAdministrator)

	var salario ConfigurationOptions
	err := BDSMLV.Find(nil).One(&salario)
	var InfoGeneralU StudentInformation
	var InfoEcoOldU Economic
	errd := MainSession.Find(bson.M{"codigo": cod}).One(&InfoGeneralU)
	UpdateDate := LastDate(InfoGeneralU)
	errd = MainSession.Update(bson.M{"codigo": cod}, &UpdateDate)

	err = EconomicSession.Find(bson.M{"id": InfoGeneralU.ID, "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}).One(&InfoEcoOldU)
	var ResultRuler string
	UpdateS := LastState(InfoEcoOldU)
	UpdateS.Salario = strconv.Itoa(salario.Salariominimo)
	ResultRuler, _ = utility.SendJSONToRuler(utility.RulerPath, "PUT", UpdateS)
	UpdateS = PostRules(UpdateS, ResultRuler)
	UpdateS.Salario = ""
	err = EconomicSession.Update(bson.M{"id": InfoGeneralU.ID, "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}, &UpdateS)
	if err != nil {
		panic(errd)
	}

	return err
}

/* functions Bonus*/

//TemplatenewEcon - create new template for the economic information of students
func TemplatenewEcon(j Economic, id bson.ObjectId, code string) Economic {
	var v XmlMatricula
	utility.GetServiceXML(&v, utility.EnrollmentService+code)

	j.EstadoProg = 0
	j.ID = bson.NewObjectId()
	j.Idc = id
	j.Periodo = time.Now().UTC().Year() /*año actual de inscripcion*/
	j.SemestreIns = utility.Semester()
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
	if strings.Compare(old.Tipoapoyo, "") == 0 || strings.Compare(old.Tipoapoyo, "A") == 0 {
		old.Tipoapoyo = "Almuerzo"
	}
	return old
}

//PostRules - Update Information economic empty
func PostRules(old Economic, Ruler string) Economic {
	if strings.Compare(Ruler, "") == 0 {
		old.EstadoProg = 0
	} else {
		old.EstadoProg = 2
	}
	old.TipoSubsidio = Ruler

	return old
}
