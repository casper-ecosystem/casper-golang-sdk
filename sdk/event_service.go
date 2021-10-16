package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/r3labs/sse/v2"
	"strings"
	"time"
)

type NodeSseChannelType int64
type NodeSseEventType int64

const (
	NodeSseChannelTypeMain NodeSseChannelType = iota
	NodeSseChannelTypeDeploys
	NodeSseChannelTypeSigs
)

func (n NodeSseChannelType) Endpoint() string {
	switch n {
	case NodeSseChannelTypeMain:
		return "/events/main"
	case NodeSseChannelTypeDeploys:
		return "/events/deploys"
	case NodeSseChannelTypeSigs:
		return "/events/sigs"
	default:
		return "unknown"
	}
}

const (
	NodeSseEventTypeAll NodeSseEventType = iota
	NodeSseEventTypeApiVersion
	NodeSseEventTypeBlockAdded
	NodeSseEventTypeDeployAccepted
	NodeSseEventTypeDeployProcessed
	NodeSseEventTypeFault
	NodeSseEventTypeFinalitySignature
	NodeSseEventTypeStep
)

func (n NodeSseEventType) String() string {
	switch n {
	case NodeSseEventTypeAll:
		return ""
	case NodeSseEventTypeApiVersion:
		return "ApiVersion"
	case NodeSseEventTypeBlockAdded:
		return "BlockAdded"
	case NodeSseEventTypeDeployAccepted:
		return "DeployAccepted"
	case NodeSseEventTypeDeployProcessed:
		return "DeployProcessed"
	case NodeSseEventTypeFault:
		return "Fault"
	case NodeSseEventTypeFinalitySignature:
		return "FinalitySignature"
	case NodeSseEventTypeStep:
		return "Step"
	default:
		return "unknown"
	}
}

type EventService struct {
	url string
}

// NewEventService create a new event service with specified URL.
// URL should be without any path for example http://159.65.118.250:9999.
func NewEventService(url string) *EventService {
	return &EventService{
		url: url,
	}
}

// AwaitEvents returns a channel which streams events specified by parameters.
func (e *EventService) AwaitEvents(channel NodeSseChannelType, eventType NodeSseEventType) (chan *EventResponse, error) {
	client := sse.NewClientWithBufferSize(fmt.Sprintf("%s%s", e.url, channel.Endpoint()), 500000)
	tempChan := make(chan *sse.Event)
	eventChan := make(chan *EventResponse)
	err := client.SubscribeChanRaw(tempChan)
	if err != nil {
		return nil, err
	}
	go func() {
		for event := range tempChan {
			if len(event.Data) == 0 {
				continue
			}
			data := string(event.Data)
			if !strings.Contains(data, eventType.String()) {
				continue
			}
			var resp EventResponse
			err := json.Unmarshal(event.Data, &resp)
			if err != nil {
				continue
			}
			eventChan <- &resp
		}
		close(eventChan)
	}()
	return eventChan, nil
}

// AwaitDeploy waits for a deploy with deployHash and returns it.
func (e *EventService) AwaitDeploy(deployHash string) (DeployProcessedEvent, error) {
	eventChan, err := e.AwaitEvents(NodeSseChannelTypeMain, NodeSseEventTypeDeployProcessed)
	if err != nil {
		return DeployProcessedEvent{}, err
	}
	for event := range eventChan {
		if event.DeployProcessed != nil && event.DeployProcessed.DeployHash == deployHash {
			return *event.DeployProcessed, nil
		}
	}
	return DeployProcessedEvent{}, errors.New("could not await deploy")
}

// AwaitNBlocks waits for n blocks and returns that block.
func (e *EventService) AwaitNBlocks(n int) (BlockAddedEvent, error) {
	eventChan, err := e.AwaitEvents(NodeSseChannelTypeMain, NodeSseEventTypeBlockAdded)
	if err != nil {
		return BlockAddedEvent{}, err
	}
	var event *EventResponse
	for i := 1; i <= n; i++ {
		event = <-eventChan
		if i == n {
			return *event.BlockAdded, nil
		}
	}
	return BlockAddedEvent{}, errors.New("could not await block")
}

// AwaitNEras waits for n eras and returns first block in that era.
func (e *EventService) AwaitNEras(n int) (BlockAddedEvent, error) {
	eventChan, err := e.AwaitEvents(NodeSseChannelTypeMain, NodeSseEventTypeBlockAdded)
	if err != nil {
		return BlockAddedEvent{}, err
	}
	currentEra := 0
	erasPassed := 0
	for event := range eventChan {
		era := event.BlockAdded.Block.Header.EraID
		if era > currentEra {
			currentEra = era
			erasPassed++
		}
		if erasPassed > n {
			return *event.BlockAdded, nil
		}
	}
	return BlockAddedEvent{}, errors.New("could not await era")
}

// AwaitUntilBlockN waits for a block with specified height and returns it.
func (e *EventService) AwaitUntilBlockN(height int) (BlockAddedEvent, error) {
	eventChan, err := e.AwaitEvents(NodeSseChannelTypeMain, NodeSseEventTypeBlockAdded)
	if err != nil {
		return BlockAddedEvent{}, err
	}
	for event := range eventChan {
		if event.BlockAdded != nil && event.BlockAdded.Block.Header.Height == height {
			return *event.BlockAdded, nil
		}
	}
	return BlockAddedEvent{}, errors.New("could not await era")
}

// AwaitUntilEraN waits for an era with specified id and returns the first block in that era.
func (e *EventService) AwaitUntilEraN(eraId int) (BlockAddedEvent, error) {
	eventChan, err := e.AwaitEvents(NodeSseChannelTypeMain, NodeSseEventTypeBlockAdded)
	if err != nil {
		return BlockAddedEvent{}, err
	}
	for event := range eventChan {
		if event.BlockAdded != nil && event.BlockAdded.Block.Header.EraID == eraId {
			return *event.BlockAdded, nil
		}
	}
	return BlockAddedEvent{}, err
}

type EventResponse struct {
	ApiVersion        *string                 `json:"ApiVersion,omitempty"`
	BlockAdded        *BlockAddedEvent        `json:"BlockAdded,omitempty"`
	DeployAccepted    *DeployAcceptedEvent    `json:"DeployAccepted,omitempty"`
	DeployProcessed   *DeployProcessedEvent   `json:"DeployProcessed,omitempty"`
	Fault             *FaultEvent             `json:"Fault"`
	FinalitySignature *FinalitySignatureEvent `json:"FinalitySignature"`
	Step              *StepEvent              `json:"Step"`
}

type BlockAddedEvent struct {
	BlockHash string        `json:"block_hash"`
	Block     BlockResponse `json:"block"`
}

type DeployAcceptedEvent struct {
	Deploy string `json:"deploy"`
}

type DeployProcessedEvent struct {
	DeployHash      string              `json:"deploy_hash"`
	Account         string              `json:"account"`
	Timestamp       time.Time           `json:"timestamp"`
	TTL             string              `json:"ttl"`
	Dependencies    []string            `json:"dependencies"`
	BlockHash       string              `json:"block_hash"`
	ExecutionResult JsonExecutionResult `json:"execution_result"`
}

type FaultEvent struct {
	EraId     int       `json:"era_id"`
	PublicKey string    `json:"public_key"`
	Timestamp time.Time `json:"timestamp"`
}

type FinalitySignatureEvent struct {
	BlockHash string `json:"block_hash"`
	EraId     int    `json:"era_id"`
	Signatue  string `json:"signatue"`
	PublicKey string `json:"public_key"`
}

type StepEvent struct {
	EraId           int    `json:"era_id"`
	ExecutionEffect string `json:"execution_effect"`
}
