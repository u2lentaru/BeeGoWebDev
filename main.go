package main

import (
	_ "BeeGoWebDev/routers"
	"os"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run("localhost:" + os.Getenv("httpport"))
}
