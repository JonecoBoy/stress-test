package reporter

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/JonecoBoy/stress-test/appContext"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

type ErrorEntry struct {
	Code  int `xml:"Code"`
	Count int `xml:"Count"`
}

type Report struct {
	TotalTimeSpent          string       `xml:"TotalTimeSpent" json:"Total time spent" csv:"Total time spent" yaml:"Total time spent" toml:"Total time spent"`
	TotalRequests           int          `xml:"TotalRequests" json:"Total requests" csv:"Total requests" yaml:"Total requests" toml:"Total requests"`
	AverageRequestTimeSpent string       `xml:"AverageRequestTimeSpent" json:"Average request time spent" csv:"Average request time spent" yaml:"Average request time spent" toml:"Average request time spent"`
	SuccessfulRequests      int          `xml:"SuccessfulRequests" json:"Successful requests" csv:"Successful requests" yaml:"Successful requests" toml:"Successful requests"`
	ErrorDistribution       []ErrorEntry `xml:"ErrorDistribution" json:"Errors Distribution" csv:"Errors Distribution" yaml:"Errors Distribution" toml:"Errors Distribution"`
}

type Reporter struct {
	Context *appContext.Context
}

func NewReporter(ctx *appContext.Context) *Reporter {
	return &Reporter{
		Context: ctx,
	}
}

func (r *Reporter) PrepareReport(fileType string) (string, error) {
	var err error
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

	percentages := make([]float64, len(r.Context.RequestTimes))
	for i, t := range r.Context.RequestTimes {
		percentages[i] = float64(t) / float64(r.Context.TotalTime) * 100
	}

	// Create a slice to hold the error entries
	var ErrorsDistribution []ErrorEntry

	// Iterate over the errors and create an entry for each error
	for _, err := range r.Context.Errors {
		// Find if the error code already exists in the slice
		found := false
		for i, entry := range ErrorsDistribution {
			if entry.Code == err {
				// If the error code exists, increment the count
				ErrorsDistribution[i].Count++
				found = true
				break
			}
		}

		// If the error code does not exist in the slice, add a new entry
		if !found {
			ErrorsDistribution = append(ErrorsDistribution, ErrorEntry{Code: err, Count: 1})
		}
	}

	report := Report{
		TotalTimeSpent:          fmt.Sprintf("%v", r.Context.TotalTime),
		TotalRequests:           r.Context.TotalRequests,
		AverageRequestTimeSpent: fmt.Sprintf("%v", average),
		SuccessfulRequests:      r.Context.SuccessfulRequests,
		ErrorDistribution:       ErrorsDistribution,
	}

	// Flatten the errorDistribution map
	var flatErrorDistribution []string
	for errCode, count := range ErrorsDistribution {
		flatErrorDistribution = append(flatErrorDistribution, fmt.Sprintf("%d: %d", errCode, count))
	}

	var data []byte

	switch fileType {
	case "json":
		data, err = json.Marshal(report)
		if err != nil {
			return "", err
		}
	case "yaml":
		data, err = yaml.Marshal(report)
		if err != nil {
			return "", err
		}
	case "toml":
		//report["Errors Distribution"] = flatErrorDistributionStr
		data, err = toml.Marshal(report)
		if err != nil {
			return "", err
		}
	case "xml":
		//report["Errors Distribution"] = flatErrorDistributionStr
		data, err = xml.MarshalIndent(report, "", "  ")
		if err != nil {
			return "", err
		}
	case "csv":
		b := &bytes.Buffer{}
		w := csv.NewWriter(b)

		// Write the non-map fields to the CSV
		err := w.Write([]string{"Total time spent", report.TotalTimeSpent})
		if err != nil {
			return "", err
		}
		err = w.Write([]string{"Total requests", fmt.Sprintf("%d", report.TotalRequests)})
		if err != nil {
			return "", err
		}
		err = w.Write([]string{"Average request time spent", report.AverageRequestTimeSpent})
		if err != nil {
			return "", err
		}
		err = w.Write([]string{"Successful requests", fmt.Sprintf("%d", report.SuccessfulRequests)})
		if err != nil {
			return "", err
		}

		// Flatten the ErrorDistribution map and write it to the CSV
		for errCode, count := range report.ErrorDistribution {
			err = w.Write([]string{fmt.Sprintf("Error %d", errCode), fmt.Sprintf("%d", count)})
			if err != nil {
				return "", err
			}
		}

		w.Flush()
		data = b.Bytes()
	default:
		report := fmt.Sprintf("Total time spent: %v\nTotal requests: %d\nSuccessful requests: %d\n",
			r.Context.TotalTime, r.Context.TotalRequests, r.Context.SuccessfulRequests)

		report += "Error distribution:\n"
		for _, errorEntry := range ErrorsDistribution {
			report += fmt.Sprintf("  %d: %d\n", errorEntry.Code, errorEntry.Count)
		}

		data = []byte(report)
	}

	//report := fmt.Sprintf("Total time spent: %v\nTotal requests: %d\nAverage request time spent: %v\nSuccessful requests: %d\n",
	//	r.Context.TotalTime, r.Context.TotalRequests, average, r.Context.SuccessfulRequests)
	//
	//report += "Error distribution:\n"
	//for errCode, count := range errorDistribution {
	//	report += fmt.Sprintf("  %d: %d\n", errCode, count)
	//}

	return string(data), nil
}

func (r *Reporter) CliReport(report string) {
	fmt.Println(report)
}

func (r *Reporter) LogToFile(data string, filename string, fileType string) error {
	// If filename doesn't have an extension, add it
	if !strings.Contains(filename, ".") {
		filename = fmt.Sprintf("%s.%s", filename, fileType)
	}

	// Write data to the file
	return ioutil.WriteFile(filename, []byte(data), 0644)

}

func (r *Reporter) GenerateHTMLReport(filename string) error {
	filename = filepath.Base(filename)
	ext := filepath.Ext(filename)
	filename = filename[0:len(filename)-len(ext)] + ".html"
	percentages := make([]float64, len(r.Context.RequestTimes))
	for i, t := range r.Context.RequestTimes {
		percentages[i] = float64(t) / float64(r.Context.TotalTime) * 100
	}
	total := time.Duration(0)
	for _, t := range r.Context.RequestTimes {
		total += t
	}
	average := total / time.Duration(len(r.Context.RequestTimes))

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
        drawAllRequestTimesChart();
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
    }`
	html += `function drawAllRequestTimesChart() {
        var data = new google.visualization.DataTable();
        data.addColumn('string', 'Request Time');
        data.addColumn('number', 'Time (ms)');
        data.addColumn('number', 'Average Time (ms)');`

	// Add the request times
	for i, requestTime := range r.Context.RequestTimes {
		html += fmt.Sprintf("data.addRow(['Request %d', %d, %f]);", i+1, requestTime.Milliseconds(), average.Seconds()*1000)
	}

	html += `
        var options = {
            title: 'All Request Times',
            seriesType: 'bars',
            series: {1: {type: 'line'}}
        };

        var chart = new google.visualization.ComboChart(document.getElementById('allRequestTimesChart'));
        chart.draw(data, options);
    }
	</script>
    </head><body>`

	html += `<div id="successErrorChart" style="width: 900px; height: 500px;"></div>`
	html += `<div id="errorDistributionChart" style="width: 900px; height: 500px;"></div>`
	html += `<div id="requestTimesChart" style="width: 900px; height: 500px;"></div>`
	html += `<div id="allRequestTimesChart" style="width: 900px; height: 500px;"></div>`
	html += "</body></html>"

	// Write HTML to the file
	return ioutil.WriteFile(filename, []byte(html), 0644)
}
