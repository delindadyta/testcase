package controllers

import (
	"bufio"
	"log"

	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	// "math"
	"io/ioutil"
	"os"

	"strings"
	"sync"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()

	// mu                 = &sync.Mutex{}
	retry              = 10
	worker             = 100
	maxDataEachProcess = 100000
	idx                = 0
	mu                 = &sync.Mutex{}
	muinsert           = &sync.Mutex{}
)

type IBaseController interface {
	// not implemented anything yet
}

type BaseController struct {
	base IBaseController
	Ctx  *orm.DataContext
}

func GetMgoValue(d tk.M, fieldName string) interface{} {
	index := strings.Index(fieldName, ".")
	if index < 0 {
		return d.Get(fieldName)
	} else {
		data := d.Get(fieldName[0:index])
		if data != nil {
			return GetMgoValue(data.(tk.M), fieldName[(index+1):len(fieldName)])
		} else {
			return nil
		}
	}
}

func (b *BaseController) GetById(m orm.IModel, id interface{}, column_name ...string) error {
	var e error
	c := b.Ctx.Connection
	column_id := "Id"
	if column_name != nil && len(column_name) > 0 {
		column_id = column_name[0]
	}
	csr, e := c.NewQuery().From(m.(orm.IModel).TableName()).Where(dbox.Eq(column_id, id)).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return e
	}
	e = csr.Fetch(m, 1, false)
	if e != nil {
		return e
	}

	return nil
}

func (b *BaseController) Delete(m orm.IModel, id interface{}, column_name ...string) error {
	column_id := "Id"
	if column_name != nil && len(column_name) > 0 {
		column_id = column_name[0]
	}
	e := b.Ctx.Connection.NewQuery().From(m.(orm.IModel).TableName()).Where(dbox.Eq(column_id, id)).Delete().Exec(nil)
	if e != nil {
		return e
	}
	return nil
}

func (b *BaseController) Update(m orm.IModel, id interface{}, column_name ...string) error {
	column_id := "Id"
	if column_name != nil && len(column_name) > 0 {
		column_id = column_name[0]
	}
	e := b.Ctx.Connection.NewQuery().From(m.(orm.IModel).TableName()).Where(dbox.Eq(column_id, id)).Update().Exec(tk.M{"data": m})
	if e != nil {
		return e
	}
	return nil
}

func (b *BaseController) Truncate(m orm.IModel) error {
	c := b.Ctx.Connection
	e := c.NewQuery().From(m.(orm.IModel).TableName()).Delete().Exec(nil)
	if e != nil {
		return e
	}

	return nil
}
func (b *BaseController) CloseDb() {
	if b.Ctx != nil {
		b.Ctx.Close()
	}
}

func (b *BaseController) WriteLog(msg interface{}) {
	log.Printf("%#v\n\r", msg)
	return
}

func PrepareConnection() (dbox.IConnection, error) {
	config := ReadConfig()
	ci := &dbox.ConnectionInfo{config["host"], config["database"], config["username"], config["password"], nil}
	c, e := dbox.NewConnection("mongo", ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}

	return c, nil
}

func ReadConfig() map[string]string {
	ret := make(map[string]string)
	file, err := os.Open(wd + "conf/app.conf")
	if err == nil {
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				break
			}

			sval := strings.Split(string(line), "=")
			ret[sval[0]] = sval[1]
		}
	} else {
		tk.Println(err.Error())
	}

	return ret
}

func (b *BaseController) GetDataSource(dataSourceFolder string) ([]os.FileInfo, string) {
	config := ReadConfig()
	source := config["datasource"]
	files, e := ioutil.ReadDir(source + "\\" + dataSourceFolder)
	if e != nil {
		tk.Println(e)
		os.Exit(0)
	}
	return files, source + "\\" + dataSourceFolder
}
