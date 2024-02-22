//go:build !selectTest || unitTest

package debug_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/debug"
)

func TestDebug(t *testing.T) {
	t.Run("Should initialize without messages.", func(t *testing.T) {
		debug.Init(true)
		assert.Equal(t, len(debug.GetAllMessages()), 0)
	})

	t.Run("Should create a new debug message.", func(t *testing.T) {
		debug.Init(true)
		newMessage := "Mock debug message"
		debug.NewMessage(newMessage)
		assert.Equal(t, len(debug.GetAllMessages()), 1)
		assert.Equal(t, debug.GetAllMessages()[0].Message, newMessage)
	})

	t.Run("Should not create any message if Debug is false.", func(t *testing.T) {
		debug.Init(false)
		debug.NewMessage("Mock debug message")
		assert.Equal(t, len(debug.GetAllMessages()), 0)
	})

	t.Run("Should get Json.", func(t *testing.T) {
		debug.Init(true)
		debug.NewMessage("Mock debug message")
		json, err := debug.GetJson()

		notEmpty := len(string(json)) > 0
		assert.Equal(t, notEmpty, true)
		assert.Nil(t, err)
	})
}
