package dto

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type User struct {
	ID              string  `json:"id"`
	Email           string  `json:"email"`
	PhoneNumber     *string `json:"phoneNumber,omitempty"`
	FirstName       string  `json:"firstName"`
	LastName        string  `json:"lastName"`
	BirthDate       string  `json:"birthDate"`
	PasswordHash    []byte  `json:"passwordHash"`
	ProfileImageURL *string `json:"profileImageURL"`

	Roles              []string `json:"roles"`
	IsDocumentVerified bool     `json:"isDocumentVerified"`
	IsEmailVerified    bool     `json:"isEmailVerified"`
	IsSuspended        bool     `json:"isSuspended"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserCreateRequest struct {
	Email                string  `json:"email"`
	PhoneNumber          *string `json:"phoneNumber"`
	FirstName            string  `json:"firstName"`
	LastName             string  `json:"lastName"`
	BirthDate            string  `json:"birthDate"`
	Password             string  `json:"password"`
	PasswordConfirmation string  `json:"passwordConfirmation"`
}

type UserUpdateRequest struct {
	Email                *string `json:"email"`
	PhoneNumber          *string `json:"phoneNumber"`
	FirstName            *string `json:"firstName"`
	LastName             *string `json:"lastName"`
	BirthDate            *string `json:"birthDate"`
	Password             *string `json:"password"`
	PasswordConfirmation *string `json:"passwordConfirmation"`

	Roles              []string `json:"roles"`
	IsDocumentVerified *bool    `json:"isDocumentVerified"`
	IsEmailVerified    *bool    `json:"isEmailVerified"`
	IsSuspended        *bool    `json:"isSuspended"`

	ProfileImageKey *string `json:"profileImageKey"`
}

type RegisterRequest struct {
	Email                string  `json:"email"`
	PhoneNumber          *string `json:"phoneNumber"`
	FirstName            string  `json:"firstName"`
	LastName             string  `json:"lastName"`
	BirthDate            string  `json:"birthDate"`
	Password             string  `json:"password"`
	PasswordConfirmation string  `json:"passwordConfirmation"`
}

type LoginRequest struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
	Password    string  `json:"password"`
}

type Token struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expiresIn"`
}

type CheckActivationCodeRequest struct {
	Code string `json:"code"`
}

type Document struct {
	ID        string  `json:"id"`
	UserID    string  `json:"userID"`
	ImageType string  `json:"imageType"`
	Status    string  `json:"status"`
	Reason    *string `json:"reason,omitempty"`
	ImageURL  string  `json:"imageURL"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetUploadDocumentDataRequest struct {
	ImageType string `json:"imageType"`
}

type CreateDocumentRequest struct {
	ImageType string `json:"imageType"`
	ObjectKey string `json:"objectKey"`
}

type CheckDocumentRequest struct {
	Status string  `json:"status"`
	Error  *string `json:"error"`
}

func ToUserResponse(m model.User) User {
	return User{
		ID:                 m.ID,
		Email:              m.Email,
		PhoneNumber:        m.PhoneNumber,
		FirstName:          m.FirstName,
		LastName:           m.LastName,
		BirthDate:          m.BirthDate,
		PasswordHash:       m.Password.Hash,
		ProfileImageURL:    m.ProfileImageURL,
		Roles:              m.Roles,
		IsDocumentVerified: m.IsDocumentVerified,
		IsEmailVerified:    m.IsEmailVerified,
		IsSuspended:        m.IsSuspended,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func FromCreateUserRequest(c *gin.Context) (model.UserCreate, error) {
	var req UserCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		return model.UserCreate{}, err
	}

	data := model.UserCreate{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		BirthDate:   req.BirthDate,
		Password: model.Password{
			Text:             &req.Password,
			TextConfirmation: &req.PasswordConfirmation,
		},
	}

	return data, nil
}

func FromUpdateRequest(c *gin.Context) (model.UserUpdate, error) {
	var req UserUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		return model.UserUpdate{}, err
	}

	return model.UserUpdate{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		BirthDate:   req.BirthDate,
		Password: model.Password{
			Text:             req.Password,
			TextConfirmation: req.PasswordConfirmation,
		},
		Roles:              req.Roles,
		IsDocumentVerified: req.IsDocumentVerified,
		IsEmailVerified:    req.IsEmailVerified,
		IsSuspended:        req.IsSuspended,
		ProfileImageKey:    req.ProfileImageKey,
	}, nil
}

func UserFilterFromCtx(c *gin.Context) (model.UserFilter, error) {
	f := model.UserFilter{}

	if v := c.Query("email"); v != "" {
		f.Email = &v
	}
	if v := c.Query("phoneNumber"); v != "" {
		f.PhoneNumber = &v
	}
	if v := c.Query("firstName"); v != "" {
		f.FirstName = &v
	}
	if v := c.Query("lastName"); v != "" {
		f.LastName = &v
	}
	if v := c.Query("isDocumentVerified"); v != "" {
		vBool, err := strconv.ParseBool(v)
		if err != nil {
			return model.UserFilter{}, model.ErrInvalidQueryParam
		}

		f.IsDocumentVerified = &vBool
	}
	if v := c.Query("isEmailVerified"); v != "" {
		vBool, err := strconv.ParseBool(v)
		if err != nil {
			return model.UserFilter{}, model.ErrInvalidQueryParam
		}

		f.IsEmailVerified = &vBool
	}
	if v := c.Query("isSuspended"); v != "" {
		vBool, err := strconv.ParseBool(v)
		if err != nil {
			return model.UserFilter{}, model.ErrInvalidQueryParam
		}

		f.IsSuspended = &vBool
	}

	p, err := pagination(c)
	if err != nil {
		return model.UserFilter{}, model.ErrInvalidQueryParam
	}

	f.Pagination = p

	return f, nil
}

func FromRegisterRequest(ctx *gin.Context) (model.UserCreate, error) {
	var req RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.UserCreate{}, err
	}

	return model.UserCreate{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		BirthDate:   req.BirthDate,
		Password: model.Password{
			Text:             &req.Password,
			TextConfirmation: &req.PasswordConfirmation,
		},
	}, nil
}

func FromLoginRequest(ctx *gin.Context) (model.Credentials, error) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.Credentials{}, err
	}

	cred := model.Credentials{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password: model.Password{
			Text: &req.Password,
		},
	}

	return cred, nil
}

func FromCheckActivationCodeRequest(c *gin.Context) (code string, err error) {
	var req CheckActivationCodeRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		return "", err
	}

	return req.Code, nil
}

func ToDocumentResponse(m model.Document) Document {
	return Document{
		ID:        m.ID,
		UserID:    m.UserID,
		ImageType: m.ImageType,
		Status:    m.Status,
		Reason:    m.Reason,
		ImageURL:  m.ImageURL,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FromCreateDocumentRequest(c *gin.Context) (imageType, objectKey string, err error) {
	var req CreateDocumentRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		return "", "", err
	}

	return req.ImageType, req.ObjectKey, nil
}

func FromGetUploadDocumentDataRequest(c *gin.Context) (imageType string, err error) {
	var req GetUploadDocumentDataRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		return "", err
	}

	return req.ImageType, nil
}

func FromCheckDocumentRequest(c *gin.Context) (status string, documentError *string, err error) {
	var req CheckDocumentRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		return "", nil, nil
	}

	return req.Status, req.Error, nil
}
