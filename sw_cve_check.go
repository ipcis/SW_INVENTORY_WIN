package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

type SoftwareInfo struct {
	Name        string `json:"Name"`
	Version     string `json:"Version"`
	InstallDate string `json:"InstallDate"`
}

type CVE struct {
	CVEID       string `json:"cve"`
	Published   string `json:"published"`
	Modified    string `json:"lastModified"`
	Description string `json:"descriptions"`
}

func main() {
	cmd := exec.Command("powershell", "Get-WmiObject -Class Win32_Product | Select-Object Name, Version, InstallDate")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Fehler beim Ausführen des PowerShell-Befehls: %v\n", err)
		return
	}

	re := regexp.MustCompile(`(?P<Name>.+?)\s+(?P<Version>(\d+(\.\d+)*)\s+)\s+(?P<InstallDate>.+)`)
	matches := re.FindAllStringSubmatch(string(output), -1)

	vendorData, err := ioutil.ReadFile("uniq_vendor.txt")
	if err != nil {
		fmt.Printf("Fehler beim Lesen der Datei 'uniq_vendor.txt': %v\n", err)
		return
	}

	var softwareList []SoftwareInfo
	for _, match := range matches {
		info := SoftwareInfo{
			Name:        match[1],
			Version:     strings.TrimSpace(match[2]),
			InstallDate: match[3],
		}
		softwareList = append(softwareList, info)

		nameLowerCase := strings.ToLower(info.Name)
		if isNameInVendorList(nameLowerCase, vendorData) {
			fmt.Printf("Gesuchter Produktname: %s - Match: Ja\n", info.Name)
			//bodystring := fetchCVEInfo(info.Name)
			bodystring := fetchCVEInfo(nameLowerCase, info.Version) // Pass both name and version
			fmt.Println(bodystring)                                 // Gibt den Body-Inhalt auf der Konsole aus
		}
	}

	jsonData, err := json.MarshalIndent(softwareList, "", "  ")
	if err != nil {
		fmt.Printf("Fehler beim Erstellen von JSON: %v\n", err)
		return
	}

	err = ioutil.WriteFile("installed_software.json", jsonData, 0644)
	if err != nil {
		fmt.Printf("Fehler beim Speichern der JSON-Daten: %v\n", err)
		return
	}

	fmt.Println("Softwareinformationen wurden in 'installed_software.json' gespeichert.")
}

func isNameInVendorList(name string, vendorData []byte) bool {
	lines := strings.Split(string(vendorData), "\n")
	for _, line := range lines {
		if strings.ToLower(line) == name {
			return true
		}
	}
	return false
}

func fetchCVEInfo(name, version string) string {
	cveURL := "https://services.nvd.nist.gov/rest/json/cves/2.0?keywordSearch=" + name + "%20" + version
	resp, err := http.Get(cveURL)
	if err != nil {
		// Handle error
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// Den Body-Inhalt in einen Byte-Slice lesen
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Fehlerbehandlung
		fmt.Println(err)
	}

	// JSON-Daten in ein übersichtliches Format umwandeln
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		// Fehlerbehandlung
		fmt.Println(err)
	}

	return prettyJSON.String()
}
