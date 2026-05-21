package app

import (
	"encoding/json"
	"fmt"
	"io"
)

type Check struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Message     string `json:"message,omitempty"`
	Remediation string `json:"remediation,omitempty"`
}

type CommandError struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Remediation string `json:"remediation,omitempty"`
}

type Result struct {
	Command        string        `json:"command"`
	Mode           string        `json:"mode"`
	Status         string        `json:"status"`
	ExitCode       int           `json:"exit_code"`
	Checks         []Check       `json:"checks"`
	PlannedActions []string      `json:"planned_actions"`
	Mutations      []string      `json:"mutations"`
	Error          *CommandError `json:"error,omitempty"`
	Results        []QueryResult `json:"results,omitempty"`
}

type QueryResult struct {
	Rank        float64 `json:"rank"`
	ID          string  `json:"id"`
	SourcePath  string  `json:"source_path"`
	Layer       string  `json:"layer"`
	Type        string  `json:"type"`
	Status      string  `json:"status"`
	Priority    string  `json:"priority"`
	Confidence  string  `json:"confidence"`
	ContextCost string  `json:"context_cost"`
	Summary     string  `json:"summary"`
	MatchReason string  `json:"match_reason"`
	GraphID     string  `json:"graph_id,omitempty"`
	GraphSource string  `json:"graph_source,omitempty"`
	EdgeType    string  `json:"edge_type,omitempty"`
	Target      string  `json:"target,omitempty"`
	Reason      string  `json:"reason,omitempty"`
	Validation  string  `json:"validation,omitempty"`
	GraphFile   string  `json:"graph_file,omitempty"`
}

func writeJSON(w io.Writer, result Result) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func writePlain(w io.Writer, result Result) error {
	if _, err := fmt.Fprintf(w, "ai-skill %s: %s (exit %d)\n", result.Command, result.Status, result.ExitCode); err != nil {
		return err
	}

	for _, check := range result.Checks {
		if check.Message == "" {
			if _, err := fmt.Fprintf(w, "- %s: %s\n", check.Name, check.Status); err != nil {
				return err
			}
			continue
		}
		if _, err := fmt.Fprintf(w, "- %s: %s - %s\n", check.Name, check.Status, check.Message); err != nil {
			return err
		}
	}

	if len(result.PlannedActions) > 0 {
		if _, err := fmt.Fprintln(w, "Planned actions:"); err != nil {
			return err
		}
		for _, action := range result.PlannedActions {
			if _, err := fmt.Fprintf(w, "- %s\n", action); err != nil {
				return err
			}
		}
	}

	if len(result.Mutations) > 0 {
		if _, err := fmt.Fprintln(w, "Mutations:"); err != nil {
			return err
		}
		for _, mutation := range result.Mutations {
			if _, err := fmt.Fprintf(w, "- %s\n", mutation); err != nil {
				return err
			}
		}
	}

	if len(result.Results) > 0 {
		if _, err := fmt.Fprintln(w, "Results:"); err != nil {
			return err
		}
		for _, item := range result.Results {
			if _, err := fmt.Fprintf(w, "- %s (%s): %s\n", item.ID, item.SourcePath, item.Summary); err != nil {
				return err
			}
		}
	}

	if result.Error != nil {
		if _, err := fmt.Fprintf(w, "Error: %s\n", result.Error.Message); err != nil {
			return err
		}
		if result.Error.Remediation != "" {
			if _, err := fmt.Fprintf(w, "Remediation: %s\n", result.Error.Remediation); err != nil {
				return err
			}
		}
	}

	return nil
}
