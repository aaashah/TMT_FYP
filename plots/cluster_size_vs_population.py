# import json
# import pandas as pd
# import matplotlib.pyplot as plt
# import seaborn as sns
# import subprocess
# from tqdm import tqdm

# log_path = "JSONlogs/output.json"
# iters = 20
# kappa_values = list(range(1, 11))
# results = []

# for kappa in tqdm(kappa_values, desc="Kappa sweep"):
#     total_final_agents = 0
#     valid_runs = 0

#     for _ in range(iters):
#         subprocess.run([
#             "./tmtSimulator",
#             "-numAgents=40",
#             "-iters=200",
#             f"-kappa={kappa}"
#         ], stdout=subprocess.DEVNULL)

#         try:
#             with open(log_path, "r") as file:
#                 game_data = json.load(file)

#             iterations = game_data.get("Iterations", [])
#             if not iterations:
#                 continue

#             final_iter = iterations[-1]
#             final_turns = final_iter.get("Turns", [])
#             if not final_turns:
#                 continue

#             final_agents = final_turns[-1].get("Agents", [])
#             if not final_agents:
#                 continue

#             total_final_agents += len(final_agents)
#             valid_runs += 1

#         except (json.JSONDecodeError, FileNotFoundError, IndexError, TypeError) as e:
#             print(f"Skipped one run due to error: {e}")
#             continue

#     avg_agents = total_final_agents / valid_runs if valid_runs else 0
#     results.append((kappa, avg_agents))

# # Plot
# df = pd.DataFrame(results, columns=["kappa", "AvgPopulation"])
# plt.figure(figsize=(10, 6))
# sns.lineplot(data=df, x="kappa", y="AvgPopulation", marker="o")
# plt.xlabel("Number of Clusters (kappa)")
# plt.ylabel("Average Population (Final Iteration)")
# plt.title("Effect of Cluster Size on Population")
# plt.grid(True)
# plt.tight_layout()
# plt.show()


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
kappa_values = list(range(1, 11))  # Cluster counts from 1 to 10
results = []

# Sweep over cluster size (kappa)
for kappa in tqdm(kappa_values, desc="Kappa sweep"):
    avg_population_ratio = 0

    for _ in range(iters):
        # Run the simulation with specified number of clusters
        subprocess.run([
            "./tmtSimulator",
            "-numAgents=40",
            "-iters=200",
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
        final_agents = final_iter.get("NumberOfAgents", 0)
        ratio = final_agents / init_agents if init_agents > 0 else 0
        avg_population_ratio += ratio

    avg_population_ratio /= iters
    results.append((kappa, avg_population_ratio))

# Create DataFrame
df = pd.DataFrame(results, columns=["Kappa", "Final/Initial Population"])

# Plot
plt.figure(figsize=(10, 6))
sns.set_context("talk")
sns.lineplot(data=df, x="Kappa", y="Final/Initial Population", marker="o")
plt.title("Effect of Cluster Size ($\\kappa$) on Population Survival")
plt.xlabel("Number of Clusters $\\kappa$")
plt.ylabel("Avg Final / Initial Population")
plt.grid(True)
plt.tight_layout()
plt.show()