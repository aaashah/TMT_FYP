import dash
from dash import dcc, html
import pandas as pd
import plotly.express as px
from dash.dependencies import Input, Output, State

# Load simulation data
df = pd.read_csv("visualisation_output/csv_data/agent_records.csv")

# Define grid dimensions
GRID_WIDTH = 70
GRID_HEIGHT = 30
CELL_SIZE = 30  # Defines pixel size of each grid cell

# Ensure consistent colors per agent
unique_agents = df["AgentID"].unique()
agent_colors = {
    agent: px.colors.qualitative.Plotly[i % len(px.colors.qualitative.Plotly)]
    for i, agent in enumerate(unique_agents)
}

# Ensure iteration numbers are integers
df["IterationNumber"] = df["IterationNumber"].astype(int)

# Get the max iteration and turn count
max_iteration = df["IterationNumber"].max()
turns_per_iteration = df.groupby("IterationNumber")["TurnNumber"].max().to_dict()

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
    Updates the grid visualization and handles movement between iterations and turns.
    """
    global turns_per_iteration

    # Determine new turn and iteration
    if next_clicks > prev_clicks:
        if current_turn < turns_per_iteration[current_iteration]:  # Move to next turn
            new_turn = current_turn + 1
            new_iteration = current_iteration
        else:  # Move to next iteration
            new_iteration = min(current_iteration + 1, max_iteration)
            new_turn = 0
    elif prev_clicks < next_clicks:
        if current_turn > 0:  # Move to previous turn
            new_turn = current_turn - 1
            new_iteration = current_iteration
        else:  # Move to previous iteration
            new_iteration = max(current_iteration - 1, 0)
            new_turn = turns_per_iteration[new_iteration]
    else:
        new_iteration = current_iteration
        new_turn = current_turn

    # Filter data for the correct iteration and turn
    filtered_df = df[
        (df["IterationNumber"] == new_iteration) & (df["TurnNumber"] == new_turn)
    ].copy()

    # Adjust positions to center agents in cells
    filtered_df["PositionX"] += 0.5
    filtered_df["PositionY"] += 0.5

    # Create scatter plot
    fig = px.scatter(
        filtered_df,
        x="PositionX",
        y="PositionY",
        color="AgentID",
        color_discrete_map=agent_colors,
        labels={"PositionX": "X Position", "PositionY": "Y Position"},
        title=f"Iteration {new_iteration} - Turn {new_turn}",
    )

    fig.update_traces(marker=dict(size=10))

    # Enforce square cells
    fig.update_layout(
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
