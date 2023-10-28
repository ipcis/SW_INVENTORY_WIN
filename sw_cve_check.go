package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
)

type SoftwareInfo struct {
	Name        string `json:"Name"`
	Version     string `json:"Version"`
	InstallDate string `json:"InstallDate"`
}

func main() {
	cmd := exec.Command("powershell", "Get-WmiObject -Class Win32_Product | Select-Object Name, Version, InstallDate")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Fehler beim Ausführen des PowerShell-Befehls: %v\n", err)
		return
	}
	fmt.Printf(string(output))

	// Verwende reguläre Ausdrücke, um Informationen aus der PowerShell-Ausgabe zu extrahieren
	re := regexp.MustCompile(`(?P<Name>.+?)\s+(?P<Version>\d+\.\d+\.\d+)\s+(?P<InstallDate>\d+)`)

	matches := re.FindAllStringSubmatch(string(output), -1)

	var softwareList []SoftwareInfo
	for _, match := range matches {
		info := SoftwareInfo{
			Name:        match[1],
			Version:     match[2],
			InstallDate: match[3],
		}
		softwareList = append(softwareList, info)
	}

	// Konvertiere die Softwareliste in JSON
	jsonData, err := json.MarshalIndent(softwareList, "", "  ")
	if err != nil {
		fmt.Printf("Fehler beim Erstellen von JSON: %v\n", err)
		return
	}

	// Speichere JSON-Daten in einer Textdatei
	err = ioutil.WriteFile("installed_software.json", jsonData, 0644)
	if err != nil {
		fmt.Printf("Fehler beim Speichern der JSON-Daten: %v\n", err)
		return
	}

	fmt.Println("Softwareinformationen wurden in 'installed_software.json' gespeichert.")
}
