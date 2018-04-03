package controllers

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/models"
	"apoyoalimentario_CRUD_API/utility"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/astaxie/beego"
)

//AdministratorController Operaciones Crud admin
type AdministratorController struct {
	beego.Controller
}

//GetStudents - optiene los estudiantes que se encuentran en un estado y una sede especifica
// @Title GetStudents
// @Description get Administrator by state and sede
// @Param	state		path 	string	true		"El estado del proceso a consultar"
// @Param	sede		path 	string	true		"la sede consultar"
// @Success 200 {object} models.Infoapoyo
// @Failure 403 :state is empty
// @router /:state/:sede [get]
func (j *AdministratorController) GetStudents() {
	var modellist models.ReportsType
	state := j.Ctx.Input.Param(":state")
	modellist.TSede = j.Ctx.Input.Param(":sede")
	modellist.Periodo = time.Now().UTC().Year()
	modellist.Semestre = utility.Semester()
	session, _ := db.GetSession()
	UserType, err := models.GetInscription(session, state, &modellist)
	defer session.Close()
	if err != nil {
		j.Data["json"] = err.Error()
	} else {
		j.Data["json"] = *UserType
	}

	j.ServeJSON()
}

//GetConfig - devuelve la configuracion del sistema
// @Title GetConfig
// @Description get configuration Administrator
// @Success 200 {string}
// @router / [get]
func (j *AdministratorController) GetConfig() {
	session, _ := db.GetSession()

	Configuration, err := models.GetConfiguration(session)

	if err != nil {
		j.Data["json"] = err.Error()
	} else {
		j.Data["json"] = *Configuration
	}

	j.ServeJSON()
}

//PutConfig - Actualiza la configuracion
// @Title PutConfig
// @Description update the configuration
// @Param	body		body 	models.Object	true		"la configuracion que se desea actualizar"
// @Success 200 {object} models.Object
// @router / [put]
func (j *AdministratorController) PutConfig() {

	var InfoConfig models.ConfigurationOptions
	resul := "Cambios guardados!"
	json.Unmarshal(j.Ctx.Input.RequestBody, &InfoConfig)
	session, _ := db.GetSession()

	erro := models.UpdateInformationConfig(session, &InfoConfig)

	if erro != nil {
		resul = erro.Error()
	}
	j.Data["json"] = resul
	j.ServeJSON()
}

//PutState - Actualiza el estado de un estudiante post verificacion
// @Title PutState
// @Description update the Infoapoyo
// @Param	code		path 	string	true		"codigo de estudiante a actualizar"
// @Success 200 {object} models.Object
// @Failure 403 :code is empty
// @router /verification/:code [put]
func (j *AdministratorController) PutState() {
	code := j.Ctx.Input.Param(":code")

	var InfoEcono models.Economic

	json.Unmarshal(j.Ctx.Input.RequestBody, &InfoEcono)
	session, _ := db.GetSession()
	err := models.UpdateStateVerificator(session, code, &InfoEcono)
	if err != nil {
		j.Data["json"] = err.Error()
	} else {
		j.Data["json"] = "ok"
	}
	defer session.Close()
	j.ServeJSON()
}

//Post - Genera reportes de estudiantes verificados correctamente segun un tipo ya especificado
// @Title GetReport
// @Description reports Administrator
// @Param	body		body 	"Opciones del reporte a generar"
// @Success 200 {object} models.Infoapoyo
// @Failure 403 body is empty
// @router /report [post]
func (j *AdministratorController) Post() {
	var modelReport models.ReportsType

	json.Unmarshal(j.Ctx.Input.RequestBody, &modelReport)

	session, _ := db.GetSession()

	//Get students in state 3 -> verified ok
	Students, err := models.GetInscription(session, "3", &modelReport)
	//Report Generic

	if err != nil {
		fmt.Printf(err.Error())
		j.Data["json"] = err.Error()
	} else {
		if modelReport.TypeReport == 1 { //Report Generic
			models.ReportsGeneric(Students, modelReport.NameSheet, &modelReport.Columnas)
		} else if modelReport.TypeReport == 2 { //Report Score final student
			models.ReportGeneral(session, Students, modelReport.NameSheet)
		} else if modelReport.TypeReport == 3 { //Sisben - ser pilo paga - Totales
			models.OthersReports(Students)
		}

		defer session.Close()
		archi, _ := ioutil.ReadFile("tempfile.xlsx")
		os.Remove("tempfile.xlsx")
		j.Ctx.Output.Body(archi)
	}
}

/*Verificador*/

//GetVerif - optiene la configuracion de la asociacion de verificadores y sedes
// @Title GetVerif
// @Description get configuration verifier
// @Success 200 {string}
// @router /verifier [get]
func (j *AdministratorController) GetVerif() {
	session, _ := db.GetSession()
	SedeVerif, err := models.GetVerifier(session)

	if err != nil {
		j.Data["json"] = err.Error()
	} else {
		j.Data["json"] = *SedeVerif
	}

	j.ServeJSON()
}

//PutVerif - Actualiza las asociaciones de verificadores y sedes
// @Title PutVerif
// @Description update the Infoapoyo
// @Success 200 {object} models.Object
// @router /verifier [put]
func (j *AdministratorController) PutVerif() {

	var InfoVerif []models.Sede

	json.Unmarshal(j.Ctx.Input.RequestBody, &InfoVerif)
	session, _ := db.GetSession()
	err := models.UpdateVerifier(session, &InfoVerif)
	if err != nil {
		j.Data["json"] = err.Error()
	} else {
		j.Data["json"] = "ok"
	}
	defer session.Close()
	j.ServeJSON()
}

//Getsede - muestra la(s) sede(s) las cuales tiene asignado un verificador
// @Title Getsede
// @Description  get sedes of verifier
// @Param	name		path 	string	true		"nombre del verificador a consultar"
// @Success 200 {string}
// @Failure 403 :name is empty
// @router /verifier/:name [get]
func (j *AdministratorController) Getsede() {
	name := j.Ctx.Input.Param(":name")
	session, _ := db.GetSession()

	SedeVerif, err := models.GetSede(session, name)

	if err != nil {
		j.Data["json"] = err.Error()
	} else {
		j.Data["json"] = *SedeVerif
	}

	j.ServeJSON()
}
