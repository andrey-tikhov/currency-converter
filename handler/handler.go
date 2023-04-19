package handler

import (
	"fmt"
	"go.uber.org/fx"
	"io"
	"my_go/controller/cb_repository"
	"my_go/controller/conversion"
	"my_go/mapper"
	"net/http"
	"time"
)

// Handler interface encapsulates external endpoints for the service
type Handler interface {
	GetCBRates(w http.ResponseWriter, req *http.Request)
	ConvertCurrency(w http.ResponseWriter, req *http.Request)
	Hello(w http.ResponseWriter, req *http.Request)
}

// Compile time check that handler implements Handler interface
var _ Handler = (*handler)(nil)

type handler struct {
	repositoryCtrl cb_repository.Controller
	conversionCtrl conversion.Controller
}

// Params is a container with dependencies for Handler interface creation
type Params struct {
	fx.In

	CBRepositoryController cb_repository.Controller
	ConversionController   conversion.Controller
}

// New is a constructor of Handler interface
func New(p Params) (Handler, error) {
	return &handler{
		repositoryCtrl: p.CBRepositoryController,
		conversionCtrl: p.ConversionController,
	}, nil
}

// Hello is GET endpoint prints "Hello World" in the http Response
func (h *handler) Hello(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Received request")
	time.Sleep(10 * time.Second)
	if req == nil || req.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	w.Header().Add("Content-Type", "text/html")
	_, err := w.Write([]byte(`hello world`))
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to process the request, err %s", err),
			http.StatusInternalServerError,
		)
	}
	fmt.Println("Completed request")
}

// GetCBRates it the POST endpoint that loads available exchange rates for the Central Bank provided in the request
// Expected json request is defined by entity.GetCBRatesRequest
// Expected json response is defined by entity.GetCBRatesResponse
func (h *handler) GetCBRates(w http.ResponseWriter, req *http.Request) {
	if req == nil || req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer req.Body.Close()
	data, err := io.ReadAll(io.LimitReader(req.Body, 1<<20))
	if err != nil {
		http.Error(w, "unable to read the body", http.StatusBadRequest)
		return
	}
	getCBRateRequest, err := mapper.BodyToGetExchangeRatesRequest(data)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("bad request, err %s", err),
			http.StatusBadRequest,
		)
		return
	}
	response, err := h.repositoryCtrl.GetCBRates(req.Context(), getCBRateRequest)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to process the request, err %s", err),
			http.StatusBadGateway,
		)
		return
	}
	getCBRateResponse, err := mapper.GetExchangeRatesResponseToBytes(response)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to process the response, err %s", err),
			http.StatusInternalServerError,
		)
		return // unreachable in tests cause response struct can always be represented as json
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(getCBRateResponse)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to process the request, err %s", err),
			http.StatusInternalServerError,
		)
		return // unreachable in tests
	}
	return
}

// ConvertCurrency is the POST endpoint to convert arbitrary amount of currency to another currency
// Expected json is defined by entity.ConvertCurrencyRequest
// Expected response is defined by entity.ConvertCurrencyResponse
func (h *handler) ConvertCurrency(w http.ResponseWriter, req *http.Request) {
	if req == nil || req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer req.Body.Close()
	data, err := io.ReadAll(io.LimitReader(req.Body, 1<<20))
	if err != nil {
		http.Error(w, "unable to read the body", http.StatusBadRequest)
		return
	}
	convertCurrencyRequest, err := mapper.BodyToConvertCurrencyRequest(data)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("bad request, err %s", err),
			http.StatusBadRequest,
		)
		return
	}
	response, err := h.conversionCtrl.Convert(req.Context(), convertCurrencyRequest)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to process the request, err %s", err),
			http.StatusBadGateway,
		)
		return
	}
	convertCurrencyResponse, err := mapper.ConvertCurrencyResponseToBytes(response)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to process the request, err %s", err),
			http.StatusInternalServerError,
		)
		return // unreachable in tests cause response struct can always be represented as json
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(convertCurrencyResponse)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to process the request, err %s", err),
			http.StatusInternalServerError,
		)
		return // unreachable in tests
	}
	return
}
