package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"
    "fmt"

    // "github.com/matheusLazaroCC-UFG/simpleGoFramework/framework"
)

// Códigos ANSI de cor
const (
    ColorReset   = "\033[0m"
    CyanBright   = "\033[96m"
    Blue         = "\033[34m"
    Green        = "\033[32m"
    Yellow       = "\033[33m"
    Red          = "\033[31m"
    White        = "\033[97m"
)

type ProcessController struct {
    service ProcessService
}

func NewProcessController(service ProcessService) *ProcessController {
    return &ProcessController{service: service}
}

func (pc *ProcessController) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/process", pc.handleRequest)
}

func (pc *ProcessController) handleRequest(w http.ResponseWriter, r *http.Request) {
    fmt.Println()
    fmt.Printf("%s╔══════════════════════════════════════════════════════════════╗%s\n", CyanBright, ColorReset)
    fmt.Printf("%s║            NOVO CICLO DE REQUISIÇÃO RECEBIDO                ║%s\n", CyanBright, ColorReset)
    fmt.Printf("%s╚══════════════════════════════════════════════════════════════╝%s\n", CyanBright, ColorReset)

    fmt.Printf("%s[Servidor] Aguardando pedido...%s\n", Blue, ColorReset)

    if r.Method != http.MethodPost {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    delayStr := r.URL.Query().Get("delay")
    delay, err := strconv.Atoi(delayStr)
    if err != nil || delay <= 0 {
        http.Error(w, "Parâmetro 'delay' inválido", http.StatusBadRequest)
        return
    }

    var body struct {
        Message string `json:"message"`
    }
    _ = json.NewDecoder(r.Body).Decode(&body)

    handlerID := time.Now().UnixNano()

    fmt.Printf("%s[Servidor] Criando tratador #%d para delay=%d segundos%s\n", Blue, handlerID, delay, ColorReset)

    go func(id int64, d int, msg string) {
        fmt.Printf("\n───────────────────────────── Tratador #%d ─────────────────────────────\n", id)
        fmt.Printf("%s[Tratador #%d] Tratando pedido...%s\n", Green, id, ColorReset)
        if msg != "" {
            fmt.Printf("%s[Tratador #%d] Mensagem: %s%s\n", Yellow, id, msg, ColorReset)
        }
        result := pc.service.Handle(d)
        fmt.Printf("%s[Tratador #%d] Respondeu: \"%s\"%s\n", White, id, result, ColorReset)
        fmt.Printf("%s[Tratador #%d] Tratador termina.%s\n", Red, id, ColorReset)
    }(handlerID, delay, body.Message)

    fmt.Printf("%s[Servidor] Respondendo requisição para tratador #%d%s\n", Blue, handlerID, ColorReset)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":   "Tratador criado",
        "handler":  handlerID,
        "info":     "O tratamento ocorre em background.",
        "message":  body.Message,
    })
}
