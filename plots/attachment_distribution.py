import json
import matplotlib.pyplot as plt
from collections import defaultdict

log_dir = "JSONlogs/output.json"
attachment_data = {"Secure": [], "Dismissive": [], "Preoccupied": [], "Fearful": []}
turn_numbers = []


with open(log_dir, "r") as file:
    GAME_DATA = json.load(file)
    turn_number = 0
    for ITER in GAME_DATA["Iterations"]:
        for TURN in ITER["Turns"]:
            agents = TURN.get("Agents", [])
            total = max(len(agents), 1)
            counter = defaultdict(int)
            for agent in agents:
                style = agent.get("AttachmentStyle", "Unknown")
                counter[style] += 1

            assert counter["Unknown"] == 0

            for style in attachment_data.keys():
                attachment_data[style].append(counter[style] / total)

            turn_numbers.append(turn_number)
            turn_number += 1


# Step 2: Plot each style over time
total_turns = len(turn_numbers)
simplified_x_ticks = range(0, total_turns + 1, 5)
color_map = {
    "Secure": "green",
    "Dismissive": "red",
    "Preoccupied": "blue",
    "Fearful": "purple",
}

plt.figure(figsize=(14, 6))
for style, color in color_map.items():
    style_prop = attachment_data[style]
    plt.plot(turn_numbers, style_prop, label=style, marker="o", color=color)

# Show all round numbers on x-axis
plt.xticks(ticks=simplified_x_ticks, rotation=-45)
plt.ylim(0, 1)
plt.ylabel("Proportion of Population")
plt.xlabel("Turn")
plt.title("Attachment Type Distribution Over Time")
plt.legend()
plt.tight_layout()
plt.show()

# stacked area plot (primer special)
plt.figure(figsize=(14, 6))
plt.stackplot(
    turn_numbers,
    attachment_data.values(),  # unpack all the arrays
    labels=attachment_data.keys(),  # use the keys as labels
    alpha=0.8,
    edgecolor="white",
)

plt.ylim(bottom=0, top=1)
plt.xlim(left=0, right=len(simplified_x_ticks))
plt.ylabel("Proportion of Population")
plt.xlabel("Turn")
plt.xticks(ticks=simplified_x_ticks, rotation=-45)
plt.title("Attachment Type Distribution Over Time")
plt.legend()
plt.show()
