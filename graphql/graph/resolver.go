package graph

import "janjiss.com/rest/users"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserService users.UserService
}
