package config

import (
	"flag"
)

type Config struct {
	NumAgents               int     `json:"NumAgents"`
	ConnectionProbability   float64 `json:"ConnectionProb"`
	PopulationRho           float64 `json:"PopulationRho"`
	InitialExpectedChildren float64 `json:"InitialExpectedChildren"`
	Debug                   bool    `json:"-"`
	Seed                    int64   `json:"-"`
}

// NewConfig parses the command line and returns a populated Config
func NewConfig() Config {
	cfg := Config{}

	flag.IntVar(&cfg.NumAgents, "numAgents", 40, "Initial number of agents")
	flag.Float64Var(&cfg.ConnectionProbability, "connectionProb", 0.35, "Probability of connections in social network")
	flag.Float64Var(&cfg.PopulationRho, "rho", 0.2, "Proportion of population required to self-sacrifice")
	flag.Float64Var(&cfg.InitialExpectedChildren, "r0", 1.9, "Initial R0 of population")
	flag.BoolVar(&cfg.Debug, "debug", false, "Log debug messages to console")
	flag.Int64Var(&cfg.Seed, "seed", 42, "Random seed for reproducibility")

	flag.Parse()

	return cfg
}
