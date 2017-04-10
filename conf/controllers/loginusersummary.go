package controllers

import (
	// "github.com/eaciit/crowd"
	// . "eaciit/bankingsalesperformance/consoleapps/helper"
	m "eaciit/bankingsalesperformance/webapps/models"
	"github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"
	"os"
	// "strconv"
	"time"
)

// LoginUserSummary
type LoginUserSummary struct {
	*BaseController
}

// Generate
func (d *LoginUserSummary) Generate(base *BaseController) {
	tk.Println("Generating Login User & List User Summary..")
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx

	// Predefine
	loginUserList := make([]tk.M, 0)

	csr, err := ctx.Connection.NewQuery().
		From("LoginUser").
		Cursor(nil)
	err = csr.Fetch(&loginUserList, 0, true)
	csr.Close()
	if err != nil {
		tk.Println(err)
		os.Exit(0)
	}

	for _, i := range loginUserList {

		doc := m.NewLoginUserSummary()
		doc.Period = i.Get("Period").(time.Time)
		doc.User_ID = i.GetString("USER_ID")
		doc.LDAP_Roles = i.GetString("LDAP_ROLES")
		doc.Peoplesoft_Country = i.GetString("PEOPLESOFT_COUNTRY")
		doc.Region = i.GetString("REGION")
		doc.Department = i.GetString("DEPARTMENT")
		doc.Module_Map = i.GetString("NEW_MODULE_MAP")
		doc.WBType = i.GetString("WB_TYPE")
		doc.TotalRequest = i.GetInt("TOTAL_REQUEST")
		doc.RM_Role = i.GetString("RM_Role")
		// // Get Data Location
		csr, err = ctx.Find(m.NewRegion(), tk.M{"where": dbox.Eq("Country", doc.Peoplesoft_Country)})
		if err != nil {
			tk.Println(err)
			os.Exit(0)
		}
		selectedCountry := m.NewRegion()
		csr.Fetch(&selectedCountry, 1, false)
		csr.Close()
		if selectedCountry.Major_Region != "" {
			doc.Major_Region = selectedCountry.Major_Region
		} else {
			doc.Major_Region = "Others"
		}

		doc.DateUpdated = time.Now().UTC()
		err = d.Ctx.Save(doc)
		if err != nil {
			tk.Println(err)
		}
	}

	userList := make([]tk.M, 0)

	csr, err = ctx.Connection.NewQuery().
		From("UserList").
		Cursor(nil)
	err = csr.Fetch(&userList, 0, true)
	csr.Close()
	if err != nil {
		tk.Println(err)
		os.Exit(0)
	}

	for _, i := range userList {
		doc := m.NewUserListSummary()
		doc.Period = i.Get("Period").(time.Time)
		doc.User_ID = i.GetInt("USER_ID")
		doc.Username = i.GetString("USER_NAME")
		doc.LDAP_Roles = i.GetString("LDAP_ROLES")
		doc.Peoplesoft_Country = i.GetString("PEOPLESOFT_COUNTRY")
		doc.Region = i.GetString("REGION")
		doc.Department = i.GetString("DEPARTMENT")
		doc.RM_Role = i.GetString("RM_Role")
		// // Get Data Location
		csr, err = ctx.Find(m.NewRegion(), tk.M{"where": dbox.Eq("Country", doc.Peoplesoft_Country)})
		if err != nil {
			tk.Println(err)
			os.Exit(0)
		}
		selectedCountry := m.NewRegion()
		csr.Fetch(&selectedCountry, 1, false)
		csr.Close()
		if selectedCountry.Major_Region != "" {
			doc.Major_Region = selectedCountry.Major_Region
		} else {
			doc.Major_Region = "Others"
		}

		doc.DateUpdated = time.Now().UTC()
		err = d.Ctx.Save(doc)
		if err != nil {
			tk.Println(err)
		}
	}
	tk.Println("Prospectst Data : COMPLETE")
}
