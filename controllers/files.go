package controllers

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/models"
	"apoyoalimentario_CRUD_API/utility"
	"os"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

//FileController Operations about Files
type FileController struct {
	beego.Controller
}

//Post - create Files
// @Title CreateFile
// @Description create Files
// @Param	body		body 	models.File	true		"body for File content"
// @Success 200 {int} models.File.Id
// @Failure 403 body is empty
// @router / [post]
func (u *FileController) Post() {

	getcode := u.Ctx.Request.MultipartForm.Value
	//get code
	code := getcode["cod"][0]
	if code != "" {
		path := utility.FileSavePath + code + "\\" + strconv.Itoa(time.Now().UTC().Year()) + "-" + strconv.Itoa(utility.Semester()) + "\\"
		err := os.MkdirAll(path, 0777)
		if err != nil {
			u.Data["json"] = err.Error()
		}
		state := ""
		session, _ := db.GetSession()
		getfiles := u.Ctx.Request.MultipartForm.File
		//get files
		for fil := range getfiles {
			if len(fil) > 0 {
				arre := getfiles[fil]
				/*Archivos pdf y menores de 500 KB se guardan en servidor y en historico BD*/
				if arre[0].Header["Content-Type"][0] == "application/pdf" && arre[0].Size < 512050 {
					u.SaveToFile(fil, path+fil+".pdf")
					models.Insertfile(session, fil, arre[0].Size, code)
				} else { /*Error al subir documento*/
					state += arre[0].Filename + "/"
				}
			}
		}
		defer session.Close()
		u.Data["json"] = state

	} else {
		u.Data["json"] = "Error Sin sesion"
	}
	u.ServeJSON()
}

//Get - get File by code
// @Title Get
// @Description get File by code
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.File
// @Failure 403 :code is empty
// @router /:code [get]
func (u *FileController) Get() {
	code := u.GetString(":code")
	session, _ := db.GetSession()
	if code != "" {
		Infofiles, err := models.GetFiles(session, code)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = Infofiles
		}
	}
	u.ServeJSON()
}
