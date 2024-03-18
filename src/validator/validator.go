package validator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"main/src/loggingMiddleware"

	"github.com/xeipuuv/gojsonschema"
)

func getSchema(request *http.Request, schema *gojsonschema.JSONLoader) {

	// Definir schema padrão como "{}", já que não existe validação (123 = '{', 125= '}')
	*schema = gojsonschema.NewStringLoader(string([]byte{123, 125}))

	//Pegando posição do arquivo pelo http da requisição
	u, err := url.Parse(request.URL.Path)
	if err != nil {
		log.Println(err.Error())
		log.Panicln("Erro ao pegar URL da request")

	}

	// Construção do caminho do arquivo usando filepath.Join
	pathScrema := filepath.Join("./src/validator/schema", path.Base(u.Path)+".json")

	//Testando se existe um arquivo para essa request
	_, err = os.Stat(pathScrema)
	if err != nil {

		file, err := ioutil.ReadFile(pathScrema)
		if err != nil {
			log.Println(err.Error())
			log.Panicln("Erro ao abrir arquivo json das schemas")
		}

		schemas := make(map[string]interface{})
		err = json.Unmarshal(file, &schemas)
		if err != nil {
			log.Println(err.Error())
			log.Panicln("Erro ao transformar json das schemas em map")
		}

		// Nome do Método
		schemaName := fmt.Sprintf("%s", request.Method)
		// Pega pelo nome do método no map a schema correspondente

		if sub_schema, ok := schemas[schemaName]; ok {
			// Passar para bytes da schema
			sub_schemaJSON, err := json.Marshal(sub_schema)
			if err != nil {
				log.Println(err.Error())
				log.Panicln("Erro ao passar a schema para bytes")

			}
			*schema = gojsonschema.NewStringLoader(string(sub_schemaJSON))

		}

	}

}
func validateSchema(lrw *loggingMiddleware.LoggingResponseWriter) error {
	// Pegando dados da requisição
	data_request := &lrw.DataRequest.Data

	// Pegando schema
	var schema gojsonschema.JSONLoader
	getSchema(lrw.Request, &schema)

	//Definir json padrão como "{}" (123 = '{', 125= '}')
	var jsonData []byte = []byte{123, 125}
	var err error
	if (*data_request) != nil {
		jsonData, err = json.Marshal(*data_request)
		if err != nil {
			log.Println(err.Error())
			log.Panicln("Erro ao transferir dados da request em json ")
		}
	}

	documentLoader := gojsonschema.NewStringLoader(string(jsonData))
	result, err := gojsonschema.Validate(schema, documentLoader)
	if err != nil {
		return fmt.Errorf("Erro ao Validar: %s", err)
	}

	if !result.Valid() {
		// Caso a requição não tenha os parametros validos

		return fmt.Errorf("Requisição invalida.\n Schema: %s", schema)

	}

	return nil
}
