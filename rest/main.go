package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"

)

type Ledger struct {
    ID   string `json:"id"`
    X    string `json:"x"`
    Y    int    `json:"y"`
    Z    int    `json:"z"`
}

type DiskSpace struct {
    ID       string `json:"id"`
    Free     int    `json:"free"`
    Capacity int    `json:"capacity"`
}

logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
    log.Fatal(err)
}
defer logFile.Close()

logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)

logLevel := os.Getenv("REST_SERVER_LOG_LEVEL")
if logLevel == "" {
    logLevel = "ERROR"
}

switch strings.ToUpper(logLevel) {
case "DEBUG":
    logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
case "INFO":
    logger.SetFlags(log.Ldate | log.Ltime)
case "WARNING":
    logger.SetFlags(log.Ldate | log.Ltime)
    logger.SetPrefix("[WARNING] ")
case "ERROR":
    logger.SetFlags(log.Ldate | log.Ltime)
    logger.SetPrefix("[ERROR] ")
default:
    logger.SetFlags(log.Ldate | log.Ltime)
    logger.SetPrefix("[UNKNOWN] ")
}

logger.Println("Server started on port 8000")

logger.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
    http.HandleFunc("/api/local/ledgers", getLedgers)
    http.HandleFunc("/api/local/ledger/", getLedger)
    http.HandleFunc("/api/local/dispace/", getDiskSpace)
    http.HandleFunc("/api/local/", handleUnknownEndpoint)

    logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal(err)
    }
    defer logFile.Close()

    logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)

    logger.Println("Server started on port 8000")

    logger.Fatal(http.ListenAndServe(":8000", nil))
}

func getLedgers(w http.ResponseWriter, r *http.Request) {
    ledgers := []string{"CH1", "CH2", "EU1", "EU2", "AP1", "AP2", "AP3", "US"}

    jsonData, err := json.Marshal(ledgers)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}

func getLedger(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/api/local/ledger/")

    filename := fmt.Sprintf("ledger_%s.json", id)
    filePath := filepath.Join(".", filename)

    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        // Create empty JSON file if it doesn't exist
        emptyData := []byte("{}")
        err := os.WriteFile(filePath, emptyData, 0644)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    jsonData, err := os.ReadFile(filePath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}

func getDiskSpace(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/api/local/dispace/")

    filename := fmt.Sprintf("diskspace_%s.json", id)
    filePath := filepath.Join(".", filename)

    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        // Create empty JSON file if it doesn't exist
        emptyData := []byte("{}")
        err := os.WriteFile(filePath, emptyData, 0644)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    jsonData, err := os.ReadFile(filePath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}

func handleUnknownEndpoint(w http.ResponseWriter, r *http.Request) {
    unknownEndpoint := r.URL.Path
    fileName := strings.ReplaceAll(unknownEndpoint, "/", "_") + ".json"

    filePath := filepath.Join(".", fileName)

    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        // Create empty JSON file if it doesn't exist
        emptyData := []byte("{}")
        err := os.WriteFile(filePath, emptyData, 0644)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    jsonData, err := os.ReadFile(filePath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}
