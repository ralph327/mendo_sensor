package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"net/http"
	"html/template"
	"log"
	"strconv"
	"strings"
	_"fmt"
)

func dteUnderscore(temp string) string {
	return strings.Replace(temp, " ", "_", -1)
}

func dteAppendLoop(data DteDatum, format string, date string) string {
	var temp string 
	for _,tempData := range data.Datum {
		if format == "dteMonthYear" {
				temp += "<div id=\""+ dteUnderscore(tempData.Date) + "\" class=\""+ format + "\"><div  class=\"fncn\" onclick=\"appendDataByDate('"+ tempData.Date +"', '"+ data.Format +"')\">"
		}else{
			temp += "<div id=\""+ dteUnderscore(date) + "_" + dteUnderscore(tempData.Date) +"\" class=\""+ format +"\"><div  class=\"fncn\" onclick=\"appendDataByDate('"+ date + " " + tempData.Date +"', '"+ data.Format +"')\">"
		}
		temp += "<p>" + tempData.Date + "</p></div>"
		
		
		freq := dteFreqSwitch(format)
		temp += dteDataLabels(freq)
		temp += "<div class=\"data\"><div class=\"label\"></div> <div class=\"tmp\">" + tempData.Temperature + "</div>" 
		temp += "<div class=\"lum\">" + tempData.Luminosity + "</div> <div class=\"mos\">" + tempData.Moisture + "</div></div>"
		
		temp += "</div>"
	}
	
	return temp
}

func history(c *gin.Context) {
	var data DteDatum
	//var temp template.HTML
	var temp string
	
	format := "dteMonthYear"
	
	data = dteAVG(format, "")
		
	temp = dteAppendLoop(data, format, "")
	
	args := gin.H{"dteMonthYear": template.HTML(temp)}
	
	c.HTML(http.StatusOK, "history.tmpl", args)
}

func appendDataByDate(c *gin.Context) {
	var data DteDatum
	var json DataReturn
	var temp string

	// Get the data passed from the AJAX call
	date := c.Query("date")
	format := c.Query("format")

	// Get the averages for a date
	data = dteAVG(format, date)
		 
	temp = dteAppendLoop(data, format, date)
	
	json.Format = data.Format
	json.HTML = temp
	
	c.JSON(200, json)
}

func getChildFormat(c *gin.Context) {
	var json ChildFormatJSON
	format := c.Query("format")
	
	json.Format = dteFormatChild(format)
	
	c.JSON(200, json)
}

func dteFreqSwitch(format string) string {
	var freq string
	switch format {
		case "dteMonthYear":
			freq = "Monthly"
		case "dteDayMonth":
			freq = "Daily"
		case "dteHourDay":
			freq = "Hourly"
		case "dteMinHour":
			freq = ""
		case "dteSecMin":
			freq = ""
	}	
	return freq
}

func dteDataLabels(frequency string) string {
	var temp string
	temp = "<div class=\"data\"><div class=\"label\">"+ frequency +" Average</div><div class=\"tmp\">Temperature</div>" 
	temp += "<div class=\"lum\">Luminosity</div> <div class=\"mos\">Moisture</div></div>"
	return temp
}

func dteFormatChild(format string) string {
	switch format {
		case "dteMonthYear":
			return "dteDayMonth"
		case "dteDayMonth":
			return "dteHourDay"
		case "dteHourDay":
			return "dteMinHour"
	}
	
	return "dteSecMin"
}

func dteRangeSwitch(format string, currDate string) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error
	
	switch format {
		case "dteMonthYear":
			rows, err = db.Raw("SELECT DISTINCT date_format(created_at, '%M %Y') FROM data ORDER BY created_at DESC;").Rows()
		case "dteDayMonth":
			rows, err = db.Raw("SELECT DISTINCT date_format(created_at, '%W the %D') FROM data WHERE ? = date_format(created_at, '%M %Y') ORDER BY created_at DESC;", currDate).Rows()
		case "dteHourDay":
			rows, err = db.Raw("SELECT DISTINCT date_format(created_at, '%l %p') FROM data WHERE ? = date_format(created_at, '%M %Y %W the %D') ORDER BY created_at DESC;", currDate).Rows()
		case "dteMinHour":
			rows, err = db.Raw("SELECT DISTINCT date_format(created_at, '%i min') FROM data WHERE ? = date_format(created_at, '%M %Y %W the %D %l %p') ORDER BY created_at DESC;", currDate).Rows()
		case "dteSecMin":
			rows, err = db.Raw("SELECT DISTINCT date_format(created_at, '%s sec') FROM data WHERE ? = date_format(created_at, '%M %Y %W the %D %l %p %i min') ORDER BY created_at DESC;", currDate).Rows()
	}
	
	return rows, err
}

func dteRange(format string, currDate string) []string {
	var dates []string
	var date string
	var rows *sql.Rows
	var err error
	
	rows, err = dteRangeSwitch(format, currDate)
	
	if err != nil {
		log.Fatal(err)
	}
	
	defer rows.Close()
	for rows.Next() {
	    rows.Scan(&date)

	    dates = append(dates, date)
	}
	
	return dates
}

func dteAVGSwitch(format string, date string) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error
	
	switch format {
		case "dteMonthYear":
			rows, err = db.Raw("SELECT AVG(temperature), AVG(luminosity), AVG(moisture) FROM data WHERE ? = date_format(created_at, '%M %Y') LIMIT 1;", date).Rows()
		case "dteDayMonth":
			rows, err = db.Raw("SELECT AVG(temperature), AVG(luminosity), AVG(moisture) FROM data WHERE ? = date_format(created_at, '%M %Y %W the %D') LIMIT 1;", date).Rows()
		case "dteHourDay":
			rows, err = db.Raw("SELECT AVG(temperature), AVG(luminosity), AVG(moisture) FROM data WHERE ? = date_format(created_at, '%M %Y %W the %D %l %p') LIMIT 1;", date).Rows()
		case "dteMinHour":
			rows, err = db.Raw("SELECT AVG(temperature), AVG(luminosity), AVG(moisture) FROM data WHERE ? = date_format(created_at, '%M %Y %W the %D %l %p %i min') LIMIT 1;", date).Rows()
		case "dteSecMin":
			rows, err = db.Raw("SELECT AVG(temperature), AVG(luminosity), AVG(moisture) FROM data WHERE ? = date_format(created_at, '%M %Y %W the %D %l %p %i min %s sec') LIMIT 1;", date).Rows()
	}
	
	return rows, err
}

// Retrieve the average for a day
func dteAVG(format string, currDate string) DteDatum {
	var data DteDatum
	var temp, lum, mos float64
	
	dates := dteRange(format,currDate)
	
	var datum = make([]DteData, len(dates))
	var tempData DteData
	
	for key, date := range dates {
		var rows *sql.Rows
		var err error
		if format == "dteMonthYear" {
			rows, err = dteAVGSwitch(format, date)
		} else {
			rows, err = dteAVGSwitch(format, currDate + " " + date)
		}

		if err != nil {
			log.Fatal(err)
		}
		
		defer rows.Close()
		for rows.Next() {
		    rows.Scan(&temp, &lum, &mos)
		}
		
		tempData.Temperature = strconv.FormatFloat(temp, 'f', 2, 64)
		tempData.Luminosity = strconv.FormatFloat(lum, 'f', 2, 64)
		tempData.Moisture = strconv.FormatFloat(mos, 'f', 2, 64)
		tempData.Date = date
		
		datum[key] = tempData
	}
	
	// Set the data information
	data.Datum = datum
	data.Format = dteFormatChild(format)

	return data	
}

