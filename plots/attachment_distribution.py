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
                
                agents = turn.get("Agents") or [] 
                total = len(agents)
                
                counter = defaultdict(int)
                for agent in agents:
                    style = agent.get("AttachmentStyle", "Unknown")
                    counter[style] += 1

                if total > 0:
                    attachment_data.append({
                        "Label": f"i{iteration_num}_t{turn_num}",
                        "Secure": counter["Secure"] / total,
                        "Dismissive": counter["Dismissive"] / total,
                        "Preoccupied": counter["Preoccupied"] / total,
                        "Fearful": counter["Fearful"] / total
                    })

# Step 2: Plot each style over time
x = [entry["Label"] for entry in attachment_data]
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
    plt.plot(x, y, label=style, marker='o', color=colors[style])

plt.xticks(rotation=90, fontsize=8)
plt.ylim(0, 1)
plt.ylabel("Proportion of Population")
plt.xlabel("Iteration_Turn")
plt.title("Attachment Type Distribution Over Time")
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.show()