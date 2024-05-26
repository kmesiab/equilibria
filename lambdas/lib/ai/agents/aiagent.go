package agents

type AgentTool string

const (
	AgentToolSearch AgentTool = "search_tool"
)

type AgentInterface interface {
	Do(input string) (string, error)
}

// AIAgent is the base struct for all AI agents.
// It defines the common fields such as role, goal, verbosity, memory
// usage, backstory, delegation allowance, and tools. This struct serves
// as the foundation for our agents, encapsulating the backstory and
// tools used, along with other configurations. It is part of the
// AgentInterface, allowing us to call Do() with an input and get an
// output completion without needing to know the internals of how the
// agent works or its backstory, etc.
type AIAgent struct {
	Role            string
	Goal            string
	Verbose         bool
	Memory          bool
	Backstory       string
	AllowDelegation bool
	Tools           []AgentTool
}
