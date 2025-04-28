import json
import os
import matplotlib.pyplot as plt

log_dir = "JSONlogs"
iterations = []

# Loop through iteration JSON files
for filename in sorted(os.listdir(log_dir)):
    if filename.startswith("iteration_") and filename.endswith(".json"):
        with open(os.path.join(log_dir, filename), "r") as f:
            data = json.load(f)
            iteration_num = data["Iteration"]
            for turn in data["Turns"]:
                iterations.append({
                    "Iteration": iteration_num,
                    "Turn": turn["TurnNumber"],
                    "AgentCount": turn["NumberOfAgents"]
                })

# Sort by iteration+turn
iterations.sort(key=lambda x: (x["Iteration"], x["Turn"]))

# Extract data for plotting
rounds = list(range(len(iterations)))
agent_counts = [entry["AgentCount"] for entry in iterations]

# Plot
plt.figure(figsize=(14, 5))
plt.plot(rounds, agent_counts, marker="o", label="Agent Count")
plt.xticks(ticks=rounds, labels=[str(r) for r in rounds], rotation=90, fontsize=8)
plt.xlabel("Rounds")
plt.ylabel("Number of Agents")
plt.title("Total Agent Population Over Time")
plt.grid(True)
plt.legend()
plt.tight_layout()
plt.show()