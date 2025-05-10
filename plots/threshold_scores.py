import json
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt

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

# Plot
plt.figure(figsize=(24, 20))
sns.set_style("white", {"axes.grid": False})
ax = sns.heatmap(
    heatmap_data,
    cmap="rocket_r",
    linewidths=0,
    linecolor="none",
)
colorbar = ax.collections[0].colorbar
colorbar.ax.tick_params(labelsize=45)  # Change tick font size
colorbar.set_label("Proportion of Threshold Reached", fontsize=40, labelpad=25)
plt.grid(False)
plt.xticks(
    ticks=range(0, len(heatmap_data.columns), 5),
    labels=heatmap_data.columns[::5],
    fontsize=45,
)
# plt.title("Agent Self-Sacrifice Decisions over Time", fontsize=40, pad=45)
plt.xlabel("Iteration", fontsize=45, labelpad=15)
plt.ylabel("Agents", fontsize=45, labelpad=25)
plt.yticks([])
plt.tight_layout()
# plt.show()
plt.savefig("figures/agentSelfSac.pdf", format="pdf")
