import json
import os
import matplotlib.pyplot as plt
from collections import defaultdict

log_dir = "JSONlogs"
attachment_data = []

# Step 1: Parse each turn's agent list and count attachment styles
for filename in sorted(os.listdir(log_dir)):
    if filename.startswith("iteration_") and filename.endswith(".json"):
        with open(os.path.join(log_dir, filename), "r") as f:
            data = json.load(f)
            iteration_num = data["Iteration"]
            for turn in data["Turns"]:
                turn_num = turn["TurnNumber"]

                # Handle missing or null agent list
                agents = turn.get("Agents")
                if agents is None:
                    agents = []

                total = len(agents)
                counter = defaultdict(int)
                for agent in agents:
                    style = agent.get("AttachmentStyle", "Unknown")
                    counter[style] += 1

                # Always append, even when total == 0, to ensure we show 0s
                attachment_data.append({
                    "Label": f"i{iteration_num}_t{turn_num}",
                    "Secure": counter["Secure"] / total if total > 0 else 0,
                    "Dismissive": counter["Dismissive"] / total if total > 0 else 0,
                    "Preoccupied": counter["Preoccupied"] / total if total > 0 else 0,
                    "Fearful": counter["Fearful"] / total if total > 0 else 0
                })

# Step 2: Plot each style over time
rounds = list(range(len(attachment_data)))
styles = ["Secure", "Dismissive", "Preoccupied", "Fearful"]
colors = {
    "Secure": "green",
    "Dismissive": "red",
    "Preoccupied": "blue",
    "Fearful": "purple"
}

plt.figure(figsize=(14, 6))
for style in styles:
    y = [entry.get(style, 0) for entry in attachment_data]
    plt.plot(rounds, y, label=style, marker='o', color=colors[style])

# Show all round numbers on x-axis
plt.xticks(ticks=rounds, labels=rounds, rotation=90, fontsize=8)

plt.ylim(0, 1)
plt.ylabel("Proportion of Population")
plt.xlabel("Rounds")
plt.title("Attachment Type Distribution Over Time")
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.show()