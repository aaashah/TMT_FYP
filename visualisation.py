import pandas as pd
import dash
from dash import dcc, html
from dash.dependencies import Input, Output
import plotly.express as px
import plotly.graph_objects as go
from datetime import datetime
import hashlib
import json

# Add these constants at the top of the file for consistent plot styling
PLOT_HEIGHT = 600  # Increased height
PLOT_WIDTH = 1000  # Increased from 700 to 1000
PLOT_MARGIN = dict(r=150)  # More space on the right for legends

# Add these global variables after the PLOT constants
last_data_hash = None
cached_data = None


# Move the data loading into a function
def load_data():
    global last_data_hash, cached_data

    base_path = "visualisation_output/csv_data"

    current_data = {
        "agent_records": pd.read_csv(f"{base_path}/agent_records.csv"),
        "infra_records": pd.read_csv(f"{base_path}/infra_records.csv"),
    }

    # Create a hash of the current data
    hash_string = ""
    for df in current_data.values():
        hash_string += df.to_json()
    current_hash = hashlib.md5(hash_string.encode()).hexdigest()

    # If the hash matches, return cached data
    if last_data_hash == current_hash and cached_data is not None:
        return cached_data

    # Otherwise, update cache and return new data
    last_data_hash = current_hash
    cached_data = current_data
    return current_data


# Load initial data
data = load_data()
agent_records = data["agent_records"]

# Initialize the Dash app
app = dash.Dash(__name__)

# Define the app layout
app.layout = html.Div(
    [
        html.H1("Agent Records"),
        dcc.Dropdown(
            id="agent-dropdown",
            options=[
                {"label": agent, "value": agent}
                for agent in agent_records["agent"].unique()
            ],
            value=agent_records["agent"].unique()[0],
        ),
        dcc.Graph(id="agent-plot"),
    ]
)


# Callback for agent status
@app.callback(Output("agent-status", "figure"), [Input("iteration-filter", "value")])
def update_agent_status(iteration):
    filtered_data = agent_records[agent_records["IterationNumber"] == iteration]

    status_counts = (
        filtered_data.groupby(["TurnNumber", "IsAlive"]).size().unstack(fill_value=0)
    )

    fig = go.Figure()

    # Safely get alive counts (True), defaulting to 0 if not present
    alive_counts = (
        status_counts[True]
        if True in status_counts.columns
        else pd.Series(0, index=status_counts.index)
    )

    # Safely get dead counts (False), defaulting to 0 if not present
    dead_counts = (
        status_counts[False]
        if False in status_counts.columns
        else pd.Series(0, index=status_counts.index)
    )

    fig.add_trace(
        go.Bar(
            x=status_counts.index, y=alive_counts, name="Alive", marker_color="green"
        )
    )

    fig.add_trace(
        go.Bar(x=status_counts.index, y=dead_counts, name="Dead", marker_color="red")
    )

    fig.update_layout(
        title="Agent Status Over Time",
        xaxis_title="Turn Number",
        yaxis_title="Number of Agents",
        barmode="stack",
        showlegend=True,
        legend=dict(yanchor="top", y=0.99, xanchor="left", x=1.05),
        height=PLOT_HEIGHT,
        width=PLOT_WIDTH,
        margin=dict(t=150, r=150, b=50, l=50),
    )

    return fig
