import json
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import numpy as np


log_dir = "JSONlogs/output.json"
data = []


with open(log_dir, "r") as file:
    GAME_DATA = json.load(file)
    turn_number = 0
    for ITER in GAME_DATA["Iterations"]:
        iteration = ITER["Iteration"]
        threshold_map: dict = ITER["AgentThresholds"]
        for agent_id, score in threshold_map.items():
            data.append((agent_id, iteration, score))


# Create DataFrame
df = pd.DataFrame(data, columns=["agent_id", "iteration", "score"])
start_times = df.groupby("agent_id")["iteration"].min()

# Sort agent IDs by start time (ascending)
ordered_agents = start_times.sort_values().index.tolist()

# Pivot again, now with sorted agent index
heatmap_data = df.pivot(index="agent_id", columns="iteration", values="score")
heatmap_data = heatmap_data.loc[ordered_agents]  # Reorder rows
heatmap_data.index = (
    heatmap_data.index.str[:8] + "…"
)  # Truncate the index and add ellipsis

# print(heatmap_data)

# Plot
plt.figure(figsize=(14, 12))
sns.set_style("white", {"axes.grid": False})
sns.heatmap(
    heatmap_data,
    cbar_kws={"label": "Proportion of Threshold Reached"},
    cmap="rocket_r",
    linewidths=0,
    linecolor="none",
)
plt.grid(False)
plt.title("Agent Self-Sacrifice Decisions over Time")
plt.xlabel("Iteration")
plt.ylabel("Agent ID")
plt.tight_layout()
# plt.show()
plt.savefig("figures/agentSelfSac.pdf", format="pdf")
