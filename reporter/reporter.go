package reporter

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/JonecoBoy/stress-test/appContext"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
	"time"
)

type Reporter struct {
	Context *appContext.Context
}

func NewReporter(ctx *appContext.Context) *Reporter {
	return &Reporter{
		Context: ctx,
	}
}

func (r *Reporter) CliReport() {
	total := time.Duration(0)
	for _, t := range r.Context.RequestTimes {
		total += t
	}
	duration := time.Duration(len(r.Context.RequestTimes))
	var average time.Duration
	if duration == 0 {
		average = 0
	} else {
		average = total / duration
	}

	fmt.Println("Total time spent:", r.Context.TotalTime)
	fmt.Println("Total requests:", r.Context.TotalRequests)
	fmt.Println("Average request time spent:", average)
	fmt.Println("Successful requests:", r.Context.SuccessfulRequests)
	//fmt.Println("Status code distribution:")
	//for code, count := range r.Context.StatusCodes {
	//	fmt.Printf("  %d: %d\n", code, count)
	//}
	// Create a map to hold the error distribution
	errorDistribution := make(map[int]int)
	for _, err := range r.Context.Errors {
		// Increment the count for this error message
		errorDistribution[err]++
	}

	fmt.Println("Error distribution:")
	for errCode, count := range errorDistribution {
		fmt.Printf("  %d: %d\n", errCode, count)
	}
}

func (r *Reporter) LogToFile(filename string, fileType string) error {
	// Prepare the report data
	total := time.Duration(0)
	for _, t := range r.Context.RequestTimes {
		total += t
	}
	average := total / time.Duration(len(r.Context.RequestTimes))

	percentages := make([]float64, len(r.Context.RequestTimes))
	for i, t := range r.Context.RequestTimes {
		percentages[i] = float64(t) / float64(r.Context.TotalTime) * 100
	}

	report := map[string]interface{}{
		"Total time spent":           r.Context.TotalTime,
		"Total requests":             r.Context.TotalRequests,
		"Average request time spent": average,
		"Successful requests":        r.Context.SuccessfulRequests,
		"Error distribution":         r.Context.Errors,
		"Requests Average time":      average,
	}
	var data []byte
	var err error

	switch fileType {
	case "json":
		data, err = json.Marshal(report)
		if err != nil {
			return err
		}
	case "yaml":
		data, err = yaml.Marshal(report)
		if err != nil {
			return err
		}
	case "toml":
		data, err = toml.Marshal(report)
		if err != nil {
			return err
		}
	case "csv":
		b := &bytes.Buffer{}
		w := csv.NewWriter(b)
		for key, value := range report {
			err := w.Write([]string{key, fmt.Sprintf("%v", value)})
			if err != nil {
				return err
			}
		}
		w.Flush()
		data = b.Bytes()
	default:
		data = []byte(fmt.Sprintf("Total time spent: %v\nTotal requests: %d\nSuccessful requests: %d\n", r.Context.TotalTime, r.Context.TotalRequests, r.Context.SuccessfulRequests))
	}

	// If filename doesn't have an extension, add it
	if !strings.Contains(filename, ".") {
		filename = fmt.Sprintf("%s.%s", filename, fileType)
	}

	// Write data to the file
	return ioutil.WriteFile(filename, data, 0644)

	//// Add the error distribution
	//data += "Error distribution:\n"
	//errorDistribution := make(map[int]int)
	//for _, err := range r.Context.Errors {
	//	errorDistribution[err]++
	//}
	//for errCode, count := range errorDistribution {
	//	data += fmt.Sprintf("  %d: %d\n", errCode, count)
	//}
	//
	//// Write data to the file
	//return ioutil.WriteFile(filename, []byte(data), 0644)
}

func (r *Reporter) GenerateHTMLReport(filename string) error {
	percentages := make([]float64, len(r.Context.RequestTimes))
	for i, t := range r.Context.RequestTimes {
		percentages[i] = float64(t) / float64(r.Context.TotalTime) * 100
	}
	total := time.Duration(0)
	for _, t := range r.Context.RequestTimes {
		total += t
	}
	//average := total / time.Duration(len(r.Context.RequestTimes))

	// Prepare the HTML report
	html := "<html><head>"
	html += `<script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
    google.charts.load('current', {'packages':['corechart']});
    google.charts.setOnLoadCallback(drawCharts);
    function drawCharts() {
        drawErrorDistributionChart();
        drawSuccessErrorChart();
        drawRequestTimesChart();
    }

    function drawErrorDistributionChart() {
        var data = google.visualization.arrayToDataTable([
        ['Error', 'Count'],`

	// Add the error distribution
	errorDistribution := make(map[int]int)
	for _, err := range r.Context.Errors {
		errorDistribution[err]++
	}
	for errCode, count := range errorDistribution {
		html += fmt.Sprintf("['%d', %d],", errCode, count)
	}

	html += `]);

        var options = {
        title: 'Error distribution'
        };

        var chart = new google.visualization.PieChart(document.getElementById('errorDistributionChart'));

        chart.draw(data, options);
    }

    function drawSuccessErrorChart() {
        var data = google.visualization.arrayToDataTable([
        ['Type', 'Count'],
        ['Successful requests', ` + fmt.Sprintf("%d", r.Context.SuccessfulRequests) + `],
        ['Error requests', ` + fmt.Sprintf("%d", len(r.Context.Errors)) + `]
        ]);

        var options = {
        title: 'Successful vs Error requests'
        };

        var chart = new google.visualization.PieChart(document.getElementById('successErrorChart'));

        chart.draw(data, options);
    }

    function drawRequestTimesChart() {
        var data = google.visualization.arrayToDataTable([
        ['Request Time', 'Count'],`

	// Calculate the min and max request times
	minRequestTime := r.Context.RequestTimes[0]
	maxRequestTime := r.Context.RequestTimes[0]
	for _, requestTime := range r.Context.RequestTimes {
		if requestTime < minRequestTime {
			minRequestTime = requestTime
		}
		if requestTime > maxRequestTime {
			maxRequestTime = requestTime
		}
	}

	// Calculate the range of each bin
	binRange := (maxRequestTime - minRequestTime) / 10

	// Initialize the bins
	bins := make([]int, 10)

	// Categorize the request times into bins
	for _, requestTime := range r.Context.RequestTimes {
		binIndex := int((requestTime - minRequestTime) / binRange)
		// Ensure that the maximum request time falls into the last bin
		if binIndex == 10 {
			binIndex = 9
		}
		bins[binIndex]++
	}

	// Add the bins to the HTML
	for i, count := range bins {
		if count > 0 { // Only add bins with a count greater than zero
			rangeStart := minRequestTime + time.Duration(i)*binRange
			rangeEnd := rangeStart + binRange
			html += fmt.Sprintf("['%dms-%dms', %d],", int(rangeStart.Milliseconds()), int(rangeEnd.Milliseconds()), count)
		}
	}

	html += `]);

        var options = {
        title: 'Request Times',
        legend: { position: 'none' },
        };

        var chart = new google.visualization.BarChart(document.getElementById('requestTimesChart'));

        chart.draw(data, options);
    }
    </script>
    </head><body>`
	html += `<div id="errorDistributionChart" style="width: 900px; height: 500px;"></div>`
	html += `<div id="successErrorChart" style="width: 900px; height: 500px;"></div>`
	html += `<div id="requestTimesChart" style="width: 900px; height: 500px;"></div>`
	html += "</body></html>"

	// Write HTML to the file
	return ioutil.WriteFile(filename, []byte(html), 0644)
}
