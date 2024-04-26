package report_generator

import (
	"time"

	"github.com/kmesiab/equilibria/lambdas/lib"
	"github.com/kmesiab/equilibria/lambdas/lib/ai"
	"github.com/kmesiab/equilibria/lambdas/lib/message"
)

type PatientReportGenerationRequest struct {
	PatientID  int64      `json:"patient_id"`
	ReportType ReportType `json:"report_type"`
}

type PatientReportGenerationResponse struct {
	PatientID string    `json:"patient_id"`
	ReportID  string    `json:"report_id"`
	ReportURL string    `json:"report_url"`
	CreatedOn time.Time `json:"created_on"`
}

type PatientReportGeneratorLambdaHandler struct {
	lib.LambdaHandler

	MemoryService     *message.MemoryService
	CompletionService ai.CompletionServiceInterface
}
