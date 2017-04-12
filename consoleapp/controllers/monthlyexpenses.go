package controllers

import (
	dbx "github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
	m "github.com/testcase/webapp/models"
)

type Me struct {
	*BaseController
}

func (d *Me) Generate(base *BaseController) {
	tk.Println("Generating Monthly Expenses ....")

	if base != nil {
		d.BaseController = base
	}

	excelName := "E:/Gopath/src/github.com/testcase/sources/TestCase.xlsx"
	xlFile, err := xlsx.OpenFile(excelName)

	if err != nil {
		tk.Println("err")
	}

	sheet := xlFile.Sheets[1]
	// tk.Println(sheet)...
	for index, row := range sheet.Rows {
		if index >= 6 {
			Description, _ := row.Cells[1].String()
			Category, _ := row.Cells[2].String()
			ProjectCost, _ := row.Cells[3].String()
			ActualCost, _ := row.Cells[4].String()

			data := m.NewMonthlyModel()
			data.Description = Description
			data.Category = Category
			data.ProjectCost = ProjectCost
			data.ActualCost = ActualCost

			tk.Println(data.TableName())
			// tk.Println(data)
			// tk.Println("#")
			existingData := new(m.MonthlyModel)
			// existingData := []m.MonthlyModel{}
			csr, err := d.Ctx.Connection.NewQuery().From(data.TableName()).Where(dbx.Eq("description", Description)).Cursor(nil)
			csr.Fetch(&existingData, 1, false)
			csr.Close()

			tk.Println(existingData.Id)
			if err == nil {
				data.Id = existingData.Id
			}
			// csr, err := d.Ctx.Connection.NewQuery().From(data.TableName()).Where(dbx.Eq("_id", Id)).Cursor(nil)
			// csr.Fetch(&existingData, 0, false)
			// csr.Close()

			err = d.Ctx.Save(data)
			if err != nil {
				tk.Println(err)
			}
			tk.Println(existingData)
		}
	}
	tk.Println("Monthly Data : COMPLETE")
	tk.Println("")
	tk.Println("List Monthly Data : ")
	// Select Data From Cost
	MonthlyList := []m.MonthlyModel{}
	csr, err := d.Ctx.Find(new(m.MonthlyModel), nil)
	if err != nil {
		tk.Println(err)
	}
	err = csr.Fetch(&MonthlyList, 0, false)
	if err != nil {
		tk.Println(err)
	}
	csr.Close()
	for _, i := range MonthlyList {
		tk.Println(i)
	}

}
