package client

import (
	"context"
	"fmt"
	"os"

	"github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	ClientCmd = &cobra.Command{
		Use: "client",
		// Aliases: []string{"insp"},
		Short: "client ",
		// Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}

	serverAddr string
)

func init() {
	ClientCmd.Flags().StringVarP(&serverAddr, "address", "", "localhost:8081", "grpc address")
}

func run() {
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("%s %s", "mock", "admin"))

	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	client := peoplepb.NewPeopleSvcClient(conn)
	res, err := client.List(ctx, &peoplepb.ListRequest{
		Filters: map[string]*peoplepb.ListText{
			"texts": {
				List: []string{"asd"},
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(res)

}
