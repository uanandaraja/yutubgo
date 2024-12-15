package main

import (
    "encoding/json"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
)

type Request struct {
    URL string `json:"url"`
}

func main() {
    http.HandleFunc("/download", handleDownload)
    http.ListenAndServe(":8080", nil)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req Request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    outputPath := filepath.Join(os.TempDir(), "output.mp3")
    cmd := exec.Command("yt-dlp",
        "--extract-audio",
        "--audio-format", "mp3",
        "--output", outputPath,
        req.URL)

    if err := cmd.Run(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "audio/mpeg")
    w.Header().Set("Content-Disposition", "attachment; filename=download.mp3")
    http.ServeFile(w, r, outputPath)
    
    defer os.Remove(outputPath)
}
