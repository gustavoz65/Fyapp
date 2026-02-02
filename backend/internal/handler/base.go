package handler

import (
	"time"

	"github.com/aws/smithy-go/middleware"
	"github.com/gustavoz65/Cashing-go/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Handler struct {
	server *server.Server
}

func NewHandler(s *server.Server) Handler {
	return Handler{server: s}
}

// HandlerFunc representa a um tipo dandler e função que processa uma requisição e retorna uma resposta ou um erro.
type HandlerFunc[Req validation.Validatable, Res any] func(c echo.Context, req Req) (Res, error)

// HandlerFuncNoContent representa um tipo de função que processa uma requisição sem retornar conteúdo na resposta.
type HandlerFuncNoContent[req validation.Validatable] func(c echo.Context, req Req) error

// ResponseHandler é responsavel por manipular a resposta de uma operação, incluindo o manuseio da resposta, obtenção da operação e adição de atributos para monitoramento.
type ResponseHandler interface {
	Handle(c echo.Context, result interface{}) error
	GetOperation() string
	AddAttributes(trx *newrelic.Transaction, result interface{})
}

// JSONResponseHandler manipula respostas no formato JSON.
type JSONResponseHandler struct {
	status int
}

func (h *JSONResponseHandler) Handle(c echo.Context, result interface{}) error {
	return c.JSON(h.status, result)
}

func (h *JSONResponseHandler) GEtOperation() string {
	return "handler"
}

func (h *JSONResponseHandler) AddAttributes(txn *newrelic.Transaction, result interface{}) {
	// Implementação específica para adicionar atributos ao New Relic, se necessário.
}

// NoContentResponseHandler manipula respostas sem conteúdo (204 No Content).
type NoContentResponseHandler struct {
	status int
}

func (h *NoContentResponseHandler) Handle(c echo.Context, result interface{}) error {
	return c.NoContent(h.status)
}

func (h *NoContentResponseHandler) GetOperation() string {
	return "handler_no_content"
}

func (h *NoContentResponseHandler) AddAttributes(txn *newrelic.Transaction, result interface{}) {
	// Implementação específica para adicionar atributos ao New Relic, se necessário.
}

type FileResponseHandler struct {
	status      int
	filename    string
	contentType string
}

func (h FileResponseHandler) Handle(c echo.Context, result interface{}) error {
	data := result.([]byte)
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+h.filename)
	return c.Blob(h.status, h.contentType, data)
}

func (h FileResponseHandler) GetOperation() string {
	return "handler_file"
}

func (h FileResponseHandler) AddAttributes(txn *newrelic.Transaction, result interface{}) {
	if txn != nil {
		// http.status_code is already set by tracing middleware
		txn.AddAttribute("file.name", h.filename)
		txn.AddAttribute("file.content_type", h.contentType)
		if data, ok := result.([]byte); ok {
			txn.AddAttribute("file.size_bytes", len(data))
		}
	}
}

// handleRequest é uma função genérica que lida com a lógica comum de manipulação de requisições, incluindo validação, logging, métricas e rastreamento.
func handleRequest[Req validation.Validatable](
	c echo.Context,
	req Req,
	handler func(c echo.Context, req Req) (interface{}, error),
	responseHandler ResponseHandler,
) error {
	start := time.Now()
	method := c.Request().Method
	path := c.Path()
	route := path

	// Get novo relic transaction
	txn := newrelic.FromContext(c.Request().Context())
	if txn != nil {
		txn.AddAttribute("handler.name", route)
		// http.method comumentado pelo middleware de rastreamento
		// http.url comumentado pelo middleware de rastreamento
		responseHandler.AddAttributes(txn, nil)
	}

	// Get context-specific logger
	loggerBuilder := middleware.GetLogger(c).With().
		Str("operation", responseHandler.GetOperation()).
		Str("method", method).
		Str("path", path).
		Str("route", route)

	// Adicionar informações específicas para FileResponseHandler
	if fileHandler, ok := responseHandler.(FileResponseHandler); ok {
		loggerBuilder = loggerBuilder.
			Str("filename", fileHandler.filename).
			Str("content_type", fileHandler.contentType)
	}

	logger := loggerBuilder.Logger()

	// user.id é adicionado pelo middleware de autenticação, se disponível

	logger.Info().Msg("handling request")

	// Validation com observability
	validationStart := time.Now()
	if err := validation.BindAndValidate(c, req); err != nil {
		validationDuration := time.Since(validationStart)

		logger.Error().
			Err(err).
			Dur("validation_duration", validationDuration).
			Msg("request validation failed")

		if txn != nil {
			txn.NoticeError(nrpkgerrors.Wrap(err))
			txn.AddAttribute("validation.status", "failed")
			txn.AddAttribute("validation.duration_ms", validationDuration.Milliseconds())
		}
		return err
	}

	validationDuration := time.Since(validationStart)
	if txn != nil {
		txn.AddAttribute("validation.status", "success")
		txn.AddAttribute("validation.duration_ms", validationDuration.Milliseconds())
	}

	logger.Debug().
		Dur("validation_duration", validationDuration).
		Msg("request validation successful")

	// Execute handler with observability
	handlerStart := time.Now()
	result, err := handler(c, req)
	handlerDuration := time.Since(handlerStart)

	if err != nil {
		totalDuration := time.Since(start)

		logger.Error().
			Err(err).
			Dur("handler_duration", handlerDuration).
			Dur("total_duration", totalDuration).
			Msg("handler execution failed")

		if txn != nil {
			txn.NoticeError(nrpkgerrors.Wrap(err))
			txn.AddAttribute("handler.status", "error")
			txn.AddAttribute("handler.duration_ms", handlerDuration.Milliseconds())
			txn.AddAttribute("total.duration_ms", totalDuration.Milliseconds())
		}
		return err
	}

	totalDuration := time.Since(start)

	// Sucesso - adicionar atributos e log
	if txn != nil {
		txn.AddAttribute("handler.status", "success")
		txn.AddAttribute("handler.duration_ms", handlerDuration.Milliseconds())
		txn.AddAttribute("total.duration_ms", totalDuration.Milliseconds())
		responseHandler.AddAttributes(txn, result)
	}

	logger.Info().
		Dur("handler_duration", handlerDuration).
		Dur("validation_duration", validationDuration).
		Dur("total_duration", totalDuration).
		Msg("request completed successfully")

	return responseHandler.Handle(c, result)
}

// Handle serve para lidar com requisições que retornam uma resposta JSON.
func Handle[Req validation.Validatable, Res any](
	h Handler,
	handler HandlerFunc[Req, Res],
	status int,
	req Req,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handleRequest(c, req, func(c echo.Context, req Req) (interface{}, error) {
			return handler(c, req)
		}, JSONResponseHandler{status: status})
	}
}

func HandleFile[Req validation.Validatable](
	h Handler,
	handler HandlerFunc[Req, []byte],
	status int,
	req Req,
	filename string,
	contentType string,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handleRequest(c, req, func(c echo.Context, req Req) (interface{}, error) {
			return handler(c, req)
		}, FileResponseHandler{
			status:      status,
			filename:    filename,
			contentType: contentType,
		})
	}
}

// HandleNoContent serve para lidar com requisições que não retornam conteúdo na resposta.
func HandleNoContent[Req validation.Validatable](
	h Handler,
	handler HandlerFuncNoContent[Req],
	status int,
	req Req,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handleRequest(c, req, func(c echo.Context, req Req) (interface{}, error) {
			err := handler(c, req)
			return nil, err
		}, NoContentResponseHandler{status: status})
	}
}
