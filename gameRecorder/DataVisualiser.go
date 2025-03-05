package gameRecorder

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// Add these constants at the top of the file
const (
	outputDir      = "visualisation_output"
	deathSymbol    = "💀"
	showLegends    = false   // Toggle for showing/hiding legends
	//showAxisLabels = true    // Keep axis labels visible
	//chartWidth     = "800px" // Increased from 500px
	//chartHeight    = "500px" // Increased from 400px
)

// CreatePlaybackHTML reads recorded data and generates an HTML file with charts
func CreatePlaybackHTML(recorder *ServerDataRecorder) {
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("Error creating output directory: %v", err)
		return
	}

	page := components.NewPage()
	page.PageTitle = "TMT Visualisation"

	iterationMap := make(map[int][]TurnRecord)
	for _, record := range recorder.TurnRecords {
		iterationMap[record.IterationNumber] = append(iterationMap[record.IterationNumber], record)
	}

	for iteration, turns := range iterationMap {
		sacrificeChart := createSacrificeChart(iteration, turns)
		if sacrificeChart != nil {
			page.AddCharts(sacrificeChart)
		}
	}

	filepath := filepath.Join(outputDir, "tmt_visualisation.html")
	f, err := os.Create(filepath)
	if err != nil {
		log.Printf("Error creating visualisation file: %v\n", err)
		return
	}
	defer f.Close()
	page.Render(f)
}

func createSacrificeChart(iteration int, turns []TurnRecord) *charts.Line {
    // Sort turns by turn number to ensure correct order
    sort.Slice(turns, func(i, j int) bool {
        return turns[i].TurnNumber < turns[j].TurnNumber
    })

    // Create a new line chart
    line := charts.NewLine()
    line.SetGlobalOptions(
        charts.WithTitleOpts(opts.Title{
            Title: fmt.Sprintf("Iteration %d - Agent Sacrifices Over Time", iteration),
            Left:  "center",
            Top:   "3%",
        }),
        charts.WithLegendOpts(opts.Legend{
            Show: opts.Bool(true),
            Left: "85%", // Moves legend to the right
            Top:  "10%", // Adjusts vertical alignment
        }),
        charts.WithTooltipOpts(opts.Tooltip{
            Show: opts.Bool(true),
        }),
        charts.WithXAxisOpts(opts.XAxis{
            Name: "Turn Number",
            NameGap: 30,
            AxisLabel: &opts.AxisLabel{
                Show: opts.Bool(true),
                Rotate: 0,
                Margin: 10,
            },
        }),
        charts.WithYAxisOpts(opts.YAxis{
            Name: "Self-Sacrifice Willingness",
            NameGap: 40,
            SplitLine: &opts.SplitLine{Show: opts.Bool(true)},
        }),
        charts.WithGridOpts(opts.Grid{
            Left: "10%",  // Moves graph slightly left to leave space for the legend
            Right: "20%", // Leaves space for the legend
        }),
        charts.WithInitializationOpts(opts.Initialization{
            Width:  "1000px", // Increased width to accommodate legend
            Height: "500px",
        }),
    )

    // Get turn numbers for X-axis
    xAxis := []int{}
    for _, turn := range turns {
        xAxis = append(xAxis, turn.TurnNumber)
    }

    // Store sacrifice willingness for each agent
    agentSacrifices := make(map[string][]opts.LineData)
    deathMarkers := make([]opts.ScatterData, 0)
    agentDeathTurn := make(map[string]bool) // Track if an agent has already been marked dead

    for _, turn := range turns {
        for _, agent := range turn.AgentRecords {
            agentID := agent.AgentID.String()
            if _, exists := agentSacrifices[agentID]; !exists {
                agentSacrifices[agentID] = make([]opts.LineData, len(turns))
            }
            agentSacrifices[agentID][turn.TurnNumber] = opts.LineData{Value: agent.SelfSacrificeWillingness}

            // If the agent is eliminated this turn, add a death marker **only once**
            if !agent.IsAlive && !agentDeathTurn[agentID] {
                deathMarkers = append(deathMarkers, opts.ScatterData{
                    Value: []interface{}{turn.TurnNumber, agent.SelfSacrificeWillingness},
                    Symbol: "rect", // Use a visible symbol instead of "pin"
                    SymbolSize: 12,
                })
                agentDeathTurn[agentID] = true // Mark that this agent has been processed
            }
        }
    }

    // Add data series for each agent
    for agentID, sacrifices := range agentSacrifices {
        line.AddSeries(
            fmt.Sprintf("Agent %s", agentID), sacrifices,
            charts.WithLineStyleOpts(opts.LineStyle{
                Width: 2,
            }),
        )
    }

    // Overlay death markers if any agents were eliminated
    if len(deathMarkers) > 0 {
        scatter := charts.NewScatter()
        scatter.AddSeries("Eliminations", deathMarkers,
            charts.WithItemStyleOpts(opts.ItemStyle{
                Color: "black",
            }),
        )
        line.Overlap(scatter)
    }

    // Set X-axis values
    line.SetXAxis(xAxis)

    return line
}

// func createSacrificePlot(recorder *ServerDataRecorder, outputDir string) {
// 	// Add safety check at the start
// 	if len(recorder.TurnRecords) == 0 {
// 		log.Println("Warning: No turn records to visualise")
// 		return
// 	}

// 	// Group turn records by iteration
// 	iterationMap := make(map[int][]TurnRecord)
// 	for _, record := range recorder.TurnRecords {
// 		iterationMap[record.IterationNumber] = append(iterationMap[record.IterationNumber], record)
// 	}

// 	// Create a page to hold all iteration charts
// 	page := components.NewPage()
// 	page.PageTitle = "Agent Sacrifices Per Iteration"

// 	// For each iteration, create a line chart
// 	for iteration, turns := range iterationMap {
// 		// Sort turns by turn number to ensure correct order
// 		sort.Slice(turns, func(i, j int) bool {
// 			return turns[i].TurnNumber < turns[j].TurnNumber
// 		})

// 		// Find first turn with agent records to initialise our agent map
// 		var initialAgentRecords []AgentRecord
// 		for _, turn := range turns {
// 			if len(turn.AgentRecords) > 0 {
// 				initialAgentRecords = turn.AgentRecords
// 				break
// 			}
// 		}

// 		if len(initialAgentRecords) == 0 {
// 			log.Printf("Warning: No agent records found in iteration %d\n", iteration)
// 			continue
// 		}

// 		// Create a new line chart with adjusted layout
// 		line := charts.NewLine()
// 		line.SetGlobalOptions(
// 			charts.WithTitleOpts(opts.Title{
// 				Title: fmt.Sprintf("Iteration %d - Agent Sacrifices Over Time", iteration),
// 				Left:  "center",  // Centers the title
// 				Top:   "2%",      // Adjusts vertical position
// 				TextStyle: &opts.TextStyle{
// 					FontSize: 18, // Larger font size for readability
// 					Bold:     true,
// 				},
// 			}),
// 			charts.WithTooltipOpts(opts.Tooltip{
// 				Show: opts.Bool(true),
// 			}),
// 			charts.WithLegendOpts(opts.Legend{
// 				Show: opts.Bool(false),
// 			}),
// 			charts.WithXAxisOpts(opts.XAxis{
// 				Name: "Turn Number",
// 				NameGap: 30, // Increases space between axis and name
// 				AxisLabel: &opts.AxisLabel{
// 					Show:   opts.Bool(true),
// 					Rotate: 0, // Ensures horizontal labels
// 					Margin: 10,
// 				},
// 			}),
// 			charts.WithYAxisOpts(opts.YAxis{
// 				Name: "Self-Sacrifice Willingness",
// 				NameGap: 40, // Moves Y-axis label further from axis
// 				SplitLine: &opts.SplitLine{Show: opts.Bool(true)}, // Adds grid lines
// 			}),
// 			charts.WithGridOpts(opts.Grid{
// 				Top:    "28%",
// 				Right:  "5%",
// 				Left:   "10%", // Add more space for Y-axis labels
// 				Bottom: "15%", // Add more space for X-axis labels
// 			}),
// 		)
// 		// Get turn numbers for X-axis
// 		xAxis := make([]int, len(turns))
// 		for i, turn := range turns {
// 			xAxis[i] = turn.TurnNumber
// 		}

// 		// Create a map of agent sacrifices over turns
// 		agentSacrifices := make(map[string][]float64)

// 		// Initialise the map with empty slices using the first found agents
// 		for _, agent := range initialAgentRecords {
// 			agentID := agent.AgentID.String()
// 			agentSacrifices[agentID] = make([]float64, len(turns))
// 		}

// 		// Fill in the sacrifices for each agent
// 		for turnIdx, turn := range turns {
// 			for _, agent := range turn.AgentRecords {
// 				agentID := agent.AgentID.String()
// 				agentSacrifices[agentID][turnIdx] = float64(agent.SelfSacrificeWillingness)
// 			}
// 		}

// 		// When adding series, we can also add death markers
// 		for agentID, sacrifices := range agentSacrifices {
// 			// Find when the agent died (if they did)
// 			var deathMarker opts.ScatterData
// 			var deathTurn int = -1

// 			// Find the turn where agent died
// 			for i, turn := range turns {
// 				for _, agent := range turn.AgentRecords {
// 					if agent.AgentID.String() == agentID {
// 						if !agent.IsAlive {
// 							deathTurn = i
// 							deathMarker = opts.ScatterData{
// 								Value:      []interface{}{xAxis[i], sacrifices[i]},
// 								Symbol:     deathSymbol,
// 								SymbolSize: 20,
// 							}
// 							break
// 						}
// 					}
// 				}
// 				if deathTurn != -1 {
// 					break
// 				}
// 			}
		
// 			// // Add the series with custom styling
// 			// line.AddSeries(agentID, generateLineItems(xAxis[:len(sacrifices)], sacrifices),
// 			// 	charts.WithLineStyleOpts(opts.LineStyle{
// 			// 		Color: teamColors[agentID],
// 			// 	}),
// 			// 	charts.WithItemStyleOpts(opts.ItemStyle{
// 			// 		Color: blue,
// 			// 	}),
// 			// )

// 			// Add death marker as a scatter plot overlay
// 			if deathTurn != -1 {
// 				scatter := charts.NewScatter()
// 				scatter.AddSeries(agentID+" Death", []opts.ScatterData{deathMarker},
// 					charts.WithItemStyleOpts(opts.ItemStyle{
// 						Color: "black",
// 					}),
// 				)
// 				line.Overlap(scatter)
// 			}
// 			// Set X-axis data
// 			line.SetXAxis(xAxis)

// 			// Add the chart to the page
// 			page.AddCharts(line)

// 			// Update file creation
// 			filepath := filepath.Join(outputDir, "agent_sacrifices.html")
// 			f, err := os.Create(filepath)
// 			if err != nil {
// 				log.Printf("Error creating score plots file: %v\n", err)
// 				return
// 			}
// 			defer f.Close()

// 			// Render the page
// 			err = page.Render(f)
// 			if err != nil {
// 				panic(err)
// 			}	
// 		}
// 	}
// }

// ExportToCSV exports the turn records to CSV files
func ExportToCSV(recorder *ServerDataRecorder, outputDir string) error {
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	var allAgentRecords []AgentRecord
	for _, turn := range recorder.TurnRecords {
		allAgentRecords = append(allAgentRecords, turn.AgentRecords...)
	}
	if err := exportStructSliceToCSV(allAgentRecords, filepath.Join(outputDir, "agent_records.csv")); err != nil {
		return fmt.Errorf("failed to export agent records: %v", err)
	}

	return nil
}

func exportStructSliceToCSV(data interface{}, filepath string) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice")
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if v.Len() == 0 {
		return nil
	}

	structType := v.Index(0).Type()
	var headers []string
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.PkgPath == "" {
			headers = append(headers, field.Name)
		}
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		var row []string
		for j := 0; j < item.Type().NumField(); j++ {
			field := item.Type().Field(j)
			if field.PkgPath == "" {
				row = append(row, fmt.Sprint(item.Field(j).Interface()))
			}
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}