package utils

import (
    "log";
    "net/http";
    "encoding/json";
)


// This is the function that will create a JSON response for our server
func JSONResponse(
    w http.ResponseWriter,
    statusCode int,
    payload interface{},
) {
    data, error := json.Marshal(payload)

    // If there is an error Parsing the struct or information being passed
    if error != nil {
        log.Printf("Error parshing the data: %v", payload)
        w.WriteHeader(500)
        return
    }

    w.Header().Add("Content-type", "application/json")
    w.WriteHeader(statusCode)
    w.Write(data)
}
