import json
import matplotlib.pyplot as plt

log_dir = "JSONlogs/output.json"
volunteers = []
required = []
turn_numbers = []
rho = None

with open(log_dir, "r") as file:
    GAME_DATA = json.load(file)
    config = GAME_DATA["Config"]
    rho = config["PopulationRho"]
    turn_number = 0
    for ITER in GAME_DATA["Iterations"]:
        for TURN in ITER["Turns"]:
            num_vols = len(TURN.get("EliminatedBySelfSacrifice", []))
            num_req = TURN["TotalRequiredEliminations"]
            volunteers.append(num_vols)
            required.append(num_req)
            turn_numbers.append(turn_number)
            turn_number += 1


total_turns = len(turn_numbers)
simplified_x_ticks = range(0, total_turns + 1, 5)

# Plot
plt.figure(figsize=(14, 5))
plt.plot(
    turn_numbers, volunteers, marker="o", label="Number of Volunteers", color="green"
)
plt.plot(turn_numbers, required, linestyle="--", label="Number Required", color="red")
plt.xticks(simplified_x_ticks, rotation=-45)
plt.xlabel("Turn")
plt.ylabel("Number of Agents")
plt.title(rf"Volunteers vs Required Eliminations per Turn ($\rho$={rho})")
plt.ylim(0, max(volunteers + required) + 1)  # Auto-scaled y-axis
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.show()
