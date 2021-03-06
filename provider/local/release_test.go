package local_test

import (
	"testing"

	"github.com/convox/praxis/types"
	"github.com/stretchr/testify/assert"
)

func TestReleaseCreateGet(t *testing.T) {
	p, err := testProvider()
	assert.NoError(t, err)
	defer testProviderCleanup(p)

	_, err = p.AppCreate("app")
	assert.NoError(t, err)

	opts := types.ReleaseCreateOptions{
		Build: "BTEST",
		Env: map[string]string{
			"APP": "app",
			"FOO": "bar",
		},
	}
	rel, err := p.ReleaseCreate("app", opts)
	assert.NoError(t, err)

	if assert.NotNil(t, rel) {
		fetched, err := p.ReleaseGet("app", rel.Id)
		assert.NoError(t, err)

		assert.EqualValues(t, rel, fetched)
	}
}

func TestReleaseList(t *testing.T) {
	p, err := testProvider()
	assert.NoError(t, err)
	defer testProviderCleanup(p)

	_, err = p.AppCreate("app")
	if !assert.NoError(t, err) {
		return
	}

	p.ReleaseCreate("app", types.ReleaseCreateOptions{Build: "B1"})
	p.ReleaseCreate("app", types.ReleaseCreateOptions{Env: map[string]string{"FOO": "bar"}})
	p.ReleaseCreate("app", types.ReleaseCreateOptions{Build: "B2"})
	p.ReleaseCreate("app", types.ReleaseCreateOptions{Build: "B3"})
	p.ReleaseCreate("app", types.ReleaseCreateOptions{Env: map[string]string{"FOO": "baz"}})
	p.ReleaseCreate("app", types.ReleaseCreateOptions{Build: "B4"})

	rs, err := p.ReleaseList("app", types.ReleaseListOptions{})

	if assert.NoError(t, err) && assert.Len(t, rs, 6) {
		assert.Equal(t, "B4", rs[0].Build)
		assert.Equal(t, map[string]string{"FOO": "baz"}, rs[0].Env)
		assert.Equal(t, "B3", rs[1].Build)
		assert.Equal(t, map[string]string{"FOO": "baz"}, rs[1].Env)
		assert.Equal(t, "B3", rs[2].Build)
		assert.Equal(t, map[string]string{"FOO": "bar"}, rs[2].Env)
		assert.Equal(t, "B2", rs[3].Build)
		assert.Equal(t, map[string]string{"FOO": "bar"}, rs[3].Env)
		assert.Equal(t, "B1", rs[4].Build)
		assert.Equal(t, map[string]string{"FOO": "bar"}, rs[4].Env)
		assert.Equal(t, "B1", rs[5].Build)
		assert.Equal(t, map[string]string(nil), rs[5].Env)
	}
}
