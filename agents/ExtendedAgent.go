package agents

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
	"github.com/google/uuid"
)


type ExtendedAgent struct {
	*agent.BaseAgent[infra.IExtendedAgent]
	Server infra.IServer
	NameID uuid.UUID

	Age int
	AgeA int // Age where mortality probability starts increasing
	AgeB int // Age where agent is definitely eliminated
	Telomere float32 // Determines lifespan decay (death probability)

	Position [2]int
	MovementPolicy string // Defines how movement is determined

	//History Tracking
	ClusterID int
	ObservedEliminationsCluster int
	ObservedEliminationsNetwork int
	Heroism                     float32 // number of times agent volunteered self-sacrifices

	// Social network and kinship group
	Network map[uuid.UUID]float32 // stores relationship strengths
	KinshipGroup        []uuid.UUID  // Descendants 

	Attachment []float32 // Attachment orientations: [anxiety, avoidance].

	// Decision-Making Parameters:
	ASP map[string]float32 // Parameters for decision-making
	PTS map[string]float32 // Parameters for behavior protocols

	Worldview uint32 // 32-bit binary representation of opinions

	// Ysterofimia (Posthumous Recognition)
	Ysterofimia float32 // Observation of self-sacrifice vs self-preservation

	Mortality bool

	MortalitySalience float32 //section in ASP module
	WorldviewValidation float32 //section in ASP module
	RelationshipValidation float32 //section in ASP module

	SelfSacrificeWillingness float32 //ASP result
}


type AgentConfig struct {
	InitSacrificeWillingness float32
	
}
//var _ infra.IExtendedAgent = (*ExtendedAgent)(nil)


func CreateExtendedAgent(server agent.IExposedServerFunctions[infra.IExtendedAgent], configParam AgentConfig, grid *infra.Grid) *ExtendedAgent {
	A := rand.Intn(25) + 40  // (40-65)
	B := A + rand.Intn(35) + 20  // Random max age (60 - 100)

	return &ExtendedAgent{
		BaseAgent: agent.CreateBaseAgent(server),
		Server:    server.(infra.IServer), // Type assert the server functions to IServer interface
		NameID:    uuid.New(),
		Attachment: []float32{rand.Float32(), rand.Float32()}, // Randomised anxiety and avoidance
		Network:    make(map[uuid.UUID]float32),
		Age:        rand.Intn(50),
		AgeA:       A,
		AgeB:       B,
		Worldview: rand.Uint32(),
		Mortality: false,
		SelfSacrificeWillingness: configParam.InitSacrificeWillingness,
		Position: [2]int{rand.Intn(grid.Width) + 1, rand.Intn(grid.Height) + 1},
	}
}

// ----------------------- Interface implementation -----------------------

func (ea *ExtendedAgent) AgentInitialised() {}

func (ea *ExtendedAgent) GetName() uuid.UUID {
	return ea.GetID()
}

func (ea *ExtendedAgent) SetName(name uuid.UUID) {
    ea.NameID = ea.GetID()
}

func (ea *ExtendedAgent) GetAge() int {
	// Beta distribution parameters (adjusted to fit UK population shape)
	return ea.Age
}

func (ea *ExtendedAgent) GetTelomere() float32 {

	if ea.Age < ea.AgeA {
		return 0.005 * float32(ea.Age) // Small increasing probability
	} else if ea.Age >= ea.AgeB {
		return 1.0 // Guaranteed death at AgeB
	} else {
		// Linearly increasing probability from AgeA to AgeB
		return float32(ea.Age-ea.AgeA) / float32(ea.AgeB-ea.AgeA)
	}
}


func (ea *ExtendedAgent) SetAge(age int) {
    ea.Age = age
}

func (ea *ExtendedAgent) GetPosition() [2]int {
	return ea.Position
}

func (ea *ExtendedAgent) Move(grid *infra.Grid) {
	newX, newY := grid.GetValidMove(ea.Position[0], ea.Position[1]) // Get a valid move
	grid.UpdateAgentPosition(ea, newX, newY)    // Update position in the grid
	ea.Position = [2]int{newX, newY}             // ✅ Assign new position
	fmt.Printf("Agent %v moved to (%d, %d)\n", ea.GetID(), newX, newY)
}

// Returns -1, 0, or 1 to move in the right direction
func getStep(current, target int) int {
	if target > current {
		return 1
	} else if target < current {
		return -1
	}
	return 0
}

func (ea *ExtendedAgent) GetAttachment() []float32 {
    return ea.Attachment
}

func (ea *ExtendedAgent) SetAttachment(attachment []float32) {
    if len(attachment) != 2 {
        panic("Attachment must have exactly two elements: [anxiety, avoidance]")
    }
    ea.Attachment = attachment
}

func randInRange(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func (ea *ExtendedAgent) GetNetwork() map[uuid.UUID]float32 {
	return ea.Network
}

func (ea *ExtendedAgent) SetNetwork(network map[uuid.UUID]float32) {
	ea.Network = network
}

// distance between two agents on grid
// func (ea *ExtendedAgent) DistanceTo(other *ExtendedAgent) float64 {
// 	return infra.Distance(ea.Position, other.Position)
// }

func (ea *ExtendedAgent) AddRelationship(otherID uuid.UUID, strength float32) {
	ea.Server.UpdateAgentRelationship(ea.GetID(), otherID, strength)
}

func (ea *ExtendedAgent) UpdateRelationship(otherID uuid.UUID, change float32) {
	ea.Server.UpdateAgentRelationship(ea.GetID(), otherID, change)
}

// Finds closest friend in social network
func (ea *ExtendedAgent) FindClosestFriend() *ExtendedAgent {
	var closestFriends []*ExtendedAgent
	minDist := math.MaxFloat64

	for friendID := range ea.GetNetwork() {
		// lookup friend in server
		agentInterface, exists := ea.Server.GetAgentByID(friendID)
		if !exists {
			continue
		}
		friend, ok := agentInterface.(*ExtendedAgent)
		if !ok {
			continue // type assertion failed
		}

		dist := distance(ea.Position, friend.Position)
		if dist < minDist {
			minDist = dist
			closestFriends = []*ExtendedAgent{friend} // start new list
		} else if dist == minDist {
			closestFriends = append(closestFriends, friend) // add equally close
		}
	}
	if len(closestFriends) == 0 {
		return nil
	}

	return closestFriends[rand.Intn(len(closestFriends))] // pick randomly
}


// Euclidean distance helper
func distance(pos1, pos2 [2]int) float64 {
	dx := float64(pos1[0] - pos2[0])
	dy := float64(pos1[1] - pos2[1])
	return math.Sqrt(dx*dx + dy*dy)
}

func (ea *ExtendedAgent) GetClusterID() int {
	return ea.ClusterID
}

func (ea *ExtendedAgent) SetClusterID(id int) {
	ea.ClusterID = id
}

// GetWorldviewBinary returns the 32-bit binary representation of the agent's worldview.
func (ea *ExtendedAgent) GetWorldviewBinary() uint32 {
	return ea.Worldview
}

func (ea *ExtendedAgent) GetHeroism() float32 {
	ea.Heroism = rand.Float32() // Randomized value for Heroism
	return ea.Heroism
}
// func (ea *ExtendedAgent) SetHeroism() {
// 	ea.Heroism = rand.Float32()
// }

func (ea *ExtendedAgent) GetYsterofimia() float32 {
	ea.Ysterofimia = rand.Float32() // Randomized value for Ysterofimia
	return ea.Ysterofimia
}

// GetMortality returns the mortality status of the agent.
func (ea *ExtendedAgent) GetMortality() bool {
	probDeath := ea.GetTelomere() // get probability of death
	randVal := rand.Float32()     // random value between 0 and 1
	ea.Mortality = randVal < probDeath // Random chance to die of natural causes
	return ea.Mortality
}

func (ea *ExtendedAgent) IsMortalitySalient() bool {
    return ea.Mortality
}

func (ea *ExtendedAgent) SetMortalitySalience(ms bool) {
    ea.Mortality = ms
}

func (ea *ExtendedAgent) GetSelfSacrificeWillingness() float32 {
	return ea.SelfSacrificeWillingness
}

// Decision-making logic
func (ea *ExtendedAgent) DecideSacrifice() float32 {
    //TO-DO: Fuzzy logic stuff

	ea.SelfSacrificeWillingness = rand.Float32() // Random willingness to sacrifice

	
    //fmt.Printf("Agent %d decision: %v\n", a.NameID, a.SacrificeChoice)
    
	// fmt.Printf("Agent %v willing to sacrifice by %s \n",
    //     ea.NameID,
    //     map[float32]string{}[ea.SelfSacrificeWillingness])
    return ea.SelfSacrificeWillingness
}

func (ea *ExtendedAgent) GetExposedInfo() infra.ExposedAgentInfo {
	return infra.ExposedAgentInfo{
		AgentUUID: ea.GetID(),
	}
}

// ----------------------- Data Recording Functions -----------------------
func (mi *ExtendedAgent) RecordAgentStatus(instance infra.IExtendedAgent) gameRecorder.AgentRecord {
	//fmt.Printf("[DEBUG] Fetching Age in RecordAgentStatus: %d for Agent %v\n", instance.GetAge(), instance.GetID()) 
	record := gameRecorder.NewAgentRecord(
		instance.GetID(),
		instance.GetAge(),
		instance.GetPosition()[0],
		instance.GetPosition()[1],
		instance.GetSelfSacrificeWillingness(),
		instance.GetAttachment()[0],
		instance.GetAttachment()[1],
		instance.GetWorldviewBinary(),
		"1",
		//instance.GetWorldviewValidation(),
		//instance.GetRelationshipValidation(),
		
	)
	return record
}
