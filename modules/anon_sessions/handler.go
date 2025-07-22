package anon_sessions

import (
	"CodeCast_backend/db"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

// POST /api/v1/sessions/anon
func CreateAnonSession(c *gin.Context) {
	// Request body contains:
	type requestBody struct {
		// In a perfect world, The frontend would be able to generate a random creator token
		CreatorToken string `json:"creator_token" binding:"required"`
		DisplayName  string `json:"display_name" binding:"required"`
	}
	var req requestBody
	// Validate the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Sets id, session code, and created_at
	id := uuid.New().String()
	sessionCode := uuid.New().String()[:8] // Short unique code
	createdAt := time.Now()

	// Insert the new session into the database
	_, err := db.DB.Exec(`INSERT INTO anon_sessions (id, session_code, creator_token, display_name, is_active, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		id, sessionCode, req.CreatorToken, req.DisplayName, true, createdAt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create session"})
		return
	}

	// Respond with the session details. We are NOT returning the creator token for security reasons.
	c.JSON(201, gin.H{
		"id":           id,
		"session_code": sessionCode,
		"display_name": req.DisplayName,
		"is_active":    true,
		"created_at":   createdAt,
	})
}

// When the user joins a session their UI will reflect.
// A new button must appear in the UI To allow them to "pull" snippets from the session.
// The user gets a random display name.
// POST /api/v1/sessions/anon/:code/join
func JoinAnonSession(c *gin.Context) {
	//Extract session code from the URL parameter.
	code := c.Param("code")

	// Check if session exists and is active
	var sessionID string
	err := db.DB.QueryRow(`SELECT id FROM anon_sessions WHERE session_code = $1 AND is_active = true`, code).Scan(&sessionID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Session not found or inactive"})
		return
	}

	// Check current participant count
	var count int
	err = db.DB.QueryRow(`SELECT COUNT(*) FROM anon_participants WHERE session_id = $1`, sessionID).Scan(&count)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to check participant count"})
		return
	}
	if count >= 20 {
		c.JSON(403, gin.H{"error": "Max amount of participants"})
		return
	}

	// Generate participant ID and random display name
	participantID := uuid.New().String()
	displayName := randomString(3)
	joinedAt := time.Now()

	// Insert the new participant into the database
	_, err = db.DB.Exec(`INSERT INTO anon_participants (id, session_id, display_name, joined_at) VALUES ($1, $2, $3, $4)`,
		participantID, sessionID, displayName, joinedAt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to join session"})
		return
	}

	// Respond with the participant details
	c.JSON(201, gin.H{
		"id":           participantID,
		"session_id":   sessionID,
		"display_name": displayName,
		"joined_at":    joinedAt,
	})
}

// Helper to generate random alphanumeric string
func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// POST /api/v1/sessions/anon/:code/snippets
// This endpoint is ONLY used by the creator of the session. Hence why it requires a creator token in the header.

func PushAnonSnippet(c *gin.Context) {
	// Extract session code from the URL parameter
	code := c.Param("code")

	// In a perfect frontend. The creator token would be store in local storage.
	// Recall that OUR FRONT END IS A VSCODE EXTENSION.

	// Authenticate creator token from header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{"error": "Missing creator token"})
		return
	}

	// Find session by code and check token
	var sessionID string
	var creatorToken string
	err := db.DB.QueryRow(`SELECT id, creator_token FROM anon_sessions WHERE session_code = $1 AND is_active = true`, code).Scan(&sessionID, &creatorToken)
	if err != nil {
		c.JSON(404, gin.H{"error": "Session not found or inactive"})
		return
	}
	if token != creatorToken {
		c.JSON(403, gin.H{"error": "Invalid creator token"})
		return
	}

	// Parse request body
	type requestBody struct {
		// In a perfect world, the frontend would be able to send the file name
		// and content of the snippet.

		FileName string `json:"file_name" binding:"required"`
		Content  string `json:"content" binding:"required"`
	}
	var req requestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Insert snippet
	snippetID := uuid.New().String()
	pushedAt := time.Now()
	_, err = db.DB.Exec(`INSERT INTO anon_snippets (id, session_id, file_name, content, pushed_at) VALUES ($1, $2, $3, $4, $5)`,
		snippetID, sessionID, req.FileName, req.Content, pushedAt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to push snippet"})
		return
	}

	c.JSON(201, gin.H{
		"id":         snippetID,
		"session_id": sessionID,
		"file_name":  req.FileName,
		"content":    req.Content,
		"pushed_at":  pushedAt,
	})
}

// GET /api/v1/sessions/anon/:code/snippets
// Pulls the most recent snippet from the anon_snippets table
func GetAnonSnippets(c *gin.Context) {
	// Extract session code from the URL parameter
	code := c.Param("code")

	// Find session by code
	var sessionID string
	err := db.DB.QueryRow(`SELECT id FROM anon_sessions WHERE session_code = $1 AND is_active = true`, code).Scan(&sessionID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Session not found or inactive"})
		return
	}

	// Fetch the most recent snippet for the session
	var id, fileName, content string
	var pushedAt time.Time
	err = db.DB.QueryRow(`
		SELECT id, file_name, content, pushed_at 
		FROM anon_snippets 
		WHERE session_id = $1 
		ORDER BY pushed_at DESC 
		LIMIT 1`, sessionID).Scan(&id, &fileName, &content, &pushedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "No snippets found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to fetch snippet"})
		}
		return
	}

	// Return the single most recent snippet
	c.JSON(200, gin.H{
		"id":        id,
		"file_name": fileName,
		"content":   content,
		"pushed_at": pushedAt,
	})
}

// POST /api/v1/sessions/anon/:code/end
// This one should literally delete anything with the session_id from the database.
func EndAnonSession(c *gin.Context) {
	code := c.Param("code")

	// Authenticate creator token from header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{"error": "Missing creator token"})
		return
	}

	// Find session by code and check token
	var sessionID string
	var creatorToken string
	err := db.DB.QueryRow(`SELECT id, creator_token FROM anon_sessions WHERE session_code = $1 AND is_active = true`, code).Scan(&sessionID, &creatorToken)
	if err != nil {
		c.JSON(404, gin.H{"error": "Session not found or inactive"})
		return
	}
	if token != creatorToken {
		c.JSON(403, gin.H{"error": "Invalid creator token"})
		return
	}

	// Delete all related records
	_, err = db.DB.Exec(`DELETE FROM anon_snippets WHERE session_id = $1`, sessionID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete snippets"})
		return
	}
	_, err = db.DB.Exec(`DELETE FROM anon_participants WHERE session_id = $1`, sessionID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete participants"})
		return
	}
	_, err = db.DB.Exec(`DELETE FROM anon_sessions WHERE id = $1`, sessionID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete session"})
		return
	}

	c.JSON(200, gin.H{"message": "Session ended and all related data deleted"})
}
