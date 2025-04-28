import json
import os
import matplotlib.pyplot as plt

# Directory containing JSON logs
log_dir = "JSONlogs"
required_eliminations_per_turn = 1  # Change this to match your current 'n'

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
                
                # Use correct key: "EliminatedBySelfSacrifice"
                v = len(turn.get("EliminatedBySelfSacrifice", []) or [])
                
                volunteer_data.append({
                    "Label": f"i{iteration_num}_t{turn_num}",
                    "Volunteers": v,
                    "Required": required_eliminations_per_turn
                })

# Prepare plot data
rounds = list(range(len(volunteer_data)))
volunteers = [entry["Volunteers"] for entry in volunteer_data]
required = [entry["Required"] for entry in volunteer_data]

# Plot
plt.figure(figsize=(14, 5))
plt.plot(rounds, volunteers, marker='o', label="Volunteers (v)", color="green")
plt.plot(rounds, required, linestyle="--", label="Required (n)", color="red")
plt.xticks(ticks=rounds, labels=rounds, rotation=90, fontsize=8)
plt.xlabel("Rounds")
plt.ylabel("Number of Agents")
plt.title("Volunteers vs Required Eliminations per Turn")
plt.ylim(0, max(volunteers + required) + 1)  # Auto-scaled y-axis
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.show()