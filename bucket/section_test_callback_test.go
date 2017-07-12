package bucket

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testKey = "testKey"
)

func Test_sectionsTestRegistry(t *testing.T) {
	assert.NotNil(t, GetSectionTestCallback(SectionTestTrue))
	assert.NotNil(t, GetSectionTestCallback(SectionTestIsNumeric))
	assert.NotNil(t, GetSectionTestCallback(SectionTestIsNotEmpty))

	assert.Nil(t, GetSectionTestCallback(testKey))

	RegisterSectionTest(testKey, func(PathSection) bool {
		return false
	})
	assert.NotNil(t, GetSectionTestCallback(testKey))
}

func parseAndAssertParseResults(t *testing.T, s string) {
	m, err := ParseSectionsTestsMap(s)
	require.NoError(t, err)
	assert.NotEmpty(t, reflect.ValueOf(m).MapKeys())

	callback, ok := m["foo"]
	assert.True(t, ok)
	assert.NotNil(t, callback)
	assert.True(t, callback.Callback("foo"))

	callback, ok = m["bar"]
	assert.True(t, ok)
	assert.NotNil(t, callback)
	assert.True(t, callback.Callback("12"))
	assert.False(t, callback.Callback("~12"))

	callback, ok = m["baz"]
	assert.True(t, ok)
	assert.NotNil(t, callback)
	assert.True(t, callback.Callback("12"))
	assert.True(t, callback.Callback("~12"))
	assert.False(t, callback.Callback("-"))

	callback, ok = m["qaz"]
	assert.False(t, ok)
}

func TestRegisterSectionTest(t *testing.T) {
	m, err := ParseSectionsTestsMap("")
	assert.Nil(t, err)
	assert.Empty(t, reflect.ValueOf(m).MapKeys())

	parseAndAssertParseResults(t, "foo:true:bar:numeric:baz:not_empty")
	parseAndAssertParseResults(t, "foo:true\nbar:numeric:baz:not_empty")
	parseAndAssertParseResults(t, "\nfoo:true\nbar:numeric:baz:not_empty")
	parseAndAssertParseResults(t, "\nfoo:true:bar:numeric\nbaz:not_empty")
	parseAndAssertParseResults(t, "\nfoo:true\nbar:numeric\nbaz:not_empty")
	parseAndAssertParseResults(t, "\nfoo:true\nbar:numeric\nbaz:not_empty\n")
}

func TestRegisterSectionTest_ErrInvalidFormat(t *testing.T) {
	_, err := ParseSectionsTestsMap("foo")
	assert.Equal(t, err, ErrInvalidFormat)

	_, err = ParseSectionsTestsMap("foo:bar:baz")
	assert.Equal(t, err, ErrInvalidFormat)
}

func TestRegisterSectionTest_ErrUnknownSectionTest(t *testing.T) {
	_, err := ParseSectionsTestsMap("foo:NOT_EISTS")
	assert.Equal(t, err, ErrUnknownSectionTest)

	_, err = ParseSectionsTestsMap("foo:true:baz:NOT_EISTS")
	assert.Equal(t, err, ErrUnknownSectionTest)
}

func TestSectionsTestsMap_String(t *testing.T) {
	m, err := ParseSectionsTestsMap("foo:true:bar:numeric:baz:not_empty")
	require.NoError(t, err)

	assert.Equal(t, "[bar: numeric, baz: not_empty, foo: true]", m.String())
}
