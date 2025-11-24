package job_usecase

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ACTSMLTrainingPayload represents the v1 ACTSML Training Payload Specification
type ACTSMLTrainingPayload struct {
	ProjectID     string                 `json:"project_id"`
	ExperimentID  string                 `json:"experiment_id"`
	Dataset       DatasetConfig          `json:"dataset"`
	Problem       ProblemConfig          `json:"problem"`
	Hyperparameters map[string]interface{} `json:"hyperparameters"`
	Preprocessing PreprocessingConfig    `json:"preprocessing"`
	Compute      ComputeConfig          `json:"compute"`
	Output       OutputConfig           `json:"output"`
	Metadata     MetadataConfig         `json:"metadata,omitempty"`
}

// DatasetConfig represents dataset configuration
type DatasetConfig struct {
	Bucket         string   `json:"bucket"`
	Path           string   `json:"path"`
	Version        string   `json:"version"`
	TargetColumn   string   `json:"target_column"`
	FeatureColumns []string `json:"feature_columns,omitempty"`
}

// ProblemConfig represents problem definition
type ProblemConfig struct {
	Type      string `json:"type"`      // classification, binary_classification, regression
	Algorithm string `json:"algorithm"` // random_forest, logistic_regression, xgboost, etc.
}

// PreprocessingConfig represents preprocessing steps
type PreprocessingConfig struct {
	MissingValues     string                 `json:"missing_values,omitempty"`     // mean, median, drop, none
	Scaling          string                 `json:"scaling,omitempty"`             // standard, minmax, robust, none
	EncodeCategorical string                 `json:"encode_categorical,omitempty"` // onehot, label, none
	TrainTestSplit   *TrainTestSplitConfig  `json:"train_test_split,omitempty"`
}

// TrainTestSplitConfig represents train/test split configuration
type TrainTestSplitConfig struct {
	TestSize    float64 `json:"test_size,omitempty"`
	Shuffle     bool    `json:"shuffle,omitempty"`
	RandomState int     `json:"random_state,omitempty"`
}

// ComputeConfig represents compute requirements
type ComputeConfig struct {
	ResourceClass string `json:"resource_class,omitempty"` // cpu_small, cpu_medium, etc.
	NumCPUs       int    `json:"num_cpus,omitempty"`
	Memory        string `json:"memory,omitempty"` // e.g., "2Gi"
}

// OutputConfig represents output configuration
type OutputConfig struct {
	Bucket string `json:"bucket"`
	Path   string `json:"path"`
}

// MetadataConfig represents additional metadata
type MetadataConfig struct {
	SubmittedBy string `json:"submitted_by,omitempty"`
	Description string `json:"description,omitempty"`
}

// ValidatePayload validates the ACTSML training payload against v1 specification
func ValidatePayload(payload json.RawMessage) (*ACTSMLTrainingPayload, error) {
	if len(payload) == 0 {
		return nil, fmt.Errorf("payload cannot be empty")
	}

	var p ACTSMLTrainingPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, fmt.Errorf("invalid JSON payload: %w", err)
	}

	// Validate required fields
	var missingFields []string

	if p.ProjectID == "" {
		missingFields = append(missingFields, "project_id")
	}
	if p.ExperimentID == "" {
		missingFields = append(missingFields, "experiment_id")
	}
	if p.Dataset.Bucket == "" {
		missingFields = append(missingFields, "dataset.bucket")
	}
	if p.Dataset.Path == "" {
		missingFields = append(missingFields, "dataset.path")
	}
	if p.Dataset.TargetColumn == "" {
		missingFields = append(missingFields, "dataset.target_column")
	}
	if p.Problem.Type == "" {
		missingFields = append(missingFields, "problem.type")
	}
	if p.Problem.Algorithm == "" {
		missingFields = append(missingFields, "problem.algorithm")
	}
	if p.Output.Bucket == "" {
		missingFields = append(missingFields, "output.bucket")
	}
	if p.Output.Path == "" {
		missingFields = append(missingFields, "output.path")
	}

	if len(missingFields) > 0 {
		return nil, fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}

	// Validate problem type
	validProblemTypes := map[string]bool{
		"classification":        true,
		"binary_classification": true,
		"regression":            true,
	}
	if !validProblemTypes[p.Problem.Type] {
		return nil, fmt.Errorf("invalid problem.type: %s (must be one of: classification, binary_classification, regression)", p.Problem.Type)
	}

	// Validate algorithm (basic check - can be extended)
	if p.Problem.Algorithm == "" {
		return nil, fmt.Errorf("problem.algorithm cannot be empty")
	}

	return &p, nil
}

