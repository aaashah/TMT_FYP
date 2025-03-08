package gameRecorder

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Output directory for grid visualization
const gridOutputDir = "visualisation_output"

func CreateGridHTML(recorder *ServerDataRecorder) {
	err := os.MkdirAll(gridOutputDir, 0755)
	if err != nil {
		log.Fatalf("Error creating output directory: %v", err)
		return
	}

	filepath := filepath.Join(gridOutputDir, "grid_visualisation.html")
	f, err := os.Create(filepath)
	if err != nil {
		log.Printf("Error creating grid visualization file: %v\n", err)
		return
	}
	defer f.Close()

	// Convert TurnRecords to JSON for JavaScript use
	gridData := make(map[int][]map[string]interface{})
	for _, turn := range recorder.TurnRecords {
		var agentPositions []map[string]interface{}
		for _, agent := range turn.AgentRecords {
			agentData := map[string]interface{}{
				"id":    agent.AgentID.String(),
				"x":     agent.PositionX,
				"y":     agent.PositionY,
				"alive": agent.IsAlive,
			}
			agentPositions = append(agentPositions, agentData)
		}
		if len(agentPositions) > 0 {
			gridData[turn.TurnNumber] = agentPositions
		}
	}

	// Convert to JSON
	jsonData, err := json.Marshal(gridData)
	if err != nil {
		log.Printf("Error converting grid data to JSON: %v\n", err)
		return
	}

	// Write HTML file
	f.WriteString(fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Agent Movement Grid</title>
			<style>
				body { font-family: Arial, sans-serif; text-align: center; }
				.grid-container { width: 80%%; margin: auto; }
				#grid-visualisation {
					width: 600px;
					height: 600px;
					border: 2px solid black;
					position: relative;
					display: grid;
					grid-template-columns: repeat(15, 1fr);  /* 15x15 grid */
					grid-template-rows: repeat(15, 1fr);
					background-color: #f0f0f0;
					background-image: 
						linear-gradient(rgba(0, 0, 0, 0.2) 1px, transparent 1px),
						linear-gradient(90deg, rgba(0, 0, 0, 0.2) 1px, transparent 1px);
					background-size: 40px 40px;
				}
				.agent {
					width: 38px;
					height: 38px;
					border-radius: 50%%;
					position: absolute;
					display: flex;
					align-items: center;
					justify-content: center;
					font-size: 18px;
					color: white;
					font-weight: bold;
					box-shadow: 0 2px 5px rgba(0, 0, 0, 0.3);
				}
			</style>
		</head>
		<body>
			<h2>Agent Movement Grid</h2>
			<div id="grid-container" class="grid-container">
				<div id="grid-visualisation"></div>
				<input type="range" min="0" max="%d" value="0" id="turn-slider" style="width: 500px;">
			</div>

			<script>
				const gridData = %s;

				function updateGridVisualisation(turn) {
					console.log("Updating grid for turn:", turn); 

					const gridContainer = document.getElementById("grid-visualisation");
					if (!gridContainer) {
						console.error("Grid container not found!");
						return;
					}

					gridContainer.innerHTML = ""; // Clear previous agents

					if (!gridData[turn] || gridData[turn].length === 0) {
						console.warn("No agent data for turn:", turn);
						return;
					}

					gridData[turn].forEach(agent => {
						let agentDiv = document.createElement("div");
						agentDiv.classList.add("agent");
						agentDiv.style.left = (agent.x * 40 + 1) + "px";  // Align within cells
						agentDiv.style.top = (agent.y * 40 + 1) + "px";
						agentDiv.style.backgroundColor = agent.alive ? "blue" : "black";
						agentDiv.innerHTML = agent.alive ? "" : "💀";

						gridContainer.appendChild(agentDiv);
					});
				}

				document.addEventListener("DOMContentLoaded", function() {
					updateGridVisualisation(0);
				});

				document.getElementById("turn-slider").addEventListener("input", function() {
					updateGridVisualisation(parseInt(this.value));
				});
			</script>
		</body>
		</html>
	`, len(recorder.TurnRecords)-1, string(jsonData)))
}