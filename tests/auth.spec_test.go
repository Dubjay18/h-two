package tests

import (
	"testing"
)

func TestAll(t *testing.T) {
	TestTokenGeneration(t)
	TestOrganizationAccessControl(t)
	TestCreateOrganizationHandler(t)
	TestRegisterUserWithDefaultOrganization(t)
	TestLoginUserSuccess(t)

}
