package server

import "github.com/google/uuid"

func (serv *TMTServer) CreateNetworkConnection(fromAgentID, toAgentID uuid.UUID, strength float32) {
	agentMap := serv.GetAgentMap()
	fromAgent, fromExists := agentMap[fromAgentID]
	toAgent, toExists := agentMap[toAgentID]

	if !fromExists || !toExists {
		return
	}

	fromAgent.AddToSocialNetwork(toAgentID, strength)

	fromAgent.PerformCreatedConnection(toAgentID)
	toAgent.ReceiveCreatedConnection(fromAgentID)
}

func (serv *TMTServer) SeverNetworkConnection(fromAgentID, toAgentID uuid.UUID) {
	agentMap := serv.GetAgentMap()
	fromAgent, fromExists := agentMap[fromAgentID]
	toAgent, toExists := agentMap[toAgentID]

	if !fromExists || !toExists {
		return
	}

	fromAgent.RemoveFromSocialNetwork(toAgentID)

	toAgent.PerformSeveredConnected(fromAgentID)
	toAgent.ReceiveSeveredConnected(fromAgentID)
}
