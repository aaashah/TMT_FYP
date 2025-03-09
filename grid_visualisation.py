import dash
from dash import dcc, html
import pandas as pd
import plotly.graph_objects as go
from dash.dependencies import Input, Output, State
import ast  # Safer parsing of string lists
import colorsys  # For generating unique colors

# Load simulation data
agent_df = pd.read_csv("visualisation_output/csv_data/agent_records.csv")
infra_df = pd.read_csv("visualisation_output/csv_data/infra_records.csv")

# Define grid dimensions
GRID_WIDTH = 70
GRID_HEIGHT = 30
CELL_SIZE = 30  # Defines pixel size of each grid cell


# Function to generate unique HSL-based colors for agents
def generate_color(index, total_agents):
    hue = (index / total_agents) % 1  # Evenly space hues
    rgb = colorsys.hsv_to_rgb(hue, 0.85, 0.85)  # Convert HSV to RGB
    return f"rgb({int(rgb[0]*255)}, {int(rgb[1]*255)}, {int(rgb[2]*255)})"


# Ensure consistent colors per agent
unique_agents = agent_df["AgentID"].unique()
agent_colors = {
    agent: generate_color(i, len(unique_agents))
    for i, agent in enumerate(unique_agents)
}

# Convert IterationNumber to int for correct sorting
agent_df["IterationNumber"] = agent_df["IterationNumber"].astype(int)
infra_df["IterationNumber"] = infra_df["IterationNumber"].astype(int)

# Parse tombstone positions from infra_records.csv
tombstone_dict = {}  # Key: (iteration, turn) → Set of tombstone positions

for _, row in infra_df.iterrows():
    iteration = row["IterationNumber"]
    turn = row["TurnNumber"]

    # Safe parsing of tombstone positions
    try:
        tombstones = (
            ast.literal_eval(row["Tombstones"])
            if isinstance(row["Tombstones"], str) and row["Tombstones"] != "[]"
            else []
        )
    except (SyntaxError, ValueError):
        tombstones = []

    # Store tombstones for all future turns
    if (iteration, turn) not in tombstone_dict:
        tombstone_dict[(iteration, turn)] = set(tombstones)

    # Carry forward tombstones from previous turns
    if turn > 0:
        prev_tombstones = tombstone_dict.get((iteration, turn - 1), set())
        tombstone_dict[
            (iteration, turn)
        ] |= prev_tombstones  # Merge previous tombstones

# Get the max iteration and turn count
max_iteration = agent_df["IterationNumber"].max()
turns_per_iteration = agent_df.groupby("IterationNumber")["TurnNumber"].max().to_dict()

# Initialize Dash app
app = dash.Dash(__name__)

# Layout
app.layout = html.Div(
    [
        html.H2("Agent Movement Grid"),
        html.Div(
            [
                html.Button(
                    "⬅️ Previous",
                    id="prev-turn",
                    n_clicks=0,
                    style={"font-size": "20px"},
                ),
                html.Span(
                    id="iteration-turn-label",
                    style={"font-size": "20px", "margin": "0 20px"},
                ),
                html.Button(
                    "Next ➡️", id="next-turn", n_clicks=0, style={"font-size": "20px"}
                ),
            ],
            style={
                "display": "flex",
                "justify-content": "center",
                "align-items": "center",
                "margin-bottom": "10px",
            },
        ),
        dcc.Graph(id="grid-plot"),
        dcc.Store(id="iteration-store", data=0),  # Track iteration
        dcc.Store(id="turn-store", data=0),  # Track turn within iteration
    ],
    style={"width": "95%", "margin": "auto"},
)


# Callback to update the grid and iteration/turn label
@app.callback(
    [
        Output("grid-plot", "figure"),
        Output("iteration-turn-label", "children"),
        Output("iteration-store", "data"),
        Output("turn-store", "data"),
    ],
    [Input("prev-turn", "n_clicks"), Input("next-turn", "n_clicks")],
    [State("iteration-store", "data"), State("turn-store", "data")],
)
def update_grid(prev_clicks, next_clicks, current_iteration, current_turn):
    """
    Handles movement through iterations and turns with proper logic.
    """
    global turns_per_iteration

    # Ensure turn structure is valid
    max_turns_in_current_iteration = turns_per_iteration.get(current_iteration, 0)

    # Handle "Next" button
    if next_clicks > prev_clicks:
        if current_turn < max_turns_in_current_iteration:
            new_turn = current_turn + 1  # Move forward one turn
            new_iteration = current_iteration
        else:
            if current_iteration < max_iteration:
                new_iteration = current_iteration + 1  # Move to next iteration
                new_turn = 0  # Reset to first turn in the new iteration
            else:
                new_iteration = max_iteration  # Stay at last iteration
                new_turn = max_turns_in_current_iteration  # Stay at last turn

    # Handle "Previous" button
    elif prev_clicks > next_clicks:
        if current_turn > 0:
            new_turn = current_turn - 1  # Move back one turn
            new_iteration = current_iteration
        else:
            if (
                current_iteration > 0
            ):  # If at turn 0, move to the last turn of the previous iteration
                new_iteration = current_iteration - 1
                new_turn = turns_per_iteration[
                    new_iteration
                ]  # Last turn of previous iteration
            else:
                new_iteration = 0  # Already at first iteration
                new_turn = 0  # Stay at turn 0

    else:
        new_iteration = current_iteration  # No clicks, stay in place
        new_turn = current_turn

    # **Filter agent data for correct iteration & turn**
    filtered_df = agent_df[
        (agent_df["IterationNumber"] == new_iteration)
        & (agent_df["TurnNumber"] == new_turn)
    ].copy()

    # Retrieve tombstones that exist up to this turn
    tombstones = list(tombstone_dict.get((new_iteration, new_turn), []))

    # Convert tombstone positions into X and Y lists
    tombstone_x = [pos[0] + 0.5 for pos in tombstones]
    tombstone_y = [pos[1] + 0.5 for pos in tombstones]

    # Initialize figure
    fig = go.Figure()

    # Add agent positions
    for agent_id in filtered_df["AgentID"].unique():
        agent_data = filtered_df[filtered_df["AgentID"] == agent_id]
        fig.add_trace(
            go.Scatter(
                x=agent_data["PositionX"] + 0.5,
                y=agent_data["PositionY"] + 0.5,
                mode="markers",
                marker=dict(size=10, color=agent_colors.get(agent_id, "gray")),
                name=f"Agent {agent_id}",
            )
        )

    # Add tombstones as text annotations
    for x, y in zip(tombstone_x, tombstone_y):
        fig.add_annotation(
            x=x,
            y=y,
            text="💀",
            showarrow=False,
            font=dict(size=15, color="black"),
        )

    # Enforce square grid layout
    fig.update_layout(
        title=f"Iteration {new_iteration} - Turn {new_turn}",
        xaxis=dict(
            range=[0, GRID_WIDTH],
            tickmode="linear",
            dtick=1,
            showgrid=True,
            gridcolor="lightgray",
            zeroline=False,
        ),
        yaxis=dict(
            range=[0, GRID_HEIGHT],
            tickmode="linear",
            dtick=1,
            showgrid=True,
            gridcolor="lightgray",
            zeroline=False,
            scaleanchor="x",
        ),
        autosize=True,
        height=GRID_HEIGHT * CELL_SIZE,
        width=GRID_WIDTH * CELL_SIZE,
        showlegend=True,
        plot_bgcolor="white",
    )

    return fig, f"Iteration {new_iteration} - Turn {new_turn}", new_iteration, new_turn


# Run the app
if __name__ == "__main__":
    app.run_server(debug=True)
