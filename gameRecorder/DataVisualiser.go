package gameRecorder

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
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
	//deathSymbol    = "💀"
	//deathSymbol = "\U0001F480"
	//showLegends    = false   // Toggle for showing/hiding legends
	//showAxisLabels = true    // Keep axis labels visible
	//chartWidth     = "800px" // Increased from 500px
	//chartHeight    = "500px" // Increased from 400px
)

// Stores agent colors permanently
var agentColorMap = make(map[string]string)

// Function to generate a consistent hex color for an agent ID
func getAgentColor(agentID string) string {
    if color, exists := agentColorMap[agentID]; exists {
        return color // Reuse existing color
    }

    // Generate a unique color based on agent ID
    hash := md5.Sum([]byte(agentID))
    hexColor := fmt.Sprintf("#%s", hex.EncodeToString(hash[:3])) // Take first 3 bytes for RGB

    agentColorMap[agentID] = hexColor
    return hexColor
}

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

	// Extract iteration keys and sort them
	iterationKeys := make([]int, 0, len(iterationMap))
	for iteration := range iterationMap {
		iterationKeys = append(iterationKeys, iteration)
	}
	sort.Ints(iterationKeys) // Sort iteration numbers in ascending order

	// Render charts in sorted order
	for _, iteration := range iterationKeys {
		turns := iterationMap[iteration]
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

			// Skip processing if the agent has already been eliminated
			if agentDeathTurn[agentID] {
				continue
			}

            if _, exists := agentSacrifices[agentID]; !exists {
                agentSacrifices[agentID] = make([]opts.LineData, len(turns))
            }
            agentSacrifices[agentID][turn.TurnNumber] = opts.LineData{Value: agent.SelfSacrificeWillingness}

            // If the agent is eliminated this turn, add a death marker **only once**
            if !agent.IsAlive && !agentDeathTurn[agentID] {
                deathMarkers = append(deathMarkers, opts.ScatterData{
                    Value: []interface{}{turn.TurnNumber, agent.SelfSacrificeWillingness},
                    Symbol: "💀", // Use a visible symbol instead of "pin"
                    SymbolSize: 12,
                })
                agentDeathTurn[agentID] = true // Mark that this agent has been processed
            }
        }
    }

    // Add data series for each agent
    for agentID, sacrifices := range agentSacrifices {
		color := getAgentColor(agentID) // Get unique color for agent
        line.AddSeries(
			fmt.Sprintf("Agent %s", agentID),
			sacrifices,
			charts.WithLineStyleOpts(opts.LineStyle{
				Width: 2,
				Color: color, // Ensure the line color stays consistent
			}),
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color: color,         // Ensure the marker (circle) color matches the line
				BorderColor: color,   // Ensure the stroke around the marker is the same
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