import json
import matplotlib.pyplot as plt

log_dir = "JSONlogs/output.json"
pop_count = []
turn_numbers = []

with open(log_dir, "r") as file:
    GAME_DATA = json.load(file)
    turn_number = 0
    for ITER in GAME_DATA["Iterations"]:
        num_alive = ITER["NumberOfAgents"]
        pop_count.append(num_alive)
        turn_numbers.append(turn_number)
        turn_number += 1

max_agents = max(pop_count)
total_turns = len(turn_numbers)
simplified_x_ticks = range(0, total_turns + 1, 5)

# Plot
plt.figure(figsize=(14, 5))
plt.xlim(0, total_turns)
plt.ylim(0, max_agents + 10)
plt.plot(turn_numbers, pop_count, marker="o", label="Agent Count")
plt.xticks(simplified_x_ticks, rotation=-45)
plt.xlabel("Iteration")
plt.ylabel("Number of Alive Agents")
plt.title("Agent Population Over Time")
plt.grid(True)
# plt.legend()
plt.tight_layout()
plt.show()
