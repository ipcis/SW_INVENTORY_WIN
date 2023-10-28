from cvelookup import CVELookup

# Ersetzen Sie 'SoftwareName' durch den Namen der installierten Software
software_name = "SoftwareName"

# Erstellen Sie eine Instanz des CVELookup
cve_lookup = CVELookup()

# Abrufen von CVE-Informationen
cve_info = cve_lookup.lookup(software_name)

if cve_info:
    print(f"CVE-Informationen für {software_name}:")
    for cve in cve_info:
        print(f"CVE-ID: {cve['id']}, Beschreibung: {cve['summary']}")
else:
    print(f"Keine CVE-Informationen für {software_name} gefunden.")
