package emotions

import (
	"github.com/kmesiab/equilibria/lambdas/lib/nrclex"
	"github.com/kmesiab/equilibria/lambdas/models"
)

type NRCLexService struct {
	client *nrclex.NRCLexClient
	repo   *nrclex.Repository
}

func NewNRCLexService(client *nrclex.NRCLexClient, repo *nrclex.Repository) *NRCLexService {
	return &NRCLexService{
		client: client,
		repo:   repo,
	}
}

func (s *NRCLexService) ProcessMessage(user *models.User, message *models.Message) (*nrclex.APIResponse, error) {

	scores, err := s.client.AnalyzeText(message.Body)

	if err != nil {

		return nil, err
	}

	emotions := &models.NrcLex{
		UserID:        user.ID,
		MessageID:     message.ID,
		Anger:         scores.EmotionScore.Anger,
		Anticipation:  scores.EmotionScore.Anticipation,
		Disgust:       scores.EmotionScore.Disgust,
		Fear:          scores.EmotionScore.Fear,
		Trust:         scores.EmotionScore.Trust,
		Joy:           scores.EmotionScore.Joy,
		Negative:      scores.EmotionScore.Negative,
		Positive:      scores.EmotionScore.Positive,
		Sadness:       scores.EmotionScore.Sadness,
		Surprise:      scores.EmotionScore.Surprise,
		VaderCompound: scores.VaderEmotionScore.Compound,
		VaderNeg:      scores.VaderEmotionScore.Neg,
		VaderNeu:      scores.VaderEmotionScore.Neu,
		VaderPos:      scores.VaderEmotionScore.Pos,
	}

	err = s.repo.Create(emotions)

	if err != nil {

		return nil, err
	}

	return scores, nil

}
