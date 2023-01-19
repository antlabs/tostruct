// antlabs, guonaihong 2022
// apache 2.0
package guesstype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试float
func TestIsFloat_Fail(t *testing.T) {
	// 反例
	for _, tc := range []string{
		"1",
		"-1",
		"0",
	} {
		assert.False(t, IsFloat(tc))
	}
}

// 测试float
func TestIsFloat(t *testing.T) {
	// 反例
	for _, tc := range []string{
		"1.1",
		"-1.1",
		"0.0",
	} {
		assert.True(t, IsFloat(tc))
	}
}

// 测试int
func TestIsInt_Fail(t *testing.T) {
	// 反例
	for _, tc := range []string{
		"1.1",
		"-1.1",
		"0.0",
	} {
		assert.False(t, IsInt(tc))
	}
}

// 测试int
func TestIsInt(t *testing.T) {
	// 正例
	for _, tc := range []string{
		"1",
		"-1",
		"0",
	} {
		assert.True(t, IsInt(tc))
	}
}

// 测试bool
func TestIsBool_fail(t *testing.T) {
	// 正例
	for _, tc := range []string{
		"hello",
		"1.1",
	} {
		assert.False(t, IsBool(tc))
	}
}

// 测试bool
func TestIsBool(t *testing.T) {
	// 正例
	for _, tc := range []string{
		"true",
		"false",
	} {
		assert.True(t, IsBool(tc))
	}
}

type testCase struct {
	data string
	need string
}

// 测试TypeOf函数
func TestTypeOf(t *testing.T) {
	for _, tc := range []testCase{
		{"1.1", "float64"},
		{"-1.1", "float64"},
		{"0.0", "float64"},
		{"0", "int"},
		{"1", "int"},
		{"-1", "int"},
		{"true", "bool"},
		{"false", "bool"},
	} {
		assert.Equal(t, TypeOf(tc.data), tc.need)
	}
}
