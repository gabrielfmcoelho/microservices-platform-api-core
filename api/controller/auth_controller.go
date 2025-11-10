package controller

import (
	"net/http"

	"github.com/gabrielfmcoelho/platform-core/bootstrap"
	"github.com/gabrielfmcoelho/platform-core/domain"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthUsecase domain.AuthUsecase
	Env         *bootstrap.Env
}

// User Login
// @Summary Login user
// @Description Authenticates a user using their email and password, then returns access and refresh tokens for session management.
// @Tags Auth User
// @ID login
// @Accept json
// @Produce json
// @Param loginRequest body domain.LoginRequest true "Login Request"
// @Success 200 {object} domain.LoginResponse "Successful login, returns access and refresh tokens"
// @Failure 400 {object} domain.ErrorResponse "Bad Request - Invalid input"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized - Incorrect email or password"
// @Failure 404 {object} domain.ErrorResponse "Not Found - User not found"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /login [post]
func (lc *AuthController) Login(c *gin.Context) {

	var request domain.LoginRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	loginResponse, err := lc.AuthUsecase.LoginUserByEmail(
		c,
		request.Email,
		request.Password,
		lc.Env.AccessTokenSecret,
		lc.Env.AccessTokenExpiryHour,
		lc.Env.RefreshTokenSecret,
		lc.Env.RefreshTokenExpiryHour,
	)

	if err != nil {
		switch err {
		case domain.ErrUserEmailNotFound:
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		case domain.ErrUserPasswordNotMatch:
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, loginResponse)
}

// Login Guest
// @Summary Login Guest
// @Description Authenticates a guest user using their IP address, then returns access and refresh tokens for session management
// @Tags Auth User
// @ID loginGuest
// @Accept json
// @Produce json
// @Success 200 {object} domain.LoginResponse "Successful login, returns access and refresh tokens"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /login-guest [post]
func (lc *AuthController) LoginGuest(c *gin.Context) {
	loginResponse, err := lc.AuthUsecase.LoginGuestUser(
		c,
		lc.Env.AccessTokenSecret,
		lc.Env.AccessTokenExpiryHour,
		lc.Env.RefreshTokenSecret,
		lc.Env.RefreshTokenExpiryHour,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, loginResponse)
}

// @Summary Forgot Password
// @Description Sends an email to the user with a link to reset their password
// @Tags Auth User
// @ID forgotPassword
// @Accept json
// @Produce json
// @Param email query string true "User email"
// @Success 200 {object} domain.SuccessResponse "Email sent successfully"
// @Failure 400 {object} domain.ErrorResponse "Bad Request - Invalid input"
// @Failure 404 {object} domain.ErrorResponse "Not Found - User not found"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /forgot-password [post]
func (lc *AuthController) ForgotPassword(c *gin.Context) {}

// @Summary Reset Password
// @Description Resets the user's password
// @Tags Auth User
// @ID resetPassword
// @Accept json
// @Produce json
// @Param email query string true "User email"
// @Param newPassword query string true "New password"
// @Success 200 {object} domain.SuccessResponse "Password reset successfully"
// @Failure 400 {object} domain.ErrorResponse "Bad Request - Invalid input"
// @Failure 404 {object} domain.ErrorResponse "Not Found - User not found"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /reset-password [post]
func (lc *AuthController) ResetPassword(c *gin.Context) {}

// @Summary Refresh Token
// @Description Refreshes the user's access token
// @Tags Auth User
// @ID refreshToken
// @Accept json
// @Produce json
// @Param refreshTokenRequest body domain.RefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} domain.RefreshTokenResponse "Access token refreshed successfully"
// @Failure 400 {object} domain.ErrorResponse "Bad Request - Invalid input"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized - Invalid refresh token"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /refresh-token [post]
func (lc *AuthController) RefreshToken(c *gin.Context) {
	var request domain.RefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	refreshResponse, err := lc.AuthUsecase.RefreshToken(
		c,
		request.RefreshToken,
		lc.Env.RefreshTokenSecret,
		lc.Env.AccessTokenSecret,
		lc.Env.AccessTokenExpiryHour,
		lc.Env.RefreshTokenExpiryHour,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Invalid or expired refresh token"})
		return
	}

	c.JSON(http.StatusOK, refreshResponse)
}
