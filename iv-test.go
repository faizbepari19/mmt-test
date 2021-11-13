package main

import (
    "fmt"
    _"html"
    "log"
    "net/http"
	"encoding/csv"
	"encoding/json"
	"os"
	"strconv"
	"errors"
)

type ScheduleRecord struct {
    FlightNo int `json:"flightno"`
    From string `json:"from"`
    To string `json:"to"`
	StartTime int `json:"start"`
	EndTime int `json:"end"`
} 

var FlighSchedue = []ScheduleRecord{}

func main() {

	generateCorpus();
	// var src , dest = "ATQ", "BLR"
	// var routes = []ScheduleRecord{}

	// routes = findRoutes(FlighSchedue, src, dest)
	

	// fmt.Println(routes)

    http.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("GET params were:", r.URL.Query())
		src := r.URL.Query().Get("src")
		dest := r.URL.Query().Get("dest")

		if src == "" || dest == "" {
			w.WriteHeader(400)
			return
		}

		routes := findRoutes(FlighSchedue, src, dest)
		fmt.Println(routes)
		jData, err := json.Marshal(routes)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
    })

    log.Fatal(http.ListenAndServe(":8081", nil))

}

func generateCorpus() {
	records := readCsvFile("./ivtest-sched.csv")
	// fmt.Println(records)

	// Loop through lines & turn into object
    for _, line := range records {
		// fmt.Println(line[0])
		fno, _ := strconv.Atoi(line[0])
		st, _ := strconv.Atoi(line[3])
		et, _ := strconv.Atoi(line[4])

        data := ScheduleRecord{
            FlightNo: fno,
            From: line[1],
            To: line[2],
            StartTime: st,
            EndTime: et,
        }

		FlighSchedue = append(FlighSchedue, data)
        // fmt.Println(data.FlightNo)
        // fmt.Println(data.From)
    }

	// fmt.Println(FlighSchedue)
}

func readCsvFile(filePath string) [][]string {
    f, err := os.Open(filePath)
    if err != nil {
        log.Fatal("Unable to read input file " + filePath, err)
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    records, err := csvReader.ReadAll()
    if err != nil {
        log.Fatal("Unable to parse file as CSV for " + filePath, err)
    }

    return records
}


func findRoutes(data []ScheduleRecord, src string, dest string) []ScheduleRecord {
	fmt.Println(src, dest)
	var routes = []ScheduleRecord{}

	for _, schedule := range data {
		if schedule.From == src && schedule.To == dest {
			//DIRECT FIGHT
			routes = append(routes, schedule)
		}

		if schedule.From == src && schedule.To != dest {
			conn_route, err := getConnectingFlight(FlighSchedue, schedule.To, dest)
			// fmt.Println(data)
			if err == nil {
				routes = append(routes, conn_route)
			}
			
		}
	}

	return routes
}

func getSchedule(data []ScheduleRecord, src string) []ScheduleRecord {
	var sch = []ScheduleRecord{}
	for _, schedule := range data {
		if schedule.From == src {
			sch = append(sch, schedule)
		}
	}

	return sch

}

https://github.com/faizbepari19/mmt-test.git


func getConnectingFlight(data []ScheduleRecord, src string, dest string) (ScheduleRecord, error) {
	for _, schedule := range data {
		if schedule.From == src && schedule.To == dest {
			//DIRECT FIGHT
			return schedule, nil
		}
	}

	return ScheduleRecord{}, errors.New("no flight found")

}
