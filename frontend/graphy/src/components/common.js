const baseUrl = "http://localhost:9090"

const Urls = {
  "folder": baseUrl + "/csv_folder", 
  "csv": baseUrl + "/csv" 
}

export function getUrl(key) {
  return Urls[key]; 
}

export function getRequestname() {
  return ""
}
