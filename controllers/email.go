package controllers

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/models"
	"encoding/json"

	"github.com/astaxie/beego"
)

//EmailController Operaciones Crud email send
type EmailController struct {
	beego.Controller
}

//GetConfig - get configuration of conection
// @Title GetConfig
// @Description get Administrator by state
// @Success 200 {object} models.Email
// @router / [get]
func (j *EmailController) GetConfig() {

	session, _ := db.GetSession()
	Configuration, err := models.SearchInfor(session)
	defer session.Close()
	if err != nil {
		j.Data["json"] = err.Error()
	} else {
		j.Data["json"] = *Configuration
	}

	j.ServeJSON()
}

//PutConfig - update the configuration
// @Title PutConfig
// @Description update the configuration
// @Param	body		body 	models.Object	true		"The body"
// @Success 200 {object} models.Object
// @router / [put]
func (j *EmailController) PutConfig() {

	var InfoConfig models.Email
	resul := "Cambios guardados!"
	json.Unmarshal(j.Ctx.Input.RequestBody, &InfoConfig)
	session, _ := db.GetSession()

	erro := models.UpdateEmailConfig(session, &InfoConfig)

	if erro != nil {
		resul = erro.Error()
	}
	j.Data["json"] = resul
	j.ServeJSON()
}

//PutEmail - email send
// @Title PutEmail
// @Description update the Infoapoyo
// @Param	code		path 	string	true		"The code you want to update"
// @Success 200 {object} models.Object
// @router /send [put]
func (j *EmailController) PutEmail() {

	var InfoToSend models.BodyEmail
	session, _ := db.GetSession()
	json.Unmarshal(j.Ctx.Input.RequestBody, &InfoToSend)
	err := models.EmailSender(&InfoToSend, session)

	j.Data["json"] = "ok"
	defer session.Close()
	if err != nil {
		j.Data["json"] = "error"
	}
	j.ServeJSON()
}

//TestEmail - Test email
// @Title TestEmail
// @Description update the Infoapoyo
// @Param	code		path 	string	true		"The code you want to update"
// @Success 200 {object} models.Object
// @router /test [put]
func (j *EmailController) TestEmail() {

	var TestToSend models.Email
	json.Unmarshal(j.Ctx.Input.RequestBody, &TestToSend)
	err := models.TestConnection(&TestToSend)

	j.Data["json"] = "Conexion exitosa!"
	if err != nil {
		j.Data["json"] = "ERROR:  " + err.Error()
	}
	j.ServeJSON()
}
