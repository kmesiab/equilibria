package facts

import (
	"github.com/kmesiab/equilibria/lambdas/lib/ai/agents/fact_agent"
	"github.com/kmesiab/equilibria/lambdas/models"
)

// ServiceInterface defines the interface for finding facts and interacting with the repository
type ServiceInterface interface {
	FindFacts(messageBody string) (*[]fact_agent.FactAgentFact, error)
	CreateFact(fact *models.Fact) error
	UpdateFact(fact *models.Fact) error
	DeleteFact(id int64) error
	FindFactByID(id int64) (*models.Fact, error)
	FindFactsByUserID(userID int64) ([]*models.Fact, error)
}
