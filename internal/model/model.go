package model

import "time"

type Role string

const (
	RoleGuest      Role = "GUEST"
	RoleUser       Role = "USER"
	RoleSpecialist Role = "SPECIALIST"
	RoleAdmin      Role = "ADMIN"
)

type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	Role         Role
	IsActive     bool
	Reputation   int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Session struct {
	ID         string
	UserID     int64
	CreatedAt  time.Time
	ExpiresAt  time.Time
	LastActive time.Time
}

type QuestionStatus string

const (
	QuestionOpen      QuestionStatus = "OPEN"
	QuestionClosed    QuestionStatus = "CLOSED"
	QuestionDuplicate QuestionStatus = "DUPLICATE"
)

type Question struct {
	ID                int64
	AuthorID          int64
	AuthorUsername    string
	CategoryID        *int64
	SpecialistID      *int64
	SpecialistUsername *string
	Title             string
	Body              string
	Status            QuestionStatus
	DuplicateOf       *int64
	ViewCount         int
	AnswerCount       int
	HasVerified       bool
	Tags              []Tag
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Comment struct {
	ID             int64
	AnswerID       int64
	AuthorID       int64
	AuthorUsername string
	Body           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Answer struct {
	ID             int64
	QuestionID     int64
	AuthorID       int64
	AuthorUsername string
	Body           string
	IsVerified     bool
	VoteScore      int
	CurrentUserVote *VoteValue
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type VoteValue string

const (
	VoteUp   VoteValue = "UP"
	VoteDown VoteValue = "DOWN"
)

type AnswerHistory struct {
	ID             int64
	AnswerID       int64
	EditorID       int64
	EditorUsername string
	Body           string
	EditedAt       time.Time
}

type QuestionHistory struct {
	ID             int64
	QuestionID     int64
	EditorID       int64
	EditorUsername string
	Title          string
	Body           string
	EditedAt       time.Time
}

type Category struct {
	ID           int64
	Name         string
	Description  string
	SpecialistID *int64
	CreatedAt    time.Time
}

type Tag struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}

type AttachmentTarget string

const (
	AttachmentTargetQuestion AttachmentTarget = "QUESTION"
	AttachmentTargetAnswer   AttachmentTarget = "ANSWER"
	AttachmentTargetComment  AttachmentTarget = "COMMENT"
)

type Attachment struct {
	ID         int64
	UploaderID int64
	TargetType AttachmentTarget
	TargetID   int64
	Filename   string
	ObjectKey  string
	MimeType   string
	SizeBytes  int64
	CreatedAt  time.Time
}

type SearchResult struct {
	ResultType    string         `json:"type"`
	AnswerID      *int64         `json:"answer_id,omitempty"`
	QuestionID    int64          `json:"question_id"`
	QuestionTitle string         `json:"question_title"`
	Snippet       string         `json:"snippet"`
	Status        QuestionStatus `json:"status"`
	AuthorID      int64          `json:"author_id"`
	CategoryID    *int64         `json:"category_id,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	Rank          float32        `json:"rank"`
}

type NotificationType string

const (
	NotifNewAnswer        NotificationType = "NEW_ANSWER"
	NotifAnswerVerified   NotificationType = "ANSWER_VERIFIED"
	NotifQuestionAssigned NotificationType = "QUESTION_ASSIGNED"
	NotifQuestionClosed   NotificationType = "QUESTION_CLOSED"
	NotifNewComment       NotificationType = "NEW_COMMENT"
)

type NotificationPayload struct {
	QuestionID    int64  `json:"question_id,omitempty"`
	QuestionTitle string `json:"question_title,omitempty"`
	AnswerID      int64  `json:"answer_id,omitempty"`
	CommentID     int64  `json:"comment_id,omitempty"`
	ActorUsername string `json:"actor_username,omitempty"`
}

type Notification struct {
	ID        int64
	UserID    int64
	Type      NotificationType
	Payload   NotificationPayload
	IsRead    bool
	CreatedAt time.Time
}

type SearchFilter struct {
	Query      string
	CategoryID *int64
	TagID      *int64
	Status     *QuestionStatus
}
