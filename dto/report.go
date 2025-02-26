package dto

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	UserID               uuid.UUID  `db:"user_id"`
	ID                   uuid.UUID  `db:"id"`
	ReportType           string     `db:"report_type"`
	OutputFilePath       *string    `db:"output_file_path"`
	DownloadURL          *string    `db:"download_url"`
	DownloadURLExpiresAt *time.Time `db:"download_url_expires_at"`
	ErrorMessage         *string    `db:"error_message"`
	CreatedAt            time.Time  `db:"created_at"`
	StartedAt            *time.Time `db:"started_at"`
	CompletedAt          *time.Time `db:"completed_at"`
	FailedAt             *time.Time `db:"failed_at"`
}
