package nrclex

type APIResponse struct {
	Text              string            `json:"text"`
	EmotionScore      EmotionScores     `json:"emotion_scores"`
	VaderEmotionScore VaderEmotionScore `json:"vader_emotion_scores"`
}

type APIErrorResponse struct {
	Error string `json:"error"`
}
