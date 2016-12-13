package util

import (
	"encoding/json"

	trello "github.com/VojtechVitek/go-trello"
)

func GetBoardLabelsIds(client *trello.Client, boardID string) (map[string]string, error) {
	type labels struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	type labelsList []labels

	data, err := client.Get("/boards/" + boardID + "/labels")
	if err != nil {
		return nil, err
	}
	var response labelsList
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	result := map[string]string{}
	for _, l := range response {
		result[l.Name] = l.ID
	}
	return result, nil
}
