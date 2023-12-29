package admin

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ggrrrr/btmt-ui/be/common/cmd"
	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/waiter"
	auth "github.com/ggrrrr/btmt-ui/be/svc-auth"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

var AdminCmd = &cobra.Command{
	Use: "admin",
	// Aliases: []string{"insp"},
	Short: "admin ",
	// Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		w = waiter.New()
		err := runNewEmail()
		if err != nil {
			fmt.Printf("error %v\n", err)
		}
	},
}

var ListCmd = &cobra.Command{
	Use: "list",
	// Aliases: []string{"insp"},
	Short: "list all email ",
	// Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		w = waiter.New()

		err := runListEmail()
		if err != nil {
			fmt.Printf("error %v\n", err)
		}

	},
}

var (
	newEmail  string
	newPasswd string
	w         waiter.Waiter
)

func init() {
	AdminCmd.Flags().StringVarP(&newEmail, "email", "e", "", "new email")
	AdminCmd.Flags().StringVarP(&newPasswd, "passwd", "p", "", "new passwd")
}

func runNewEmail() error {
	defer func() {
		w.Wait()
		fmt.Println("Wait")
	}()
	defer func() {
		f := w.CancelFunc()
		f()
		fmt.Println("Cancel")
	}()
	ctx, app, err := prepCli()
	if err != nil {
		return err
	}
	err = app.UserCreate(ctx, ddd.AuthPasswd{
		Email:       newEmail,
		Passwd:      newPasswd,
		SystemRoles: []string{"admin"},
		Status:      ddd.StatusEnabled,
	})
	return err
}
func runUpdateEmail() error {
	defer func() {
		w.Wait()
		fmt.Println("Wait")
	}()
	defer func() {
		w.CancelFunc()()
		fmt.Println("Cancel")
	}()
	ctx, app, err := prepCli()
	if err != nil {
		return err
	}
	err = app.UserChangePasswd(ctx, newEmail, newPasswd, newPasswd)
	return err
}

func prepCli() (context.Context, app.App, error) {
	var cfg config.AppConfig
	logger.Init(logger.Config{
		Level:  cmd.GlobalFlags.LogLevel,
		Format: "console",
	})
	hostname, _ := os.Hostname()
	ctx := roles.CtxWithAuthInfo(context.Background(), roles.CreateAdminUser(
		"root",
		roles.Device{
			DeviceInfo: fmt.Sprintf("%v@%v", os.Getenv("USER"), hostname),
			RemoteAddr: "localhost",
		},
	))
	err := config.InitConfig(&cfg)
	if err != nil {
		return nil, nil, err
	}
	// InitApp/
	app, err := auth.InitApp(ctx, w, cfg)
	if err != nil {
		return nil, nil, err
	}
	return ctx, app, nil
}

func runListEmail() error {
	defer func() {
		w.Wait()
		fmt.Println("Wait")
	}()
	defer func() {
		f := w.CancelFunc()
		f()
		fmt.Println("Cancel")
	}()

	ctx, app, err := prepCli()
	if err != nil {
		return err
	}
	res, err := app.UserList(ctx)
	for _, a := range res.Payload() {
		fmt.Printf("%v \n", a)
	}
	return err
}
