package controllers

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/models"
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"

	//Libreria conexion mongoDB
	_ "gopkg.in/mgo.v2"
)

//EconomicController Operations CRUD Information economic "Apoyo Alimentario"
type EconomicController struct {
	beego.Controller
}

//GetState - get State of student for intro in plataform
// @Title GetState
// @Description get State of student for intro in plataform
// @Param	code		path 	string	true		"El id del estudiante a consultar"
// @Success 200 {object} models.Infoapoyo
// @Failure 403 :code is empty
// @router /state/:code  [get]
func (j *EconomicController) GetState() {
	code := j.GetString(":code")
	session, _ := db.GetSession()
	if strings.Compare(code, "") != 0 {
		state := models.GetStatus(session, code)
		j.Data["json"] = state
	}
	defer session.Close()
	j.ServeJSON()
}

//Get - get Information of student in BD by code
// @Title Get
// @Description get Information of student in BD by code
// @Param	code	path 	string	true		"El codigo del estudiante a consultar informacion economica"
// @Success 200 {object} models.Infoapoyo
// @Failure 403 :code is empty
// @router /:code [get]
func (j *EconomicController) Get() {
	Code := j.GetString(":code")

	session, _ := db.GetSession()
	defer session.Close()

	if strings.Compare(Code, "") != 0 {
		Infoapoyo, err := models.GetInformationEconomic(session, Code)
		if err != nil {
			j.Data["json"] = err.Error()
		} else {
			j.Data["json"] = *Infoapoyo
		}
	}
	j.ServeJSON()
}

//Put - update the Information economic of student
// @Title Put
// @Description update the Infoapoyo
// @Param	code		path 	string	true		"The code you want to update"
// @Success 200 {object} models.Object
// @Failure 403 :code is empty
// @router /:code [put]
func (j *EconomicController) Put() {
	Codigo := j.Ctx.Input.Param(":code")

	var InfoEcono models.Economic
	resul := "update success!"
	json.Unmarshal(j.Ctx.Input.RequestBody, &InfoEcono)
	session, _ := db.GetSession()

	keyFileDelete, err := models.UpdateInformationEconomic(session, &InfoEcono, Codigo)
	if len(*keyFileDelete) > 0 {
		models.Deletefile(session, Codigo, keyFileDelete)
	}
	if err != nil {
		resul = err.Error()
	}
	j.Data["json"] = resul
	defer session.Close()
	j.ServeJSON()
}

//LastPut - update the Information last step of student
// @Title LastPut
// @Description update the Information last step of student
// @Param	code		path 	string	true		"The code you want to update"
// @Success 200 {object} models.Object
// @Failure 403 :code is empty
// @router /state/:code [put]
func (j *EconomicController) LastPut() {
	code := j.Ctx.Input.Param(":code")

	var InfoEcono models.Economic

	json.Unmarshal(j.Ctx.Input.RequestBody, &InfoEcono)
	session, _ := db.GetSession()

	if strings.Compare(code, "") != 0 {
		keyFileDelete, erro := models.UpdateInformationEconomic(session, &InfoEcono, code)
		if len(*keyFileDelete) > 0 {
			models.Deletefile(session, code, keyFileDelete)
		}
		if erro != nil {
			j.Data["json"] = erro
		} else {
			values, err := models.GetRequiredFiles(session, code)
			if err != nil {
				j.Data["json"] = err.Error()
			} else {
				files, err := models.Completefile(session, code, values)
				if err != nil {
					j.Data["json"] = *files
				} else {
					err = models.UpdateState(session, code)
					if err != nil {
						j.Data["json"] = err.Error()
					} else {
						j.Data["json"] = *files
					}
				}
			}
		}
	}
	defer session.Close()
	j.ServeJSON()
}
