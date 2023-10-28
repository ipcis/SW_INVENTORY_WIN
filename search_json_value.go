package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	// Pfad zur JSON-Datei
	jsonFilePath := "./nvdcpematch-1.0.json/nvdcpematch-1.0.json"

	// String, nach dem gesucht werden soll
	searchString := "cpe:2.3:a:tightvnc:tightvnc:*:*:*:*:*:*:*:*" // Ersetze "Suchtext" durch den gewünschten Suchbegriff

	// JSON-Datei öffnen
	file, err := os.Open(jsonFilePath)
	if err != nil {
		fmt.Printf("Fehler beim Öffnen der JSON-Datei: %v\n", err)
		return
	}
	defer file.Close()

	// JSON-Datei in einen Byte-Slice einlesen
	byteValue, _ := ioutil.ReadAll(file)

	// JSON-Struktur dekodieren
	var data map[string]interface{}
	if err := json.Unmarshal(byteValue, &data); err != nil {
		fmt.Printf("Fehler beim Dekodieren der JSON-Daten: %v\n", err)
		return
	}

	// Den gewünschten String in den JSON-Daten suchen
	if findStringInJSON(data, searchString) {
		fmt.Printf("Der Suchbegriff '%s' wurde in der JSON-Datei gefunden.\n", searchString)
	} else {
		fmt.Printf("Der Suchbegriff '%s' wurde nicht in der JSON-Datei gefunden.\n", searchString)
	}
}

// Rekursive Funktion zum Suchen eines Strings in JSON-Daten
func findStringInJSON(data interface{}, searchString string) bool {
	switch val := data.(type) {
	case string:
		return strings.Contains(val, searchString)
	case map[string]interface{}:
		for _, v := range val {
			if findStringInJSON(v, searchString) {
				return true
			}
		}
	case []interface{}:
		for _, v := range val {
			if findStringInJSON(v, searchString) {
				return true
			}
		}
	}
	return false
}
