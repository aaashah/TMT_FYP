# A(Life and Death): A Computational Model of Terror Management Theory

## Overview

This repository contains the codebase for the self-organising multi-agent simulation described in the A(Life and Death) paper. The simulation is designed to model the effects of mortality salience on behviour, based on Terror Management Theory (Greenberg, Pyszczynski, & Solomon, 1986). Agents vary by attachment style and make decisions to self-sacrifice based on internal drives for worldview validation, relationship validation, and survival.

## Core Features

- **Attachment Styles**: Secure, Dismissive, Preoccupied, Fearful — based on Bartholomew and Horowitz (1991).
- **ASM Modules**: Self-sacrifice decisions are driven by:
  - **Mortality Salience**
  - **Worldview Validation**
  - **Relationship Validation**
- **Elimination Logic**:
  - Natural death via aging (telomere).
  - Voluntary self-sacrifice.
  - Involuntary elimination if minimum eliminations not met.

## Project Structure

```
TMT_Attachment
├── agents/             # Agent implementations for each attachment style
├── config/             # Simulation configuration
├── gameRecorder/       # Logging utilities for recording simulation data
├── infra/              # Core types, interfaces, and grid logic
├── JSONlogs/           # JSON files containing simulation logs per turn
├── plots/              # Python scripts for analyzing and visualizing results
├── server/             # Main server loop and elimination mechanics
├── grid_visualiser.py  # Agent grid visualisation tool
├── main.go             # Entry point for running the simulation
```

## Running the Simulation

Ensure prerequisites are installed, then run:

```bash
go run main.go
```

From the root directory.

JSON logs will be output to the `JSONlogs/` folder.

## Plotting and Visualisation

Python plotting scripts are provided in the `plots/` directory.

Visualisation of the agent simulation on the grid environment is also provided in the Python script `grid_visualiser.py`.
