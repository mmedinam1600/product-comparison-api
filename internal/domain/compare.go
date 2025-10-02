package domain

// CompareRequest representa la solicitud de comparación de productos
type CompareRequest struct {
	Ids    []string  `json:"ids" binding:"required,min=1"`
	Fields *[]string `json:"fields,omitempty"`
}

// Metric define el tipo de métrica para la comparación de campos
type Metric string

const (
	LowerIsBetter  Metric = "lower_is_better"
	HigherIsBetter Metric = "higher_is_better"
	TrueIsBetter   Metric = "true_is_better"
)

// DiffField representa las diferencias de un campo específico entre productos
type DiffField struct {
	Values map[string]interface{} `json:"values"`
	Metric *Metric                `json:"metric,omitempty"`
	Best   []string               `json:"best"`
}

// CompareResult contiene el resultado de la comparación
type CompareResult struct {
	Items        []Item               `json:"items"`
	SharedFields []string             `json:"shared_fields"`
	Diff         map[string]DiffField `json:"diff"`
}

// ComparePolicy contiene la configuración de la comparación aplicada
type ComparePolicy struct {
	EffectiveMode      string   `json:"effective_mode"`
	ComparabilityScore float64  `json:"comparability_score"`
	Warnings           []string `json:"warnings,omitempty"`
}

// Metadata contiene metadatos adicionales de la comparación
type Metadata struct {
	Order           []string      `json:"order"`
	RequestedFields *[]string     `json:"requested_fields,omitempty"`
	ResolvedFields  []string      `json:"resolved_fields"`
	ComparePolicy   ComparePolicy `json:"compare_policy"`
	Currency        string        `json:"currency"`
	Version         string        `json:"version"`
}
