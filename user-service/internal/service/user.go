package service

import (
	"context"
	"log/slog"
	"time"

	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/pkg/security"
	"carsharing/user-service/internal/validation"

	"github.com/go-playground/validator/v10"
)

type UserService struct {
	log                   *slog.Logger
	validate              *validator.Validate
	userRepo              UserRepository
	docRepo               DocumentRepository
	objectStorage         ObjectStorage
	documentAnalyzer      DocumentAnalyzer
	publisher             Publisher
	activationCodeStorage ActivationCodeStorage
	mailer                Mailer
}

func NewUserService(
	log *slog.Logger,
	validate *validator.Validate,
	userRepo UserRepository,
	docRepo DocumentRepository,
	objectStorage ObjectStorage,
	documentAnalyzer DocumentAnalyzer,
	publisher Publisher,
	activationCodeStorage ActivationCodeStorage,
	mailer Mailer,
) *UserService {
	return &UserService{
		log:                   pkglog.WithComponent(log, "service.UserService"),
		validate:              validate,
		userRepo:              userRepo,
		docRepo:               docRepo,
		objectStorage:         objectStorage,
		documentAnalyzer:      documentAnalyzer,
		publisher:             publisher,
		activationCodeStorage: activationCodeStorage,
		mailer:                mailer,
	}
}

func (s *UserService) SignIn(ctx context.Context, creds validation.Credentials) (string, error) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(s.log, "SignIn"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, creds); err != nil {
		return "", err
	}

	filter := model.UserFilter{}
	if creds.Email != nil {
		filter.Email = creds.Email
	} else {
		filter.PhoneNumber = creds.PhoneNumber
	}

	user, err := s.userRepo.FindOne(ctx, filter)
	if err != nil {
		logger.Error("repo: finding user by id", pkglog.Err(err))

		return "", err
	}

	if err := security.CheckStringHash(creds.Password, user.PasswordHash); err != nil {
		logger.Warn("invalid user password hash", pkglog.Err(model.ErrBcrypt))

		return "", model.ErrUnauthenticated
	}

	return user.ID, nil
}

func (s *UserService) Create(ctx context.Context, data validation.UserCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Create"), utils.MetadataFromCtx(ctx))

	return s.insertUser(ctx, log, data)
}

func (s *UserService) Register(ctx context.Context, data validation.UserCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Register"), utils.MetadataFromCtx(ctx))

	return s.insertUser(ctx, log, data)
}

func (s *UserService) insertUser(ctx context.Context, log *slog.Logger, data validation.UserCreate) (string, error) {
	if err := validation.ValidateInput(s.validate, data); err != nil {
		return "", err
	}

	hash, err := security.HashString(data.Password)
	if err != nil {
		log.Error("hashing password", pkglog.Err(err))

		return "", model.ErrBcrypt
	}

	user := model.User{
		Email:        data.Email,
		PhoneNumber:  data.PhoneNumber,
		FirstName:    data.FirstName,
		LastName:     data.LastName,
		BirthDate:    data.BirthDate,
		PasswordHash: hash,
		Roles:        []sharedmodel.Role{sharedmodel.RoleUser},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	id, err := s.userRepo.Insert(ctx, user)
	if err != nil {
		log.Error("repo: inserting user", pkglog.Err(err))

		return "", err
	}

	if err := s.publisher.PublishUserCreated(ctx, id); err != nil {
		log.Error("event: publishing user created", pkglog.Err(err))
	}

	return id, nil
}

func (s *UserService) Get(ctx context.Context, id string) (model.User, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Get"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return model.User{}, err
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		log.Error("repo: finding user by id", pkglog.Err(err))

		return model.User{}, err
	}

	if err := s.resolveProfileImageURL(ctx, log, &user); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *UserService) List(ctx context.Context, filter validation.UserFilter) ([]model.User, error) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(s.log, "List"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, filter); err != nil {
		return nil, err
	}

	if filter.Pagination == nil {
		filter.Pagination = &sharedmodel.Pagination{
			Limit:  sharedmodel.DefaultPaginationLimit,
			Offset: sharedmodel.DefaultPaginationOffset,
		}
	}

	repoFilter := model.UserFilter{
		Email:              filter.Email,
		PhoneNumber:        filter.PhoneNumber,
		FirstName:          filter.FirstName,
		LastName:           filter.LastName,
		IsDocumentVerified: filter.IsDocumentVerified,
		IsEmailVerified:    filter.IsEmailVerified,
		IsSuspended:        filter.IsSuspended,
		Pagination:         filter.Pagination,
	}

	users, err := s.userRepo.Find(ctx, repoFilter)
	if err != nil {
		logger.Error("repo: listing users", pkglog.Err(err))

		return nil, err
	}

	for i := range users {
		if err := s.resolveProfileImageURL(ctx, logger, &users[i]); err != nil {
			return nil, err
		}
	}

	return users, nil
}

func (s *UserService) Update(ctx context.Context, id string, data validation.UserUpdate) error {
	logger := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Update"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return err
	}

	repoUpdate := model.UserUpdate{
		Email:              data.Email,
		PhoneNumber:        data.PhoneNumber,
		FirstName:          data.FirstName,
		LastName:           data.LastName,
		BirthDate:          data.BirthDate,
		ProfileImageKey:    data.ProfileImageKey,
		IsDocumentVerified: data.IsDocumentVerified,
		IsEmailVerified:    data.IsEmailVerified,
		IsSuspended:        data.IsSuspended,
		UpdatedAt:          time.Now(),
	}
	if data.Roles != nil {
		roles := make([]sharedmodel.Role, len(data.Roles))
		for i, r := range data.Roles {
			roles[i] = sharedmodel.Role(r)
		}
		repoUpdate.Roles = roles
	}

	if data.Password != nil {
		hash, err := security.HashString(*data.Password)
		if err != nil {
			logger.Error("hashing password", pkglog.Err(err))

			return model.ErrBcrypt
		}
		repoUpdate.PasswordHash = hash
	}

	if err := s.userRepo.Update(ctx, id, repoUpdate); err != nil {
		logger.Error("repo: updating user", pkglog.Err(err))

		return err
	}

	isSecurityUpdate := data.Password != nil || len(data.Roles) > 0 || data.IsSuspended != nil
	if err := s.publisher.PublishUserUpdated(ctx, id, isSecurityUpdate); err != nil {
		logger.Error("event: publishing user updated", pkglog.Err(err))
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	logger := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Delete"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		logger.Error("repo: deleting user", pkglog.Err(err))

		return err
	}

	if err := s.publisher.PublishUserDeleted(ctx, id); err != nil {
		logger.Error("event: publishing user deleted", pkglog.Err(err))
	}

	return nil
}

func (s *UserService) GetUserProfileImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetUserProfileImageUploadData"), utils.MetadataFromCtx(ctx))

	data, err := s.objectStorage.GetUserProfileImageUploadData(ctx)
	if err != nil {
		logger.Error("object storage: getting upload data", pkglog.Err(err))

		return sharedmodel.ImageUploadData{}, err
	}

	return data, nil
}

func (s *UserService) resolveProfileImageURL(ctx context.Context, log *slog.Logger, user *model.User) error {
	if user.ProfileImage.Key == "" {
		return nil
	}

	imageURL, err := s.objectStorage.GetImageURL(ctx, user.ProfileImage.Key)
	if err != nil {
		log.Error("object storage: getting image url", pkglog.Err(err))

		return err
	}
	user.ProfileImage.URL = imageURL

	return nil
}
