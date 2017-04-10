package controllers

import (
	// "github.com/eaciit/crowd"
	. "eaciit/bankingsalesperformance/consoleapps/helper"
	m "eaciit/bankingsalesperformance/webapps/models"
	"github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"
	"os"
	"strconv"
	"time"
)

// AccountPlanSummary
type AccountPlanSummary struct {
	*BaseController
}

// Generate
func (d *AccountPlanSummary) Generate(base *BaseController) {
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx

	// Predefine
	planList := make([]tk.M, 0)
	pipes := []tk.M{}
	projection := tk.M{}
	grouping := tk.M{}
	err := Deserialize(`
		{"$project":{
	        "AP_Approved_Date":"$AP_Approved_Date",
	        "Status":
	        {"$cond":[{ "$eq": [ "$AP_Status", "APPROVED"] },"Approved", 
	            {"$cond":[{ "$or":[{"$eq": [ "$AP_Status", "APPROVAL ON GOING"]},{"$eq": [ "$AP_Status", "DRAFT"]}] },
	        "Initiated", "UNKNOWN"]}]},
	        "Year":{"$year":"$AP_Approved_Date"},
	        "Month":{"$month":"$AP_Approved_Date"},
	        "Location":"$GAM_Location",
	        "Group_ID":"$Group_ID"
	    }}
	`, &projection)
	err = Deserialize(`
		{"$group":{
        "_id":{
            "group_id":"$Group_ID",
            "location":"$Location",
            "year":"$Year",
            "month":"$Month"
        },
        "approved":{"$sum":{"$cond":[{"$eq":["$Status","Approved"]},1,0]}},
        "initiated":{"$sum":{"$cond":[{"$eq":["$Status","Initiated"]},1,0]}},
        "unknown":{"$sum":{"$cond":[{"$eq":["$Status","UNKNOWN"]},1,0]}},
        "submission":{"$sum":1}
    }}
	`, &grouping)
	pipes = append(pipes, projection)
	pipes = append(pipes, grouping)
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.group_id": 1, "_id.location": 1, "_id.year": 1, "_id.month": 1}})

	csr, err := ctx.Connection.NewQuery().Command("pipe", pipes).
		From("AccountPlan").
		Cursor(nil)
	err = csr.Fetch(&planList, 0, true)
	csr.Close()
	if err != nil {
		tk.Println(err)
		os.Exit(0)
	}

	for _, i := range planList {
		doc := m.NewAccountPlanSummary()
		id := i.Get("_id").(tk.M)
		year := id.GetInt("year")
		month := id.GetInt("month")
		doc.Period = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		doc.Year, _ = strconv.Atoi(doc.Period.Format("2006"))
		doc.MonthYear, _ = strconv.Atoi(doc.Period.Format("200601"))
		doc.Group_ID = id.GetInt("group_id")
		doc.Country = id.GetString("location")
		// Get Data Location
		csr, err = ctx.Find(m.NewRegion(), tk.M{"where": dbox.Eq("Country", doc.Country)})
		if err != nil {
			tk.Println(err)
			os.Exit(0)
		}
		selectedCountry := m.NewRegion()
		csr.Fetch(&selectedCountry, 1, false)
		csr.Close()
		doc.Region = selectedCountry.Region
		doc.Major_Region = selectedCountry.Major_Region

		// Check wheter its BCA Group or not
		csr, err = ctx.Find(m.NewClient(), tk.M{"where": dbox.And(dbox.Eq("Client_Group_ID", strconv.Itoa(doc.Group_ID)), dbox.Eq("Client_BCA_Flag", "Y"))})
		if err != nil {
			tk.Println(err)
			os.Exit(0)
		}
		selectedClient := m.NewClient()
		csr.Fetch(&selectedClient, 1, false)
		csr.Close()
		if selectedClient.Client_ID == "" {
			doc.IsBCAGroup = 1 //Set as active
		}
		doc.Approved = i.GetInt("approved")
		doc.Initiated = i.GetInt("initiated")
		doc.Unknown = i.GetInt("unknown")
		doc.Submission = i.GetInt("submission")

		doc.DateUpdated = time.Now().UTC()
		err = d.Ctx.Save(doc)
		if err != nil {
			tk.Println(err)
		}
	}
	tk.Println("Generating Account Plan Summary..")
	tk.Println("Account Plant Data : COMPLETE")
}
