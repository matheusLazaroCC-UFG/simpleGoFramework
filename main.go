package main

import (
    "github.com/matheusLazaroCC-UFG/simpleGoFramework/framework"
    "github.com/matheusLazaroCC-UFG/simpleGoFramework/job"
    "log"
)

func main() {
    // Configurações básicas
    config := &framework.Config{
        Port: "3000",
    }

    // Cria a aplicação principal
    app := framework.NewApp(config)

    // Instancia o serviço de "jobs" (tarefas em background)
    jobService := job.NewJobService()

    // Cria o controller de "jobs" e injeta o serviço
    jobController := job.NewJobController(jobService)

    // Registra o controller no mini-framework
    app.RegisterController(jobController)

    log.Println("Servidor rodando em http://localhost:3000")
    log.Println("Use POST /jobs para criar tarefas, GET /jobs e GET /jobs/{id} para consultar...")

    // Inicia o servidor
    app.Start()
}
