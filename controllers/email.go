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

//GetConfig - optiene la configuracion del correo
// @Title GetConfig
// @Description optiene la configuracion del correo
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

//PutConfig - Actualiza la configuracion del correo electronico
// @Title PutConfig
// @Description Actualiza la configuracion del correo electronico
// @Param	body		body 	models.Object	true		"Parametros de configuracion correo, contraseña, protocolo, seguridad, etc..."
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

//PutEmail - Envia el correo electronico a un estudiante
// @Title PutEmail
// @Description Envia el correo electronico a un estudiante
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

//TestEmail - Prueba de la configuracion del correo
// @Title TestEmail
// @Description Prueba de la configuracion del correo
// @Success 200 {object} models.Object
// @router /test [put]
func (j *EmailController) TestEmail() {

	var TestToSend models.Email
	json.Unmarshal(j.Ctx.Input.RequestBody, &TestToSend)
	err := models.TestConnection(&TestToSend)
	/*
	   return 1 to conection success
	   return "(text of error) conection fail"
	*/
	j.Data["json"] = 1
	if err != nil {
		j.Data["json"] = "ERROR:  " + err.Error()
	}
	j.ServeJSON()
}
