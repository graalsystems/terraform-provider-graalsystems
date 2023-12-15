package graalsystems

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider_BuildApiWithCredentials(t *testing.T) {

	ctx := context.Background()
	api, err := buildApi(ctx, "http://172.27.0.1:4200/api/v1", "http://localhost:8089", "version", "platform-vincent-internal", "vdevillers", "devillerspwd", "", "", "")
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	assert.NotNil(t, api)
	projects, response, err := api.ProjectAPI.FindProjects(ctx).XTenant("platform-vincent-internal").Execute()
	fmt.Println("Err", err)
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response", response)
	if assert.NotNil(t, projects) {

	}
}
