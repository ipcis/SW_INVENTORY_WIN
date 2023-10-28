# SW_INVENTORY_WIN
Different code to extract SW inventory informations

```
   31  unzip nvdcve-1.1-2023.json.zip
   32  ls
   33  head nvdcve-1.1-*.json
   34  cat nvdcve-1.1-*.json
   35  grep cpe23Uri nvdcve-1.1-*.json
   36  grep cpe23Uri nvdcve-1.1-*.json | cut -d ':'
   37  grep cpe23Uri nvdcve-1.1-*.json | cut -d ':' -f1
   38  grep cpe23Uri nvdcve-1.1-*.json | cut -d ':' -f2
   39  grep cpe23Uri nvdcve-1.1-*.json | cut -d ':' -f3
   40  grep cpe23Uri nvdcve-1.1-*.json | cut -d ':' -f4
   41  grep cpe23Uri nvdcve-1.1-*.json | cut -d ':' -f5
   42  grep cpe23Uri nvdcve-1.1-*.json | cut -d ':' -f5 | sort -u
   43  grep cpe23Uri nvdcve-1.1-*.json | cut -d ':' -f5 | sort -u > uniq_vendor_2022.txt
   44  cat uniq_vendor_2022.txt
   45  cat uniq_vendor_2022.txt | grep microsoft
   46  cat uniq_vendor_2022.txt | grep juniper
   47  cat uniq_vendor_2022.txt | grep junipdfs
   48  cat uniq_vendor_2022.txt | grep tightvnc
   49  cat uniq_vendor_2022.txt | wc -l
   50  cat uniq_vendor_2022.txt | grep wireguard
   51  cat uniq_vendor_2022.txt | grep vs_codecover
```


folgendes fehlt noch:

- einige programme sind im wmic nicht drin: koennte mit exe search nachgebessert werden
- oben in der json sollte der hostname, os und os build number rein
- os buildnumber sollte cve check durchgefuehrt werden und eol check
- cve output sollte als json erfolgen
- 
