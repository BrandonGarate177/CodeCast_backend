package anon_sessions

import "time"

type AnonSession struct {
	ID           string    `db:"id"`
	SessionCode  string    `db:"session_code"`
	CreatorToken string    `db:"creator_token"`
	DisplayName  string    `db:"display_name"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
}

type AnonParticipant struct {
	ID          string    `db:"id"`
	SessionID   string    `db:"session_id"`
	DisplayName string    `db:"display_name"`
	JoinedAt    time.Time `db:"joined_at"`
}

type AnonSnippet struct {
	ID        string    `db:"id"`
	SessionID string    `db:"session_id"`
	FileName  string    `db:"file_name"`
	Content   string    `db:"content"`
	PushedAt  time.Time `db:"pushed_at"`
}
