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

//GetState - optiene el estado del estudiante en el programa para determinar el ingreso a la plataforma
// @Title GetState
// @Description  optiene el estado del estudiante en el programa para determinar el ingreso a la plataforma
// @Param	code		path 	string	true		"El codigo del estudiante a consultar"
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

//Get - optiene la informacion socioeconomica del estudiante por codigo
// @Title Get
// @Description optiene la informacion socioeconomica del estudiante
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

//Put - Actualiza la informacion socioeconomica del estudiante
// @Title Put
// @Description Actualiza la informacion socioeconomica del estudiante
// @Param	code		path 	string	true		"el codigo del estudiante a actualizar"
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

//LastPut - Comprueba y actualiza la informacion y archivos en el historial
// @Title LastPut
// @Description Comprueba y actualiza la informacion y archivos en el historial ultimo paso para cambiar de estado al estudiante
// @Param	code		path 	string	true		"el codigo del estudiante a actualizar"
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
			j.Data["json"] = nil
		} else {
			values, err := models.GetRequiredFiles(session, code)
			if err != nil {
				j.Data["json"] = nil
			} else {
				err := models.Completefile(session, code, values)
				if err != nil {
					j.Data["json"] = nil
				} else {
					err = models.UpdateState(&InfoEcono, session, code)
					if err != nil {
						j.Data["json"] = nil
					} else {
						j.Data["json"] = InfoEcono
					}
				}
			}
		}
	}
	defer session.Close()
	j.ServeJSON()
}
