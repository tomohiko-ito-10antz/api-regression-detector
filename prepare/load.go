package prepare

import (
	"encoding/json"
	"os"

	"github.com/Jumpaku/api-regression-detector/db"
)

func Load(file *os.File) (tables db.Tables, err error) {
	decoder := json.NewDecoder(file)
	decoder.UseNumber()
	if err := decoder.Decode(&tables); err != nil {
		return nil, err
	}
	return tables, nil
}
