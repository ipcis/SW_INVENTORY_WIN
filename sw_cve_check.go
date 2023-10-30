package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

//go:embed vendors.txt
var vendorData embed.FS

//go:embed products.txt
var productData embed.FS

type SoftwareInfo struct {
	DisplayName    string `json:"DisplayName"`
	DisplayVersion string `json:"DisplayVersion"`
	Publisher      string `json:"Publisher"`
	InstallDate    string `json:"InstallDate"`
}

type CVE struct {
	CVEID       string `json:"cve"`
	Published   string `json:"published"`
	Modified    string `json:"lastModified"`
	Description string `json:"descriptions"`
}

func main() {
	/*
		vendorData, err := vendorData.ReadFile("vendors.txt")
		if err != nil {
			fmt.Printf("Fehler beim Lesen der Datei 'vendors.txt': %v\n", err)
			return
		}

		productData, err := productData.ReadFile("products.txt")
		if err != nil {
			fmt.Printf("Fehler beim Lesen der Datei 'products.txt': %v\n", err)
			return
		}
	*/

	/*
		// Erstelle eine Map, die Software-Bezeichnungen auf CPE-Strings abbildet
		softwareMap := map[string]string{
			"vlc media player":         "cpe:2.3:a:videolan:vlc_media_player",
			"notepad\\+\\+":            "cpe:2.3:a:notepad-plus-plus:notepad\\+\\+",
			"7-zip":                    "cpe:2.3:a:7-zip:7-zip",
			"tightvnc":                 "cpe:2.3:a:tightvnc:tightvnc",
			"python.*core interpreter": "cpe:2.3:a:python:python",
			"oracle vm virtualbox": 	"cpe:2.3:a:python:python",
			// Füge hier weitere Software-Bezeichnungen und CPE-Strings hinzu
		}*/

	softwareMap := map[string]string{
		"vlc media player": "cpe:2.3:a:videolan:vlc_media_player",
		//"notepad\\+\\+":            "cpe:2.3:a:notepad%2B%2B:notepad%2B%2B", matched nicht auf cpe kein plan warum https://services.nvd.nist.gov/rest/json/cves/2.0?keywordSearch=notepad%2B%2B
		"7-zip":                    "cpe:2.3:a:7-zip:7-zip",
		"tightvnc":                 "cpe:2.3:a:tightvnc:tightvnc",
		"python.*core interpreter": "cpe:2.3:a:python:python",
		"oracle vm virtualbox":     "cpe:2.3:a:oracle:vm_virtualbox",
		"google chrome":            "cpe:2.3:a:google:chrome",
		//"Microsoft Office (.*?) (\\d{4})": "cpe:2.3:a:microsoft:office",
		"adobe acrobat reader":     "cpe:2.3:a:adobe:acrobat_reader",
		"skype":                    "cpe:2.3:a:skype:skype",
		"adobe photoshop":          "cpe:2.3:a:adobe:photoshop",
		"dropbox":                  "cpe:2.3:a:dropbox:dropbox",
		"teamviewer":               "cpe:2.3:a:teamviewer:teamviewer",
		"spotify":                  "cpe:2.3:a:spotify:spotify",
		"itunes":                   "cpe:2.3:a:apple:itunes",
		"windows media player":     "cpe:2.3:a:microsoft:windows_media_player",
		"microsoft edge":           "cpe:2.3:a:microsoft:edge",
		"discord":                  "cpe:2.3:a:discord:discord",
		"slack":                    "cpe:2.3:a:slack:slack",
		"microsoft teams":          "cpe:2.3:a:microsoft:teams",
		"whatsapp desktop":         "cpe:2.3:a:whatsapp:whatsapp_desktop",
		"evernote":                 "cpe:2.3:a:evernote:evernote",
		"obs studio":               "cpe:2.3:a:obs:obs_studio",
		"quicktime player":         "cpe:2.3:a:apple:quicktime_player",
		"filezilla":                "cpe:2.3:a:filezilla:filezilla",
		"onedrive":                 "cpe:2.3:a:microsoft:onedrive",
		"microsoft onenote":        "cpe:2.3:a:microsoft:onenote",
		"thunderbird":              "cpe:2.3:a:mozilla:thunderbird",
		"utorrent":                 "cpe:2.3:a:utorrent:utorrent",
		"bittorrent":               "cpe:2.3:a:bittorrent:bittorrent",
		"inkscape":                 "cpe:2.3:a:inkscape:inkscape",
		"zoom":                     "cpe:2.3:a:zoom:zoom",
		"microsoft paint":          "cpe:2.3:a:microsoft:paint",
		"microsoft visio":          "cpe:2.3:a:microsoft:visio",
		"libreoffice":              "cpe:2.3:a:libreoffice:libreoffice",
		"openoffice":               "cpe:2.3:a:openoffice:openoffice",
		"google drive":             "cpe:2.3:a:google:drive",
		"visual studio":            "cpe:2.3:a:microsoft:visual_studio",
		"adobe premiere pro":       "cpe:2.3:a:adobe:premiere_pro",
		"blender":                  "cpe:2.3:a:blender:blender",
		"steam":                    "cpe:2.3:a:steam:steam",
		"virtualbox":               "cpe:2.3:a:oracle:virtualbox",
		"gimp":                     "cpe:2.3:a:gimp:gimp",
		"audacity":                 "cpe:2.3:a:audacity:audacity",
		"adobe after effects":      "cpe:2.3:a:adobe:after_effects",
		"microsoft excel":          "cpe:2.3:a:microsoft:excel",
		"adobe premiere elements":  "cpe:2.3:a:adobe:premiere_elements",
		"adobe indesign":           "cpe:2.3:a:adobe:indesign",
		"adobe lightroom":          "cpe:2.3:a:adobe:lightroom",
		"autodesk autocad":         "cpe:2.3:a:autodesk:autocad",
		"visual studio code":       "cpe:2.3:a:microsoft:visual_studio_code",
		"sketchup":                 "cpe:2.3:a:sketchup:sketchup",
		"adobe illustrator":        "cpe:2.3:a:adobe:illustrator",
		"corel draw":               "cpe:2.3:a:corel:draw",
		"final cut pro":            "cpe:2.3:a:apple:final_cut_pro",
		"adobe audition":           "cpe:2.3:a:adobe:audition",
		"adobe animate":            "cpe:2.3:a:adobe:animate",
		"google earth":             "cpe:2.3:a:google:earth",
		"winamp":                   "cpe:2.3:a:winamp:winamp",
		"opera browser":            "cpe:2.3:a:opera:browser",
		"camtasia":                 "cpe:2.3:a:camtasia:camtasia",
		"coreldraw":                "cpe:2.3:a:coreldraw:coreldraw",
		"avid pro tools":           "cpe:2.3:a:avid:pro_tools",
		"android studio":           "cpe:2.3:a:android:studio",
		"oracle sql developer":     "cpe:2.3:a:oracle:sql_developer",
		"winrar":                   "cpe:2.3:a:winrar:winrar",
		"firefox":                  "cpe:2.3:a:mozilla:firefox",
		"chrome":                   "cpe:2.3:a:google:chrome",
		"edge":                     "cpe:2.3:a:microsoft:edge",
		"atom":                     "cpe:2.3:a:atom:atom",
		"eclipse":                  "cpe:2.3:a:eclipse:eclipse",
		"intellij idea":            "cpe:2.3:a:intellij:idea",
		"mysql":                    "cpe:2.3:a:mysql:mysql",
		"postgresql":               "cpe:2.3:a:postgresql:postgresql",
		"mongodb":                  "cpe:2.3:a:mongodb:mongodb",
		"adobe acrobat pro":        "cpe:2.3:a:adobe:acrobat_pro",
		"sublime text":             "cpe:2.3:a:sublime:text",
		"picasa":                   "cpe:2.3:a:picasa:picasa",
		"winzip":                   "cpe:2.3:a:winzip:winzip",
		"google sketchup":          "cpe:2.3:a:google:sketchup",
		"autodesk 3ds max":         "cpe:2.3:a:autodesk:3ds_max",
		"visual studio community":  "cpe:2.3:a:microsoft:visual_studio_community",
		"microsoft powerpoint":     "cpe:2.3:a:microsoft:powerpoint",
		"microsoft outlook":        "cpe:2.3:a:microsoft:outlook",
		"autodesk maya":            "cpe:2.3:a:autodesk:maya",
		"avid media composer":      "cpe:2.3:a:avid:media_composer",
		"cinema 4d":                "cpe:2.3:a:maxon:cinema_4d",
		"autodesk inventor":        "cpe:2.3:a:autodesk:inventor",
		"visual studio enterprise": "cpe:2.3:a:microsoft:visual_studio_enterprise",
		"adobe dreamweaver":        "cpe:2.3:a:adobe:dreamweaver",
		"lightworks":               "cpe:2.3:a:lightworks:lightworks",
		"autodesk revit":           "cpe:2.3:a:autodesk:revit",
		"unity":                    "cpe:2.3:a:unity:unity",
		"adobe premiere rush":      "cpe:2.3:a:adobe:premiere_rush",
		"autodesk autocad lt":      "cpe:2.3:a:autodesk:autocad_lt",
		"mysql workbench":          "cpe:2.3:a:mysql:workbench",
		"eclipse ide":              "cpe:2.3:a:eclipse:ide",
		"adobe captivate":          "cpe:2.3:a:adobe:captivate",
		"corel painter":            "cpe:2.3:a:corel:painter",
		"microsoft word":           "cpe:2.3:a:microsoft:word",
	}

	// PowerShell-Befehl
	powershellCmd := "Get-ItemProperty HKLM:\\Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\* | Select-Object DisplayName, DisplayVersion, Publisher, InstallDate | ConvertTo-Json"

	// PowerShell-Befehl ausführen
	cmd := exec.Command("powershell", "-Command", powershellCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Fehler beim Ausführen des PowerShell-Befehls:", err)
		return
	}

	// JSON-Ausgabe parsen
	var softwareInfo []SoftwareInfo
	err = json.Unmarshal(out.Bytes(), &softwareInfo)
	if err != nil {
		fmt.Println("Fehler beim Parsen der JSON-Ausgabe:", err)
		return
	}

	// Ergebnis ausgeben
	for _, info := range softwareInfo {
		fmt.Println("DisplayName:", info.DisplayName)
		fmt.Println("DisplayVersion:", info.DisplayVersion)
		fmt.Println("Publisher:", info.Publisher)

		//fmt.Println("Publisher to lower:", getFirstWordLower(info.Publisher))

		// Suche nach Übereinstimmungen und erhalte die CPE-Strings in einer Map
		cpeMap := findCPE(info.DisplayName, softwareMap)

		cveInfo := ""

		// Gib die gefundenen CPE-Strings aus
		for software, cpe := range cpeMap {
			cveInfo = ""
			fmt.Printf("----------------------------------------\n")
			fmt.Printf("Software: %s\nCPE: %s:%s\n", software, cpe, info.DisplayVersion)
			fmt.Printf("----------------------------------------\n")
			cveInfo = fetchCVEInfo(cpe, info.DisplayVersion)
			fmt.Printf(cveInfo)

		}

		fmt.Println("InstallDate:", info.InstallDate)
		fmt.Println()
	}

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

func isNameInProductList(name string, productData []byte) bool {
	lines := strings.Split(string(productData), "\n")
	for _, line := range lines {
		if strings.ToLower(line) == name {
			return true
		}
	}
	return false
}

func fetchCVEInfo(name, version string) string {
	//cveURL := "https://services.nvd.nist.gov/rest/json/cves/2.0?keywordSearch=" + name + "%20" + version

	cveURL := "https://services.nvd.nist.gov/rest/json/cves/2.0?cpeName=" + name + ":" + version + "&isVulnerable"
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

func getFirstWordLower(s string) string {
	parts := strings.Fields(s)
	if len(parts) > 0 {
		return strings.ToLower(parts[0])
	}
	return "" // or any other default value
}

func findCPE(text string, softwareMap map[string]string) map[string]string {
	cpeMap := make(map[string]string)

	for software, cpe := range softwareMap {
		// Escape des Suchstrings, um Sonderzeichen korrekt zu behandeln
		//escapedSoftware := regexp.QuoteMeta(software)

		//fmt.Println("----- SW:" + strings.ToLower(escapedSoftware))
		//fmt.Println("----- SW:" + strings.ToLower(text))

		//re := regexp.MustCompile(strings.ToLower(escapedSoftware))
		re := regexp.MustCompile(strings.ToLower(software))
		if re.MatchString(strings.ToLower(text)) {
			cpeMap[software] = cpe
		}
	}

	return cpeMap
}
