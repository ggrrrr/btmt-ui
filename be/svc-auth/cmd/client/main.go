package client

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/authpb"
)

var (
	email  string
	passwd string

	serverAddr string
	// opts       []grpc.DialOption
)

var ClientCmd = &cobra.Command{
	Use: "client",
	// Aliases: []string{"insp"},
	Short: "client ",
	// Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	ClientCmd.Flags().StringVarP(&serverAddr, "address", "", "localhost:8081", "grpc address")
	ClientCmd.Flags().StringVarP(&email, "email", "e", "", "email")
	ClientCmd.Flags().StringVarP(&passwd, "password", "p", "", "password")

}

func run() {
	ctx := context.Background()
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	client := authpb.NewAuthSvcClient(conn)
	resp, err := client.LoginPasswd(context.Background(), &authpb.LoginPasswdRequest{
		Email:    email,
		Password: passwd,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp.Payload.Token)

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("%s %s", roles.AuthSchemeBeaerer, resp.Payload.Token))
	_, err = client.ValidateToken(ctx, &authpb.ValidateTokenRequest{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
