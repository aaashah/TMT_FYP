import json
import os
import matplotlib.pyplot as plt

# Directory containing JSON logs
log_dir = "JSONlogs"
required_eliminations_per_turn = 2  # Change this to match your current 'n'

# Container for parsed data
volunteer_data = []

# Load and parse JSON files
for filename in sorted(os.listdir(log_dir)):
    if filename.startswith("iteration_") and filename.endswith(".json"):
        with open(os.path.join(log_dir, filename), "r") as f:
            data = json.load(f)
            iteration_num = data["Iteration"]
            for turn in data["Turns"]:
                turn_num = turn["TurnNumber"]
                
                # Use correct key: "SelfSacrificedAgents"
                v = len(turn.get("SelfSacrificedAgents", []) or [])
                
                volunteer_data.append({
                    "Label": f"i{iteration_num}_t{turn_num}",
                    "Volunteers": v,
                    "Required": required_eliminations_per_turn
                })

# Prepare plot data
x_labels = [entry["Label"] for entry in volunteer_data]
volunteers = [entry["Volunteers"] for entry in volunteer_data]
required = [entry["Required"] for entry in volunteer_data]

# Plot
plt.figure(figsize=(14, 5))
plt.plot(x_labels, volunteers, marker='o', label="Volunteers (v)", color="green")
plt.plot(x_labels, required, linestyle="--", label="Required (n)", color="red")
plt.xticks(rotation=90, fontsize=8)
plt.xlabel("Iteration_Turn")
plt.ylabel("Number of Agents")
plt.title("Volunteers vs Required Eliminations per Turn")
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.show()