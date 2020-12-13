package utils

type ContextKey string

const (
	UserIDKey       ContextKey = "userID"
	RestrictModeKey ContextKey = "restricted_mode"
	OrgRoleKey      ContextKey = "org_role"
)
