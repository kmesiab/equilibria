package report_generator

import (
	"encoding/json"
	"errors"
)

type ReportType string

const (
	ReportTypeTrigger   ReportType = "trigger"
	ReportTypeParent    ReportType = "parent"
	ReportTypeMinor     ReportType = "minor"
	ReportTypeAdult     ReportType = "adult"
	ReportTypeTherapist ReportType = "therapist"
	ReportTypeProgress  ReportType = "progress"
)

func (r *ReportType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	reportType := ReportType(s)

	switch reportType {
	case ReportTypeTrigger,
		ReportTypeParent,
		ReportTypeMinor,
		ReportTypeAdult,
		ReportTypeTherapist,
		ReportTypeProgress:
		*r = reportType

		return nil
	default:
		return errors.New("invalid ReportType value")
	}
}
