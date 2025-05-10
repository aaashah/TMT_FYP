package config

import (
	"flag"
	"math"
)

type Config struct {
	GridWidth               int     `json:"GridWidth"`
	GridHeight              int     `json:"GridHeight"`
	NumAgents               int     `json:"NumAgents"`
	DismissiveProp          float64 `json:"DismissiveProp"`
	FearfulProp             float64 `json:"FearfulProp"`
	PreoccupiedProp         float64 `json:"ProccupiedProp"`
	SecureProp              float64 `json:"SecureProp"`
	NumIterations           int     `json:"NumIterations"`
	NumTurns                int     `json:"NumTurns"`
	NumClusters             int     `json:"NumClusters"`
	ConnectionProbability   float64 `json:"ConnectionProb"`
	PopulationRho           float64 `json:"PopulationRho"`
	InitialExpectedChildren float64 `json:"InitialExpectedChildren"`
	MinExpectedChildren     float64 `json:"MinExpectedChildren"`
	MaxExpectedChildren     float64 `json:"MaxExpectedChildren"`
	MutationRate            float64 `json:"Mu"`
	ASPThreshold            float64 `json:"ASPThreshold"`
	Debug                   bool    `json:"-"`
	Seed                    int64   `json:"-"`
}

// NewConfig parses the command line and returns a populated Config
func NewConfig() Config {
	cfg := Config{}

	flag.IntVar(&cfg.GridWidth, "width", 50, "Width of Grid World")
	flag.IntVar(&cfg.GridHeight, "height", 50, "Height of Grid World")
	flag.IntVar(&cfg.NumAgents, "numAgents", 40, "Initial number of agents")
	flag.Float64Var(&cfg.DismissiveProp, "dismissive", 0.25, "Initial proportion of dismissive agents")
	flag.Float64Var(&cfg.FearfulProp, "fearful", 0.25, "Initial proportion of fearful agents")
	flag.Float64Var(&cfg.PreoccupiedProp, "preoccupied", 0.25, "Initial proportion of preoccupied agents")
	flag.Float64Var(&cfg.SecureProp, "secure", 0.25, "Initial proportion of secure agents")
	flag.IntVar(&cfg.NumIterations, "iters", 100, "Number of iterations")
	flag.IntVar(&cfg.NumTurns, "turns", 10, "Initial number of turns")
	flag.IntVar(&cfg.NumClusters, "kappa", 3, "Number of agent clusters")
	flag.Float64Var(&cfg.ConnectionProbability, "connectionProb", 0.35, "Probability of connections in social network")
	flag.Float64Var(&cfg.PopulationRho, "rho", 0.2, "Proportion of population required to self-sacrifice")
	flag.Float64Var(&cfg.InitialExpectedChildren, "init_r0", 2.0, "Initial R0 of population")
	flag.Float64Var(&cfg.MinExpectedChildren, "min_r0", 1.9, "Minimum R0 of population")
	flag.Float64Var(&cfg.MaxExpectedChildren, "max_r0", 2.1, "Maximum R0 of population")
	flag.Float64Var(&cfg.MutationRate, "mu", 0.2, "Mutation rate of spawned children")
	flag.Float64Var(&cfg.ASPThreshold, "tau", 0.5, "Threshold for ASP decision")
	flag.BoolVar(&cfg.Debug, "debug", false, "Log debug messages to console")
	flag.Int64Var(&cfg.Seed, "seed", 42, "Random seed for reproducibility")

	flag.Parse()

	epsilon := 0.05
	if math.Abs(cfg.DismissiveProp+cfg.FearfulProp+cfg.PreoccupiedProp+cfg.SecureProp-1) > epsilon {
		panic("Proportion of attachment types do not sum to 1.0")
	}

	return cfg
}
