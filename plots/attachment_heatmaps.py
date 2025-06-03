import os
import json
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import subprocess
from tqdm import tqdm

log_path = "JSONlogs/output.json"
attachment_types = ["secure", "dismissive", "preoccupied", "fearful"]
plot_dir = "figures"
os.makedirs(plot_dir, exist_ok=True)

for style in tqdm(attachment_types, desc="Attachment Type Sweep"):
    # Build the correct CLI arguments
    props = {
        "secure": 0.0,
        "dismissive": 0.0,
        "preoccupied": 0.0,
        "fearful": 0.0,
    }
    props[style] = 1.0

    # Run the simulator
    subprocess.run([
        "./tmtSimulator",
        "-numAgents=40",
        "-iters=100",
        "-seed=42",  # seed for reproducibility
        f"-tau=0.5",
        f"-secure={props['secure']}",
        f"-dismissive={props['dismissive']}",
        f"-preoccupied={props['preoccupied']}",
        f"-fearful={props['fearful']}"
    ], stdout=subprocess.DEVNULL)

    # Load JSON output
    with open(log_path, "r") as file:
        game_data = json.load(file)

    data = []
    for ITER in game_data.get("Iterations", []):
        iteration = ITER["Iteration"]
        threshold_map = ITER.get("AgentThresholds", {})
        for agent_id, score in threshold_map.items():
            data.append((agent_id, iteration, score))

    if not data:
        continue

    df = pd.DataFrame(data, columns=["agent_id", "iteration", "score"])
    start_times = df.groupby("agent_id")["iteration"].min()
    ordered_agents = start_times.sort_values().index.tolist()
    heatmap_data = df.pivot(index="agent_id", columns="iteration", values="score")
    heatmap_data = heatmap_data.loc[ordered_agents]
    heatmap_data.index = heatmap_data.index.str[:8] + "…"

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
    colorbar.ax.tick_params(labelsize=45)
    colorbar.set_label("Proportion of Threshold Reached", fontsize=40, labelpad=25)
    plt.grid(False)
    plt.xticks(
        ticks=range(0, len(heatmap_data.columns), 5),
        labels=heatmap_data.columns[::5],
        fontsize=45,
    )
    plt.xlabel("Iteration", fontsize=45, labelpad=15)
    plt.ylabel("Agents", fontsize=45, labelpad=25)
    plt.yticks([])
    plt.tight_layout()
    plt.savefig(f"{plot_dir}/agentSelfSac_{style}.pdf", format="pdf")
    plt.close()