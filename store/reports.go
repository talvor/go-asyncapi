package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/talvor/asyncapi/dto"
)

type ReportStore struct {
	db *sqlx.DB
}

func NewReportStore(db *sql.DB) *ReportStore {
	return &ReportStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

func (s *ReportStore) Create(ctx context.Context, userID uuid.UUID, reportType string) (*dto.Report, error) {
	const dml = `INSERT INTO reports (user_id, report_type) VALUES ($1, $2) RETURNING *`

	var report dto.Report
	if err := s.db.GetContext(ctx, &report, dml, userID, reportType); err != nil {
		return nil, fmt.Errorf("failed to insert report for user %s: %w", userID, err)
	}
	return &report, nil
}

func (s *ReportStore) Update(ctx context.Context, report *dto.Report) (*dto.Report, error) {
	const dml = `UPDATE reports SET 
	              output_file_path = $1, 
	              download_url = $2, 
	              download_url_expires_at = $3, 
	              error_message = $4, 
	              started_at = $5, 
	              completed_at = $6, 
	              failed_at = $7 
	            WHERE user_id = $8 AND id = $9 RETURNING *`

	var updatedReport dto.Report
	if err := s.db.GetContext(ctx, &updatedReport, dml,
		report.OutputFilePath,
		report.DownloadURL,
		report.DownloadURLExpiresAt,
		report.ErrorMessage,
		report.StartedAt,
		report.CompletedAt,
		report.FailedAt,
		report.UserID,
		report.ID,
	); err != nil {
		return nil, fmt.Errorf("failed to update report %s for user %s: %w", report.ID, report.UserID, err)
	}
	return &updatedReport, nil
}

func (s *ReportStore) ByPrimaryKey(ctx context.Context, userID, reportID uuid.UUID) (*dto.Report, error) {
	const query = `SELECT * FROM reports WHERE user_id = $1 AND id = $2`

	var report dto.Report
	if err := s.db.GetContext(ctx, &report, query, userID, reportID); err != nil {
		return nil, fmt.Errorf("failed to get report %s for user %s: %w", reportID, userID, err)
	}
	return &report, nil
}
