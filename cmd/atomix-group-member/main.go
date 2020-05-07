package main

import (
	"context"
	"fmt"
	atomix "github.com/atomix/go-client/pkg/client"
	"github.com/atomix/go-client/pkg/client/membership"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func main() {
	cmd := &cobra.Command{
		Use:  "atomix-group-member",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			member := args[0]
			test, _ := cmd.Flags().GetString("test")
			controller, _ := cmd.Flags().GetString("controller")
			namespace, _ := cmd.Flags().GetString("namespace")
			groupName, _ := cmd.Flags().GetString("group")
			client, err := atomix.New(
				controller,
				atomix.WithMemberID(member),
				atomix.WithNamespace(namespace),
				atomix.WithScope(test))
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			group, err := client.GetMembershipGroup(ctx, groupName)
			cancel()
			if err != nil {
				return err
			}
			ch := make(chan membership.Membership)
			err = group.Watch(context.Background(), ch)
			if err != nil {
				return err
			}
			for event := range ch {
				fmt.Println(event)
			}
			return nil
		},
	}
	cmd.Flags().StringP("controller", "c", "", "the controller address")
	cmd.Flags().StringP("namespace", "n", "", "the test namespace")
	cmd.Flags().StringP("test", "t", "", "the test name")
	cmd.Flags().StringP("group", "g", "", "the group name")
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
