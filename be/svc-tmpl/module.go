package tmpl

import (
	"context"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/blob/awss3"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/rest"
)

type Module struct{}

var _ (system.Module) = (*Module)(nil)

func (*Module) Name() string {
	return "svc-tmpl"
}

func (*Module) Startup(ctx context.Context, s system.Service) (err error) {
	return Root(ctx, s)
}

func Root(ctx context.Context, s system.Service) error {
	logger.Info().Msg("svc-tmpl")

	blobClient, err := awss3.NewClient("test-bucket-1", awsclient.AwsConfig{
		Region:   "us-east-1",
		Endpoint: "http://localhost:4566",
	})
	if err != nil {
		return err
	}

	// fileName := "glass-mug-variant.png"
	// beerFile, err := os.Open(fileName)

	// pushInfo, err := blobClient.Push(s.Waiter().Context(), "localhost", "images/beer1",
	// 	&blob.BlobInfo{
	// 		Type:        "image",
	// 		ContentType: "image/png",
	// 		Name:        "beer.png",
	// 		Owner:       "me",
	// 	},
	// 	beerFile,
	// )
	// fmt.Printf("asdasdasdasd %+v \n", pushInfo)

	a, _ := app.New(app.WithBlobFetcher(blobClient))

	restApp := rest.New(a)
	s.Mux().Mount("/tmpl", restApp.Router())

	if s.Mux() == nil {
		return fmt.Errorf("system.Mux is nil")
	}
	return nil

}
