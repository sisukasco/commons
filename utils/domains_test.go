package utils_test

import (
	"testing"

	"github.com/sisukasco/commons/utils"

	"github.com/stretchr/testify/assert"
)

func TestMatchDomains(t *testing.T) {
	assert.True(t, utils.MatchDomain("localhost", "localhost"))

	assert.True(t, utils.MatchDomain("localhost,some.domain.com", "localhost"))

	assert.False(t, utils.MatchDomain("localhost,some.domain.com", "domain.com"))

	assert.True(t, utils.MatchDomain("*.domain.com", "sub.domain.com"))

	assert.False(t, utils.MatchDomain("*.domain.com", "domain.com"))

	assert.False(t, utils.MatchDomain("*.domain.com", "ww.tt.domain.com"))

	assert.True(t, utils.MatchDomain("*.domain.com,*.*.domain.com", "ww.tt.domain.com"))

	assert.True(t, utils.MatchDomain("www.domain.com,domain.com,other-domain.com", "domain.com"))
}
