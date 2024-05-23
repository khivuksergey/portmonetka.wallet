package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/khivuksergey/portmonetka.common"
	serviceerror "github.com/khivuksergey/portmonetka.wallet/error"
	"github.com/khivuksergey/portmonetka.wallet/internal/core/port/service"
	"github.com/khivuksergey/portmonetka.wallet/internal/model"
	"github.com/khivuksergey/webserver/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type WalletHandler struct {
	walletService service.WalletService
	logger        logger.Logger
	validate      *validator.Validate
}

func NewWalletHandler(services *service.Manager, logger logger.Logger) *WalletHandler {

	return &WalletHandler{
		walletService: services.Wallet,
		logger:        logger,
		validate:      model.GetWalletValidator(),
	}
}

// GetWallets retrieves user's wallets.
//
// @Tags Wallet
// @Summary Get user's wallets
// @Description Gets user's wallets
// @ID get-wallets
// @Accept json
// @Produce json
// @Param userId path uint64 true "Authorized user ID"
// @Success 200 {object} model.Response "Wallets retrieved"
// @Failure 422 {object} model.Response "Unprocessable entity"
// @Router /users/{userId}/wallets [get]
func (u WalletHandler) GetWallets(c echo.Context) error {
	requestUuid := c.Get(common.RequestUuidKey).(string)
	userId := c.Get("userId").(uint64)

	wallets, err := u.walletService.GetWalletsByUserId(userId)
	if err != nil {
		return common.NewUnprocessableEntityError(serviceerror.CannotGetWallets, err)
	}

	u.logger.Info(logger.LogMessage{
		Action:      "GetWallets",
		Message:     "Wallets retrieved",
		UserId:      &userId,
		RequestUuid: requestUuid,
	})

	return c.JSON(http.StatusOK, model.Response{
		Message:     "Wallets retrieved",
		Data:        wallets,
		RequestUuid: requestUuid,
	})
}

// CreateWallet creates a new wallet for user.
//
// @Tags Wallet
// @Summary Create a new wallet
// @Description Creates a new wallet with the provided information
// @ID create-wallet
// @Accept json
// @Produce json
// @Param userId path uint64 true "Authorized user ID"
// @Param wallet body model.WalletCreateDTO true "Wallet object to be created"
// @Success 201 {object} model.Response "Wallet created"
// @Failure 400 {object} model.Response "Bad request"
// @Failure 422 {object} model.Response "Unprocessable entity"
// @Router /users/{userId}/wallets [post]
func (u WalletHandler) CreateWallet(c echo.Context) error {
	requestUuid := c.Get(common.RequestUuidKey).(string)
	userId := c.Get("userId").(uint64)

	walletCreateDTO, err := u.bindWalletCreateDtoValidate(c, userId)
	if err != nil {
		return common.NewValidationError(serviceerror.InvalidInputData, err)
	}

	wallet, err := u.walletService.CreateWallet(walletCreateDTO)
	if err != nil {
		return common.NewUnprocessableEntityError(serviceerror.CannotCreateWallet, err)
	}

	u.logger.Info(logger.LogMessage{
		Action:      "CreateWallet",
		Message:     "Wallet created",
		UserId:      &wallet.UserId,
		Data:        map[string]uint64{"id": wallet.Id},
		RequestUuid: requestUuid,
	})

	return c.JSON(http.StatusCreated, model.Response{
		Message:     "Wallet created",
		Data:        wallet,
		RequestUuid: requestUuid,
	})
}

// UpdateWallet updates the wallet.
//
// @Tags Wallet
// @Summary Update wallet
// @Description Updates wallet's properties
// @ID update-wallet
// @Accept json
// @Produce json
// @Param userId path uint64 true "Authorized user ID"
// @Param wallet body model.WalletUpdateDTO true "Wallet update attributes"
// @Success 200 {object} model.Response "Wallet updated"
// @Failure 400 {object} model.Response "Bad request"
// @Failure 422 {object} model.Response "Unprocessable entity"
// @Router /users/{userId}/wallets/{walletId} [patch]
func (u WalletHandler) UpdateWallet(c echo.Context) error {
	requestUuid := c.Get(common.RequestUuidKey).(string)
	userId := c.Get("userId").(uint64)

	walletUpdateDTO, err := u.bindWalletUpdateDtoValidate(c, userId)
	if err != nil {
		return common.NewValidationError(serviceerror.InvalidInputData, err)
	}

	wallet, err := u.walletService.UpdateWallet(walletUpdateDTO)
	if err != nil {
		return common.NewUnprocessableEntityError(serviceerror.CannotUpdateWallet, err)
	}

	u.logger.Info(logger.LogMessage{
		Action:      "UpdateWallet",
		Message:     "Wallet updated",
		UserId:      &userId,
		Data:        map[string]uint64{"id": wallet.Id},
		RequestUuid: requestUuid,
	})

	return c.JSON(http.StatusOK, model.Response{
		Message:     "Wallet updated",
		Data:        wallet,
		RequestUuid: requestUuid,
	})
}

// DeleteWallet deletes the wallet by ID.
//
// @Tags Wallet
// @Summary Delete wallet
// @Description Deletes wallet by the provided wallet ID
// @ID delete-wallet
// @Accept json
// @Produce json
// @Param userId path uint64 true "Authorized user ID"
// @Param wallet body model.WalletDeleteDTO true "Wallet delete request"
// @Success 204 {string} string "No content"
// @Failure 400 {object} model.Response "Bad request"
// @Failure 422 {object} model.Response "Unprocessable entity"
// @Router /users/{userId}/wallets/{walletId} [delete]
func (u WalletHandler) DeleteWallet(c echo.Context) error {
	requestUuid := c.Get(common.RequestUuidKey).(string)
	userId := c.Get("userId").(uint64)

	walletDeleteDTO, err := u.bindWalletDeleteDtoValidate(c, userId)
	if err != nil {
		return common.NewValidationError(serviceerror.InvalidInputData, err)
	}

	if err := u.walletService.DeleteWallet(walletDeleteDTO); err != nil {
		return common.NewUnprocessableEntityError(serviceerror.CannotDeleteWallet, err)
	}

	u.logger.Info(logger.LogMessage{
		Action:      "DeleteWallet",
		Message:     "Wallet deleted",
		UserId:      &userId,
		Data:        map[string]uint64{"id": walletDeleteDTO.Id},
		RequestUuid: requestUuid,
	})

	return c.NoContent(http.StatusNoContent)
}

func (u WalletHandler) bindWalletCreateDtoValidate(c echo.Context, userId uint64) (*model.WalletCreateDTO, error) {
	walletCreateDTO := new(model.WalletCreateDTO)
	if err := c.Bind(walletCreateDTO); err != nil {
		return nil, err
	}
	if err := u.validate.Struct(walletCreateDTO); err != nil {
		return nil, err
	}
	walletCreateDTO.UserId = userId
	return walletCreateDTO, nil
}

func (u WalletHandler) bindWalletUpdateDtoValidate(c echo.Context, userId uint64) (*model.WalletUpdateDTO, error) {
	walletUpdateDTO := new(model.WalletUpdateDTO)
	if err := c.Bind(walletUpdateDTO); err != nil {
		return nil, err
	}
	if err := u.validate.Struct(walletUpdateDTO); err != nil {
		return nil, err
	}
	walletUpdateDTO.UserId = userId
	walletId, _ := strconv.ParseUint(c.Param("walletId"), 10, 64)
	walletUpdateDTO.Id = walletId
	return walletUpdateDTO, nil
}

func (u WalletHandler) bindWalletDeleteDtoValidate(c echo.Context, userId uint64) (*model.WalletDeleteDTO, error) {
	walletDeleteDTO := new(model.WalletDeleteDTO)
	if err := c.Bind(walletDeleteDTO); err != nil {
		return nil, err
	}
	if err := u.validate.Struct(walletDeleteDTO); err != nil {
		return nil, err
	}
	walletDeleteDTO.UserId = userId
	walletId, _ := strconv.ParseUint(c.Param("walletId"), 10, 64)
	walletDeleteDTO.Id = walletId
	return walletDeleteDTO, nil
}
