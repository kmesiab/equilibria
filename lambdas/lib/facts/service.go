package facts

import (
	"fmt"

	"github.com/kmesiab/equilibria/lambdas/lib/ai"
	"github.com/kmesiab/equilibria/lambdas/lib/ai/agents"
	"github.com/kmesiab/equilibria/lambdas/lib/ai/agents/fact_agent"
	"github.com/kmesiab/equilibria/lambdas/models"
)

// Service implements the Service interface
type Service struct {
	factFinderAgent agents.AgentInterface
	repo            *Repository
}

func NewService(
	serviceRepo *Repository,
	completionService ai.CompletionServiceInterface,
) *Service {

	return &Service{
		repo:            serviceRepo,
		factFinderAgent: fact_agent.NewFactAgent(completionService),
	}
}

func (s *Service) FindFacts(messageBody string) (*[]fact_agent.FactAgentFact, error) {

	response, err := s.factFinderAgent.Do(messageBody)

	if err != nil {

		return nil, err
	}

	facts, err := fact_agent.ParseResponse(response)

	// Re-wrap this to indentify this as a parsing error
	if err != nil {

		return nil, fmt.Errorf("could not parse response from OpenAI: %s", response)
	}

	return facts, nil
}

func (s *Service) CreateFact(fact *models.Fact) error {
	return s.repo.Create(fact)
}

func (s *Service) UpdateFact(fact *models.Fact) error {
	return s.repo.Update(fact)
}

func (s *Service) DeleteFact(id int64) error {
	return s.repo.Delete(id)
}

func (s *Service) FindFactByID(id int64) (*models.Fact, error) {
	return s.repo.FindByID(id)
}

func (s *Service) FindFactsByUserID(userID int64) ([]*models.Fact, error) {
	return s.repo.FindByUserID(userID)
}
