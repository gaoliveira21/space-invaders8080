package debug

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ResponseError struct {
	StatusCode int
	Message    string
}

type MemoryDumpResponse struct {
	Data []string
}

func getMemoryDump(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	memory, err := os.ReadFile(".dump/memory")
	if err != nil {
		res := &ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Could not read memory dump",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	hexMemoryDump := []string{}
	for _, v := range memory {
		hexMemoryDump = append(hexMemoryDump, fmt.Sprintf("%.2X", v))
	}

	res := &MemoryDumpResponse{
		Data: hexMemoryDump,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func Start() {
	http.HandleFunc("GET /memory-dump", getMemoryDump)

	http.ListenAndServe(":8080", nil)
}
