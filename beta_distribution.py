import numpy as np
import matplotlib.pyplot as plt
import scipy.stats as stats

# Beta distribution parameters
alpha = 1.4  # Controls peak
beta = 2.15  # Controls skewness
min_age = 0
max_age = 100

# E(x) = alpha/alpha + beta
# Generate a beta-distributed sample
samples = np.random.beta(alpha, beta, 10000)

# Scale the samples to the age range
ages = min_age + samples * (max_age - min_age)

# Plot the histogram
plt.figure(figsize=(8, 5))
plt.hist(ages, bins=50, density=True, alpha=0.6, color="b", edgecolor="black")

# Plot the theoretical PDF
x = np.linspace(min_age, max_age, 1000)
pdf = stats.beta.pdf((x - min_age) / (max_age - min_age), alpha, beta) / (
    max_age - min_age
)
plt.plot(x, pdf, "r-", label="Beta PDF")

plt.xlabel("Age")
plt.ylabel("Density")
plt.title("Beta Distribution for Agent Age Initialization")
plt.legend()
plt.grid()
plt.show()
