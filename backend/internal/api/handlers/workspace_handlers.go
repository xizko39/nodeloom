package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xizko39/nodeloom/internal/workspace"
)

// Initialize the workspace service
var workspaceService *workspace.SupabaseService

func InitWorkspaceHandlers(ws *workspace.SupabaseService) {
	workspaceService = ws
}

// CreateWorkspace handles creating a new workspace
func CreateWorkspace(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workspace, err := workspaceService.CreateWorkspace(uuid.New(), req.Name) // Pass a dummy user ID for now
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workspace"})
		return
	}

	c.JSON(http.StatusCreated, workspace)
}

// GetWorkspace handles fetching a specific workspace by ID
func GetWorkspace(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	workspace, err := workspaceService.GetWorkspace(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
		return
	}

	c.JSON(http.StatusOK, workspace)
}

// GetWorkspaces handles fetching all workspaces
func GetWorkspaces(c *gin.Context) {
	workspaces, err := workspaceService.GetAllWorkspaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch workspaces"})
		return
	}

	c.JSON(http.StatusOK, workspaces)
}

// UpdateWorkspace handles updating a specific workspace by ID
func UpdateWorkspace(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workspace, err := workspaceService.UpdateWorkspace(id, req.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
		return
	}

	c.JSON(http.StatusOK, workspace)
}

// DeleteWorkspace handles deleting a specific workspace by ID
func DeleteWorkspace(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	err = workspaceService.DeleteWorkspace(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddNode handles adding a new node to a workspace
func AddNode(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	var req struct {
		Type     workspace.NodeType `json:"type" binding:"required"`
		Label    string             `json:"label" binding:"required"`
		Position workspace.Position `json:"position" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node, err := workspaceService.AddNode(workspaceID, req.Type, req.Label, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add node"})
		return
	}

	c.JSON(http.StatusCreated, node)
}

// RemoveNode handles removing a specific node from a workspace
func RemoveNode(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	nodeID, err := uuid.Parse(c.Param("nodeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID"})
		return
	}

	err = workspaceService.RemoveNode(workspaceID, nodeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddEdge handles adding an edge between nodes in a workspace
func AddEdge(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	var req struct {
		Source uuid.UUID `json:"source" binding:"required"`
		Target uuid.UUID `json:"target" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	edge, err := workspaceService.AddEdge(workspaceID, req.Source, req.Target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add edge"})
		return
	}

	c.JSON(http.StatusCreated, edge)
}

// RemoveEdge handles removing an edge from a workspace
func RemoveEdge(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	edgeID, err := uuid.Parse(c.Param("edgeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid edge ID"})
		return
	}

	err = workspaceService.RemoveEdge(workspaceID, edgeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Edge not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
