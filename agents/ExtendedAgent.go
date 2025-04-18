package agents

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"sort"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
	"github.com/google/uuid"
)

type ExtendedAgent struct {
	*agent.BaseAgent[infra.IExtendedAgent]
	Server infra.IServer
	
	Telomere *infra.Telomere

	Position       infra.PositionVector
	MovementPolicy string // Defines how movement is determined

	//History Tracking
	ClusterID                   int
	ObservedEliminationsCluster int
	ObservedEliminationsNetwork int
	Heroism                     int // number of times agent volunteered self-sacrifices

	// Social network and kinship group
	network      map[uuid.UUID]float32 // stores relationship strengths
	kinshipGroup []uuid.UUID           // Descendants
	parent1      uuid.UUID
	parent2      uuid.UUID

	Attachment infra.Attachment // Attachment orientations: [anxiety, avoidance].

	// Decision-Making Parameters:
	//ASP map[string]float32 // Parameters for decision-making
	//PTS map[string]float32 // Parameters for behavior protocols
	PTW infra.PTSParams // Parameters for PTS

	worldview uint32 // 32-bit binary representation of opinions

	// Ysterofimia (Posthumous Recognition)
	Ysterofimia *infra.Ysterofimia // Observation of self-sacrifice vs self-preservation

	AgentIsAlive bool // True if agent is alive

	MortalitySalience      float32 //section in ASP module
	WorldviewValidation    float32 //section in ASP module
	RelationshipValidation float32 //section in ASP module

	SelfSacrificeWillingness float32 //ASP result
}

// type AgentConfig struct {
// 	InitSacrificeWillingness float32
// }

//var _ infra.IExtendedAgent = (*ExtendedAgent)(nil)

func CreateExtendedAgent(server infra.IServer, grid *infra.Grid, parent1ID uuid.UUID, parent2ID uuid.UUID, worldview uint32) *ExtendedAgent {
	A := rand.Intn(25) + 40     // (40-65)
	B := A + rand.Intn(35) + 20 // Random max age (60 - 100)

	return &ExtendedAgent{
		BaseAgent:     agent.CreateBaseAgent(server),
		Server:        server,                                                               // Type assert the server functions to IServer interface
		Attachment:    infra.Attachment{Anxiety: rand.Float32(), Avoidance: rand.Float32()}, // Randomised anxiety and avoidance
		Heroism:       0,                                                                    //start at 0 increment if chose to self-sacrifice
		network:       make(map[uuid.UUID]float32),
		parent1: 	   parent1ID,
		parent2: 	   parent2ID,
		Telomere:      infra.NewTelomere(rand.Intn(50), A, B, 0.5),
		worldview:     worldview,
		Ysterofimia:   infra.NewYsterofimia(),
		AgentIsAlive:  true,
		Position:      infra.PositionVector{X: rand.Intn(grid.Width) + 1, Y: rand.Intn(grid.Height) + 1},
	}
}

// ----------------------- Interface implementation -----------------------

func (ea *ExtendedAgent) AgentInitialised() {}

func (ea *ExtendedAgent) GetName() uuid.UUID {
	return ea.GetID()
}

func (ea *ExtendedAgent) GetAge() int {
	// Beta distribution parameters (adjusted to fit UK population shape)
	return ea.Telomere.GetAge()
}

func (ea *ExtendedAgent) GetTelomere() float32 {
	return ea.Telomere.GetProbabilityOfDeath()
}

func (ea *ExtendedAgent) IncrementAge() {
	ea.Telomere.IncrementAge()
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

func (ea *ExtendedAgent) GetNetwork() map[uuid.UUID]float32 {
	return ea.network
}

func (ea *ExtendedAgent) AddRelationship(otherID uuid.UUID, strength float32) {
	ea.Server.UpdateAgentRelationship(ea.GetID(), otherID, strength)
}

func (ea *ExtendedAgent) RemoveRelationship(otherID uuid.UUID) {
	delete(ea.network, otherID)
}

func (ea *ExtendedAgent) UpdateRelationship(otherID uuid.UUID, change float32) {
	ea.Server.UpdateAgentRelationship(ea.GetID(), otherID, change)
}

// Finds closest friend in social network
func (ea *ExtendedAgent) FindClosestFriend() infra.IExtendedAgent {
	var closestFriends []infra.IExtendedAgent
	minDist := math.Inf(1)

	for friendID := range ea.network {
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
	return ea.worldview
}

// func (ea *ExtendedAgent) SetWorldviewBinary(worldview uint32) {
// 	ea.Worldview = worldview
// }

func (ea *ExtendedAgent) AddDescendant(childID uuid.UUID) {
	ea.kinshipGroup = append(ea.kinshipGroup, childID)
}

func (ea *ExtendedAgent) GetParents() (uuid.UUID, uuid.UUID) {
	return ea.parent1, ea.parent2
}

// func (ea *ExtendedAgent) SetParents(parent1, parent2 uuid.UUID) {
// 	ea.Parent1 = parent1
// 	ea.Parent2 = parent2
// }

func (ea *ExtendedAgent) GetYsterofimia() *infra.Ysterofimia {
	return ea.Ysterofimia
}

func (ea *ExtendedAgent) MarkAsDead() {
	ea.AgentIsAlive = false
}

func (ea *ExtendedAgent) IsAlive() bool {
	return ea.AgentIsAlive
}

func (ea *ExtendedAgent) RelativeAgeToNetwork() float32 {
	var totalAge int
	var numAgentsNetwork float32
	age := ea.GetAge()

	for friendID := range ea.network {
		friend, ok := ea.Server.GetAgentByID(friendID)
		if ok && friendID != ea.GetID() {
			totalAge += friend.GetAge()
			numAgentsNetwork++
		}
	}
	if numAgentsNetwork == 0 || age == 0 {
		return 0
	}

	averageAge := float32(totalAge) / numAgentsNetwork
	diff := float32(age) - averageAge
	if diff <= 0 {
		return 0
	}
	return diff / float32(age)
}

func (ea *ExtendedAgent) GetMemorialProximity(grid *infra.Grid) float32 {
	agentMap := ea.Server.GetAgentMap()
	selfPosition := ea.GetPosition()
	clusterID := ea.GetClusterID()
	memorials := append(grid.Tombstones, grid.Temples...)

	if len(memorials) == 0 {
		return 0 // no memorials
	}

	//numerator - distance from self to all memorials
	var selfMemorialDistanceSum float64
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
			score := worldviewAlignment(ea.worldview, otherAgent.GetWorldviewBinary())
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
	for friendID := range ea.network {
		if other, ok := agentMap[friendID]; ok {
			score := worldviewAlignment(ea.worldview, other.GetWorldviewBinary())
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

func (ea *ExtendedAgent) GetEstrangement() float32 {
	kin := ea.kinshipGroup
	network := ea.network

	if len(kin) == 0 {
		return 0.0 // no descendants
	}

	connectedDescendants := 0
	for _, descendantsID := range kin {
		if _, ok := network[descendantsID]; ok {
			connectedDescendants++
		}
	}

	return float32(connectedDescendants) / float32(len(kin))
}

func (ea *ExtendedAgent) GetProSocialEsteem() float32 {
	network := ea.network
	if len(network) == 0 {
		return 0.0 // No neighbors, no esteem
	}
	sumEsteem := float32(0.0)
	for _, esteem := range network {
		sumEsteem += esteem
	}

	return sumEsteem / float32(len(network))
}

func (ea *ExtendedAgent) GetHeroismTendency() float32 {
	agentMap := ea.Server.GetAgentMap()
	network := ea.network
	if len(network) == 0 {
		return 0.0 // No neighbors, no tendency
	}

	heroismScores := []int{}

	for id := range network {
		if agent, ok := agentMap[id]; ok {
			heroismScores = append(heroismScores, agent.GetHeroism())
		}
	}

	// Sort heroism scores in ascending order
	sort.Ints(heroismScores)

	selfHeroism := ea.GetHeroism()
	index := sort.SearchInts(heroismScores, selfHeroism)

	return float32(index) / float32(len(heroismScores))
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
	ysterofimia := ea.GetYsterofimia().ComputeYsterofimia() // compute ysterofimia

	return infra.W5*cpr + infra.W6*npr + infra.W7*ysterofimia
}

func (ea *ExtendedAgent) ComputeRelationshipValidation() float32 {
	//w8, w9, w10 := float32(0.25), float32(0.25), float32(0.5) // tweak

	est := ea.GetEstrangement()                // compute EST
	pse := ea.GetProSocialEsteem()             // compute PSE
	heroismTendency := ea.GetHeroismTendency() // compute heroism tendency

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
	ea.network[id] = change
}

func (ea *ExtendedAgent) GetPTSParams() infra.PTSParams {
	return ea.PTW
}

func (ea *ExtendedAgent) UpdateEsteem(friendID uuid.UUID, isCheck bool) {
	currentEsteem := ea.network[friendID]
	if isCheck {
		ea.network[friendID] = currentEsteem + ea.PTW.Alpha*(1-currentEsteem)
	} else {
		ea.network[friendID] = currentEsteem - ea.PTW.Beta*(currentEsteem)
	}
}

func (ea *ExtendedAgent) CreateWellbeingCheckMessage() *infra.WellbeingCheckMessage {
	return &infra.WellbeingCheckMessage{
		BaseMessage: ea.CreateBaseMessage(),
	}
}


func (ea *ExtendedAgent) CreateReplyMessage() *infra.ReplyMessage {
	return &infra.ReplyMessage{
		BaseMessage: ea.CreateBaseMessage(),
	}
}

func (ea *ExtendedAgent) HandleWellbeingCheckMessage(msg *infra.WellbeingCheckMessage) {
	fmt.Printf("Agent %v received wellbeing check from %v\n", ea.GetID(), msg.Sender)
	// not receiving??
	if rand.Float32() < ea.PTW.ReplyProb {
		reply := infra.ReplyMessage{BaseMessage: ea.CreateBaseMessage()}
		ea.SendMessage(&reply, msg.Sender)
		//fmt.Printf("Agent %v sending reply message to %v\n", ea.GetID(), msg.Sender)

		//then update alpha
		//ea.UpdateEsteem(msg.Sender, true, ea.PTW.Alpha, ea.PTW.Beta)
	}
}

func (ea *ExtendedAgent) HandleReplyMessage(msg *infra.ReplyMessage) {
	// update alpha
	//ea.UpdateEsteem(msg.Sender, true, ea.PTW.Alpha, ea.PTW.Beta)
	ea.SignalMessagingComplete()
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
func (ea *ExtendedAgent) RecordAgentJSON(instance infra.IExtendedAgent) gameRecorder.JSONAgentRecord {
	return gameRecorder.JSONAgentRecord{
		ID:                  ea.GetID().String(),
		IsAlive:             ea.IsAlive(),
		Age:                 ea.GetAge(),
		//AttachmentStyle:     ea.AttachmentStyle.String(),
		AttachmentAnxiety:   ea.Attachment.Anxiety,
		AttachmentAvoidance: ea.Attachment.Avoidance,
		ClusterID:           ea.ClusterID,
		Position:            gameRecorder.Position{X: ea.Position.X, Y: ea.Position.Y},
		Worldview:           ea.worldview,
		Heroism:             ea.Heroism,
		//MortalitySalience:   ea.ComputeMortalitySalience(),
		//WorldviewValidation: ea.ComputeWorldviewValidation(),
		//RelationshipValidation: ea.ComputeRelationshipValidation(),
		//ASPDecison: ea.GetASPDecision(nil).String(),
	}
}