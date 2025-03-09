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

# Get iteration numbers
df["IterationNumber"] = df["IterationNumber"].astype(int)  # Ensure it's an integer

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
        dcc.Store(id="turn-store", data=df["TurnNumber"].min()),  # Store current turn
    ],
    style={"width": "95%", "margin": "auto"},
)


# Callback to update the grid and iteration/turn label
@app.callback(
    [
        Output("grid-plot", "figure"),
        Output("iteration-turn-label", "children"),
        Output("turn-store", "data"),
    ],
    [Input("prev-turn", "n_clicks"), Input("next-turn", "n_clicks")],
    [State("turn-store", "data")],
)
def update_grid(prev_clicks, next_clicks, current_turn):
    # Determine new turn number
    total_turns = df["TurnNumber"].max()
    new_turn = current_turn + (1 if next_clicks > prev_clicks else -1)

    # Keep within valid range
    new_turn = max(df["TurnNumber"].min(), min(new_turn, total_turns))

    # Get corresponding iteration number
    current_iteration = df[df["TurnNumber"] == new_turn]["IterationNumber"].iloc[0]

    filtered_df = df[df["TurnNumber"] == new_turn].copy()

    # Center agents within cells
    filtered_df["PositionX"] += 0.5
    filtered_df["PositionY"] += 0.5

    # Ensure unique agent positions per turn
    filtered_df = filtered_df.drop_duplicates(subset=["AgentID", "TurnNumber"])

    fig = px.scatter(
        filtered_df,
        x="PositionX",
        y="PositionY",
        color="AgentID",
        color_discrete_map=agent_colors,
        labels={"PositionX": "X Position", "PositionY": "Y Position"},
        title=f"Iteration {current_iteration} - Turn {new_turn}",
    )

    fig.update_traces(marker=dict(size=10))  # Adjust marker size

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
            scaleanchor="x",  # Forces square cells
        ),
        autosize=True,
        height=GRID_HEIGHT * CELL_SIZE,
        width=GRID_WIDTH * CELL_SIZE,
        showlegend=True,
        plot_bgcolor="white",
    )

    return fig, f"Iteration {current_iteration} - Turn {new_turn}", new_turn


# Run the app
if __name__ == "__main__":
    app.run_server(debug=True)
