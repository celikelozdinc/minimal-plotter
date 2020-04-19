package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	statistics "github.com/gonum/stat"
	plot "gonum.org/v1/plot"
	plotter "gonum.org/v1/plot/plotter"
	vg "gonum.org/v1/plot/vg"
)

// Solution struct stores results of the associated experiments
type Solution struct {
	solutionType                string
	experimentToRestoreDuration map[int]float64
	experimentToMemoryFootprint map[int][]float64
	allRestoreDuration          []float64
	meanRestoreDuration         float64
	stdDevRestoreDuration       float64
	allMemoryFootprint          []float64
	meanMemoryFootprint         float64
	stdMemoryFootprint          float64
}

// returns rows & columns
// as slice of slices of strings
func readExperiments() [][]string {
	fileHandler, err := os.Open("data/6000Msg.csv")
	defer fileHandler.Close()
	if err != nil {
		log.Panic("Can not open file")
	}

	experiments, err := csv.NewReader(fileHandler).ReadAll()

	if err != nil {
		log.Panic("Can not read experiments from csv file")
	}

	return experiments
}

// CSV File schema:
// ReplicaSet,Solution,Experiment,SmocId,RestoreDurationInSec,VmPeak,VmSize,VmHWM,VmRSS,VmData,DeltaMemoryUsageInKbFromTop
func parseExperiments(experiments [][]string, database map[string]*Solution) {

	for row := range experiments {
		if row == 0 {
			// Header row
			continue
		}
		// distributed || centralized || conventional
		solution := experiments[row][1]
		// 1...10
		experiment, _ := strconv.Atoi(experiments[row][2])

		// smoc1...5
		smocID := experiments[row][3]

		if smocID == "smoc5" {
			// store RestoreDuration
			restoreDuration, _ := strconv.ParseFloat(experiments[row][4], 64)
			database[solution].experimentToRestoreDuration[experiment] = restoreDuration
		} else {
			// store MemoryFootprint
			// uses VmRSS metric
			memoryFootprint, _ := strconv.ParseFloat(experiments[row][8], 64)
			database[solution].experimentToMemoryFootprint[experiment] = append(database[solution].experimentToMemoryFootprint[experiment], memoryFootprint)
		}

	}

}

func calculateStatistics(database map[string]*Solution) {
	// Statistics for duration
	for experiment := range database["distributed"].experimentToRestoreDuration {
		database["distributed"].allRestoreDuration = append(database["distributed"].allRestoreDuration, database["distributed"].experimentToRestoreDuration[experiment])
	}
	database["distributed"].meanRestoreDuration, database["distributed"].stdDevRestoreDuration = statistics.MeanStdDev(database["distributed"].allRestoreDuration, nil)

	for experiment := range database["centralized"].experimentToRestoreDuration {
		database["centralized"].allRestoreDuration = append(database["centralized"].allRestoreDuration, database["centralized"].experimentToRestoreDuration[experiment])
	}
	database["centralized"].meanRestoreDuration, database["centralized"].stdDevRestoreDuration = statistics.MeanStdDev(database["centralized"].allRestoreDuration, nil)

	for experiment := range database["conventional"].experimentToRestoreDuration {
		database["conventional"].allRestoreDuration = append(database["conventional"].allRestoreDuration, database["conventional"].experimentToRestoreDuration[experiment])
	}
	database["conventional"].meanRestoreDuration, database["conventional"].stdDevRestoreDuration = statistics.MeanStdDev(database["conventional"].allRestoreDuration, nil)

	// Statistics for memory footprint
	for experiment := range database["distributed"].experimentToMemoryFootprint {
		// each of experiment includes: [20164 15432 16140 27464]
		var total float64 = 0
		// calculate total memory footprint of experiment
		for calculation := range database["distributed"].experimentToMemoryFootprint[experiment] {
			total += database["distributed"].experimentToMemoryFootprint[experiment][calculation]
		}
		database["distributed"].allMemoryFootprint = append(database["distributed"].allMemoryFootprint, total)
	}
	database["distributed"].meanMemoryFootprint, database["distributed"].stdMemoryFootprint = statistics.MeanStdDev(database["distributed"].allMemoryFootprint, nil)

	for experiment := range database["centralized"].experimentToMemoryFootprint {
		// each of experiment includes: [20164 15432 16140 27464]
		var total float64 = 0
		// calculate total memory footprint of experiment
		for calculation := range database["centralized"].experimentToMemoryFootprint[experiment] {
			total += database["centralized"].experimentToMemoryFootprint[experiment][calculation]
		}
		database["centralized"].allMemoryFootprint = append(database["centralized"].allMemoryFootprint, total)
	}
	database["centralized"].meanMemoryFootprint, database["centralized"].stdMemoryFootprint = statistics.MeanStdDev(database["centralized"].allMemoryFootprint, nil)

	for experiment := range database["conventional"].experimentToMemoryFootprint {
		// each of experiment includes: [20164 15432 16140 27464]
		var total float64 = 0
		// calculate total memory footprint of experiment
		for calculation := range database["conventional"].experimentToMemoryFootprint[experiment] {
			total += database["conventional"].experimentToMemoryFootprint[experiment][calculation]
		}
		database["conventional"].allMemoryFootprint = append(database["conventional"].allMemoryFootprint, total)
	}
	database["conventional"].meanMemoryFootprint, database["conventional"].stdMemoryFootprint = statistics.MeanStdDev(database["conventional"].allMemoryFootprint, nil)

}

func plotRestoreDuration(database map[string]*Solution) {
	var values plotter.Values
	values = append(values, database["distributed"].meanRestoreDuration)
	values = append(values, database["centralized"].meanRestoreDuration)
	values = append(values, database["conventional"].meanRestoreDuration)
	labels := []string{"Distributed", "Centralized", "Conventional"}

	// Create a vertical BarChart
	plot, err := plot.New()
	if err != nil {
		log.Panic(err)
	}

	barChart, err := plotter.NewBarChart(values, 0.5*vg.Centimeter)

	if err != nil {
		log.Panic(err)
	}

	//mean95, err := plotutil.NewErrorPoints(plotutil.MeanAndConf95,)

	plot.Title.Text = " "
	plot.X.Label.Text = "#replicas=4"
	plot.Y.Label.Text = "Restore Duration(sec)"
	plot.Add(barChart)
	plot.NominalX(labels...)

	err = plot.Save(3*vg.Inch, 3*vg.Inch, "data/Restore_Duration.png")
	if err != nil {
		log.Panic(err)
	}

}

func plotMemoryFootprint(database map[string]*Solution) {
	var values plotter.Values
	values = append(values, database["distributed"].meanMemoryFootprint)
	values = append(values, database["centralized"].meanMemoryFootprint)
	values = append(values, database["conventional"].meanMemoryFootprint)
	labels := []string{"Distributed", "Centralized", "Conventional"}

	// Create a vertical BarChart
	plot, err := plot.New()
	if err != nil {
		log.Panic(err)
	}

	barChart, err := plotter.NewBarChart(values, 0.5*vg.Centimeter)

	if err != nil {
		log.Panic(err)
	}

	//mean95, err := plotutil.NewErrorPoints(plotutil.MeanAndConf95,)

	plot.Title.Text = " "
	plot.X.Label.Text = "#replicas=4"
	plot.Y.Label.Text = "Memory Footprint(KiB)"
	plot.Add(barChart)
	plot.NominalX(labels...)

	err = plot.Save(3*vg.Inch, 3*vg.Inch, "data/Memory_Footprint.png")
	if err != nil {
		log.Panic(err)
	}

}



func main() {
	database := make(map[string]*Solution)
	database["distributed"] = &Solution{
		solutionType:                "distributed",
		experimentToRestoreDuration: make(map[int]float64),
		experimentToMemoryFootprint: make(map[int][]float64),
		allRestoreDuration:          make([]float64, 0),
		meanRestoreDuration:         0,
		stdDevRestoreDuration:       0,
		allMemoryFootprint:          make([]float64, 0),
		meanMemoryFootprint:         0,
		stdMemoryFootprint:          0,
	}
	database["centralized"] = &Solution{
		solutionType:                "centralized",
		experimentToRestoreDuration: make(map[int]float64),
		experimentToMemoryFootprint: make(map[int][]float64),
		allRestoreDuration:          make([]float64, 0),
		meanRestoreDuration:         0,
		stdDevRestoreDuration:       0,
		allMemoryFootprint:          make([]float64, 0),
		meanMemoryFootprint:         0,
		stdMemoryFootprint:          0,
	}
	database["conventional"] = &Solution{
		solutionType:                "conventional",
		experimentToRestoreDuration: make(map[int]float64),
		experimentToMemoryFootprint: make(map[int][]float64),
		allRestoreDuration:          make([]float64, 0),
		meanRestoreDuration:         0,
		stdDevRestoreDuration:       0,
		allMemoryFootprint:          make([]float64, 0),
		meanMemoryFootprint:         0,
		stdMemoryFootprint:          0,
	}

	experiments := readExperiments()
	parseExperiments(experiments, database)
	calculateStatistics(database)
	plotRestoreDuration(database)
	plotMemoryFootprint(database)
}
