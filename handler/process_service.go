package handler

import (
    "fmt"
    "time"
)

// ProcessService simula um tratamento de pedido
type ProcessService interface {
    Handle(delay int) string
}

type processServiceImpl struct{}

// NewProcessService cria uma nova instância do serviço
func NewProcessService() ProcessService {
    return &processServiceImpl{}
}

// Handle simula uma tarefa demorada (bloqueia por N segundos)
func (s *processServiceImpl) Handle(delay int) string {
    time.Sleep(time.Duration(delay) * time.Second)
    return fmt.Sprintf("Processamento concluído em %d segundos", delay)
}
