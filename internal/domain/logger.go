package domain

import (
	"database/sql"
	"time"
)

type SqlNullString struct {
	sql.NullString
}

func NewSqlNullString(s string) SqlNullString {
	if s == "" {
		return SqlNullString{NullString: sql.NullString{String: "", Valid: false}}
	}
	return SqlNullString{NullString: sql.NullString{String: s, Valid: true}}
}

const LOG_TABLE = "request_logs"

type RequestLog struct {
	ID         uint          `gorm:"primaryKey;autoIncrement;column:id"`
	Timestamp  time.Time     `gorm:"not null;column:timestamp"`
	IP         string        `gorm:"size:45;not null;column:ip"`
	Method     string        `gorm:"size:10;not null;column:method"`
	Path       string        `gorm:"size:255;not null;column:path"`
	QueryParam string        `gorm:"type:text;column:query_param"`
	UserAgent  string        `gorm:"type:text;column:user_agent"`
	Headers    string        `gorm:"type:jsonb;column:headers"`
	Body       SqlNullString `gorm:"type:text;column:body"`
	Response   SqlNullString `gorm:"type:text;column:response"`
	StatusCode int           `gorm:"not null;column:status_code"`
	CreatedAt  time.Time     `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt  time.Time     `gorm:"autoUpdateTime;column:updated_at"`
	Latency    int64         `gorm:"column:latency_ms"` // latency dalam ms
}

// TableName custom untuk GORM
func (RequestLog) TableName() string {
	return LOG_TABLE
}
