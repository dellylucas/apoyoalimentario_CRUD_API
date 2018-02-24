package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:AdministratorController"] = append(beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:AdministratorController"],
		beego.ControllerComments{
			Method: "GetConfig",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:AdministratorController"] = append(beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:AdministratorController"],
		beego.ControllerComments{
			Method: "PutConfig",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:AdministratorController"] = append(beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:AdministratorController"],
		beego.ControllerComments{
			Method: "GetStudents",
			Router: `/:state/:sede`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:AdministratorController"] = append(beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:AdministratorController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:user`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:AdministratorController"] = append(beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:AdministratorController"],
		beego.ControllerComments{
			Method: "GetReport",
			Router: `/report/:state/:sede`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:AdministratorController"] = append(beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:AdministratorController"],
		beego.ControllerComments{
			Method: "PutState",
			Router: `/verification/:code`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:EconomicController"] = append(beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:EconomicController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:code`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:EconomicController"] = append(beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:EconomicController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:code`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:EconomicController"] = append(beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:EconomicController"],
		beego.ControllerComments{
			Method: "GetState",
			Router: `/state/:code`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:EconomicController"] = append(beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:EconomicController"],
		beego.ControllerComments{
			Method: "LastPut",
			Router: `/state/:code`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:FileController"] = append(beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:FileController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:FileController"] = append(beego.GlobalControllerRouter["apoyoalimentario_CRUD_API/controllers:FileController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:code`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
