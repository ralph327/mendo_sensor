package main

func ctof(c float64) (f float64) {
     f = (c * 1.8) + 32
     return f
}

// trigger a watering
func trigger(d *Data) {
	if d.Temperature >= 60 && d.Temperature <= 90 && d.Moisture > 0 && d.Moisture < 30 && d.Luminosity < 25 {
		   d.Triggered = true
	   }else{
		   d.Triggered = false
	   }
}
