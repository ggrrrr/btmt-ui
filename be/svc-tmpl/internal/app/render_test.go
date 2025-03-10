package app

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	templv1 "github.com/ggrrrr/btmt-ui/be/common/templ/v1"
	tmplpb "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

type mockBlobStore struct {
}

// Fetch implements blob.Store.
func (m *mockBlobStore) Fetch(ctx context.Context, tenant string, blobId blob.BlobId) (blob.BlobReader, error) {
	return blob.BlobReader{}, nil
}

// Head implements blob.Store.
func (m *mockBlobStore) Head(ctx context.Context, tenant string, blobId blob.BlobId) (blob.BlobMD, error) {
	return blob.BlobMD{}, nil
}

// List implements blob.Store.
func (m *mockBlobStore) List(ctx context.Context, tenant string, blobId blob.BlobId) ([]blob.ListResult, error) {
	return []blob.ListResult{}, nil

}

// Push implements blob.Store.
func (m *mockBlobStore) Push(ctx context.Context, tenant string, blobId blob.BlobId, blobInfo blob.BlobMD, reader io.ReadSeeker) (blob.BlobId, error) {
	return blob.BlobId{}, nil
}

var _ (blob.Store) = (*mockBlobStore)(nil)

func TestRender(t *testing.T) {

	// TODO use this as test template body
	_ = `# Header from: {{ .Person.Name }}
# Items.key1: {{ .Items.key1.Item1 }} [{{ .Items.key1.Item }}]
# Lists.list1: {{range index .Lists "list1"}}
  * {{.}}{{end}}
---
# table: {{ .Tables.table1.Name }}
{{range .Tables.table1.Headers}}| {{ . }} {{end}} |
------------------------------------
{{range .Tables.table1.Rows}}{{ range .}}| {{ . }} {{ end}} | 
{{end}}
------------------------------------
end.

{{ renderImg "imageName" }}
`

	rootCtx := roles.CtxWithAuthInfo(context.Background(), roles.CreateSystemAdminUser("localhost", "admin", app.Device{}))

	testApp := &Application{
		tracer:       tracer.Tracer(otelScope),
		blobStore:    &mockBlobStore{},
		appPolices:   roles.NewAppPolices(),
		imagesFolder: blob.BlobId{},
		stateStore:   &state.MockStore{},
		repo:         nil,
	}
	mapData, err := structpb.NewStruct(map[string]any{"mapKey_1": "val 1"})
	require.NoError(t, err)

	actualHTml, err := testApp.RenderHtml(rootCtx, &tmplpb.RenderRequest{
		Data: &templv1.Data{
			Items: map[string]*structpb.Struct{
				"item_key_1": mapData,
			},
		},
		Body: `Hi,
some data {{ .Items.item_key_1.mapKey_1 }}
{{ renderImg "imageName" }}
`,
	})

	require.NoError(t, err)
	require.Equal(t, "Hi,\nsome data val 1\n<img src=\"http://localhost:8010/tmpl/image/imageName/resized\" ></img>\n", actualHTml)

	fmt.Printf("%#v \n", actualHTml)

}
