package controllers

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/models"
	"encoding/json"

	"github.com/astaxie/beego"
)

//AdministratorController Operaciones Crud admin
type AdministratorController struct {
	beego.Controller
}

// @Title Get
// @Description get Administrator by user
// @Param	user		path 	string	true		"El estado del proceso a consultar"
// @Success 200 {object} models.Infoapoyo
// @Failure 403 :user is empty
// @router /:user [get]
func (j *AdministratorController) Get() {
	user := j.Ctx.Input.Param(":user")
	session, _ := db.GetSession()

	UserType, err := models.GetTypeUser(session, user)

	if err != nil {
		j.Data["json"] = err.Error()
	} else {
		j.Data["json"] = UserType
	}

	j.ServeJSON()
}

// @Title Get
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

// @Title Get
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
