package tasks

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/matheusLazaroCC-UFG/simpleGoFramework/framework"
)

// TaskController expõe endpoints para processos concorrentes
type TaskController struct {
    service TaskService
}

// NewTaskController injeta o TaskService
func NewTaskController(service TaskService) *TaskController {
    return &TaskController{service: service}
}

// RegisterRoutes implementa a interface Controller, registrando rotas
func (tc *TaskController) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/primes", tc.handlePrimes)
}

// handlePrimes lida com GET /primes?start=1&end=1000&workers=5
func (tc *TaskController) handlePrimes(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
        return
    }

    // Extrai parâmetros de query
    query := r.URL.Query()
    startStr := query.Get("start")
    endStr := query.Get("end")
    workersStr := query.Get("workers")

    // Converte para int
    start, _ := strconv.Atoi(startStr)
    end, _ := strconv.Atoi(endStr)
    workers, _ := strconv.Atoi(workersStr)

    if start == 0 && end == 0 {
        // Se não foi passado nada, define um exemplo
        start = 1
        end = 100
        workers = 5
    } else if workers == 0 {
        workers = 4 // default
    }

    primes := tc.service.FindPrimesInRange(start, end, workers)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "start":   start,
        "end":     end,
        "workers": workers,
        "primes":  primes,
    })
}
