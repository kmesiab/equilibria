package nrclex

type EmotionScores struct {
	Anger        float64 `json:"anger"`
	Anticipation float64 `json:"anticipation"`
	Disgust      float64 `json:"disgust"`
	Fear         float64 `json:"fear"`
	Trust        float64 `json:"trust"`
	Joy          float64 `json:"joy"`
	Negative     float64 `json:"negative"`
	Positive     float64 `json:"positive"`
	Sadness      float64 `json:"sadness"`
	Surprise     float64 `json:"surprise"`
}
