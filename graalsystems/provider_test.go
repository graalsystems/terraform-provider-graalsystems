package graalsystems

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProvider_BuildApi(t *testing.T) {

	ctx := context.Background()
	api, err := buildApi(ctx, "http://172.24.240.1:4200/api/v1", "http://172.24.240.1:8089", "version", "platform-vincent-internal", "vdevillers", "devillerspwd")
	assert.Nil(t, err)
	projects, _, _ := api.ProjectApi.FindProjects(ctx).Execute()

	if assert.NotNil(t, projects) {

	}
}
