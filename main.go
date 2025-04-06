package main

import (
    "log"
    "net/http"
    "strconv"

    "github.com/matheusLazaroCC-UFG/simpleGoFramework/framework"
    "github.com/matheusLazaroCC-UFG/simpleGoFramework/tasks"
)

func main() {
    // Configuração básica
    config := &framework.Config{
        Port: "3000",
    }

    // Cria a aplicação principal (mini-framework)
    app := framework.NewApp(config)

    // Instancia nosso serviço de tarefas
    taskService := tasks.NewTaskService()

    // Cria o controller de tarefas, injeta o serviço
    taskController := tasks.NewTaskController(taskService)

    // Registra as rotas definidas pelo nosso controller
    app.RegisterController(taskController)

    log.Println("Para testar, acesse por exemplo:")
    log.Println("GET http://localhost:3000/primes?start=1&end=1000&workers=5")

    // Inicia o servidor
    app.Start()
}
