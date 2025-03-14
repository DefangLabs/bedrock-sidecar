package bedrock

import (
	"encoding/json"
	"fmt"
	"os"
)

type ModelMap map[string]string

func NewModelMap() (ModelMap, error) {
	envVarName := "MODEL_NAME_MAP"
	modelNameMap := map[string]string{}
	// read MODEL_NAME_MAP as json from env
	err := json.Unmarshal([]byte(os.Getenv(envVarName)), &modelNameMap)
	if err != nil {
		return nil, fmt.Errorf("%w: unable to unmarshal %s", err, envVarName)
	}

	return modelNameMap, nil
}

func (m ModelMap) BedrockModelID(openAIModel string) string {
	mappedModelID, ok := m[openAIModel]
	if !ok {
		return openAIModel
	}
	return mappedModelID
}
