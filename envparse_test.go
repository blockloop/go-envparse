package envparse

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldIgnoreCommentLines(t *testing.T) {
	str := `# HELLO=world`
	assert.NoError(t, ParseString(str))
	assert.Equal(t, "", os.Getenv("HELLO"))
}

func TestShouldSet(t *testing.T) {
	str := `HELLO=world`
	assert.NoError(t, ParseString(str))
	assert.Equal(t, "world", os.Getenv("HELLO"))
}

func TestShouldTrimWhitespaces(t *testing.T) {
	str := `   HELLO=world`
	assert.NoError(t, ParseString(str))
	assert.Equal(t, "world", os.Getenv("HELLO"))
}

func TestShouldIgnoreEmptyLines(t *testing.T) {
	str := `
    HELLO=world
    WORLD=hello


    `

	assert.NoError(t, ParseString(str))
	assert.Equal(t, "world", os.Getenv("HELLO"))
	assert.Equal(t, "hello", os.Getenv("WORLD"))
}

func TestShouldErrorIfLineFormatIsInvalid(t *testing.T) {
	str := `HELLO`
	assert.Error(t, ParseString(str))
}

func TestShouldSplitOnFirstEqSign(t *testing.T) {
	str := `HELLO=hello=world`
	assert.NoError(t, ParseString(str))
	assert.Equal(t, "hello=world", os.Getenv("HELLO"))
}

func TestShouldErrorIfFileDoesNotExist(t *testing.T) {
	f := "not/a/real/file/path"
	assert.Error(t, ParseFile(f))
}
