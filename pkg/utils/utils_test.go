package utils_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/marcofranssen/terraform-provider-dexidp/pkg/utils"
)

func TestListStringValuesToSlice(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	tfList, _ := types.ListValueFrom(ctx, types.StringType, []string{"foo", "bar"})
	list := utils.ListStringValuesToSlice(tfList)

	assert.Len(list, 2)
	assert.ElementsMatch(list, []string{"foo", "bar"})
}
