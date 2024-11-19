// Main application from the server. Here will the be initial setup of the
// server and the health check path with the corresponding response.
package main

import (
    "log";
    "time";
    "fmt";
    "errors";
    "net/http";
    "encoding/json";

    ls "jacobitosuperstar/LoanSizing/internal/loan_sizer";
    ff "jacobitosuperstar/LoanSizing/internal/financial_formulas";
)

const PORT = ":8000";


type HealthCheckResponse struct {
    Now time.Time   `json:"now"`
}

type Response struct {
    Message string `json:"message"`
}


func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /", handleRoot)
    // TODO: Change this path to /health/ later.
    mux.HandleFunc("POST /loan_sizer", handleLoanSizer)

    log.Printf("Server listening in the port %s", PORT)
    err := http.ListenAndServe(PORT, mux)
    if err != nil {
        log.Fatal(err)
    }
}


// handleRoot handles the request for the root direction of the server. Returns
// an OK http status and the current time of the server.
func handleRoot(
    w http.ResponseWriter,
    r *http.Request,
) {
    response := HealthCheckResponse{
        Now: time.Now(),
    }
    JSONResponse(w, http.StatusOK, response)
}

// handleLoanSizer hadles the post request with the information to size the
// loan and if everything is correct, returns the json representation of the
// LoanSizer struct
func handleLoanSizer(
    w http.ResponseWriter,
    r *http.Request,
) {
    var loan_sizer ls.LoanSizer
    err := json.NewDecoder(r.Body).Decode(&loan_sizer)

    if err != nil {
        response := Response{
            Message: "Invalid request body",
        }
        JSONResponse(w, http.StatusBadRequest, response)
        return
    }

    loan_sizer, err = ls.InitLoanSizer(loan_sizer)
    if err != nil {
        var validationError *ff.ValidationError
        var response Response

        if errors.Is(err, validationError) {
            response = Response{
                Message: fmt.Sprintf("Validation Error: %v", err),
            }
            JSONResponse(w, http.StatusBadRequest, response)
        } else {
            response = Response{
                Message: "Internal Server Error",
            }
            log.Println(err)
            JSONResponse(w, http.StatusInternalServerError, response)
        }
        return
    }

    JSONResponse(w, http.StatusOK, loan_sizer)
    return
}
