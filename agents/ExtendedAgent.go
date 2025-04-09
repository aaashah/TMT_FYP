package agents

import (
	"math/bits"
	"math/rand"

	//"fmt"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
	"github.com/google/uuid"
)

type ExtendedAgent struct {
	*agent.BaseAgent[infra.IExtendedAgent]
	Server infra.IServer
	// NameID uuid.UUID

	Age int
	// AgeA     int     // Age where mortality probability starts increasing
	// AgeB     int     // Age where agent is definitely eliminated
	// Telomere float32 // Determines lifespan decay (death probability)
	Telomere infra.Telomere

	Position       infra.PositionVector
	MovementPolicy string // Defines how movement is determined

	//History Tracking
	ClusterID                   int
	ObservedEliminationsCluster int
	ObservedEliminationsNetwork int
	Heroism                     int // number of times agent volunteered self-sacrifices

	// Social network and kinship group
	Network      map[uuid.UUID]float32 // stores relationship strengths
	KinshipGroup []uuid.UUID           // Descendants

	Attachment infra.Attachment // Attachment orientations: [anxiety, avoidance].

	// Decision-Making Parameters:
	ASP map[string]float32 // Parameters for decision-making
	PTS map[string]float32 // Parameters for behavior protocols

	Worldview uint32 // 32-bit binary representation of opinions

	// Ysterofimia (Posthumous Recognition)
	Ysterofimia float32 // Observation of self-sacrifice vs self-preservation

	//Mortality bool
	AgentIsAlive bool // True if agent is alive

	MortalitySalience      float32 //section in ASP module
	WorldviewValidation    float32 //section in ASP module
	RelationshipValidation float32 //section in ASP module

	SelfSacrificeWillingness float32 //ASP result
}

type AgentConfig struct {
	InitSacrificeWillingness float32
}

//var _ infra.IExtendedAgent = (*ExtendedAgent)(nil)

func CreateExtendedAgent(server infra.IServer, configParam AgentConfig, grid *infra.Grid) *ExtendedAgent {
	A := rand.Intn(25) + 40     // (40-65)
	B := A + rand.Intn(35) + 20 // Random max age (60 - 100)

	return &ExtendedAgent{
		BaseAgent:                agent.CreateBaseAgent(server),
		Server:                   server,                                                               // Type assert the server functions to IServer interface
		Attachment:               infra.Attachment{Anxiety: rand.Float32(), Avoidance: rand.Float32()}, // Randomised anxiety and avoidance
		Heroism:                  0,                                                                    //start at 0 increment if chose to self-sacrifice
		Network:                  make(map[uuid.UUID]float32),
		Age:                      rand.Intn(50),
		Telomere:                 infra.NewTelomere(A, B, 0.5),
		Worldview:                rand.Uint32(),
		Ysterofimia:              rand.Float32(),
		AgentIsAlive:                   true,
		SelfSacrificeWillingness: configParam.InitSacrificeWillingness,
		Position:                 infra.PositionVector{X: rand.Intn(grid.Width) + 1, Y: rand.Intn(grid.Height) + 1},
	}
}

// const (
// 	// ASP weights
// 	w1 = float32(0.25)
// 	w2 = float32(0.25)
// 	w3 = float32(0.25)
// 	w4 = float32(0.25)
// 	w5 = float32(0.25)
// 	w6 = float32(0.25)
// 	w7 = float32(0.5)
// 	w8 = float32(0.25)
// 	w9 = float32(0.25)
// 	w10 = float32(0.5)
// )

const MaxFloat32 = float32(3.4028235e+38) // largest float32 value
// ----------------------- Interface implementation -----------------------

func (ea *ExtendedAgent) AgentInitialised() {}

func (ea *ExtendedAgent) GetName() uuid.UUID {
	return ea.GetID()
}

func (ea *ExtendedAgent) GetAge() int {
	// Beta distribution parameters (adjusted to fit UK population shape)
	return ea.Age
}

func (ea *ExtendedAgent) GetTelomere() float32 {

	return ea.Telomere.GetProbabilityOfDeath(ea.Age)

}

func (ea *ExtendedAgent) IncrementAge() {
	ea.Age++
}

func (ea *ExtendedAgent) GetPosition() infra.PositionVector {
	return ea.Position
}

func (ea *ExtendedAgent) SetPosition(newPos infra.PositionVector) {
	ea.Position = newPos
}

func (ea *ExtendedAgent) GetAttachment() infra.Attachment {
	return ea.Attachment
}

func randInRange(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

// func (ea *ExtendedAgent) GetNetwork() map[uuid.UUID]float32 {
// 	return ea.Network
// }

func (ea *ExtendedAgent) AddRelationship(otherID uuid.UUID, strength float32) {
	ea.Server.UpdateAgentRelationship(ea.GetID(), otherID, strength)
}

func (ea *ExtendedAgent) UpdateRelationship(otherID uuid.UUID, change float32) {
	ea.Server.UpdateAgentRelationship(ea.GetID(), otherID, change)
}

// Finds closest friend in social network
func (ea *ExtendedAgent) FindClosestFriend() infra.IExtendedAgent {
	var closestFriends []infra.IExtendedAgent
	minDist := MaxFloat32

	for friendID := range ea.Network {
		// lookup friend in server
		agentInterface, exists := ea.Server.GetAgentByID(friendID)
		if !exists {
			continue
		}

		dist := ea.Position.Dist(agentInterface.GetPosition())
		if dist < minDist {
			minDist = dist
			closestFriends = []infra.IExtendedAgent{agentInterface} // start new list
		} else if dist == minDist {
			closestFriends = append(closestFriends, agentInterface) // add equally close
		}
	}
	if len(closestFriends) == 0 {
		return nil
	}

	return closestFriends[rand.Intn(len(closestFriends))] // pick randomly
}

func (ea *ExtendedAgent) GetClusterID() int {
	return ea.ClusterID
}

func (ea *ExtendedAgent) SetClusterID(id int) {
	ea.ClusterID = id
}

func (ea *ExtendedAgent) IncrementClusterEliminations(n int) {
	ea.ObservedEliminationsCluster += n
}

func (ea *ExtendedAgent) IncrementNetworkEliminations(n int) {
	ea.ObservedEliminationsNetwork += n
}

func (ea *ExtendedAgent) IncrementHeroism() {
	ea.Heroism++
}

func (ea *ExtendedAgent) GetHeroism() int {
	return ea.Heroism
}

// GetWorldviewBinary returns the 32-bit binary representation of the agent's worldview.
func (ea *ExtendedAgent) GetWorldviewBinary() uint32 {
	return ea.Worldview
}

// func (ea *ExtendedAgent) GetYsterofimia() float32 {
// 	return ea.Ysterofimia
// }

func (ea *ExtendedAgent) MarkAsDead() {
	ea.AgentIsAlive = false
}

func (ea *ExtendedAgent) IsAlive() bool {
	return ea.AgentIsAlive
}

func (ea *ExtendedAgent) RelativeAgeToNetwork() float32 {
	var totalAge int
	var numAgentsNetwork int

	for friendID := range ea.Network {
		friend, ok := ea.Server.GetAgentByID(friendID)
		if ok && friendID != ea.GetID() {
			totalAge += friend.GetAge()
			numAgentsNetwork++
		}
	}
	if numAgentsNetwork == 0 || ea.Age == 0 {
		return 0
	}
	averageAge := totalAge / numAgentsNetwork
	diff := ea.Age - averageAge
	if diff <= 0 {
		return 0
	}
	return float32(diff) / float32(ea.Age)
}

func (ea *ExtendedAgent) GetMemorialProximity(grid *infra.Grid) float32 {
	agentMap := ea.Server.GetAgentMap()
	selfPosition := ea.GetPosition()
	clusterID := ea.GetClusterID()
	//memorials := []infra.PositionVector{}
	memorials := append(grid.Tombstones, grid.Temples...)

	// for tomb := range grid.Tombstones {
	// 	memorials = append(memorials, infra.PositionVector{X: tomb[0], Y: tomb[1]})
	// }
	// for temples := range grid.Temples {
	// 	memorials = append(memorials, infra.PositionVector{X: temples[0], Y: temples[1]})
	// }
	


	if len(memorials) == 0 {
		return 0 // no memorials
	}

	//numerator - distance from self to all memorials
	var selfMemorialDistanceSum float32
	for _, mem := range memorials {
		selfMemorialDistanceSum += selfPosition.Dist(mem)
	}

	//denominator- distance from memorials and distance from cluster agents to memorials
	denominator := selfMemorialDistanceSum
	for _, otherAgent := range agentMap {
		if otherAgent.GetID() == ea.GetID() {
			continue // skip self
		}
		if otherAgent.GetClusterID() == clusterID {
			otherPosition := otherAgent.GetPosition()
			for _, mem := range memorials {
				denominator += otherPosition.Dist(mem)
			}
		}
	}

	return float32(selfMemorialDistanceSum / denominator)
}

func worldviewAlignment(a, b uint32) float32 {
	// XNOR the numbers to find aligned bits
	alignedBits := ^(a ^ b)

	// Count the differing bits (Hamming weight)
	alignedBitCount := bits.OnesCount32(alignedBits)

	// Divide by 32 to get average bit alignment
	return float32(alignedBitCount) / 32.0
}

func (ea *ExtendedAgent) GetCPR() float32 {
	// compute cluster profiles
	agentMap := ea.Server.GetAgentMap()
	clusterID := ea.GetClusterID()
	clusterAlignments := []float32{}
	for _, otherAgent := range agentMap {
		if otherAgent.GetID() == ea.GetID() {
			continue // skip self
		}
		if otherAgent.GetClusterID() == clusterID {
			score := worldviewAlignment(ea.Worldview, otherAgent.GetWorldviewBinary())
			clusterAlignments = append(clusterAlignments, score)
		}
	}

	if len(clusterAlignments) == 0 {
		return 0
	}
	// compute average alignment
	var totalAlignment float32
	for _, score := range clusterAlignments {
		totalAlignment += score
	}
	return float32(totalAlignment) / float32(len(clusterAlignments))
}

func (ea *ExtendedAgent) GetNPR() float32 {
	// compute network profiles
	networkAlignments := []float32{}
	agentMap := ea.Server.GetAgentMap()
	for friendID := range ea.Network {
		if other, ok := agentMap[friendID]; ok {
			score := worldviewAlignment(ea.Worldview, other.GetWorldviewBinary())
			networkAlignments = append(networkAlignments, score)
		}
	}
	if len(networkAlignments) == 0 {
		return 0
	}
	// compute average alignment
	var totalAlignment float32
	for _, score := range networkAlignments {
		totalAlignment += score
	}
	return float32(totalAlignment) / float32(len(networkAlignments))
}

func (ea *ExtendedAgent) ComputeMortalitySalience(grid *infra.Grid) float32 {
	//w1, w2, w3, w4 := float32(0.25), float32(0.25), float32(0.25), float32(0.25) // tweak

	ce := float32(ea.ObservedEliminationsCluster)
	ne := float32(ea.ObservedEliminationsNetwork)
	ra := float32(ea.RelativeAgeToNetwork())
	mp := float32(ea.GetMemorialProximity(grid))

	return infra.W1*ce + infra.W2*ne + infra.W3*ra + infra.W4*mp
}

func (ea *ExtendedAgent) ComputeWorldviewValidation() float32 {
	//w5, w6, w7 := float32(0.25), float32(0.25), float32(0.5) // tweak

	cpr := ea.GetCPR()
	npr := ea.GetNPR() // compute NPR
	ysterofimia := ea.Ysterofimia

	return infra.W5*cpr + infra.W6*npr + infra.W7*ysterofimia
}

func (ea *ExtendedAgent) ComputeRelationshipValidation() float32 {
	//w8, w9, w10 := float32(0.25), float32(0.25), float32(0.5) // tweak

	est := rand.Float32()             // compute EST
	pse := rand.Float32()             // compute PSE
	heroismTendency := rand.Float32() // compute heroism tendency

	return infra.W8*est + infra.W9*pse + infra.W10*heroismTendency
}

func (ea *ExtendedAgent) GetSelfSacrificeWillingness() float32 {
	return ea.SelfSacrificeWillingness
}

// Decision-making logic
func (ea *ExtendedAgent) GetASPDecision(grid *infra.Grid) infra.ASPDecison {
	threshold := float32(0.75) //random threshold

	ms := ea.ComputeMortalitySalience(grid)
	wv := ea.ComputeWorldviewValidation()
	rv := ea.ComputeRelationshipValidation()

	sum := 0
	for _, score := range []float32{ms, wv, rv} {
		if score > threshold {
			sum += 1
		} else {
			sum -= 1
		}
	}

	if sum > 0 {
		return infra.SELF_SACRIFICE // Self-sacrifice
	} else if sum < 0 {
		return infra.NOT_SELF_SACRIFICE // Reject self-sacrifice
	} else {
		return infra.INACTION // No action
	}

}


func (ea *ExtendedAgent) GetExposedInfo() infra.ExposedAgentInfo {
	return infra.ExposedAgentInfo{
		AgentUUID: ea.GetID(),
	}
}

func (ea *ExtendedAgent) UpdateSocialNetwork(id uuid.UUID, change float32) {
	ea.Network[id] = change
}

func (ea *ExtendedAgent) HandleWellbeingCheckMessage(msg *infra.WellbeingCheckMessage) {
	// depend on attachment
}

func (ea *ExtendedAgent) HandleReplyMessage(msg *infra.ReplyMessage) {
	// depend on attachment
}

// ----------------------- Data Recording Functions -----------------------
func (mi *ExtendedAgent) RecordAgentStatus(instance infra.IExtendedAgent) gameRecorder.AgentRecord {
	//fmt.Printf("[DEBUG] Fetching Age in RecordAgentStatus: %d for Agent %v\n", instance.GetAge(), instance.GetID())
	agentPos := instance.GetPosition()
	agentAttach := instance.GetAttachment()
	record := gameRecorder.NewAgentRecord(
		instance.GetID(),
		instance.GetAge(),
		agentPos.X,
		agentPos.Y,
		instance.GetSelfSacrificeWillingness(),
		agentAttach.Anxiety,
		agentAttach.Avoidance,
		instance.GetWorldviewBinary(),
		"1",
		//instance.GetWorldviewValidation(),
		//instance.GetRelationshipValidation(),

	)
	return record
}
