package prepare

import (
	"encoding/json"
	"os"
)

func ReadTablesFrom(file *os.File) (tables Tables, err error) {
	decoder := json.NewDecoder(file)
	decoder.UseNumber()
	if err := decoder.Decode(&tables); err != nil {
		return nil, err
	}
	return tables, nil
}

func WriteSqlTo(sql string, file *os.File) (err error) {
	_, err = file.WriteString(sql)
	return err
}
