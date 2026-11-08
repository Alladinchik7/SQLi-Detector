package report

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Alladinchik/slqi-detector/internal/scanner"
)

// internal/report
type Summary struct {
	StartedAt  time.Time         `json:"started_at"`
	FinishedAt time.Time         `json:"finished_at"`
	Targets    int               `json:"targets_count"`
	Tasks      int               `json:"tasks_count"`
	Findings   []scanner.Finding `json:"findings"`
}

func CreateSummary(started, finished time.Time, targets, tasks int, findings []scanner.Finding) *Summary {
	return &Summary{
		StartedAt:  started,
		FinishedAt: finished,
		Targets:    targets,
		Tasks:      tasks,
		Findings:   findings,
	}
}

func (s *Summary) Write(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if path == "" {
		_, err := os.Stdout.Write(append(data, '\n'))
		return err
	}
	return os.WriteFile(path, data, 0644)
}
