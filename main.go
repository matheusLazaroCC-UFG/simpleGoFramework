package main

import (
    "github.com/matheusLazaroCC-UFG/simpleGoFramework/framework"
    "github.com/matheusLazaroCC-UFG/simpleGoFramework/user"
)

func main() {
    // Configurações para a aplicação
    config := &framework.Config{
        Port: "3000",
    }

    // Cria a aplicação principal
    app := framework.NewApp(config)

    // Instancia o serviço de usuário com, por exemplo, 3 workers
    userService := user.NewUserService(3)

    // Cria o controller e injeta o serviço
    userController := user.NewUserController(userService)

    // Registra o controller no nosso "framework"
    app.RegisterController(userController)

    // Inicia o servidor
    app.Start()
}
