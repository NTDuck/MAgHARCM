package print

import (
	"errors"
	"fmt"
	"io"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

func Events(events *adk.AsyncIterator[*adk.AgentEvent]) {
	for {
		event, ok := events.Next()
		if !ok {
			break
		}
		if event == nil {
			continue
		}
		if event.Err != nil {
			fmt.Println("Error:", event.Err)
			continue
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
						fmt.Println("[ERROR]", err)
					}
					break
				}
				Msg(msg)
			}
			messageOutput.MessageStream.Close()
		} else if messageOutput.Message != nil {
			Msg(messageOutput.Message)
		}
	}
	fmt.Println()
}

func Msg(msg *schema.Message) {
	if msg == nil {
		return
	}
	if msg.ReasoningContent != "" {
		fmt.Print(msg.ReasoningContent)
	}
	if msg.Content != "" {
		fmt.Print(msg.Content)
	}
}
