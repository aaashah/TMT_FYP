import json
import pandas as pd
import random
import seaborn as sns
import matplotlib.pylab as plt
import os

log_dir = "JSONlogs/output.json"
iters = 1
data = []

for tau in [0, 0.5, 1.0]:
    for rho in [0, 0.5, 1.0]:
        pop_change = 0
        for _ in range(iters):
            os.system(f"go run main.go -numAgents=10 -iters=10 -rho={rho} -tau={tau}")

            with open(log_dir, "r") as file:
                GAME_DATA = json.load(file)
                CONFIG = GAME_DATA["Config"]
                init_agents = CONFIG["NumAgents"]
                FINAL_ITER = GAME_DATA["Iterations"][-1]
                FINAL_TURN = FINAL_ITER["Turns"][-1]
                final_agents = FINAL_TURN["NumberOfAgents"]
                pop_change += final_agents / init_agents

        pop_change /= iters
        data.append((tau, rho, pop_change))

# Create a DataFrame
df = pd.DataFrame(data, columns=["tau", "rho", "value"])
# Pivot into 2D form (rows = rho, columns = tau)
pivot = df.pivot(index="rho", columns="tau", values="value")

pivot.to_pickle("figures/rho_vs_tau.pkl")

plt.figure()
sns.heatmap(pivot)
plt.show()
