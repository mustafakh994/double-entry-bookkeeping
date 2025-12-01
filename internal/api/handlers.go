package api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type createAccountRequest struct {
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

func (server *Server) createAccount(ctx echo.Context) error {
	var req createAccountRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	account, err := server.service.CreateAccount(ctx.Request().Context(), req.Balance, req.Currency)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `param:"id"`
}

func (server *Server) getAccount(ctx echo.Context) error {
	var req getAccountRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	account, err := server.service.GetAccount(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, account)
}

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
}

func (server *Server) createTransfer(ctx echo.Context) error {
	var req transferRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	result, err := server.service.Transfer(ctx.Request().Context(), req.FromAccountID, req.ToAccountID, req.Amount)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, result)
}

func errorResponse(err error) echo.Map {
	return echo.Map{"error": err.Error()}
}
