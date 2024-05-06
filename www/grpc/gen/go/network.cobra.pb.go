// Code generated by protoc-gen-cobra. DO NOT EDIT.

package pactus

import (
	client "github.com/NathanBaulch/protoc-gen-cobra/client"
	flag "github.com/NathanBaulch/protoc-gen-cobra/flag"
	iocodec "github.com/NathanBaulch/protoc-gen-cobra/iocodec"
	cobra "github.com/spf13/cobra"
	grpc "google.golang.org/grpc"
	proto "google.golang.org/protobuf/proto"
)

func NetworkClientCommand(options ...client.Option) *cobra.Command {
	cfg := client.NewConfig(options...)
	cmd := &cobra.Command{
		Use:   cfg.CommandNamer("Network"),
		Short: "Network service client",
		Long:  "Network service provides RPCs for retrieving information about the network.",
	}
	cfg.BindFlags(cmd.PersistentFlags())
	cmd.AddCommand(
		_NetworkGetNetworkInfoCommand(cfg),
		_NetworkGetNodeInfoCommand(cfg),
	)
	return cmd
}

func _NetworkGetNetworkInfoCommand(cfg *client.Config) *cobra.Command {
	req := &GetNetworkInfoRequest{}

	cmd := &cobra.Command{
		Use:   cfg.CommandNamer("GetNetworkInfo"),
		Short: "GetNetworkInfo RPC client",
		Long:  "GetNetworkInfo retrieves information about the overall network.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.UseEnvVars {
				if err := flag.SetFlagsFromEnv(cmd.Parent().PersistentFlags(), true, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Network"); err != nil {
					return err
				}
				if err := flag.SetFlagsFromEnv(cmd.PersistentFlags(), false, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Network", "GetNetworkInfo"); err != nil {
					return err
				}
			}
			return client.RoundTrip(cmd.Context(), cfg, func(cc grpc.ClientConnInterface, in iocodec.Decoder, out iocodec.Encoder) error {
				cli := NewNetworkClient(cc)
				v := &GetNetworkInfoRequest{}

				if err := in(v); err != nil {
					return err
				}
				proto.Merge(v, req)

				res, err := cli.GetNetworkInfo(cmd.Context(), v)

				if err != nil {
					return err
				}

				return out(res)

			})
		},
	}

	cmd.PersistentFlags().BoolVar(&req.OnlyOnline, cfg.FlagNamer("OnlyOnline"), false, "Only returns the peers with online status")

	return cmd
}

func _NetworkGetNodeInfoCommand(cfg *client.Config) *cobra.Command {
	req := &GetNodeInfoRequest{}

	cmd := &cobra.Command{
		Use:   cfg.CommandNamer("GetNodeInfo"),
		Short: "GetNodeInfo RPC client",
		Long:  "GetNodeInfo retrieves information about a specific node in the network.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.UseEnvVars {
				if err := flag.SetFlagsFromEnv(cmd.Parent().PersistentFlags(), true, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Network"); err != nil {
					return err
				}
				if err := flag.SetFlagsFromEnv(cmd.PersistentFlags(), false, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Network", "GetNodeInfo"); err != nil {
					return err
				}
			}
			return client.RoundTrip(cmd.Context(), cfg, func(cc grpc.ClientConnInterface, in iocodec.Decoder, out iocodec.Encoder) error {
				cli := NewNetworkClient(cc)
				v := &GetNodeInfoRequest{}

				if err := in(v); err != nil {
					return err
				}
				proto.Merge(v, req)

				res, err := cli.GetNodeInfo(cmd.Context(), v)

				if err != nil {
					return err
				}

				return out(res)

			})
		},
	}

	return cmd
}
