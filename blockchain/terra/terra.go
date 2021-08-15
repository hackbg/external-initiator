// might cover all cosmos/tendermint chains later

package terra

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/external-initiator/store"
	"github.com/smartcontractkit/external-initiator/subscriber"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tidwall/gjson"
)

const Name = "terra"

type TerraParams struct {
	ContractAddress string `json:"contract_address"`
	AccountAddress  string `json:"account_address"`
	// FcdUrl          string
}

type manager struct {
	endpointName    string
	contractAddress string
	accountAddress  string
	// fcdUrl          string
	subscriber subscriber.ISubscriber
}

func createManager(sub store.Subscription) (*manager, error) {
	conn, err := subscriber.NewSubscriber(sub.Endpoint)
	if err != nil {
		return nil, err
	}

	return &manager{
		endpointName:    sub.EndpointName,
		contractAddress: sub.Terra.ContractAddress,
		accountAddress:  sub.Terra.AccountAddress,
		subscriber:      conn,
	}, nil
}

func (tm *manager) Stop() {
	// TODO!
}

func (tm *manager) query(ctx context.Context, address, query string, t interface{}) error {
	// TODO! potentially use Tendermint http client
	url := fmt.Sprintf("%s/wasm/contracts/%s/store?query_msg=%s", os.Getenv("TERRA_URL"), address, query)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var decoded map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return err
	}

	if err := json.Unmarshal(decoded["result"], &t); err != nil {
		return err
	}

	return nil
}

func (tm *manager) subscribe(ctx context.Context, queryFilter string, handler func(event EventRecords)) error {
	responses := make(chan json.RawMessage)
	filter := []string{queryFilter}
	params, err := json.Marshal(filter)
	if err != nil {
		return err
	}

	err = tm.subscriber.Subscribe(ctx, "subscribe", "unsubscribe", params, responses)
	if err != nil {
		return err
	}

	go func() {
		for {
			resp, ok := <-responses
			if !ok {
				return
			}

			events, err := extractEvents(resp)
			if err != nil {
				logger.Error(err)
				continue
			}
			eventRecords, err := parseEvents(events)
			if err != nil {
				logger.Error(err)
				continue
			}
			if eventRecords != nil {
				handler(*eventRecords)
			}
		}
	}()

	return nil
}

func extractEvents(data json.RawMessage) ([]types.Event, error) {
	value := gjson.Get(string(data), "data.value.TxResult.result.events") // TODO! this parsing should be improved

	var events []types.Event
	err := json.Unmarshal([]byte(value.Raw), &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func parseEvents(events []types.Event) (*EventRecords, error) {
	var eventRecords EventRecords

	for _, event := range events {
		switch event.Type {

		case "wasm-new_round":
			roundIdStr, err := getAttributeValue(event.Attributes, "round_id")
			if err != nil {
				return nil, err
			}
			roundId, err := strconv.Atoi(roundIdStr)
			if err != nil {
				return nil, err
			}
			startedBy, err := getAttributeValue(event.Attributes, "started_by")
			if err != nil {
				return nil, err
			}
			var startedAt uint64
			startedAtStr, err := getAttributeValue(event.Attributes, "started_at")
			if err == nil {
				value, err := strconv.Atoi(startedAtStr)
				if err != nil {
					return nil, err
				}
				startedAt = uint64(value)
			}
			eventRecords.NewRound = append(eventRecords.NewRound, EventNewRound{
				RoundId:   uint32(roundId),
				StartedBy: Addr(startedBy),
				StartedAt: startedAt,
			})

		case "wasm-submission_received":
			submissionStr, err := getAttributeValue(event.Attributes, "submission")
			if err != nil {
				return nil, err
			}
			submission := new(big.Int)
			submission, _ = submission.SetString(submissionStr, 10)

			roundIdStr, err := getAttributeValue(event.Attributes, "round_id")
			if err != nil {
				return nil, err
			}
			roundId, err := strconv.Atoi(roundIdStr)
			if err != nil {
				return nil, err
			}
			oracle, err := getAttributeValue(event.Attributes, "oracle")
			if err != nil {
				return nil, err
			}

			eventRecords.SubmissionReceived = append(eventRecords.SubmissionReceived, EventSubmissionReceived{
				Oracle:     Addr(oracle),
				Submission: Value{*submission},
				RoundId:    uint32(roundId),
			})

		case "wasm-answer_updated":
			roundIdStr, err := getAttributeValue(event.Attributes, "round_id")
			if err != nil {
				return nil, err
			}
			roundId, err := strconv.Atoi(roundIdStr)
			if err != nil {
				return nil, err
			}
			valueStr, err := getAttributeValue(event.Attributes, "current")
			if err != nil {
				return nil, err
			}
			value := new(big.Int)
			value, _ = value.SetString(valueStr, 10)

			eventRecords.AnswerUpdated = append(eventRecords.AnswerUpdated, EventAnswerUpdated{
				Value:   Value{*value},
				RoundId: uint32(roundId),
			})

		case "wasm-oracle_permissions_updated":
			addedStr, err := getAttributeValue(event.Attributes, "added")
			if err != nil {
				return nil, err
			}
			var added []string
			err = json.Unmarshal([]byte(addedStr), &added)
			if err != nil {
				return nil, err
			}
			for _, oracle := range added {
				eventRecords.OraclePermissionsUpdated = append(eventRecords.OraclePermissionsUpdated, EventOraclePermissionsUpdated{
					Oracle: Addr(oracle),
					Bool:   true,
				})
			}

			removedStr, err := getAttributeValue(event.Attributes, "removed")
			if err != nil {
				return nil, err
			}
			var removed []string
			err = json.Unmarshal([]byte(removedStr), &removed)
			if err != nil {
				return nil, err
			}
			for _, oracle := range removed {
				eventRecords.OraclePermissionsUpdated = append(eventRecords.OraclePermissionsUpdated, EventOraclePermissionsUpdated{
					Oracle: Addr(oracle),
					Bool:   false,
				})
			}
		}
	}

	return &eventRecords, nil
}

func getAttributeValue(attributes []types.EventAttribute, attributeKey string) (string, error) {
	for _, attr := range attributes {
		if string(attr.Key) == attributeKey {
			return string(attr.Value), nil
		}
	}

	return "", fmt.Errorf("attribute key %s does not exist", attributeKey)
}
