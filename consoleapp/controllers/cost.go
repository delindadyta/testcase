package controllers

import (
	// "github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
	m "github.com/testcase/webapp/models"
)

type Cost struct {
	*BaseController
}

func (d *Cost) Generate(base *BaseController) {
	tk.Println("Generating Cost Summary..")
	if base != nil {
		d.BaseController = base
	}

	excelFileName := "E:/Gopath/src/github.com/testcase/sources/TestCase.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		tk.Println(err)
	}

	sheet := xlFile.Sheets[2]
	tk.Println(sheet)
	for index, row := range sheet.Rows {
		if index >= 1 {
			Category, _ := row.Cells[1].String()
			Cost, _ := row.Cells[2].Float()

			data := m.NewCostModel()
			data.Category = Category
			data.Cost = Cost

			err = d.Ctx.Save(data)
			if err != nil {
				tk.Println(err)
			}
			tk.Println(data)
		}

	}

	tk.Println("Cost Data : COMPLETE")
	tk.Println("")
	tk.Println("List Cost Data : ")
	// Select Data From Cost
	CostList := []m.CostModel{}
	csr, err := d.Ctx.Find(new(m.CostModel), nil)
	if err != nil {
		tk.Println(err)
	}
	err = csr.Fetch(&CostList, 0, false)
	if err != nil {
		tk.Println(err)
	}
	csr.Close()
	for _, i := range CostList {
		tk.Println(i)
	}
}
