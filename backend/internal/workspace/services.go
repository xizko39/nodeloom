package workspace

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/xizko39/nodeloom/internal/database"
)

// Define custom error messages
var (
	ErrWorkspaceNotFound = fmt.Errorf("workspace not found")
	ErrNodeNotFound      = fmt.Errorf("node not found")
	ErrEdgeNotFound      = fmt.Errorf("edge not found")
)

// SupabaseService handles workspace operations using Supabase
type SupabaseService struct {
	client *database.SupabaseClient
}

// NewSupabaseService initializes a new Supabase-based WorkspaceService
func NewSupabaseService(client *database.SupabaseClient) *SupabaseService {
	return &SupabaseService{
		client: client,
	}
}

// CreateWorkspace creates a new workspace in Supabase
func (s *SupabaseService) CreateWorkspace(userID uuid.UUID, name string) (*Workspace, error) {
	workspace := Workspace{
		ID:    uuid.New(),
		Name:  name,
		Nodes: []Node{},
		Edges: []Edge{},
	}

	body, status, err := s.client.Request("POST", "workspaces", workspace)
	if err != nil {
		return nil, err
	}

	if status != http.StatusCreated {
		log.Printf("Supabase returned status %d: %s", status, string(body))
		return nil, fmt.Errorf("failed to create workspace: %s", string(body))
	}

	var insertedWorkspaces []Workspace
	err = json.Unmarshal(body, &insertedWorkspaces)
	if err != nil {
		return nil, err
	}

	if len(insertedWorkspaces) == 0 {
		return nil, fmt.Errorf("no workspace was inserted")
	}

	return &insertedWorkspaces[0], nil
}

// GetWorkspace retrieves a workspace with its nodes and edges from Supabase
func (s *SupabaseService) GetWorkspace(id uuid.UUID) (*Workspace, error) {
	// Fetch workspace
	workspaceEndpoint := fmt.Sprintf("workspaces?id=eq.%s", id.String())
	body, status, err := s.client.Request("GET", workspaceEndpoint, nil)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		log.Printf("Supabase returned status %d: %s", status, string(body))
		return nil, fmt.Errorf("failed to get workspace: %s", string(body))
	}

	var workspaces []Workspace
	err = json.Unmarshal(body, &workspaces)
	if err != nil {
		return nil, err
	}

	if len(workspaces) == 0 {
		return nil, ErrWorkspaceNotFound
	}

	workspace := &workspaces[0]

	// Fetch associated nodes
	nodesEndpoint := fmt.Sprintf("nodes?workspace_id=eq.%s", id.String())
	body, status, err = s.client.Request("GET", nodesEndpoint, nil)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		log.Printf("Supabase returned status %d: %s", status, string(body))
		return nil, fmt.Errorf("failed to get nodes: %s", string(body))
	}

	var nodes []Node
	err = json.Unmarshal(body, &nodes)
	if err != nil {
		return nil, err
	}
	workspace.Nodes = nodes

	// Fetch associated edges
	edgesEndpoint := fmt.Sprintf("edges?workspace_id=eq.%s", id.String())
	body, status, err = s.client.Request("GET", edgesEndpoint, nil)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		log.Printf("Supabase returned status %d: %s", status, string(body))
		return nil, fmt.Errorf("failed to get edges: %s", string(body))
	}

	var edges []Edge
	err = json.Unmarshal(body, &edges)
	if err != nil {
		return nil, err
	}
	workspace.Edges = edges

	return workspace, nil
}

// GetAllWorkspaces retrieves all workspaces from Supabase
func (s *SupabaseService) GetAllWorkspaces() ([]Workspace, error) {
	body, status, err := s.client.Request("GET", "workspaces", nil)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		log.Printf("Supabase returned status %d: %s", status, string(body))
		return nil, fmt.Errorf("failed to fetch workspaces: %s", string(body))
	}

	var workspaces []Workspace
	err = json.Unmarshal(body, &workspaces)
	if err != nil {
		return nil, err
	}

	return workspaces, nil
}

// UpdateWorkspace updates a workspace's name in Supabase
func (s *SupabaseService) UpdateWorkspace(id uuid.UUID, name string) (*Workspace, error) {
	update := map[string]string{"name": name}

	body, status, err := s.client.Request("PATCH", fmt.Sprintf("workspaces?id=eq.%s", id.String()), update)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		log.Printf("Supabase returned status %d: %s", status, string(body))
		return nil, fmt.Errorf("failed to update workspace: %s", string(body))
	}

	var updatedWorkspaces []Workspace
	err = json.Unmarshal(body, &updatedWorkspaces)
	if err != nil {
		return nil, err
	}

	if len(updatedWorkspaces) == 0 {
		return nil, ErrWorkspaceNotFound
	}

	return &updatedWorkspaces[0], nil
}

// DeleteWorkspace deletes a workspace from Supabase
func (s *SupabaseService) DeleteWorkspace(id uuid.UUID) error {
	body, status, err := s.client.Request("DELETE", fmt.Sprintf("workspaces?id=eq.%s", id.String()), nil)
	if err != nil {
		return err
	}

	if status != http.StatusOK && status != http.StatusNoContent {
		log.Printf("Supabase returned status %d: %s", status, string(body))
		return fmt.Errorf("failed to delete workspace: %s", string(body))
	}

	return nil
}

// AddNode adds a new node to a workspace in Supabase
func (s *SupabaseService) AddNode(workspaceID uuid.UUID, nodeType NodeType, label string, position Position) (*Node, error) {
	node := Node{
		ID:       uuid.New(),
		Type:     nodeType,
		Label:    label,
		Data:     make(map[string]interface{}),
		Position: position,
	}

	body, status, err := s.client.Request("POST", "nodes", node)
	if err != nil {
		return nil, err
	}

	if status != http.StatusCreated {
		log.Printf("Supabase returned status %d: %s", status, string(body))
		return nil, fmt.Errorf("failed to add node: %s", string(body))
	}

	var insertedNodes []Node
	err = json.Unmarshal(body, &insertedNodes)
	if err != nil {
		return nil, err
	}

	if len(insertedNodes) == 0 {
		return nil, fmt.Errorf("no node was inserted")
	}

	return &insertedNodes[0], nil
}

// RemoveNode removes a node from a workspace in Supabase
func (s *SupabaseService) RemoveNode(workspaceID, nodeID uuid.UUID) error {
	endpoint := fmt.Sprintf("nodes?id=eq.%s&workspace_id=eq.%s", nodeID.String(), workspaceID.String())
	body, status, err := s.client.Request("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	if status != http.StatusOK && status != http.StatusNoContent {
		log.Printf("Supabase returned status %d: %s", status, string(body))
		return fmt.Errorf("failed to remove node: %s", string(body))
	}

	return nil
}

// AddEdge adds a new edge to a workspace in Supabase
func (s *SupabaseService) AddEdge(workspaceID, sourceID, targetID uuid.UUID) (*Edge, error) {
	edge := Edge{
		ID:     uuid.New(),
		Source: sourceID,
		Target: targetID,
	}

	body, status, err := s.client.Request("POST", "edges", edge)
	if err != nil {
		return nil, err
	}

	if status != http.StatusCreated {
		log.Printf("Supabase returned status %d: %s", status, string(body))
		return nil, fmt.Errorf("failed to add edge: %s", string(body))
	}

	var insertedEdges []Edge
	err = json.Unmarshal(body, &insertedEdges)
	if err != nil {
		return nil, err
	}

	if len(insertedEdges) == 0 {
		return nil, fmt.Errorf("no edge was inserted")
	}

	return &insertedEdges[0], nil
}

// RemoveEdge removes an edge from a workspace in Supabase
func (s *SupabaseService) RemoveEdge(workspaceID, edgeID uuid.UUID) error {
	endpoint := fmt.Sprintf("edges?id=eq.%s&workspace_id=eq.%s", edgeID.String(), workspaceID.String())
	body, status, err := s.client.Request("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	if status != http.StatusOK && status != http.StatusNoContent {
		log.Printf("Supabase returned status %d: %s", status, string(body))
		return fmt.Errorf("failed to remove edge: %s", string(body))
	}

	return nil
}
