package print

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

func Events(events *adk.AsyncIterator[*adk.AgentEvent]) {
	var currentAgent string
	var inReasoning bool

	for {
		event, ok := events.Next()
		if !ok {
			break
		}
		if event == nil {
			continue
		}
		if event.Err != nil {
			fmt.Printf("\n\033[31m[ERROR]\033[0m %v\n", event.Err)
			continue
		}

		if event.AgentName != "" && event.AgentName != currentAgent {
			currentAgent = event.AgentName
			fmt.Printf("\n\033[1;34m=== Agent: %s ===\033[0m\n", currentAgent)
		}

		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}

		messageOutput := event.Output.MessageOutput
		if messageOutput.IsStreaming && messageOutput.MessageStream != nil {
			for {
				msg, err := messageOutput.MessageStream.Recv()
				if err != nil {
					if !errors.Is(err, io.EOF) {
						fmt.Printf("\n\033[31m[STREAM ERROR]\033[0m %v\n", err)
					}
					break
				}
				streamMsg(msg, &inReasoning)
			}
			messageOutput.MessageStream.Close()
		} else if messageOutput.Message != nil {
			printStaticMsg(messageOutput.Message)
		}
	}
	if inReasoning {
		fmt.Print("\033[0m\n")
	}
	fmt.Println()
}

func streamMsg(msg *schema.Message, inReasoning *bool) {
	if msg == nil {
		return
	}

	if msg.ReasoningContent != "" {
		if !*inReasoning {
			*inReasoning = true
			fmt.Print("\033[2m\033[36m[Thinking] ")
		}
		fmt.Print(msg.ReasoningContent)
	}

	if msg.Content != "" {
		if *inReasoning {
			*inReasoning = false
			fmt.Print("\033[0m\n")
		}
		fmt.Print(msg.Content)
	}

	if len(msg.ToolCalls) > 0 {
		if *inReasoning {
			*inReasoning = false
			fmt.Print("\033[0m\n")
		}
		for _, tc := range msg.ToolCalls {
			fmt.Printf("\n\033[1;33m[Tool Call]\033[0m \033[33m%s\033[0m(%s)\n", tc.Function.Name, tc.Function.Arguments)
		}
	}
}

func printStaticMsg(msg *schema.Message) {
	if msg == nil {
		return
	}
	if msg.Role == schema.Tool {
		preview := msg.Content
		if len(preview) > 300 {
			preview = preview[:300] + "... (truncated)"
		}
		fmt.Printf("\033[1;32m[Tool Result %s]\033[0m %s\n", msg.ToolCallID, strings.TrimSpace(preview))
		return
	}
	if msg.ReasoningContent != "" {
		fmt.Printf("\033[2m\033[36m[Thinking]\n%s\033[0m\n", msg.ReasoningContent)
	}
	if len(msg.ToolCalls) > 0 {
		for _, tc := range msg.ToolCalls {
			var formattedArgs string
			var js any
			if json.Unmarshal([]byte(tc.Function.Arguments), &js) == nil {
				formattedBytes, _ := json.Marshal(js)
				formattedArgs = string(formattedBytes)
			} else {
				formattedArgs = tc.Function.Arguments
			}
			fmt.Printf("\033[1;33m[Tool Call]\033[0m \033[33m%s\033[0m(%s)\n", tc.Function.Name, formattedArgs)
		}
	}
	if msg.Content != "" {
		fmt.Print(msg.Content)
	}
}

func Msg(msg *schema.Message) {
	printStaticMsg(msg)
}
