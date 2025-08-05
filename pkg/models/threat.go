package models

type ThreatResult struct {
	*Result            // Herda do Result existente
	Technique  string  `gorm:"column:technique" json:"technique"`
	Confidence float64 `gorm:"column:confidence" json:"confidence"`
	Risk       string  `gorm:"column:risk" json:"risk"`
	BaseDomain string  `gorm:"column:base_domain" json:"base_domain"`
	Similarity float64 `gorm:"column:similarity" json:"similarity"`
}

func (*ThreatResult) TableName() string {
	return "threat_results"
}
