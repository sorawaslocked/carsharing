package auth

import (
	sharedmodel "carsharing/shared/model"

	usersvc "carsharing/protos/gen/service/user"
)

// ownerExtractFn extracts the target user ID from the request so the interceptor
// can check whether the caller is operating on their own resource.
type ownerExtractFn func(req any) (userID string, ok bool)

type methodPolicy struct {
	public       bool
	allowedRoles []sharedmodel.Role
	ownerExtract ownerExtractFn
}

// Duck-typing interfaces — avoids importing concrete proto request types here.
type idCarrier interface{ GetId() string }
type userIDCarrier interface{ GetUserId() string }

func extractByID(req any) (string, bool) {
	c, ok := req.(idCarrier)
	if !ok {
		return "", false
	}
	id := c.GetId()
	return id, id != ""
}

func extractByUserID(req any) (string, bool) {
	c, ok := req.(userIDCarrier)
	if !ok {
		return "", false
	}
	id := c.GetUserId()
	return id, id != ""
}

var privilegedRoles = []sharedmodel.Role{sharedmodel.RoleAdmin, sharedmodel.RoleUserManager}

func buildPolicies() map[string]methodPolicy {
	return map[string]methodPolicy{
		// Public — no authentication required.
		usersvc.HealthService_Health_FullMethodName: {public: true},
		usersvc.UserService_Register_FullMethodName: {public: true},
		usersvc.UserService_SignIn_FullMethodName:   {public: true},

		// Privileged roles only.
		usersvc.UserService_CreateUser_FullMethodName:    {allowedRoles: privilegedRoles},
		usersvc.UserService_ListUsers_FullMethodName:     {allowedRoles: privilegedRoles},
		usersvc.UserService_CheckDocument_FullMethodName: {allowedRoles: privilegedRoles},

		// Privileged roles OR the resource owner.
		usersvc.UserService_GetUser_FullMethodName:    {allowedRoles: privilegedRoles, ownerExtract: extractByID},
		usersvc.UserService_UpdateUser_FullMethodName: {allowedRoles: privilegedRoles, ownerExtract: extractByID},
		usersvc.UserService_DeleteUser_FullMethodName: {allowedRoles: privilegedRoles, ownerExtract: extractByID},
		usersvc.UserService_ListDocuments_FullMethodName: {
			allowedRoles: privilegedRoles,
			ownerExtract: extractByUserID,
		},

		// Any authenticated user.
		usersvc.UserService_SendActivationCode_FullMethodName:        {},
		usersvc.UserService_CheckActivationCode_FullMethodName:       {},
		usersvc.UserService_GetProfileImageUploadData_FullMethodName: {},
		usersvc.UserService_CreateDocument_FullMethodName:            {},
		usersvc.UserService_GetUploadDocumentData_FullMethodName:     {},
		usersvc.UserService_StreamDocumentAnalyzed_FullMethodName:    {},
	}
}
