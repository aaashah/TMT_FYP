package gameRecorder

import (
	//"encoding/csv"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"

	//"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// Add these constants at the top of the file
const (
	outputDir      = "visualisation_output"
	deathSymbol    = "💀"
	showLegends    = false   // Toggle for showing/hiding legends
	showAxisLabels = true    // Keep axis labels visible
	chartWidth     = "800px" // Increased from 500px
	chartHeight    = "500px" // Increased from 400px
)

// CreatePlaybackHTML reads recorded data and generates an HTML file with charts
func CreatePlaybackHTML(recorder *ServerDataRecorder) {
	// Create output directory
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("Error creating output directory: %v", err)
		return
	}

	// Create a single page to hold all charts
	page := components.NewPage()
	page.PageTitle = "TMT Visualisation"

	// Add custom CSS for layout
	page.PageTitle = `
		<title>TMT Visualisation</title>
		<style>
			.chart-container { 
				display: flex; 
				flex-wrap: wrap; 
				justify-content: space-between; 
				margin: 20px;
			}
			.chart { 
				width: 48%; 
				min-width: 800px;  // Match chartWidth
				margin-bottom: 40px;  // Increased spacing between charts
			}
		</style>
		<div class="chart-container">
	`

	// Group turn records by iteration
	iterationMap := make(map[int][]TurnRecord)
	for _, record := range recorder.TurnRecords {
		iterationMap[record.IterationNumber] = append(iterationMap[record.IterationNumber], record)
	}

	// // For each iteration, create side-by-side charts
	// for iteration, turns := range iterationMap {
	// 	scoreChart := createScoreChart(iteration, turns)
	// 	contributionChart := createContributionChart(iteration, turns)

	// 	// Add both charts to the page
	// 	page.AddCharts(scoreChart, contributionChart)
	// }

	// Create the output file
	filepath := filepath.Join(outputDir, "tmt_visualisation.html")
	f, err := os.Create(filepath)
	if err != nil {
		log.Printf("Error creating visualisation file: %v\n", err)
		return
	}
	defer f.Close()

	// Render the page
	err = page.Render(f)
	if err != nil {
		panic(err)
	}
}

func createSacrificeChart (iteration int, turns []TurnRecord) *charts.Line {
	// Sort turns by turn number to ensure correct order
	sort.Slice(turns, func(i, j int) bool {
		return turns[i].TurnNumber < turns[j].TurnNumber
	})

	// Find first turn with agent records
	var initialAgentRecords []AgentRecord
	for _, turn := range turns {
		if len(turn.AgentRecords) > 0 {
			initialAgentRecords = turn.AgentRecords
			break
		}
	}

	if len(initialAgentRecords) == 0 {
		log.Printf("Warning: No agent records found in iteration %d\n", iteration)
		return nil
	}

	// Create a new line chart with adjusted layout
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("Iteration %d - Agent Sacrifices over Time", iteration),
			Top:   "5%",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: opts.Bool(true),
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(showLegends),
			Top:  "15%",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name:    "Turn Number",
			NameGap: 30,
			AxisLabel: &opts.AxisLabel{
				Show: opts.Bool(showAxisLabels),
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name:    "Score",
			NameGap: 30,
			AxisLabel: &opts.AxisLabel{
				Show: opts.Bool(showAxisLabels),
			},
		}),
		charts.WithGridOpts(opts.Grid{
			Top:          "25%",
			Right:        "5%",
			Left:         "10%",
			Bottom:       "15%",
			ContainLabel: opts.Bool(true),
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  chartWidth,
			Height: chartHeight,
		}),
	)

	// Get turn numbers for X-axis
	xAxis := make([]int, len(turns))
	for i, turn := range turns {
		xAxis[i] = turn.TurnNumber
	}

	// Create a map of agent sacrifices over turns
	agentSacrifices := make(map[string][]float64)
	for _, agent := range initialAgentRecords {
		agentID := agent.AgentID.String()
		agentSacrifices[agentID] = make([]float64, len(turns))
	}

	// Add series and death markers
	for agentID, sacrifices := range agentSacrifices {
		var deathMarker opts.ScatterData
		var deathTurn int = -1

		// Find death turn
		for i, turn := range turns {
			for _, agent := range turn.AgentRecords {
				if agent.AgentID.String() == agentID && !agent.IsAlive {
					deathTurn = i
					deathMarker = opts.ScatterData{
						Value:      []interface{}{xAxis[i], sacrifices[i]},
						Symbol:     "💀",
						SymbolSize: 20,
					}
					break
				}
			}
			if deathTurn != -1 {
				break
			}

			// Add the series
			// line.AddSeries(agentID, generateLineItems(xAxis[:50], 5),
			// charts.WithLineStyleOpts(opts.LineStyle{
			// Color: teamColors[agentID],
			// }),
			// )

			// Add death marker
			if deathTurn != -1 {
				scatter := charts.NewScatter()
				scatter.AddSeries(agentID+" Death", []opts.ScatterData{deathMarker})
				line.Overlap(scatter)
			}
		}
		// Add the series
		// line.AddSeries(agentID, generateLineItems(xAxis[:50], 5),
		// 	charts.WithLineStyleOpts(opts.LineStyle{
		// 		Color: teamColors[agentID],
		// 	}),
		// )

		// Add death marker
		if deathTurn != -1 {
			scatter := charts.NewScatter()
			scatter.AddSeries(agentID+" Death", []opts.ScatterData{deathMarker})
			line.Overlap(scatter)
		}
	}
	line.SetXAxis(xAxis)
	return line
}

func createSacrificePlot(recorder *ServerDataRecorder, outputDir string) {
	// Add safety check at the start
	if len(recorder.TurnRecords) == 0 {
		log.Println("Warning: No turn records to visualise")
		return
	}

	// Group turn records by iteration
	iterationMap := make(map[int][]TurnRecord)
	for _, record := range recorder.TurnRecords {
		iterationMap[record.IterationNumber] = append(iterationMap[record.IterationNumber], record)
	}

	// Create a page to hold all iteration charts
	page := components.NewPage()
	page.PageTitle = "Agent Sacrifices Per Iteration"

	// For each iteration, create a line chart
	for iteration, turns := range iterationMap {
		// Sort turns by turn number to ensure correct order
		sort.Slice(turns, func(i, j int) bool {
			return turns[i].TurnNumber < turns[j].TurnNumber
		})

		// Find first turn with agent records to initialise our agent map
		var initialAgentRecords []AgentRecord
		for _, turn := range turns {
			if len(turn.AgentRecords) > 0 {
				initialAgentRecords = turn.AgentRecords
				break
			}
		}

		if len(initialAgentRecords) == 0 {
			log.Printf("Warning: No agent records found in iteration %d\n", iteration)
			continue
		}

		// Create a new line chart with adjusted layout
		line := charts.NewLine()
		line.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{
				Title: fmt.Sprintf("Iteration %d - Agent Sacrifices over Time", iteration),
				Top:   "5%", // Move title to top with some padding
			}),
			charts.WithTooltipOpts(opts.Tooltip{
				Show: opts.Bool(true),
			}),
			charts.WithLegendOpts(opts.Legend{
				Show: opts.Bool(false),
			}),
			charts.WithXAxisOpts(opts.XAxis{
				Name:    "Turn Number",
				NameGap: 30, // Add gap between axis and name
				AxisLabel: &opts.AxisLabel{
					Show:   opts.Bool(true),
					Margin: 20, // Add margin to labels
				},
			}),
			charts.WithYAxisOpts(opts.YAxis{
				Name:    "Score",
				NameGap: 30, // Add gap between axis and name
			}),
			charts.WithGridOpts(opts.Grid{
				Top:    "28%",
				Right:  "5%",
				Left:   "10%", // Add more space for Y-axis labels
				Bottom: "15%", // Add more space for X-axis labels
			}),
		)
		// Get turn numbers for X-axis
		xAxis := make([]int, len(turns))
		for i, turn := range turns {
			xAxis[i] = turn.TurnNumber
		}

		// Create a map of agent sacrifices over turns
		agentSacrifices := make(map[string][]float64)

		// Initialize the map with empty slices using the first found agents
		for _, agent := range initialAgentRecords {
			agentID := agent.AgentID.String()
			agentSacrifices[agentID] = make([]float64, len(turns))
		}

		// Fill in the sacrifices for each agent
		for turnIdx, turn := range turns {
			for _, agent := range turn.AgentRecords {
				agentID := agent.AgentID.String()
				agentSacrifices[agentID][turnIdx] = float64(agent.SelfSacrificeWillingness)
			}
		}

		// When adding series, we can also add death markers
		for agentID, sacrifices := range agentSacrifices {
			// Find when the agent died (if they did)
			var deathMarker opts.ScatterData
			var deathTurn int = -1

			// Find the turn where agent died
			for i, turn := range turns {
				for _, agent := range turn.AgentRecords {
					if agent.AgentID.String() == agentID {
						if !agent.IsAlive {
							deathTurn = i
							deathMarker = opts.ScatterData{
								Value:      []interface{}{xAxis[i], sacrifices[i]},
								Symbol:     deathSymbol,
								SymbolSize: 20,
							}
							break
						}
					}
				}
				if deathTurn != -1 {
					break
				}
			}
		
			// // Add the series with custom styling
			// line.AddSeries(agentID, generateLineItems(xAxis[:len(sacrifices)], sacrifices),
			// 	charts.WithLineStyleOpts(opts.LineStyle{
			// 		Color: teamColors[agentID],
			// 	}),
			// 	charts.WithItemStyleOpts(opts.ItemStyle{
			// 		Color: blue,
			// 	}),
			// )

			// Add death marker as a scatter plot overlay
			if deathTurn != -1 {
				scatter := charts.NewScatter()
				scatter.AddSeries(agentID+" Death", []opts.ScatterData{deathMarker},
					charts.WithItemStyleOpts(opts.ItemStyle{
						Color: "black",
					}),
				)
				line.Overlap(scatter)
			}
			// Set X-axis data
			line.SetXAxis(xAxis)

			// Add the chart to the page
			page.AddCharts(line)

			// Update file creation
			filepath := filepath.Join(outputDir, "agent_sacrifices.html")
			f, err := os.Create(filepath)
			if err != nil {
				log.Printf("Error creating score plots file: %v\n", err)
				return
			}
			defer f.Close()

			// Render the page
			err = page.Render(f)
			if err != nil {
				panic(err)
			}	
		}
	}
}

// ExportToCSV exports the turn records to CSV files, creating separate files for different data types
func ExportToCSV(recorder *ServerDataRecorder, outputDir string) error {
	// Create output directory if it doesn't exist
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Export agent records (flattened from turn records)
	var allAgentRecords []AgentRecord
	for _, turn := range recorder.TurnRecords {
		allAgentRecords = append(allAgentRecords, turn.AgentRecords...)
	}
	if err := exportStructSliceToCSV(allAgentRecords, filepath.Join(outputDir, "agent_records.csv")); err != nil {
		return fmt.Errorf("failed to export agent records: %v", err)
	}

	// Export infra records (flattened from turn records) - with filtering
	var allInfraRecords []InfraRecord
	for _, turn := range recorder.TurnRecords {
		// Only include records where either turn or iteration is non-zero
		if turn.InfraRecord.TurnNumber != 0 || turn.InfraRecord.IterationNumber != 0 {
			allInfraRecords = append(allInfraRecords, turn.InfraRecord)
		}
	}
	if err := exportStructSliceToCSV(allInfraRecords, filepath.Join(outputDir, "infra_records.csv")); err != nil {
		return fmt.Errorf("failed to export infra records: %v", err)
	}

	return nil
}

func exportStructSliceToCSV(data interface{}, filepath string) error {
	// Get the slice value and type
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice")
	}

	// Create the file
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// If the slice is empty, return early
	if v.Len() == 0 {
		return nil
	}

	// Get the type of the struct
	structType := v.Index(0).Type()

	// Write header
	var headers []string
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}
		headers = append(headers, field.Name)
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write data
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		var row []string
		for j := 0; j < item.Type().NumField(); j++ {
			field := item.Type().Field(j)
			// Skip unexported fields
			if field.PkgPath != "" {
				continue
			}
			fieldValue := item.Field(j)
			// Convert the field value to string based on its type
			var strValue string
			switch fieldValue.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				strValue = strconv.FormatInt(fieldValue.Int(), 10)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				strValue = strconv.FormatUint(fieldValue.Uint(), 10)
			case reflect.Float32, reflect.Float64:
				strValue = strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
			case reflect.Bool:
				strValue = strconv.FormatBool(fieldValue.Bool())
			case reflect.String:
				strValue = fieldValue.String()
			default:
				// For complex types, use fmt.Sprint
				strValue = fmt.Sprint(fieldValue.Interface())
			}
			row = append(row, strValue)
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}