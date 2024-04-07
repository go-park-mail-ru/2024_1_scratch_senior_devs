package elasticsearch

import "encoding/json"

func processValue(value interface{}) interface{} {
	switch typedValue := value.(type) {
	case map[string]interface{}:
		outputSlice := make([]map[string]interface{}, 0)
		for key, val := range typedValue {
			outputSlice = append(outputSlice, map[string]interface{}{"key": key, "value": processValue(val)})
		}
		return outputSlice
	case []interface{}:
		outputSlice := make([]interface{}, 0)
		for _, val := range typedValue {
			outputSlice = append(outputSlice, processValue(val))
		}
		return outputSlice
	default:
		return value
	}
}

func getElasticData(inputJSON []byte) (interface{}, error) {
	var inputMap map[string]interface{}
	if err := json.Unmarshal(inputJSON, &inputMap); err != nil {
		return []byte{}, err
	}

	return processValue(inputMap), nil
}
