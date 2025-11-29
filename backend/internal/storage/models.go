package storage

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User represents a platform user
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"` // bcrypt hash
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Functions  []Function  `gorm:"foreignKey:UserID" json:"functions,omitempty"`
	APIKeys    []APIKey    `gorm:"foreignKey:UserID" json:"api_keys,omitempty"`
	Executions []Execution `gorm:"foreignKey:UserID" json:"executions,omitempty"`
}

// Function represents a user-uploaded cloud function
type Function struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Runtime     string    `gorm:"not null" json:"runtime"` // nodejs, python, go
	Code        string    `gorm:"type:text;not null" json:"code"`
	EntryPoint  string    `gorm:"default:index.handler" json:"entry_point"`
	MemoryMB    int       `gorm:"default:128" json:"memory_mb"`
	TimeoutSec  int       `gorm:"default:30" json:"timeout_sec"`
	Status      string    `gorm:"default:active" json:"status"` // active, inactive, error
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User       User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Executions []Execution `gorm:"foreignKey:FunctionID" json:"executions,omitempty"`
}

// Execution represents a single function execution
type Execution struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	FunctionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"function_id"`
	Status     string         `gorm:"default:pending" json:"status"` // pending, running, success, failed
	Input      datatypes.JSON `gorm:"type:jsonb" json:"input"`
	Output     datatypes.JSON `gorm:"type:jsonb" json:"output"`
	Error      string         `gorm:"type:text" json:"error,omitempty"`
	Logs       string         `gorm:"type:text" json:"logs"`
	DurationMS int64          `json:"duration_ms"`
	MemoryUsed int            `json:"memory_used"` // in MB
	StartedAt  *time.Time     `json:"started_at,omitempty"`
	CompletedAt *time.Time    `json:"completed_at,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`

	User     User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Function Function `gorm:"foreignKey:FunctionID" json:"function,omitempty"`
}

// APIKey represents an API key for function invocation
type APIKey struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	Key       string    `gorm:"uniqueIndex;not null" json:"key"` // hashed
	Prefix    string    `gorm:"not null" json:"prefix"`          // first 8 chars for display
	LastUsed  *time.Time `json:"last_used,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BeforeCreate hook for User
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook for Function
func (f *Function) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook for Execution
func (e *Execution) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook for APIKey
func (k *APIKey) BeforeCreate(tx *gorm.DB) error {
	if k.ID == uuid.Nil {
		k.ID = uuid.New()
	}
	return nil
}
