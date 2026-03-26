package engine

import (
	"fmt"
	"strings"

	"github.com/blackms/ExplainableEngine/internal/models"
)

// NarrativeLevel controls the detail of the generated text.
type NarrativeLevel string

const (
	LevelBasic    NarrativeLevel = "basic"
	LevelAdvanced NarrativeLevel = "advanced"
)

// NarrativeLanguage controls the output language.
type NarrativeLanguage string

const (
	LangEN NarrativeLanguage = "en"
	LangIT NarrativeLanguage = "it"
)

// NarrativeResult is the output of the narrative engine.
type NarrativeResult struct {
	ExplanationID   string            `json:"explanation_id"`
	Level           NarrativeLevel    `json:"level"`
	Language        NarrativeLanguage `json:"language"`
	Narrative       string            `json:"narrative"`
	ConfidenceLevel string            `json:"confidence_level"`
	HasMissingData  bool              `json:"has_missing_data"`
}

// GenerateNarrative creates a human-readable explanation from an ExplainResponse.
func GenerateNarrative(resp *models.ExplainResponse, level NarrativeLevel, lang NarrativeLanguage) (*NarrativeResult, error) {
	if level != LevelBasic && level != LevelAdvanced {
		return nil, fmt.Errorf("unsupported narrative level: %q", level)
	}
	if lang != LangEN && lang != LangIT {
		return nil, fmt.Errorf("unsupported narrative language: %q", lang)
	}

	confLevel := confidenceLevel(resp.Confidence, lang)
	hasMissing := resp.MissingImpact > 0.1

	var narrative string
	switch level {
	case LevelBasic:
		narrative = generateBasic(resp, confLevel, lang)
	case LevelAdvanced:
		narrative = generateAdvanced(resp, confLevel, lang, hasMissing)
	}

	return &NarrativeResult{
		ExplanationID:   resp.ID,
		Level:           level,
		Language:        lang,
		Narrative:       narrative,
		ConfidenceLevel: confLevel,
		HasMissingData:  hasMissing,
	}, nil
}

// confidenceLevel returns a human-readable confidence label.
func confidenceLevel(confidence float64, lang NarrativeLanguage) string {
	if lang == LangIT {
		switch {
		case confidence >= 0.8:
			return "alta"
		case confidence >= 0.5:
			return "moderata"
		default:
			return "bassa"
		}
	}
	// English (default)
	switch {
	case confidence >= 0.8:
		return "high"
	case confidence >= 0.5:
		return "moderate"
	default:
		return "low"
	}
}

// driverPercentage looks up the percentage for a driver name from the breakdown items.
func driverPercentage(name string, breakdown []models.BreakdownItem) float64 {
	for _, b := range breakdown {
		if b.Label == name {
			return b.Percentage
		}
	}
	return 0
}

func generateBasic(resp *models.ExplainResponse, confLevel string, lang NarrativeLanguage) string {
	if len(resp.TopDrivers) == 0 {
		if lang == LangIT {
			return fmt.Sprintf("Il punteggio %s è %.2f. La confidenza è %s.", resp.Target, resp.FinalValue, confLevel)
		}
		return fmt.Sprintf("The %s score is %.2f. Confidence is %s.", resp.Target, resp.FinalValue, confLevel)
	}

	topDriver := resp.TopDrivers[0]
	pct := driverPercentage(topDriver.Name, resp.Breakdown)

	if lang == LangIT {
		return fmt.Sprintf("Il punteggio %s è %.2f, guidato principalmente da %s (%.1f%%). La confidenza è %s.",
			resp.Target, resp.FinalValue, topDriver.Name, pct, confLevel)
	}
	return fmt.Sprintf("The %s score is %.2f, primarily driven by %s (%.1f%%). Confidence is %s.",
		resp.Target, resp.FinalValue, topDriver.Name, pct, confLevel)
}

func generateAdvanced(resp *models.ExplainResponse, confLevel string, lang NarrativeLanguage, hasMissing bool) string {
	var sb strings.Builder

	confPct := resp.Confidence * 100

	if lang == LangIT {
		sb.WriteString(fmt.Sprintf("Il punteggio %s è %.2f con confidenza %s (%.1f%%).",
			resp.Target, resp.FinalValue, confLevel, confPct))
	} else {
		sb.WriteString(fmt.Sprintf("The %s score is %.2f with %s confidence (%.1f%%).",
			resp.Target, resp.FinalValue, confLevel, confPct))
	}

	// Driver lines
	drivers := resp.TopDrivers
	if len(drivers) > 3 {
		drivers = drivers[:3]
	}

	if len(drivers) > 0 {
		sb.WriteString("\n\n")
		if lang == LangIT {
			sb.WriteString("Fattori principali:")
		} else {
			sb.WriteString("Key drivers:")
		}

		for _, d := range drivers {
			pct := driverPercentage(d.Name, resp.Breakdown)
			sb.WriteString("\n")
			if lang == LangIT {
				sb.WriteString(fmt.Sprintf("- %s: contributo del %.1f%% (impatto: %.2f)", d.Name, pct, d.Impact))
			} else {
				sb.WriteString(fmt.Sprintf("- %s: %.1f%% contribution (impact: %.2f)", d.Name, pct, d.Impact))
			}
		}
	}

	// Missing data warning
	if hasMissing {
		missingPct := resp.MissingImpact * 100
		sb.WriteString("\n\n")
		if lang == LangIT {
			sb.WriteString(fmt.Sprintf("Nota: il %.1f%% dei dati di input è mancante, il che potrebbe influire sull'affidabilità.", missingPct))
		} else {
			sb.WriteString(fmt.Sprintf("Note: %.1f%% of input data is missing, which may affect reliability.", missingPct))
		}
	}

	return sb.String()
}
