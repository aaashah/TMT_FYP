import dash
from dash import dcc, html
import pandas as pd
import plotly.graph_objects as go
from dash.dependencies import Input, Output, State
import ast  # Safer parsing of string lists
import colorsys  # For generating unique colors
from dash import callback_context

# Load simulation data
agent_df = pd.read_csv("visualisation_output/csv_data/agent_records.csv")
infra_df = pd.read_csv("visualisation_output/csv_data/infra_records.csv")

# Define grid dimensions
GRID_WIDTH = 70
GRID_HEIGHT = 30
CELL_SIZE = 30  # Defines pixel size of each grid cell


# Function to generate unique HSL-based colors for agents
def generate_color(index, total_agents):
    hue = (index / total_agents) % 1
    rgb = colorsys.hsv_to_rgb(hue, 0.85, 0.85)
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

# Parse tombstone positions
tombstone_dict = {}
for _, row in infra_df.iterrows():
    iteration, turn = row["IterationNumber"], row["TurnNumber"]
    try:
        tombstones = (
            ast.literal_eval(row["Tombstones"])
            if isinstance(row["Tombstones"], str) and row["Tombstones"] != "[]"
            else []
        )
    except (SyntaxError, ValueError):
        tombstones = []
    if (iteration, turn) not in tombstone_dict:
        tombstone_dict[(iteration, turn)] = set(tombstones)
    if turn > 0:
        prev_tombstones = tombstone_dict.get((iteration, turn - 1), set())
        tombstone_dict[(iteration, turn)] |= prev_tombstones

# Get max iteration and turn count
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
                html.Button(
                    "▶ Play",
                    id="play-pause",
                    n_clicks=0,
                    style={"font-size": "20px", "margin": "0 10px"},
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
        html.Div(
            id="iteration-turn-label",
            style={"textAlign": "center", "font-size": "20px"},
        ),
        dcc.Graph(id="grid-plot"),
        dcc.Store(id="iteration-store", data=0),
        dcc.Store(id="turn-store", data=0),
        dcc.Store(id="animation-state", data=True),  # ✅ Start Paused
        dcc.Interval(
            id="animation-interval", interval=1000, n_intervals=0, disabled=False
        ),  # ✅ Start Disabled
    ],
    style={"width": "95%", "margin": "auto"},
)


# Play/Pause button callback
@app.callback(
    [
        Output("animation-interval", "disabled"),
        Output("play-pause", "children"),
        Output("animation-state", "data"),
    ],
    [Input("play-pause", "n_clicks")],
    [State("animation-state", "data")],
)
def toggle_animation(n_clicks, is_playing):
    if is_playing:
        return True, "▶ Play", False  # ✅ Start Paused
    else:
        return False, "⏸ Pause", True  # ✅ Play when clicked


# Grid update callback
@app.callback(
    [
        Output("grid-plot", "figure"),
        Output("iteration-turn-label", "children"),
        Output("iteration-store", "data"),
        Output("turn-store", "data"),
    ],
    [
        Input("prev-turn", "n_clicks"),
        Input("next-turn", "n_clicks"),
        Input("animation-interval", "n_intervals"),
    ],
    [State("iteration-store", "data"), State("turn-store", "data")],
)
def update_grid(prev_clicks, next_clicks, n_intervals, current_iteration, current_turn):
    global turns_per_iteration

    # Ensure turn structure is valid
    max_turns_in_current_iteration = turns_per_iteration.get(current_iteration, 0)

    # ✅ Detect which button was clicked
    ctx = callback_context
    if not ctx.triggered:
        triggered_id = None
    else:
        triggered_id = ctx.triggered[0]["prop_id"].split(".")[0]

    # ✅ Handle "Next" button & Animation
    if triggered_id in ["next-turn", "animation-interval"]:
        if current_turn < max_turns_in_current_iteration:
            new_turn = current_turn + 1
            new_iteration = current_iteration
        else:
            if current_iteration < max_iteration:
                new_iteration = current_iteration + 1
                new_turn = 0
            else:
                new_iteration = max_iteration
                new_turn = max_turns_in_current_iteration

    # ✅ Handle "Previous" button
    elif triggered_id == "prev-turn":
        if current_turn > 0:
            new_turn = current_turn - 1
            new_iteration = current_iteration
        else:
            if current_iteration > 0:
                new_iteration = current_iteration - 1
                new_turn = turns_per_iteration.get(new_iteration, 0)
            else:
                new_iteration = 0
                new_turn = 0

    else:
        new_iteration, new_turn = current_iteration, current_turn

    # ✅ Filter agent data
    filtered_df = agent_df[
        (agent_df["IterationNumber"] == new_iteration)
        & (agent_df["TurnNumber"] == new_turn)
    ].copy()

    # ✅ Retrieve tombstones
    tombstones = list(tombstone_dict.get((new_iteration, new_turn), []))

    # ✅ Convert tombstone positions
    tombstone_x = [pos[0] + 0.5 for pos in tombstones]
    tombstone_y = [pos[1] + 0.5 for pos in tombstones]

    # ✅ Initialize figure
    fig = go.Figure()

    # ✅ Add agent positions
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

    # ✅ Add tombstones
    for x, y in zip(tombstone_x, tombstone_y):
        fig.add_annotation(
            x=x, y=y, text="💀", showarrow=False, font=dict(size=15, color="black")
        )

    # ✅ Enforce square grid layout
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
