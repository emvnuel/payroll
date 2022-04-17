package main

import (
	"github.com/emvnuel/payroll/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/payroll", controllers.GetPayroll)
	r.Run() // listen and serve on 0.0.0.0:8080
}
