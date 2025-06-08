import os
import json
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import subprocess
from tqdm import tqdm

# Constants
log_path = "JSONlogs/output.json"
iters = 20
kappa_values = list(range(1, 11))
attachment_styles = ["Secure", "Dismissive", "Preoccupied", "Fearful"]

# Data storage
data = []

# Sweep over cluster size (kappa)
for kappa in tqdm(kappa_values, desc="Kappa sweep"):
    style_ratios = {style: 0 for style in attachment_styles}
    total_ratio = 0

    for _ in range(iters):
        subprocess.run([
            "./tmtSimulator",
            "-numAgents=40",
            "-iters=200",
            "-seed=42",
            f"-kappa={kappa}"
        ], stdout=subprocess.DEVNULL)

        # Load JSON output
        with open(log_path, "r") as file:
            game_data = json.load(file)

        config = game_data.get("Config", {})
        init_agents = config.get("NumAgents", 40)
        iterations = game_data.get("Iterations", [])
        if not iterations:
            continue

        final_iter = iterations[-1]
        final_turns = final_iter.get("Turns", [])
        if not final_turns:
            continue

        final_turn = final_turns[-1]
        final_agents = final_turn.get("Agents", [])
        if not final_agents:
            continue

        # Overall final / initial ratio
        final_total = len(final_agents)
        ratio = final_total / init_agents
        total_ratio += ratio

        # Count by attachment style
        style_counts = {style: 0 for style in attachment_styles}
        for agent in final_agents:
            style = agent.get("AttachmentStyle")
            if style in style_counts:
                style_counts[style] += 1

        # Update style-specific ratios
        for style in attachment_styles:
            style_ratio = style_counts[style] / init_agents
            style_ratios[style] += style_ratio

    # Store averaged data for this kappa
    avg_total_ratio = total_ratio / iters
    data.append((kappa, "All", avg_total_ratio))
    for style in attachment_styles:
        avg_style_ratio = style_ratios[style] / iters
        data.append((kappa, style, avg_style_ratio))

# Create DataFrame
df = pd.DataFrame(data, columns=["Kappa", "AttachmentStyle", "AvgFinal/Initial"])

# Plot
plt.figure(figsize=(10, 6))
sns.set_context("talk")
sns.lineplot(data=df, x="Kappa", y="AvgFinal/Initial", hue="AttachmentStyle", marker="o")
plt.title("Effect of Cluster Size ($\\kappa$) on Population Survival (by Attachment Style)")
plt.xlabel("Number of Clusters $\\kappa$")
plt.ylabel("Avg Final / Initial Population")
plt.grid(True)
plt.tight_layout()
plt.show()