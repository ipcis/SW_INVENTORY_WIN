package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ExeFileInformation struct {
	ProductVersion string `json:"product_version"`
	LastModified   string `json:"last_modified"`
	ProductName    string `json:"product_name"`
}

func findExeFilesWithDetails(dirs []string) map[string]ExeFileInformation {
	exeFiles := make(map[string]ExeFileInformation)

	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(info.Name(), ".exe") {
				productVersion, lastModified, productName, err := getFileDetails(path)
				if err == nil {
					exeInfo := ExeFileInformation{
						ProductVersion: productVersion,
						LastModified:   lastModified,
						ProductName:    productName,
					}
					exeFiles[path] = exeInfo
				}
			}

			return nil
		})

		if err != nil {
			fmt.Printf("Fehler bei der Verzeichnisdurchsuchung für %s: %v\n", dir, err)
		}
	}

	return exeFiles
}

func getFileDetails(filePath string) (string, string, string, error) {
	cmd1 := exec.Command("powershell", "(Get-Command '"+filePath+"').FileVersionInfo.ProductVersion")
	output1, err := cmd1.CombinedOutput()
	if err != nil {
		return "", "", "", err
	}
	productVersion := strings.TrimSpace(string(output1))

	cmd2 := exec.Command("powershell", "(Get-Command '"+filePath+"').CreationTime")
	output2, err := cmd2.CombinedOutput()
	if err != nil {
		return "", "", "", err
	}
	lastModified := strings.TrimSpace(string(output2))

	cmd3 := exec.Command("powershell", "(Get-Command '"+filePath+"').FileVersionInfo.ProductName")
	output3, err := cmd3.CombinedOutput()
	if err != nil {
		return "", "", "", err
	}
	productName := strings.TrimSpace(string(output3))

	return productVersion, lastModified, productName, nil
}

func main() {
	/*
		directories := []string{
			"C:\\Program Files",
			"C:\\Program Files (x86)",
			// Füge hier weitere Verzeichnisse hinzu, die du durchsuchen möchtest.
			// Zum Beispiel: "D:\\MyDirectory"
		}*/

	//directories := []string{"C:\\Program Files\\VideoLAN"}
	directories := []string{"C:\\Program Files\\"}

	exeFilesWithDetails := findExeFilesWithDetails(directories)

	jsonData, err := json.MarshalIndent(exeFilesWithDetails, "", "  ")
	if err != nil {
		fmt.Printf("Fehler beim Erstellen von JSON: %v\n", err)
		return
	}

	fmt.Printf("JSON-Daten:\n%s\n", string(jsonData))
}
