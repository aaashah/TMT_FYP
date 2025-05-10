import json
import pandas as pd
import seaborn as sns
import matplotlib.pylab as plt
import subprocess
from tqdm import tqdm

log_dir = "JSONlogs/output.json"
iters = 20
data = []

for tau in tqdm([0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0]):
    for rho in tqdm([0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0]):
        pop_change = 0
        for _ in tqdm(range(iters)):
            subprocess.run(
                [
                    "./tmtSimulator",
                    "-numAgents=40",
                    "-iters=200",
                    f"-rho={rho}",
                    f"-tau={tau}",
                ]
            )

            with open(log_dir, "r") as file:
                GAME_DATA = json.load(file)
                CONFIG = GAME_DATA["Config"]
                init_agents = CONFIG["NumAgents"]
                FINAL_ITER = GAME_DATA["Iterations"][-1]
                final_agents = FINAL_ITER["NumberOfAgents"]
                pop_change += final_agents / init_agents

        pop_change /= iters
        data.append((tau, rho, min(pop_change, 3)))

# Create a DataFrame
df = pd.DataFrame(data, columns=["tau", "rho", "value"])
# Pivot into 2D form (rows = rho, columns = tau)
pivot = df.pivot(index="rho", columns="tau", values="value")

pivot.to_pickle("figures/rho_vs_tau.pkl")

# pivot = pd.read_pickle("figures/rho_vs_tau.pkl")

plt.figure()
sns.heatmap(pivot)
plt.gca().invert_yaxis()  # Bottom to top
plt.xlabel(r"$\tau$")
plt.ylabel(r"$\rho$")
plt.show()
