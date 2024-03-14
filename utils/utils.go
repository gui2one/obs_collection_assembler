package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadJsonToStringMap(json_path string) map[string]interface{} {
	bytes, err := os.ReadFile(json_path)
	if err != nil {
		return nil
	}

	var json_map map[string]interface{}
	err = json.Unmarshal(bytes, &json_map)
	if err != nil {
		fmt.Println("Problem UnMarshaling JSON file")
		return nil
	}

	return json_map

}

func ReadJsonToString(json_path string) (string, error) {
	bytes, err := os.ReadFile(json_path)
	if err != nil {
		return "", err
	}

	var text string = string(bytes)

	return text, nil

}

func Hey() {
	fmt.Println("hey there ! What's up ?")
}
