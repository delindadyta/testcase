package main

import (
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	. "github.com/testcase/consoleapp/controllers"
	"os"
	"runtime"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	tk.Println("Starting the app..\n")

	conn, e := PrepareConnection()
	if e != nil {
		tk.Println(e)
	} else {
		base := new(BaseController)
		base.Ctx = orm.New(conn)
		defer base.Ctx.Close()
		new(Cost).Generate(base)

	}

	tk.Println("Application Close..")
}
