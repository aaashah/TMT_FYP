package agents

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/MattSScott/TMT_SOMAS/gameRecorder"
	"github.com/MattSScott/TMT_SOMAS/infra"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/google/uuid"
)

type ExtendedAgent struct {
	*agent.BaseAgent[infra.IExtendedAgent]
	infra.IServer

	telomere *infra.Telomere

	position infra.PositionVector

	//History Tracking
	clusterID int
	//ClusterHistory			    []int // the cluster ID the agent belonged to at each turn
	// clusterSizeHistory          []int // the size of the cluster at each turn
	// networkSizeHistory          []int // the size of the network at each turn
	observedEliminationsCluster int
	clusterEliminationTolerance int
	observedEliminationsNetwork int
	networkEliminationTolerance int
	heroism                     int // number of times agent volunteered self-sacrifices

	// Social network and kinship group
	network       map[uuid.UUID]float32 // stores relationship strengths
	networkLength int
	kinshipGroup  []uuid.UUID // Descendants
	parent1       uuid.UUID
	parent2       uuid.UUID

	attachment infra.Attachment // Attachment orientations: [anxiety, avoidance].

	// Decision-Making Parameters:
	//ASP map[string]float32 // Parameters for decision-making
	//PTS map[string]float32 // Parameters for behavior protocols
	PTW infra.PTSParams // Parameters for PTS

	worldview *infra.Worldview

	// ysterofimia (Posthumous Recognition)
	ysterofimia *infra.Ysterofimia // Observation of self-sacrifice vs self-preservation

	agentIsAlive bool // True if agent is alive

	//MortalitySalience      float32 //section in ASP module
	//WorldviewValidation    float32 //section in ASP module
	//RelationshipValidation float32 //section in ASP module
}

func CreateExtendedAgent(server infra.IServer, parent1ID uuid.UUID, parent2ID uuid.UUID, worldview *infra.Worldview) *ExtendedAgent {
	initAgents := float64(server.GetInitNumberAgents())

	clusTol := 0.5 * initAgents
	netTol := 0.5 * initAgents

	ag := &ExtendedAgent{
		BaseAgent:                   agent.CreateBaseAgent(server),
		IServer:                     server,                                                               // Type assert the server functions to IServer interface
		attachment:                  infra.Attachment{Anxiety: rand.Float32(), Avoidance: rand.Float32()}, // Randomised anxiety and avoidance
		heroism:                     0,                                                                    //start at 0 increment if chose to self-sacrifice
		network:                     make(map[uuid.UUID]float32),
		parent1:                     parent1ID,
		parent2:                     parent2ID,
		telomere:                    infra.NewTelomere(),
		worldview:                   worldview,
		ysterofimia:                 infra.NewYsterofimia(),
		clusterEliminationTolerance: int(clusTol),
		networkEliminationTolerance: int(netTol),
		agentIsAlive:                true,
		position:                    infra.PositionVector{X: rand.Intn(infra.GRID_WIDTH), Y: rand.Intn(infra.GRID_HEIGHT)},
	}

	ag.UpdateSocialNetwork(ag.GetID(), 0.5)
	return ag
}

// ----------------------- Interface implementation -----------------------

func (ea *ExtendedAgent) AgentInitialised() {}

func (ea *ExtendedAgent) GetName() uuid.UUID {
	return ea.GetID()
}

func (ea *ExtendedAgent) GetAge() int {
	return ea.telomere.GetAge()
}

func (ea *ExtendedAgent) GetTelomere() float64 {
	return ea.telomere.GetProbabilityOfDeath()
}

func (ea *ExtendedAgent) IncrementAge() {
	ea.telomere.IncrementAge()
}

func (ea *ExtendedAgent) GetPosition() infra.PositionVector {
	return ea.position
}

func (ea *ExtendedAgent) SetPosition(newPos infra.PositionVector) {
	ea.position = newPos
}

func (ea *ExtendedAgent) GetAttachment() infra.Attachment {
	return ea.attachment
}

func randInRange(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func (ea *ExtendedAgent) GetNetwork() map[uuid.UUID]float32 {
	return ea.network
}

func (ea *ExtendedAgent) AddRelationship(otherID uuid.UUID, strength float32) {
	ea.UpdateAgentRelationship(ea.GetID(), otherID, strength)
}

func (ea *ExtendedAgent) RemoveRelationship(otherID uuid.UUID) {
	delete(ea.network, otherID)
}

func (ea *ExtendedAgent) UpdateRelationship(otherID uuid.UUID, change float32) {
	ea.UpdateAgentRelationship(ea.GetID(), otherID, change)
}

// Finds closest friend in social network
func (ea *ExtendedAgent) FindClosestFriend() infra.IExtendedAgent {
	var closestFriends []infra.IExtendedAgent
	minDist := math.Inf(1)

	for friendID := range ea.network {
		// lookup friend in server
		agentInterface, exists := ea.GetAgentByID(friendID)
		if !exists {
			continue
		}

		dist := ea.position.Dist(agentInterface.GetPosition())
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
	return ea.clusterID
}

func (ea *ExtendedAgent) SetClusterID(id int) {
	ea.clusterID = id
}

// func (ea *ExtendedAgent) AppendClusterHistory(clusterID int, clusterSize int) {
// 	//ea.ClusterHistory = append(ea.ClusterHistory, clusterID)
// 	ea.clusterSizeHistory = append(ea.clusterSizeHistory, clusterSize)
// }

func (ea *ExtendedAgent) IncrementClusterEliminations(n int) {
	ea.observedEliminationsCluster += n
}

// func (ea *ExtendedAgent) AppendNetworkSizeHistory(networkSize int) {
// 	ea.networkSizeHistory = append(ea.networkSizeHistory, networkSize)
// }

func (ea *ExtendedAgent) IncrementNetworkEliminations(n int) {
	ea.observedEliminationsNetwork += n
}

func (ea *ExtendedAgent) SetPreEliminationNetworkLength(length int) {
	ea.networkLength = length
}

func (ea *ExtendedAgent) IncrementHeroism() {
	ea.heroism++
}

func (ea *ExtendedAgent) GetHeroism() int {
	return ea.heroism
}

// GetWorldviewBinary returns the 32-bit binary representation of the agent's worldview.
func (ea *ExtendedAgent) GetWorldview() *infra.Worldview {
	return ea.worldview
}

func (ea *ExtendedAgent) UpdateWorldview(trend float64, seasonal int) {
	ea.worldview.UpdateWorldview(trend, seasonal)
}

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
	return ea.ysterofimia
}

func (ea *ExtendedAgent) MarkAsDead() {
	ea.agentIsAlive = false
}

func (ea *ExtendedAgent) IsAlive() bool {
	return ea.agentIsAlive
}

// take the total number of eliminations you've ever seen (1)
// divide it by an agent-specific tolerance (2)
func (ea *ExtendedAgent) ClusterEliminations() float32 {
	propElim := float32(ea.observedEliminationsCluster) / float32(ea.clusterEliminationTolerance)
	return min(propElim, 1.0)
}

// take the total number of eliminations you've ever seen (1)
// divide it by an agent-specific tolerance (2)
func (ea *ExtendedAgent) NetworkEliminations() float32 {
	propElim := float32(ea.observedEliminationsNetwork) / float32(ea.networkEliminationTolerance)
	return min(propElim, 1.0)
}

func (ea *ExtendedAgent) RelativeAgeToNetwork() float32 {
	thisAge := ea.GetAge()
	ages := make([]int, 0)
	networkSize := len(ea.network)
	for agentID := range ea.network {
		if agent, ok := ea.GetAgentByID(agentID); ok {
			agentAge := agent.GetAge()
			ages = append(ages, agentAge)
		}
	}

	sort.Ints(ages)
	agePos := sort.SearchInts(ages, thisAge) + 1
	return float32(agePos) / float32(networkSize)
}

func (ea *ExtendedAgent) GetMemorialProximity(grid *infra.Grid) float32 {
	selfPosition := ea.GetPosition()
	clusterID := ea.GetClusterID()
	memorials := append(grid.Tombstones, grid.Temples...)

	totalMemorialInfluence := 0.0
	for _, mem := range memorials {
		distToMem := selfPosition.Dist(mem)
		totalMemorialInfluence += 1 / distToMem
	}

	totalClusterInfluence := 0.0
	for _, ag := range ea.GetAgentMap() {
		if ag.GetClusterID() != clusterID || ag.GetID() == ea.GetID() {
			continue
		}
		otherPosition := ag.GetPosition()
		distToAgent := selfPosition.Dist(otherPosition)
		totalClusterInfluence += 1 / distToAgent
	}

	if totalClusterInfluence == 0 && totalMemorialInfluence == 0 {
		return 0
	}

	return float32(totalClusterInfluence) / (float32(totalClusterInfluence) + float32(totalMemorialInfluence))
}

// func worldviewAlignment(a, b uint32) float32 {
// 	// XNOR the numbers to find aligned bits
// 	alignedBits := ^(a ^ b)

// 	// Count the differing bits (Hamming weight)
// 	alignedBitCount := bits.OnesCount32(alignedBits)

// 	// Divide by 32 to get average bit alignment
// 	return float32(alignedBitCount) / 32.0
// }

func (ea *ExtendedAgent) GetCPR() float32 {
	// compute cluster profiles
	agentMap := ea.GetAgentMap()
	clusterID := ea.GetClusterID()
	clusterAlignments := []float64{}
	for _, otherAgent := range agentMap {
		if otherAgent.GetID() == ea.GetID() {
			continue // skip self
		}
		if otherAgent.GetClusterID() == clusterID {
			// score := worldviewAlignment(ea.worldview, otherAgent.GetWorldviewBinary())
			score := ea.worldview.CompareWorldviews(otherAgent.GetWorldview())
			clusterAlignments = append(clusterAlignments, score)
		}
	}

	if len(clusterAlignments) == 0 {
		return 0
	}
	// compute average alignment
	var totalAlignment float64
	for _, score := range clusterAlignments {
		totalAlignment += score
	}
	return float32(totalAlignment) / float32(len(clusterAlignments))
}

func (ea *ExtendedAgent) GetNPR() float32 {
	// compute network profiles
	networkAlignments := []float64{}
	agentMap := ea.GetAgentMap()
	for friendID := range ea.network {
		if other, ok := agentMap[friendID]; ok {
			// score := worldviewAlignment(ea.worldview, other.GetWorldviewBinary())
			score := ea.worldview.CompareWorldviews(other.GetWorldview())
			networkAlignments = append(networkAlignments, score)
		}
	}
	if len(networkAlignments) == 0 {
		return 0
	}
	// compute average alignment
	var totalAlignment float64
	for _, score := range networkAlignments {
		totalAlignment += score
	}
	return float32(totalAlignment) / float32(len(networkAlignments))
}

// prop. links agent cut vs links cut to you
// prop. links created vs links created to you
func (ea *ExtendedAgent) GetEstrangement() float32 {
	kin := ea.kinshipGroup
	network := ea.network
	//fmt.Print("Number of kin: ", len(kin), " ")

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
	agentMap := ea.GetAgentMap()
	selfHeroism := ea.GetHeroism()
	network := ea.network

	heroismScores := []int{selfHeroism}

	for id := range network {
		if agent, ok := agentMap[id]; ok {
			heroismScores = append(heroismScores, agent.GetHeroism())
		}
	}

	// Sort heroism scores in ascending order
	sort.Ints(heroismScores)
	index := sort.SearchInts(heroismScores, selfHeroism) + 1

	return float32(index) / float32(len(heroismScores))
}

func (ea *ExtendedAgent) ComputeMortalitySalience(grid *infra.Grid) float32 {
	//w1, w2, w3, w4 := float32(0.25), float32(0.25), float32(0.25), float32(0.25) // tweak

	ce := float32(ea.ClusterEliminations())
	ne := float32(ea.NetworkEliminations())
	ra := float32(ea.RelativeAgeToNetwork())
	mp := float32(ea.GetMemorialProximity(grid))
	// fmt.Printf("Agent %v MS Scores: CE=%.2f, NE=%.2f, RA=%.2f, MP=%.2f\n", ea.GetID(), ce, ne, ra, mp)
	//ea.MortalitySalience = infra.W1*ce + infra.W2*ne + infra.W3*ra + infra.W4*mp

	return infra.W1*ce + infra.W2*ne + infra.W3*ra + infra.W4*mp
}

func (ea *ExtendedAgent) ComputeWorldviewValidation() float32 {
	//w5, w6, w7 := float32(0.25), float32(0.25), float32(0.5) // tweak

	cpr := ea.GetCPR()
	npr := ea.GetNPR()                                      // compute NPR
	ysterofimia := ea.GetYsterofimia().ComputeYsterofimia() // compute ysterofimia
	//fmt.Printf("Agent %v WV Scores: CPR=%.2f, NPR=%.2f, Ysterofimia=%.2f\n", ea.GetID(), cpr, npr, ysterofimia)
	//ea.WorldviewValidation = infra.W5*cpr + infra.W6*npr + infra.W7*ysterofimia

	return infra.W5*cpr + infra.W6*npr + infra.W7*ysterofimia
}

func (ea *ExtendedAgent) ComputeRelationshipValidation() float32 {
	//w8, w9, w10 := float32(0.25), float32(0.25), float32(0.5) // tweak

	est := ea.GetEstrangement()                // compute EST
	pse := ea.GetProSocialEsteem()             // compute PSE
	heroismTendency := ea.GetHeroismTendency() // compute heroism tendency
	fmt.Printf("Agent %v RV Scores: EST=%.2f, PSE=%.2f, HeroismTendency=%.2f\n", ea.GetID(), est, pse, heroismTendency)
	//ea.RelationshipValidation = infra.W8*est + infra.W9*pse + infra.W10*heroismTendency

	return infra.W8*est + infra.W9*pse + infra.W10*heroismTendency
}

// Decision-making logic
func (ea *ExtendedAgent) GetASPDecision(grid *infra.Grid) infra.ASPDecison {
	threshold := float32(0.2) //random threshold

	ms := ea.ComputeMortalitySalience(grid)
	wv := ea.ComputeWorldviewValidation()
	rv := ea.ComputeRelationshipValidation()

	// Debug log
	fmt.Printf("Agent %v ASP Scores: MS=%.2f, WV=%.2f, RV=%.2f\n", ea.GetID(), ms, wv, rv)

	sum := 0
	for _, score := range []float32{ms, wv, rv} {
		if score > threshold {
			sum += 1
		} else {
			sum -= 1
		}
	}

	if sum > 0 {
		// fmt.Printf("✅ Agent %v decided to SELF-SACRIFICE\n", ea.GetID())
		ea.IncrementHeroism()
		return infra.SELF_SACRIFICE // Self-sacrifice
	} else if sum < 0 {
		return infra.NOT_SELF_SACRIFICE // Reject self-sacrifice
	} else {
		return infra.INACTION // No action
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
	//fmt.Printf("Agent %v received wellbeing check from %v\n", ea.GetID(), msg.Sender)
	if rand.Float32() < ea.PTW.ReplyProb {
		reply := ea.CreateReplyMessage()
		//ea.SendMessage(reply, msg.Sender)
		ea.SendSynchronousMessage(reply, msg.Sender)
		//fmt.Printf("Agent %v sending reply message to %v\n", ea.GetID(), msg.Sender)

		//then update alpha
		//fmt.Printf("Agent esteem before: %f\n", ea.network[msg.Sender])
		ea.UpdateEsteem(msg.Sender, true)
		//fmt.Printf("Agent esteem after: %f\n", ea.network[msg.Sender])
	}
}

func (ea *ExtendedAgent) HandleReplyMessage(msg *infra.ReplyMessage) {
	// update alpha
	ea.UpdateEsteem(msg.Sender, true)
	ea.SignalMessagingComplete()
}

// ----------------------- Data Recording Functions -----------------------

func (ea *ExtendedAgent) RecordAgentJSON(instance infra.IExtendedAgent) gameRecorder.JSONAgentRecord {
	styleMap := map[infra.AttachmentType]string{
		infra.DISMISSIVE:  "Dismissive",
		infra.FEARFUL:     "Fearful",
		infra.PREOCCUPIED: "Preoccupied",
		infra.SECURE:      "Secure",
	}

	return gameRecorder.JSONAgentRecord{
		ID:                  ea.GetID().String(),
		IsAlive:             ea.IsAlive(),
		Age:                 ea.GetAge(),
		AttachmentStyle:     styleMap[ea.attachment.Type],
		AttachmentAnxiety:   ea.attachment.Anxiety,
		AttachmentAvoidance: ea.attachment.Avoidance,
		ClusterID:           ea.clusterID,
		Position:            gameRecorder.Position{X: ea.position.X, Y: ea.position.Y},
		// Worldview:           ea.worldview,
		Heroism: ea.heroism,
		//MortalitySalience:      ea.MortalitySalience,
		//WorldviewValidation:    ea.WorldviewValidation,
		//RelationshipValidation: ea.RelationshipValidation,
		//ASPDecison: 		    ea.GetASPDecision(nil),
	}
}
