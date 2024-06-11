package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_HasRole(t *testing.T) {
	t.Parallel()

	doTestHasRole(t, 0, 0, true)
	doTestHasRole(t, 0, 1, false)
	doTestHasRole(t, 0, 5, false)

	doTestHasRole(t, 3, 0, true)
	doTestHasRole(t, 3, 1, true)
	doTestHasRole(t, 3, 2, true)
	doTestHasRole(t, 3, 3, true)
	doTestHasRole(t, 3, 4, false)
	doTestHasRole(t, 3, 5, false)
}

func doTestHasRole(t *testing.T, role int, testRole int, hasRole bool) {
	user := User{Role: role}
	assert.Equal(t, user.HasRole(testRole), hasRole)
}
