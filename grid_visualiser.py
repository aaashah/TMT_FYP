import json
import os
import colorsys
import plotly.graph_objects as go
import dash
from dash import dcc, html, callback_context
from dash.dependencies import Input, Output, State

# --- Constants ---
LOG_DIR = "JSONlogs"
GRID_WIDTH = 70
GRID_HEIGHT = 30
CELL_SIZE = 30


# --- Helper Functions ---
def generate_color(index, total_agents):
    hue = (index / total_agents) % 1
    rgb = colorsys.hsv_to_rgb(hue, 0.85, 0.85)
    return f"rgb({int(rgb[0]*255)}, {int(rgb[1]*255)}, {int(rgb[2]*255)})"


# --- Load JSON Logs ---
iteration_logs = {}
max_iteration = -1
turns_per_iteration = {}

for filename in sorted(os.listdir(LOG_DIR)):
    if filename.startswith("iteration_") and filename.endswith(".json"):
        with open(os.path.join(LOG_DIR, filename), "r") as f:
            data = json.load(f)
            iteration = data["Iteration"]
            iteration_logs[iteration] = data["Turns"]
            turns_per_iteration[iteration] = len(data["Turns"]) - 1
            max_iteration = max(max_iteration, iteration)

# --- Agent Colors ---
unique_agent_ids = set()
for turns in iteration_logs.values():
    for turn in turns:
        for agent in turn.get("Agents") or []:
            unique_agent_ids.add(agent["ID"])

agent_colors = {
    agent_id: generate_color(i, len(unique_agent_ids))
    for i, agent_id in enumerate(sorted(unique_agent_ids))
}

# --- Dash App Setup ---
app = dash.Dash(__name__)

app.layout = html.Div(
    [
        html.H2("Agent Movement Grid (from JSON Logs)"),
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
        dcc.Store(id="animation-state", data=True),
        dcc.Interval(
            id="animation-interval", interval=1000, n_intervals=0, disabled=False
        ),
    ]
)


# --- Callbacks ---
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
    return (True, "▶ Play", False) if is_playing else (False, "⏸ Pause", True)


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
    ctx = callback_context
    triggered_id = ctx.triggered[0]["prop_id"].split(".")[0] if ctx.triggered else None
    max_turns = turns_per_iteration.get(current_iteration, 0)

    if triggered_id in ["next-turn", "animation-interval"]:
        if current_turn < max_turns:
            current_turn += 1
        elif current_iteration < max_iteration:
            current_iteration += 1
            current_turn = 0
    elif triggered_id == "prev-turn":
        if current_turn > 0:
            current_turn -= 1
        elif current_iteration > 0:
            current_iteration -= 1
            current_turn = turns_per_iteration.get(current_iteration, 0)

    turn_data = iteration_logs[current_iteration][current_turn]
    agents = turn_data.get("Agents", [])
    tombstones = turn_data.get("TombstoneLocations") or []
    temples = turn_data.get("TempleLocations") or []

    fig = go.Figure()

    for agent in agents:
        agent_id = agent["ID"]
        x = agent["Position"]["X"]
        y = agent["Position"]["Y"]
        fig.add_trace(
            go.Scatter(
                x=[x],
                y=[y],
                mode="markers",
                marker=dict(size=10, color=agent_colors.get(agent_id, "gray")),
                name=f"Agent {agent_id}",
            )
        )

    for t in tombstones:
        fig.add_annotation(
            x=t["X"], y=t["Y"], text="💀", showarrow=False, font=dict(size=15)
        )

    for t in temples:
        fig.add_annotation(
            x=t["X"], y=t["Y"], text="🏛️", showarrow=False, font=dict(size=15)
        )

    fig.update_layout(
        title=f"Iteration {current_iteration} - Turn {current_turn}",
        xaxis=dict(
            range=[-0.5, GRID_WIDTH + 1],
            tickvals=[i + 0.5 for i in range(GRID_WIDTH + 1)],
            ticktext=[str(i) for i in range(GRID_WIDTH + 1)],
            showgrid=True,
            gridcolor="lightgray",
        ),
        yaxis=dict(
            range=[-0.5, GRID_HEIGHT + 1],
            tickvals=[i + 0.5 for i in range(GRID_HEIGHT + 1)],
            ticktext=[str(i) for i in range(GRID_HEIGHT + 1)],
            showgrid=True,
            gridcolor="lightgray",
        ),
        height=GRID_HEIGHT * CELL_SIZE,
        width=GRID_WIDTH * CELL_SIZE,
        showlegend=True,
        plot_bgcolor="white",
        autosize=True,
    )

    return (
        fig,
        f"Iteration {current_iteration} - Turn {current_turn}",
        current_iteration,
        current_turn,
    )


# --- Run App ---
if __name__ == "__main__":
    app.run_server(debug=True)
