package main

import (
	"./controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/task/", &controllers.TaskController{}, "get:ListTasks;post:NewTask")
	beego.Router("/task/:id:int", &controllers.TaskController{}, "get:GetTask;put:UpdateTask")
	beego.Run()

	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole, `{"level":2,"color":true}`)
	log.Debug("this is debug message")
}
