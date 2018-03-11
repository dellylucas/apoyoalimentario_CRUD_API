package controllers

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/models"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/astaxie/beego"
)

//AdministratorController Operaciones Crud admin
type AdministratorController struct {
	beego.Controller
}

//GetStudents - get Administrator by state
// @Title GetStudents
// @Description get Administrator by state
// @Param	state		path 	string	true		"El estado del proceso a consultar"
// @Param	sede		path 	string	true		"El estado del proceso a consultar"
// @Success 200 {object} models.Infoapoyo
// @Failure 403 :state is empty
// @router /:state/:sede [get]
func (j *AdministratorController) GetStudents() {
	state := j.Ctx.Input.Param(":state")
	sedeChecker := j.Ctx.Input.Param(":sede")
	session, _ := db.GetSession()
	UserType, err := models.GetInscription(session, state, sedeChecker)
	if err != nil {
		j.Data["json"] = err.Error()
	} else {
		j.Data["json"] = UserType
	}

	j.ServeJSON()
}

//GetConfig - get configuration Administrator
// @Title GetConfig
// @Description get configuration Administrator
// @Success 200 {string}
// @router / [get]
func (j *AdministratorController) GetConfig() {
	session, _ := db.GetSession()

	Message, err := models.GetConfiguration(session)

	if err != nil {
		j.Data["json"] = err.Error()
	} else {
		j.Data["json"] = Message
	}

	j.ServeJSON()
}

//PutConfig - update the configuration
// @Title PutConfig
// @Description update the configuration
// @Param	body		body 	models.Object	true		"The body"
// @Success 200 {object} models.Object
// @router / [put]
func (j *AdministratorController) PutConfig() {

	var InfoConfig models.ConfigurationOptions
	resul := "Cambios guardados!"
	json.Unmarshal(j.Ctx.Input.RequestBody, &InfoConfig)
	session, _ := db.GetSession()

	erro := models.UpdateInformationConfig(session, InfoConfig)

	if erro != nil {
		resul = erro.Error()
	}
	j.Data["json"] = resul
	defer session.Close()
	j.ServeJSON()
}

//PutState - update the Infoapoyo
// @Title PutState
// @Description update the Infoapoyo
// @Param	code		path 	string	true		"The code you want to update"
// @Success 200 {object} models.Object
// @Failure 403 :code is empty
// @router /verification/:code [put]
func (j *AdministratorController) PutState() {
	code := j.Ctx.Input.Param(":code")

	var InfoEcono models.Economic

	json.Unmarshal(j.Ctx.Input.RequestBody, &InfoEcono)
	session, _ := db.GetSession()
	err := models.UpdateStateVerificator(session, code, InfoEcono)
	if err != nil {
		j.Data["json"] = err.Error()
	} else {
		j.Data["json"] = "ok"
	}
	defer session.Close()
	j.ServeJSON()
}

//Post - reports Administrator
// @Title GetReport
// @Description reports Administrator
// @Param	body		body 	"body for File content"
// @Success 200 {object} models.Infoapoyo
// @Failure 403 body is empty
// @router /report [post]
func (j *AdministratorController) Post() {
	var modelReport models.ReportsType

	json.Unmarshal(j.Ctx.Input.RequestBody, &modelReport)

	session, _ := db.GetSession()

	UserType, err := models.GetInscription(session, "3", modelReport.TSede)
	//Report Generic
	if modelReport.TypeReport == 1 {
		models.ReportsGeneric(UserType, modelReport.NameSheet, modelReport.Columnas)
	} else if modelReport.TypeReport == 2 {
		//Report Score
		models.ReportGeneral(UserType, modelReport.NameSheet)
	}

	archi, _ := ioutil.ReadFile("tempfile.xlsx")
	//os.Remove("tempfile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
		j.Data["json"] = err.Error()
	} else {
		j.Ctx.Output.Body(archi)
	}
}

/*Verificador*/

//GetVerif - get configuration Administrator
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
		j.Data["json"] = SedeVerif
	}

	j.ServeJSON()
}

//PutVerif - update the verifier
// @Title PutVerif
// @Description update the Infoapoyo
// @Param	code		path 	string	true		"The code you want to update"
// @Success 200 {object} models.Object
// @router /verifier [put]
func (j *AdministratorController) PutVerif() {

	var InfoVerif []models.Sede

	json.Unmarshal(j.Ctx.Input.RequestBody, &InfoVerif)
	session, _ := db.GetSession()
	err := models.UpdateVerifier(session, InfoVerif)
	if err != nil {
		j.Data["json"] = err.Error()
	} else {
		j.Data["json"] = "ok"
	}
	defer session.Close()
	j.ServeJSON()
}

//Getsede - get configuration Administrator
// @Title Getsede
// @Description get configuration verifier
// @Param	name		path 	string	true		"El estado del proceso a consultar"
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
		j.Data["json"] = SedeVerif
	}

	j.ServeJSON()
}
