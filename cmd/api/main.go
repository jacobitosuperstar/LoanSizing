// Main application from the server. Here will the be initial setup of the
// server and the health check path with the corresponding response.
package main


import (
    "log";
    "time";
    "net/http";
    utils "jacobitosuperstar/LoanSizing/internal/utils";
)

const PORT = ":8000";


type HealthCheckResponse struct {
    Now time.Time   `json:"now"`
}


func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /", handleRoot)
    // TODO: Change this path to /health/ later.

    log.Printf("Server listening in the port %s", PORT)
    err := http.ListenAndServe(PORT, mux)
    if err != nil {
        log.Fatal(err)
    }
}


func handleRoot(
    w http.ResponseWriter,
    r *http.Request,
) {
    response := HealthCheckResponse{
        Now: time.Now(),
    }
    utils.JSONResponse(w, http.StatusOK, response)
}
