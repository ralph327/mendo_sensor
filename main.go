package main

import (
	"github.com/tarm/serial"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"
	"strings"
	"net/http"
	"strconv"
)

var db gorm.DB

type Data struct {
	gorm.Model
	Temperature	float64	`json:"temperature"`
	Luminosity	float64	`json:"luminosity"`
	Moisture		float64	`json:"moisture"`
	Triggered		bool		`json:"triggered"`
}

type ChildFormatJSON struct {
	Format		string	`json:"format"`
}

type DataReturn struct {
	Format		string	`json:"format"`
	HTML			string 	`json:"html"`
}

type DteDatum struct {
	Format		string	`json:"format"`
	Datum 		[]DteData	`json:"datum"`
}

type DteData struct {
	Date			string	`json:"date"`
	Temperature	string	`json:"temperature"`
	Luminosity	string	`json:"luminosity"`
	Moisture		string	`json:"moisture"`
}

func main() {
	/************************
	 * Set up the variables
	 ***********************/
	var reads []string
	var n,len_reads int
	var err error
	buffer := make([]byte, 128)

	
	/************************	 
	 * Set up the sensor
	 ***********************/

	sensor_config := &serial.Config{Name: "/dev/ttyACM0", Baud: 9600}
	
	sensor, err := serial.OpenPort(sensor_config)
	
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error Opening Port")
	}

	/************************	 
	 * Set up the database
	 ***********************/	
	db, err = gorm.Open("mysql", "mendo_sensor:DBh2oREADer@/mendo_sensor?charset=utf8&parseTime=True&loc=Local")
	db.SingularTable(true)
	db.AutoMigrate(&Data{})
	
	/************************	 
	 * Set up the server
	 ***********************/
	r := gin.Default()
	
	r.LoadHTMLGlob("tmpl/*")
	
	r.StaticFile("/style.css", "./css/style.css")
	r.Static("/images", "./images")
	r.Static("/scripts", "./scripts")
	
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})
	
	r.GET("/read", func(c *gin.Context) {
		n, err = sensor.Read(buffer)
		
		if err != nil {
			log.Fatal(err)
		}
		
		reads = strings.Split(string(buffer[:n]) , ":")
		len_reads = len(reads)
				
		// Length of reads will be at least four when 
		// sensor.Reads the appropriate amount of info
		if len_reads == 4 || len_reads == 7 {
			var data Data
			celsius,_ := strconv.ParseFloat(reads[1],64);
			
			data.Moisture, _ = strconv.ParseFloat(reads[0], 64);
			data.Temperature = ctof(celsius)
			data.Luminosity, _ = strconv.ParseFloat(reads[2],64);
			
			// Check values to see if a watering
			// was triggered
			trigger(&data)
			
			// Insert into the database
			db.Create(&data)
			
			c.JSON(200, data)
		}
	})
	
	// The history page, allows drill down by date granularity
	r.GET("/history", history)
	
	r.GET("/appendDataByDate", appendDataByDate)
	
	r.GET("/getChildFormat", getChildFormat)
	
	// Prime the buffer
	for x:=0; x<10; x++ {
		
		n, err = sensor.Read(buffer)
		
		if err != nil {
			log.Fatal(err)
		}
		
		/*
		err = sensor.Flush()
		
		if err != nil{
			log.Fatal(err)
		}
		*/
	}
	
	r.Run(":5555")
}
