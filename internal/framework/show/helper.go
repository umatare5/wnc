package show

import (
	"encoding/json"
	"fmt"
	"log"
)

func printJson(data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(jsonData))
}
