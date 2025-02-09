package endpoints

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/SuperMatch/model"
	"github.com/SuperMatch/model/dto"
	"github.com/SuperMatch/pkg/db"
	"github.com/SuperMatch/pkg/db/dao"
	Service "github.com/SuperMatch/service"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRegistrationHandler godoc
//
//	@Summary		userRegister
//	@Description	User register API
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			userRegister	body		dto.UserRegistrationDTO	true	"userRegister"
//	@Success		200				{string}	string					"user created."
//	@Failure		400				{string}	string					"Bad request"
//	@Failure		500				{string}	string					"Internal Server Error"
//	@Router			/user/register [post]
func UserRegistrationHandler(c *gin.Context) {
	var userData dto.UserRegistrationDTO

	if err := c.BindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in parsing request."})
		return
	}

	user := model.User{
		FirstName:  userData.FirstName,
		LastName:   userData.LastName,
		Mobile:     userData.Mobile,
		Code:       userData.Code,
		Email:      userData.Email,
		Password:   userData.Password,
		IsActive:   true,
		SignUpType: "password",
	}

	userDao := &dao.UserDao{
		Connection: *db.GlobalOrm,
	}

	users, err := userDao.FindByMobileEmail(user.Mobile, user.Email)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "db connection error"})
		return
	}

	if len(users) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user already exist"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	user.Password = string(hashedPassword)
	user, _ = userDao.Insert(user)

	//just for test purpose only
	var profileDTO dto.UserProfile

	userProfileService := Service.NewUserProfileService()
	actualUserProfile, err := userProfileService.GetUserProfileFromDB(int(user.ID))

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	actualUserProfile.UserId = int(user.ID)
	profileDTO, err = userProfileService.CreateUserProfile(actualUserProfile, profileDTO)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in creating user.", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": profileDTO, "message": "user created successfully."})
}

// LoginUser godoc
//
//	@Summary		userLogin
//	@Description	User Username and Password login API
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			userLogin	body		dto.Login	true	"userLogin"
//	@Success		200			{string}	string		token
//	@Failure		400			{string}	string		"Bad request"
//	@Failure		500			{string}	string		"Internal Server Error"
//	@Router			/user/login [post]
func LoginUser(c *gin.Context) {
	var newUser dto.Login
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	if strings.Contains(newUser.Username, " ") || strings.Contains(newUser.Username, "?") || strings.Contains(newUser.Password, " ") || strings.Contains(newUser.Password, "?") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username or passwrod is missing."})
		return
	}

	userRepo := dao.UserDao{
		Connection: *db.GlobalOrm,
	}

	actualUser, err := userRepo.FindByEmail(newUser.Username)

	if err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "No such user found."})
		return
	}

	if strings.ToLower(actualUser.SignUpType) != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "password didn't exists for this account. please sign in using " + actualUser.SignUpType + " method."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(actualUser.Password), []byte(newUser.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "username or password is wrong."})
		return
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "No such user Found."})
		return
	}

	JWTService := Service.JWTImpl{}
	token, expiresAt, err := JWTService.GenerateToken(int(actualUser.ID), actualUser.Email, true)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error in generating auth token"})
		return
	}

	err = JWTService.SaveToken(actualUser, token, expiresAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error in saving auth token", "error": err.Error()})
	}

	c.Header("expires_at", expiresAt.String())
	c.Header("token", token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GoogleLoginHandler godoc
//
//	@Summary		GoogleLoginHandler
//	@Description	google login API
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			idToken	body		dto.UserIDToken	true	"ID token"
//	@Success		200		{string}	string			token
//	@Failure		400		{string}	string			"Bad request"
//	@Failure		500		{string}	string			"Internal Server Error"
//	@Router			/user/google/login [post]
func GoogleLoginHandler(c *gin.Context) {

	var idToken dto.UserIDToken
	if err := c.BindJSON(&idToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	loginService := Service.NewLoginService()
	resp, err := loginService.GoogleLogin(idToken.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "error in GoogleLogin function", "error": err.Error()})
		return
	}

	userRepo := dao.UserDao{
		Connection: *db.GlobalOrm,
	}

	userDetails, err := userRepo.FindByEmail(resp.Email)
	if err == nil {
		if strings.ToLower(userDetails.SignUpType) != "social" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "social signin didn't exists for this account. please sign in using " + userDetails.SignUpType + " method.", "error": err.Error()})
			return
		}

		loginService := Service.NewLoginService()
		token, expiresAt, err := loginService.UserSignIN(userDetails)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error in user sign in", "error": err.Error()})
			return
		}

		c.Header("expires_at", expiresAt.String())
		c.Header("token", token)
		c.JSON(http.StatusOK, gin.H{"token": token, "user_id": userDetails.ID})

	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		user := loginService.GetUserDetailsFromGoogle(resp)
		userID, token, expiresAt, err := loginService.UserSignUP(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error in user sign up", "error": err.Error()})
			return
		}

		c.Header("expires_at", expiresAt.String())
		c.Header("token", token)
		c.JSON(http.StatusOK, gin.H{"token": token, "user_id": userID})

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user is not registered", "error": err.Error()})
		return
	}

}

// SendOTP godoc
//
//	@Summary		Send OTP
//	@Description	API for sending verification sms to user with OTP
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			userData	body		dto.UserData	true	"userData"
//	@Success		200			{string}	string			"OTP sent successfully"
//	@Failure		400			{string}	string			"Bad request"
//	@Failure		500			{string}	string			"Internal Server Error"
//	@Router			/user/send/otp [post]
func SendOTP(c *gin.Context) {
	var userData dto.UserData
	if err := c.BindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	loginService := Service.NewLoginService()
	_, err := loginService.SendOTPService(userData.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send the OTP for verification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// VerifyOTP godoc
//
//	@Summary		Verify OTP
//	@Description	API for verification of the OTP
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			verifyUser	body		dto.VerifyUser	true	"verify user data"
//	@Success		200			{string}	string			token
//	@Failure		400			{string}	string			"Bad request"
//	@Failure		500			{string}	string			"Internal Server Error"
//	@Router			/user/verify/otp [post]
func VerifyOTP(c *gin.Context) {
	var verifyUser dto.VerifyUser
	if err := c.BindJSON(&verifyUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	loginService := Service.NewLoginService()
	isOtpVerified, err := loginService.VerifyOTPService(verifyUser.User.PhoneNumber, verifyUser.OTP)
	if err != nil || !isOtpVerified {
		c.JSON(http.StatusBadRequest, gin.H{"message": "OTP not approved, user verification failed", "error": err.Error()})
		return
	}

	if isOtpVerified {
		token, expiresAt, err := loginService.CheckUserExistOrNot(verifyUser.User.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "error checking user existence", "error": err.Error()})
			return
		}
		c.Header("expires_at", expiresAt.String())
		c.Header("token", token)
		c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully", "token": token})
	}

}

// SendVerificationEmail godoc
//
//	@Summary		Send Verification Email
//	@Description	API for sending verification email to user
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			verifyEmail	body		dto.UserEmail	true	"send user verification email"
//	@Success		200			{string}	string			"verification email sent successfully."
//	@Failure		400			{string}	string			"Bad request"
//	@Failure		500			{string}	string			"Internal Server Error"
//	@Router			/user/verify/email [post]
func SendVerificationEmail(c *gin.Context) {
	var verifyEmail dto.UserEmail
	if err := c.BindJSON(&verifyEmail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	loginService := Service.NewLoginService()
	err := loginService.SendVerificationEmail(verifyEmail.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error sending verification email", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification email sent successfully."})
}

// VerifyEmail godoc
//
//	@Summary		Verify Email
//	@Description	API for verify email of the user
//	@Tags			Authentication
//	@Produce		json
//	@Param			verification_code	query		string	true	"verification code"
//	@Success		200					{string}	string	"email verified successfully."
//	@Failure		400					{string}	string	"Bad request"
//	@Failure		500					{string}	string	"Internal Server Error"
//	@Router			/user/verify/email [get]
func VerifyEmail(c *gin.Context) {
	verificationCode, ok := c.GetQuery("verification_code")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verificationCode is required."})
		return
	}

	loginService := Service.NewLoginService()
	err := loginService.VerifyEmailService(verificationCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in verifying email. ", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully."})
}

// DeleteUserHandler godoc
//
//	@Summary		Delete user handler
//	@Description	API for deleting user by updating the email. Just for testing purpose.
//	@Tags			For Testing Purpose Only
//	@Accept			json
//	@Produce		json
//	@Param			user_id	header		int				true	"User ID"
//	@Param			email	body		dto.UserEmail	true	"new email which will replace the old email"
//	@Success		200		{string}	string			"user deleted successfully."
//	@Failure		400		{string}	string			"Bad request"
//	@Failure		500		{string}	string			"Internal Server Error"
//	@Router			/user/delete [post]
func DeleteUserHandler(c *gin.Context) {
	ID := c.Request.Header.Get("user_id")
	userID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	var email dto.UserEmail
	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	loginService := Service.NewLoginService()
	err = loginService.DeleteUser(userID, email.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error deleting user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully."})
}
