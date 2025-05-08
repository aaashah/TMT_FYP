import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import numpy as np

# Example: simulate agent score data
np.random.seed(42)
data = []

for agent_id in range(100):
    agent_name = f"agent_{agent_id}"
    start = np.random.randint(0, 20)
    length = np.random.randint(10, 50)
    for t in range(start, start + length):
        score = np.random.rand()
        data.append((agent_name, t, score))

# Create DataFrame
df = pd.DataFrame(data, columns=["agent_id", "iteration", "score"])

# Pivot the data to get agents as rows, iterations as columns
heatmap_data = df.pivot(index="agent_id", columns="iteration", values="score")

# Plot
plt.figure(figsize=(14, 12))
sns.heatmap(heatmap_data, cmap="viridis", cbar_kws={'label': 'Score'}, linewidths=0.1, linecolor='gray')
plt.title("Agent Scores Over Time")
plt.xlabel("Iteration")
plt.ylabel("Agent ID")
plt.tight_layout()
plt.show()
