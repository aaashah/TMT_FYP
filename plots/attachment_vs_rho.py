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
rho_values = [round(i * 0.1, 1) for i in range(11)]
attachment_styles = ["Secure", "Dismissive", "Preoccupied", "Fearful"]

results = []

# Sweep over rho values
for rho in tqdm(rho_values, desc="Rho sweep"):
    style_counts = {style: 0 for style in attachment_styles}
    total_agents = 0

    for _ in range(iters): 
        # Run the simulation with specified rho
        subprocess.run([
            "./tmtSimulator",
            "-numAgents=40",
            "-iters=200",
            f"-rho={rho}"
        ], stdout=subprocess.DEVNULL)

        # Load JSON output
        with open(log_path, "r") as file:
            game_data = json.load(file)

        iterations = game_data.get("Iterations", [])
        if not iterations:
            continue

        final_iter = iterations[-1]
        final_turns = final_iter.get("Turns", [])
        if not final_turns:
            continue

        final_turn = final_turns[-1]
        final_agents = final_turn.get("Agents")
        if not final_agents:
            continue

        total_agents += len(final_agents)

        for agent in final_agents:
            style = agent.get("AttachmentStyle")
            if style in style_counts:
                style_counts[style] += 1

    # Average and store results
    for style in attachment_styles:
        avg = style_counts[style] / iters
        results.append((rho, style, avg))

# Create DataFrame
df = pd.DataFrame(results, columns=["rho", "AttachmentStyle", "Count"])

# Plot
plt.figure(figsize=(10, 6))
sns.lineplot(data=df, x="rho", y="Count", hue="AttachmentStyle", marker="o")
plt.xlabel(r"$\rho$")
plt.ylabel("Avg Count of Agents (Final Iteration)")
plt.title("Attachment Style Survival vs rho")
plt.grid(True)
plt.tight_layout()
plt.show()