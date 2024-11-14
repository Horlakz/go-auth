package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/horlakz/go-auth/dto"
	"github.com/horlakz/go-auth/payload/request"
	"github.com/horlakz/go-auth/payload/response"
	"github.com/horlakz/go-auth/service"
	"github.com/horlakz/go-auth/validator"
)

type authHandler struct {
	authService service.AuthServiceInterface
	validator   validator.AuthValidator
}

type AuthHandlerInterface interface {
	Login(c *fiber.Ctx) error
	RefreshAccessToken(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	ForgotPassword(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
}

func NewAuthHandler(authService service.AuthServiceInterface) AuthHandlerInterface {
	return &authHandler{authService: authService}
}

func (handler *authHandler) Login(c *fiber.Ctx) error {
	var resp response.Response
	var loginDto dto.LoginDto

	loginRequest := new(request.LoginRequest)

	if err := c.BodyParser(loginRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Invalid request"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if _, err := handler.validator.LoginValidate(*loginRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		resp.Error = err
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	loginDto.Email = loginRequest.Email
	loginDto.Password = loginRequest.Password

	loginResp, err := handler.authService.Login(loginDto)

	if err != nil {
		resp.Status = 400
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = 200
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = loginResp

	return c.JSON(resp)
}

func (handler *authHandler) Register(c *fiber.Ctx) error {
	var resp response.Response
	var registerDto dto.RegisterDto

	registerRequest := new(request.RegisterRequest)

	if err := c.BodyParser(registerRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Invalid request"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if vEs, err := handler.validator.RegisterValidate(*registerRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		resp.Data = vEs
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	registerDto.FullName = registerRequest.FullName
	registerDto.Email = registerRequest.Email
	registerDto.Password = registerRequest.Password
	registerDto.Role = registerRequest.Role

	if err := handler.authService.Register(registerDto); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusCreated
	resp.Message = "Registeration is Successful"

	return c.JSON(resp)
}

func (handler *authHandler) RefreshAccessToken(c *fiber.Ctx) error {
	var resp response.Response

	refreshAccessTokenRequest := new(request.RefreshAccessTokenRequest)

	if err := c.BodyParser(refreshAccessTokenRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Invalid request"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	token, err := handler.authService.RefreshAccessToken(refreshAccessTokenRequest.RefreshToken)

	if err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = map[string]interface{}{"token": token}

	return c.JSON(resp)
}

func (handler *authHandler) ForgotPassword(c *fiber.Ctx) error {
	var resp response.Response

	emailRequest := new(request.EmailRequest)

	if err := c.BodyParser(emailRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Invalid request"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if vEs, err := handler.validator.EmailValidate(*emailRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		resp.Error = vEs
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if err := handler.authService.ForgotPassword(emailRequest.Email); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = "Reset Password Code sent to email"

	return c.JSON(resp)
}

func (handler *authHandler) ResetPassword(c *fiber.Ctx) error {
	var resp response.Response
	var resetPasswordDto dto.ResetPasswordDto

	resetPasswordRequest := new(request.ResetPasswordRequest)

	if err := c.BodyParser(resetPasswordRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Invalid request"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if vEs, err := handler.validator.ResetPasswordValidate(*resetPasswordRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		resp.Error = vEs
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resetPasswordDto.Code = resetPasswordRequest.Code
	resetPasswordDto.Email = resetPasswordRequest.Email
	resetPasswordDto.Password = resetPasswordRequest.Password

	if err := handler.authService.ResetPassword(resetPasswordDto); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	return c.JSON(resp)
}
