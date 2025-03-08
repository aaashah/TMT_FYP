import dash
from dash import dcc, html
import pandas as pd
import plotly.express as px
from dash.dependencies import Input, Output
import hashlib

# Constants for styling
PLOT_HEIGHT = 600
PLOT_WIDTH = 600  # Keep square grid
PLOT_MARGIN = dict(r=150)

# Global cache for data
last_data_hash = None
cached_data = None


# Function to load and cache data
def load_data():
    global last_data_hash, cached_data

    base_path = "visualisation_output/csv_data"
    current_data = {
        "agent_records": pd.read_csv(f"{base_path}/agent_records.csv"),
        "infra_records": pd.read_csv(f"{base_path}/infra_records.csv"),
    }

    # Generate hash of the current data
    hash_string = "".join(df.to_json() for df in current_data.values())
    current_hash = hashlib.md5(hash_string.encode()).hexdigest()

    # If the data hasn't changed, use cached data
    if last_data_hash == current_hash and cached_data is not None:
        return cached_data

    last_data_hash = current_hash
    cached_data = current_data
    return current_data


# Load initial data
data = load_data()
agent_records = data["agent_records"]

# Ensure consistent colors per agent
unique_agents = agent_records["AgentID"].unique()
agent_colors = {
    agent: px.colors.qualitative.Plotly[i % len(px.colors.qualitative.Plotly)]
    for i, agent in enumerate(unique_agents)
}

# Initialize Dash app
app = dash.Dash(__name__)

# Layout
app.layout = html.Div(
    [
        html.H2("Agent Movement Grid"),
        dcc.Slider(
            id="turn-slider",
            min=agent_records["TurnNumber"].min(),
            max=agent_records["TurnNumber"].max(),
            value=agent_records["TurnNumber"].min(),
            marks={
                i: str(i)
                for i in range(
                    agent_records["TurnNumber"].min(),
                    agent_records["TurnNumber"].max() + 1,
                    5,
                )
            },
            step=1,
        ),
        dcc.Graph(id="grid-plot"),
    ]
)


# Callback to update the grid
@app.callback(Output("grid-plot", "figure"), [Input("turn-slider", "value")])
def update_grid(turn):
    filtered_df = agent_records[agent_records["TurnNumber"] == turn].copy()

    # Adjust positions to center agents inside cells
    filtered_df["PositionX"] += 0.5
    filtered_df["PositionY"] += 0.5

    fig = px.scatter(
        filtered_df,
        x="PositionX",
        y="PositionY",
        color="AgentID",
        color_discrete_map=agent_colors,
        labels={"PositionX": "X Position", "PositionY": "Y Position"},
        title=f"Agent Movement at Turn {turn}",
    )

    fig.update_traces(marker=dict(size=12))  # Adjust marker size

    # Set grid size correctly
    grid_size = (
        max(agent_records["PositionX"].max(), agent_records["PositionY"].max()) + 1
    )

    fig.update_layout(
        xaxis=dict(
            range=[0, grid_size],
            tickmode="linear",
            dtick=1,
            scaleanchor="y",  # Keeps square cells
            showgrid=True,
            gridcolor="lightgray",
            zeroline=False,
        ),
        yaxis=dict(
            range=[0, grid_size],
            tickmode="linear",
            dtick=1,
            scaleanchor="x",  # Keeps square cells
            showgrid=True,
            gridcolor="lightgray",
            zeroline=False,
        ),
        showlegend=True,
        plot_bgcolor="white",
        height=PLOT_HEIGHT,  # Keep plot square
        width=PLOT_WIDTH,  # Make grid visually balanced
    )

    return fig


# Run the app
if __name__ == "__main__":
    app.run_server(debug=True)
