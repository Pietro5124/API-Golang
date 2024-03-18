package main

import (
	"fmt"
	"log"
	"main/src/api/api1"

	"main/src/loggingMiddleware"

	"main/src/validator"
	"net/http"

	"github.com/gorilla/mux"
)

var porta int = 2000

func main() {

	r := mux.NewRouter()

	// Adiciona o middleware
	r.Use(Middleware)
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	// Define as rotas da API
	apiRouter.HandleFunc("/API1", api1.Main)

	// Inicia o servidor na porta 2000
	fmt.Printf("\nServiço iniciado - Porta: %d\n  ", porta)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", porta), r))

}
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		lrw := loggingMiddleware.CreateLoggingResponseWriter(response, request)

		// Antes da solicitação
		if request.Method != "OPTIONS" {
			// Validação para a request
			if validator.Main(lrw) /*Vc pode fazer qualquer validação necessaria nessa parte (Permissão, schema e etc..) */ {
				// Executa o handler da próxima etapa
				next.ServeHTTP(response, lrw.UpdateRequest(request))
			}
		}

		// Executa o próximo manipulador

		// Depois da solicitação log
		fmt.Printf("\n Foi")
		// fmt.Printf("\n  [%s] %s (%s) - %d \n", time.Now().Format("2006-01-02 15:04:05"), request.Method, request.URL.Path, request.Response.StatusCode)
		// // fmt.Print(response.Header())

	})
}
