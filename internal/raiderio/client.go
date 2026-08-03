package raiderio

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetCharacter(name string, realm string, region string) (*Character, error) {

	url := fmt.Sprintf(
		"https://raider.io/api/v1/characters/profile?region=%s&realm=%s&name=%s&fields=mythic_plus_scores_by_season:current,gear",
		region,
		realm,
		name,
	)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var character Character

	err = json.NewDecoder(response.Body).Decode(&character)
	if err != nil {
		return nil, err
	}

	return &character, nil
}