import json
import os
import colorsys
import plotly.graph_objects as go
import dash
from dash import dcc, html, callback_context
from dash.dependencies import Input, Output, State

# --- Constants ---
LOG_DIR = "JSONlogs/output.json"
GRID_WIDTH = 50
GRID_HEIGHT = 50
CELL_SIZE = 30


# --- Load JSON Logs ---
iteration_logs = {}
max_iteration = -1
turns_per_iteration = {}


with open(LOG_DIR, "r") as file:
    GAME_DATA = json.load(file)
    turn_number = 0
    for ITER in GAME_DATA["Iterations"]:
        iter_num = ITER["Iteration"]
        TURN_DATA = ITER["Turns"]
        iteration_logs[iter_num] = TURN_DATA
        turns_per_iteration[iter_num] = len(TURN_DATA) - 1
        max_iteration = max(max_iteration, iter_num)

# --- Agent Colors ---

color_map = {
    "Secure": "green",
    "Dismissive": "red",
    "Preoccupied": "blue",
    "Fearful": "purple",
}

agent_colors = {}
for turns in iteration_logs.values():
    for turn in turns:
        for agent in turn.get("Agents") or []:
            agent_id = agent["ID"]
            attch = agent["AttachmentStyle"]
            agent_colors[agent_id] = color_map[attch]


# Legend for color map
style_map = [
    html.Div(
        children=[
            html.Div(
                style={
                    "width": "15px",
                    "height": "15px",
                    "background-color": color,
                    "display": "inline-block",
                    "margin-right": "10px",
                }
            ),
            html.H3(label),
        ],
        style={"margin": "10px", "display": "flex", "align-items": "center"},
    )
    for label, color in color_map.items()
]


# --- Dash App Setup ---
app = dash.Dash(__name__)

app.layout = html.Div(
    [
        html.H1("Animated Agent Movement", style={"textAlign": "center"}),
        html.Div(
            [
                html.Button(
                    "⬅️ Previous",
                    id="prev-turn",
                    n_clicks=0,
                    style={"font-size": "20px"},
                ),
                html.Button(
                    "▶️ Play",
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
        html.Div(
            [
                # dcc.Graph(figure=style_map),
                html.Div(style_map, style={"flexDirection": "column"}),
                dcc.Graph(id="grid-plot"),
                dcc.Store(id="iteration-store", data=0),
                dcc.Store(id="turn-store", data=0),
                dcc.Store(id="animation-state", data=True),
                dcc.Interval(
                    id="animation-interval", interval=100, n_intervals=0, disabled=False
                ),
            ],
            id="plot-body",
            style={
                "display": "flex",
                "flexDirection": "row",
                "justifyContent": "center",
                "alignContent": "center",
                "alignItems": "center",
            },
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
    return (True, "▶️ Play", False) if is_playing else (False, "⏸️ Pause", True)


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

    fig.update_layout(
        title=f"Iteration {current_iteration} - Turn {current_turn}",
        xaxis=dict(
            range=[-1, GRID_WIDTH + 1],
            tickvals=[i for i in range(GRID_WIDTH + 1)],
            ticktext=[str(i) for i in range(GRID_WIDTH + 1)],
            showgrid=True,
            gridcolor="lightgray",
        ),
        yaxis=dict(
            range=[-1, GRID_HEIGHT + 1],
            tickvals=[i for i in range(GRID_HEIGHT + 1)],
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

    return (
        fig,
        f"Iteration {current_iteration} - Turn {current_turn}",
        current_iteration,
        current_turn,
    )


# --- Run App ---
if __name__ == "__main__":
    app.run(debug=True)
