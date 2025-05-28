import json
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import subprocess
from tqdm import tqdm

# Constants
log_path = "JSONlogs/output.json"
iters = 20
kappa_values = list(range(1, 11))  # κ from 1 to 10
results = []

# Loop over kappa
for kappa in tqdm(kappa_values, desc="Kappa sweep"):
    total_volunteers = 0
    total_turns = 0

    for _ in range(iters):
        subprocess.run([
            "./tmtSimulator",
            "-numAgents=40",
            "-iters=200",
            f"-kappa={kappa}"
        ], stdout=subprocess.DEVNULL)

        with open(log_path, "r") as file:
            game_data = json.load(file)

        # Traverse all iterations and turns
        for iteration in game_data.get("Iterations", []):
            for turn in iteration.get("Turns", []):
                vols = turn.get("NumVolunteers", 0)
                total_volunteers += vols
                total_turns += 1

    avg_volunteers = total_volunteers / total_turns if total_turns > 0 else 0
    results.append((kappa, avg_volunteers))

# Create DataFrame
df = pd.DataFrame(results, columns=["kappa", "AverageVolunteers"])

# Plot
plt.figure(figsize=(10, 6))
sns.lineplot(data=df, x="kappa", y="AverageVolunteers", marker="o")
plt.xlabel("Number of Clusters $\\kappa$")
plt.ylabel("Avg Voluntary Self-Sacrifices per Turn")
plt.title("Effect of Cluster Size ($\\kappa$) on Volunteering Rates")
plt.grid(True)
plt.tight_layout()
plt.show()