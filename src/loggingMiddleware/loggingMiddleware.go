package loggingMiddleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type LoggingResponseWriter struct {
	ResponseW   http.ResponseWriter
	Request     *http.Request
	DataRequest dataRequest
}
type dataRequest struct {
	Data interface{}
	// Aqui vc pode colocar qualquer parametro necessario como nivel de permissão, tokens de acesso entre outras coisas
}

func CreateLoggingResponseWriter(response http.ResponseWriter, request *http.Request) *LoggingResponseWriter {
	lrw := &LoggingResponseWriter{}

	// Pegando dados da requisição
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Headers", "*")
	response.Header().Set("Access-Control-Allow-Methods", "*")

	lrw.WriteRequest(response, request)

	return lrw
}
func (lrw *LoggingResponseWriter) WriteRequest(response http.ResponseWriter, request *http.Request) {

	lrw.DataRequest.getValuesRequest(request)
	lrw.ResponseW = response
	lrw.Request = request

}
func (lrw *LoggingResponseWriter) UpdateRequest(request *http.Request) *http.Request {

	jsonData, err := json.Marshal(lrw.DataRequest)
	if err != nil {
		log.Println(err.Error())
		log.Panicln("Falha em atualizar a request")
	}

	request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))
	request.ContentLength = int64(len(jsonData))
	request.Header.Set("Content-Type", "application/json")

	return request
}

func (dataRequest *dataRequest) getValuesRequest(request *http.Request) {
	request.ParseMultipartForm(640000)

	// Pegando todos os bytes que tem no corpo da request(JSON,XML.Texto simples, etc...)
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		log.Panicln("Falha ao pegar o corpo da request")
	}
	// Teste para ver se tem algo no corpo do request
	if len(body) != 0 {
		err = json.Unmarshal(body, &dataRequest)

		if err != nil {
			err = json.Unmarshal(body, &dataRequest.Data)
			if err != nil {
				log.Println(err.Error())
				log.Panicln("Falha ao pegar dados do corpo da request")
			}
		}

	} else if request.MultipartForm != nil {
		// Adicionar todos os valores do form-data na solicitação
		//Pegando strings
		err := request.ParseForm()
		if err == nil {
			for key, values := range request.PostForm {
				if len(values) > 0 && values[0] != "" {
					(*dataRequest).Data.(map[string]interface{})[key] = values[0]
				}
			}
		}

	} else {

		// Adicionar todos os valores dos args na solicitação
		for key, values := range request.URL.Query() {
			if len(values) > 0 && values[0] != "" {
				(*dataRequest).Data.(map[string]interface{})[key] = values[0]
			}
		}
	}

}

func (lrw *LoggingResponseWriter) Header() http.Header {
	return lrw.ResponseW.Header()
}

func (lrw *LoggingResponseWriter) Write(b []byte) (int, error) {
	return lrw.ResponseW.Write(b)
}

// func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {

// 	lrw.ResponseW.WriteHeader(statusCode)
// }
