package gameRecorder

import (
	//"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

	//"reflect"
	"sort"
	//"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// Add these constants at the top of the file
const (
	//deathSymbol    = "💀"
	showLegends    = false   // Toggle for showing/hiding legends
	showAxisLabels = true    // Keep axis labels visible
	chartWidth     = "800px" // Increased from 500px
	chartHeight    = "500px" // Increased from 400px
)

// CreatePlaybackHTML generates visualizations for the recorded game data
func CreatePlaybackHTML(recorder *ServerDataRecorder) {
	// Create output directory if it doesn't exist
	outputDir := "visualization_output"
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Printf("Error creating output directory: %v\n", err)
		return
	}

	// Create a single page to hold all charts
	page := components.NewPage()
	page.PageTitle = "Game Visualization"

	// Add custom CSS for layout
	page.PageTitle = `
		<title>Game Visualization</title>
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
	filepath := filepath.Join(outputDir, "game_visualization.html")
	f, err := os.Create(filepath)
	if err != nil {
		log.Printf("Error creating visualization file: %v\n", err)
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
			Title: fmt.Sprintf("Iteration %d - Agent Scores over Time", iteration),
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

	// Create a map of agent scores over turns
	agentScores := make(map[string][]float64)
	for _, agent := range initialAgentRecords {
		agentID := agent.AgentID.String()
		agentScores[agentID] = make([]float64, len(turns))
	}

	// Add series and death markers
	for agentID, scores := range agentScores {
		var deathMarker opts.ScatterData
		var deathTurn int = -1

		// Find death turn
		for i, turn := range turns {
			for _, agent := range turn.AgentRecords {
				if agent.AgentID.String() == agentID && !agent.IsAlive {
					deathTurn = i
					deathMarker = opts.ScatterData{
						Value:      []interface{}{xAxis[i], scores[i]},
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