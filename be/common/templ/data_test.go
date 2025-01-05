package templ

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	templv1 "github.com/ggrrrr/btmt-ui/be/common/templ/v1"
)

func TestConverter(t *testing.T) {
	tests := []struct {
		name  string
		prepF func(t *testing.T) *templv1.Data
		data  templData
	}{
		{
			name: "1",
			prepF: func(t *testing.T) *templv1.Data {
				items, err := structpb.NewStruct(map[string]any{"k1": "v1"})
				require.NoError(t, err)
				return &templv1.Data{
					Items: map[string]*structpb.Struct{
						"item1": items,
					},
				}
			},
			data: templData{
				Items: map[string]any{
					"item1": map[string]any{
						"k1": "v1",
					},
				},
			},
		},
		{
			name: "2",
			prepF: func(t *testing.T) *templv1.Data {
				return &templv1.Data{}
			},
			data: templData{Items: map[string]any{}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fromData := tc.prepF(t)
			actual := fromV1(fromData)
			assert.Equal(t, tc.data, actual)

		})
	}
}
