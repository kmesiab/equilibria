package fact_agent

import "fmt"

const backstory = `
Analyze the following text message and identify any facts relevant to mental health management.  

Not all messages will include facts worth remembering.

If there are important facts, condense each into a summary while maintaining as much relevant context as possible. If a message is instructive, include it. 

If there are no important facts, do not include it. Additionally, provide a section explaining why each fact is relevant to the 
therapist. 

Respond in JSON format as follows:

%s
`

func getBackStory() string {
	return fmt.Sprintf(backstory, getExampleJSONResponseText())
}

func getExampleJSONResponseText() string {
	return fmt.Sprintf(
		jsonTextBlockTemplate,
		exampleResponseText,
		exampleNoResponseJson,
		fullExampleResponseText,
	)
}

const jsonTextBlockTemplate = "```json\n%s\n```" +
	`If there are no facts worth remembering, respond with:` +
	"\n```json\n%s\n```" + // exampleNoResponseJson

	`**Example: **
	Text Message: "I just started a new job and moved to a new city. I'm feeling 
	overwhelmed by the changes, but I'm trying to stay positive. Also, I turned 19 
	last month!"

	JSON Response:` +
	"\n```json\n%s\n```" // fullExampleResponseText

const exampleNoResponseJson = `[]`
const exampleResponseText = `
[
	{
		"fact": "Summarized fact",
		"reasoning": "Explanation of why this is relevant to the therapist"
	},
	{
		"fact": "Summarized fact",
		"reasoning": "Explanation of why this is relevant to the therapist"
	}
]
`
const fullExampleResponseText = `
[
	{
		"fact": "Started new job",
		"reasoning": "Major life events can be a source of new stressors"
	},
	{
		"fact": "Moved to a new city",
		"reasoning": "Relocation can impact mental health due to changes in environment and support systems"
	},
	{
		"fact": "Age 19",
		"reasoning": "Age is relevant to the type of mental health care"
	}
]`
