package api1

import (
	"log"
	"main/src/loggingMiddleware"
	"runtime/debug"

	"net/http"
)

func Main(response http.ResponseWriter, request *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("PANIC")
			http.Error(response, "Erro interno", http.StatusInternalServerError)
		} else {
			debug.FreeOSMemory()
		}

	}()
	var resp interface{}
	var err error

	// Aqui vc pega os dados da request e coloca nos struct lrw
	lrw := loggingMiddleware.CreateLoggingResponseWriter(response, request)

	// Verifica se o método HTTP é POST
	if request.Method == "GET" {
		resp, err = get(lrw)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if request.Method == "POST" {
		resp, err = post(lrw)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if request.Method == "PUT" {
		resp, err = put(lrw)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Caso não exista esse metodo na API
		// Um exemplo disso é o method DELETE, já que não foi implementado nessa API
		// Retornamos
		http.Error(response, "Método inválido", http.StatusMethodNotAllowed)
		return
	}

	// Se caso exista uma resposta
	if resp != nil {

		// Nesse caso eu vou definir o Content-Type como application/json e escrever o corpo da response
		response.Header().Set("ContentType", "application/json")
		_, err = response.Write(resp.([]byte))
		if err != nil {
			log.Panicln("Erro ao carregar resposta ")

		}

	}

}
func get(request *loggingMiddleware.LoggingResponseWriter) (interface{}, error) {
	// Aqui vc pode pegar os parametros enviados na request
	// Dados := request.DataRequest.Data
	// Os dados enviados na request podem ser de qualquer tipo como map[string]interface{}, []map[string]interface{} entre outros....
	// vc pode definir o tipo de Dados utilizando request.DataRequest.Data .(map[string]interface{}), se caso a interface não for desse tipo o valor retornado será nil

	// Efetuar a operação de GET
	log.Panicln("Mensagem enviada no PANIC")
	return nil, nil
}
func post(request *loggingMiddleware.LoggingResponseWriter) (interface{}, error) {
	// Efetuar a operação de POST
	return nil, nil
}

func put(request *loggingMiddleware.LoggingResponseWriter) (interface{}, error) {
	// Efetuar a operação de PUT
	return nil, nil
}
