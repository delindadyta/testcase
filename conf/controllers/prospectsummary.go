package controllers

import (
	// "github.com/eaciit/crowd"
	. "eaciit/bankingsalesperformance/consoleapps/helper"
	m "eaciit/bankingsalesperformance/webapps/models"
	"os"
	"strconv"
	"time"

	"github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"
)

// ProspectSummary
type ProspectSummary struct {
	*BaseController
}

// Generate
func (d *ProspectSummary) Generate(base *BaseController) {
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx

	// Predefine
	prospectList, prospectConvertedList := make([]tk.M, 0), make([]tk.M, 0)
	pipes := []tk.M{}
	projection := tk.M{}
	grouping := tk.M{}
	err := Deserialize(`
		{"$project":{
	        "Created_On":"$Created_On",
	        "Status":"$Prospect_Status",
	        "Year":{"$year":"$Created_On"},
	        "Month":{"$month":"$Created_On"},
	        "DiffWeeks": {"$divide": [{ "$subtract": [ "$Last_Updated", "$Created_On" ]}, 604800000]},
	        "Location":"$Domicile_Country",
	        "Group_ID":"$Linked_Group_ID",
			"Segment":"$Prospect_Segment"
	    }}
	`, &projection)
	tk.Println(err)
	err = Deserialize(`
		{"$group":{
	        "_id":{
	            "year":"$Year",
	            "month":"$Month",
	            "group_id":"$Group_ID",
	            "location":"$Location",
				"segment":"$Segment"
	        },
	        "ct0006":{"$sum":{"$cond":[{"$lte":["$DiffWeeks",6]},1,0]}},
	        "ct0712":{"$sum":{"$cond":[{"$and":[{"$gt":["$DiffWeeks",6]},{"$lte":["$DiffWeeks",12]}]},1,0]}},
	        "ct1318":{"$sum":{"$cond":[{"$and":[{"$gt":["$DiffWeeks",12]},{"$lte":["$DiffWeeks",18]}]},1,0]}},
	        "ct1924":{"$sum":{"$cond":[{"$and":[{"$gt":["$DiffWeeks",18]},{"$lte":["$DiffWeeks",24]}]},1,0]}},
	        "ct2500":{"$sum":{"$cond":[{"$gt":["$DiffWeeks",24]},1,0]}},
	        "created":{"$sum":1},
	        "qualified":{"$sum":{"$cond":[{"$eq":["$Status","Qualified"]},1,0]}},
	        "lead":{"$sum":{"$cond":[{"$eq":["$Status","Lead"]},1,0]}},
	        "converted":{"$sum":{"$cond":[{"$eq":["$Status","Converted"]},1,0]}},
	        "disqualified":{"$sum":{"$cond":[{"$eq":["$Status","Disqualified"]},1,0]}}
	    }}
	`, &grouping)
	pipes = append(pipes, projection)
	pipes = append(pipes, grouping)
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.year": 1, "_id.month": 1}})
	csr, err := ctx.Connection.NewQuery().Command("pipe", pipes).
		From("Prospect").
		Cursor(nil)
	err = csr.Fetch(&prospectList, 0, true)
	csr.Close()
	if err != nil {
		tk.Println(err)
		os.Exit(0)
	}

	// Get Prospect Converted List
	err = Deserialize(`
		{"$project":{
	        "Updated_On":"$Last_Updated",
	        "Status":"$Prospect_Status",
	        "Year":{"$year":"$Last_Updated"},
	        "Month":{"$month":"$Last_Updated"},
	        "Location":"$Domicile_Country",
	        "Group_ID":"$Linked_Group_ID"
	    }}
	`, &projection)
	err = Deserialize(`
		{"$group":{
	        "_id":{
	            "year":"$Year",
	            "month":"$Month",
	            "group_id":"$Group_ID",
	            "location":"$Location"
	        },
	        "converted":{"$sum":{"$cond":[{"$eq":["$Status","Converted"]},1,0]}},
	    }}
	`, &grouping)
	pipes = []tk.M{}
	pipes = append(pipes, projection)
	pipes = append(pipes, grouping)
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.year": 1, "_id.month": 1}})
	csr, err = ctx.Connection.NewQuery().Command("pipe", pipes).
		From("Prospect").
		Cursor(nil)
	err = csr.Fetch(&prospectConvertedList, 0, true)
	csr.Close()
	if err != nil {
		tk.Println(err)
		os.Exit(0)
	}

	Created, Converted, Lead, Qualified, Disqualified, Open := 0, 0, 0, 0, 0, 0
	PrevMonthCreated, PrevMonthConverted, PrevMonthDisqualified := 0, 0, 0
	for _, i := range prospectList {
		doc := m.NewProspectSummary()
		id := i.Get("_id").(tk.M)
		year := id.GetInt("year")
		month := id.GetInt("month")
		if month == 1 {
			Created = 0
			Converted = 0
			Lead = 0
			Qualified = 0
			Disqualified = 0
			Open = 0
			PrevMonthCreated = 0
			PrevMonthConverted = 0
			PrevMonthDisqualified = 0
		}
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
		if selectedCountry.Major_Region != "" {
			doc.Major_Region = selectedCountry.Major_Region
		} else {
			doc.Major_Region = "Others"
		}

		doc.Ct0006 = i.GetInt("ct0006")
		doc.Ct0712 = i.GetInt("ct0712")
		doc.Ct1318 = i.GetInt("ct1318")
		doc.Ct1924 = i.GetInt("ct1924")
		doc.Ct2500 = i.GetInt("ct2500")

		doc.Created = i.GetInt("created")
		doc.Qualified = i.GetInt("qualified")
		doc.Lead = i.GetInt("lead")
		doc.Converted = i.GetInt("converted")
		doc.Disqualified = i.GetInt("disqualified")
		doc.Open = PrevMonthCreated + doc.Created - PrevMonthConverted - PrevMonthDisqualified
		PrevMonthCreated = doc.Created
		PrevMonthConverted = doc.Converted
		PrevMonthDisqualified = doc.Disqualified
		// Add to YTD value
		Created += doc.Created
		Qualified += doc.Qualified
		Lead += doc.Lead
		Converted += doc.Converted
		Disqualified += doc.Disqualified
		Open += doc.Open

		doc.CreatedYtd = Created
		doc.QualifiedYtd = Qualified
		doc.LeadYtd = Lead
		doc.ConvertedYtd = Converted
		doc.DisqualifiedYtd = Disqualified
		doc.OpenYtd = Open
		doc.DateUpdated = time.Now().UTC()
		err = d.Ctx.Save(doc)
		if err != nil {
			tk.Println(err)
		}
	}

	// Gen Prospect Converted
	Created, Converted, Lead, Qualified, Disqualified, Open = 0, 0, 0, 0, 0, 0
	PrevMonthCreated, PrevMonthConverted, PrevMonthDisqualified = 0, 0, 0
	for _, i := range prospectConvertedList {
		doc := m.NewProspectConvertedSummary()
		id := i.Get("_id").(tk.M)
		year := id.GetInt("year")
		month := id.GetInt("month")
		if month == 1 {
			Created = 0
			Converted = 0
			Lead = 0
			Qualified = 0
			Disqualified = 0
			Open = 0
			PrevMonthCreated = 0
			PrevMonthConverted = 0
			PrevMonthDisqualified = 0
		}
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
		if selectedCountry.Major_Region != "" {
			doc.Major_Region = selectedCountry.Major_Region
		} else {
			doc.Major_Region = "Others"
		}

		doc.Converted = i.GetInt("converted")

		Converted += doc.Converted

		doc.ConvertedYtd = Converted
		doc.DateUpdated = time.Now().UTC()
		err = d.Ctx.Save(doc)
		if err != nil {
			tk.Println(err)
		}
	}
	tk.Println("Generating Prospects Summary..")
	tk.Println("Prospectst Data : COMPLETE")
}
