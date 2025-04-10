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

func BlockchainClientCommand(options ...client.Option) *cobra.Command {
	cfg := client.NewConfig(options...)
	cmd := &cobra.Command{
		Use:   cfg.CommandNamer("Blockchain"),
		Short: "Blockchain service client",
		Long:  "Blockchain service defines RPC methods for interacting with the blockchain.",
	}
	cfg.BindFlags(cmd.PersistentFlags())
	cmd.AddCommand(
		_BlockchainGetBlockCommand(cfg),
		_BlockchainGetBlockHashCommand(cfg),
		_BlockchainGetBlockHeightCommand(cfg),
		_BlockchainGetBlockchainInfoCommand(cfg),
		_BlockchainGetConsensusInfoCommand(cfg),
		_BlockchainGetAccountCommand(cfg),
		_BlockchainGetValidatorCommand(cfg),
		_BlockchainGetValidatorByNumberCommand(cfg),
		_BlockchainGetValidatorAddressesCommand(cfg),
		_BlockchainGetPublicKeyCommand(cfg),
		_BlockchainGetTxPoolContentCommand(cfg),
	)
	return cmd
}

func _BlockchainGetBlockCommand(cfg *client.Config) *cobra.Command {
	req := &GetBlockRequest{}

	cmd := &cobra.Command{
		Use:   cfg.CommandNamer("GetBlock"),
		Short: "GetBlock RPC client",
		Long:  "GetBlock retrieves information about a block based on the provided request parameters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.UseEnvVars {
				if err := flag.SetFlagsFromEnv(cmd.Parent().PersistentFlags(), true, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain"); err != nil {
					return err
				}
				if err := flag.SetFlagsFromEnv(cmd.PersistentFlags(), false, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain", "GetBlock"); err != nil {
					return err
				}
			}
			return client.RoundTrip(cmd.Context(), cfg, func(cc grpc.ClientConnInterface, in iocodec.Decoder, out iocodec.Encoder) error {
				cli := NewBlockchainClient(cc)
				v := &GetBlockRequest{}

				if err := in(v); err != nil {
					return err
				}
				proto.Merge(v, req)

				res, err := cli.GetBlock(cmd.Context(), v)

				if err != nil {
					return err
				}

				return out(res)

			})
		},
	}

	cmd.PersistentFlags().Uint32Var(&req.Height, cfg.FlagNamer("Height"), 0, "The height of the block to retrieve.")
	flag.EnumVar(cmd.PersistentFlags(), &req.Verbosity, cfg.FlagNamer("Verbosity"), "The verbosity level for block information.")

	return cmd
}

func _BlockchainGetBlockHashCommand(cfg *client.Config) *cobra.Command {
	req := &GetBlockHashRequest{}

	cmd := &cobra.Command{
		Use:   cfg.CommandNamer("GetBlockHash"),
		Short: "GetBlockHash RPC client",
		Long:  "GetBlockHash retrieves the hash of a block at the specified height.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.UseEnvVars {
				if err := flag.SetFlagsFromEnv(cmd.Parent().PersistentFlags(), true, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain"); err != nil {
					return err
				}
				if err := flag.SetFlagsFromEnv(cmd.PersistentFlags(), false, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain", "GetBlockHash"); err != nil {
					return err
				}
			}
			return client.RoundTrip(cmd.Context(), cfg, func(cc grpc.ClientConnInterface, in iocodec.Decoder, out iocodec.Encoder) error {
				cli := NewBlockchainClient(cc)
				v := &GetBlockHashRequest{}

				if err := in(v); err != nil {
					return err
				}
				proto.Merge(v, req)

				res, err := cli.GetBlockHash(cmd.Context(), v)

				if err != nil {
					return err
				}

				return out(res)

			})
		},
	}

	cmd.PersistentFlags().Uint32Var(&req.Height, cfg.FlagNamer("Height"), 0, "The height of the block to retrieve the hash for.")

	return cmd
}

func _BlockchainGetBlockHeightCommand(cfg *client.Config) *cobra.Command {
	req := &GetBlockHeightRequest{}

	cmd := &cobra.Command{
		Use:   cfg.CommandNamer("GetBlockHeight"),
		Short: "GetBlockHeight RPC client",
		Long:  "GetBlockHeight retrieves the height of a block with the specified hash.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.UseEnvVars {
				if err := flag.SetFlagsFromEnv(cmd.Parent().PersistentFlags(), true, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain"); err != nil {
					return err
				}
				if err := flag.SetFlagsFromEnv(cmd.PersistentFlags(), false, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain", "GetBlockHeight"); err != nil {
					return err
				}
			}
			return client.RoundTrip(cmd.Context(), cfg, func(cc grpc.ClientConnInterface, in iocodec.Decoder, out iocodec.Encoder) error {
				cli := NewBlockchainClient(cc)
				v := &GetBlockHeightRequest{}

				if err := in(v); err != nil {
					return err
				}
				proto.Merge(v, req)

				res, err := cli.GetBlockHeight(cmd.Context(), v)

				if err != nil {
					return err
				}

				return out(res)

			})
		},
	}

	cmd.PersistentFlags().StringVar(&req.Hash, cfg.FlagNamer("Hash"), "", "The hash of the block to retrieve the height for.")

	return cmd
}

func _BlockchainGetBlockchainInfoCommand(cfg *client.Config) *cobra.Command {
	req := &GetBlockchainInfoRequest{}

	cmd := &cobra.Command{
		Use:   cfg.CommandNamer("GetBlockchainInfo"),
		Short: "GetBlockchainInfo RPC client",
		Long:  "GetBlockchainInfo retrieves general information about the blockchain.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.UseEnvVars {
				if err := flag.SetFlagsFromEnv(cmd.Parent().PersistentFlags(), true, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain"); err != nil {
					return err
				}
				if err := flag.SetFlagsFromEnv(cmd.PersistentFlags(), false, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain", "GetBlockchainInfo"); err != nil {
					return err
				}
			}
			return client.RoundTrip(cmd.Context(), cfg, func(cc grpc.ClientConnInterface, in iocodec.Decoder, out iocodec.Encoder) error {
				cli := NewBlockchainClient(cc)
				v := &GetBlockchainInfoRequest{}

				if err := in(v); err != nil {
					return err
				}
				proto.Merge(v, req)

				res, err := cli.GetBlockchainInfo(cmd.Context(), v)

				if err != nil {
					return err
				}

				return out(res)

			})
		},
	}

	return cmd
}

func _BlockchainGetConsensusInfoCommand(cfg *client.Config) *cobra.Command {
	req := &GetConsensusInfoRequest{}

	cmd := &cobra.Command{
		Use:   cfg.CommandNamer("GetConsensusInfo"),
		Short: "GetConsensusInfo RPC client",
		Long:  "GetConsensusInfo retrieves information about the consensus instances.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.UseEnvVars {
				if err := flag.SetFlagsFromEnv(cmd.Parent().PersistentFlags(), true, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain"); err != nil {
					return err
				}
				if err := flag.SetFlagsFromEnv(cmd.PersistentFlags(), false, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain", "GetConsensusInfo"); err != nil {
					return err
				}
			}
			return client.RoundTrip(cmd.Context(), cfg, func(cc grpc.ClientConnInterface, in iocodec.Decoder, out iocodec.Encoder) error {
				cli := NewBlockchainClient(cc)
				v := &GetConsensusInfoRequest{}

				if err := in(v); err != nil {
					return err
				}
				proto.Merge(v, req)

				res, err := cli.GetConsensusInfo(cmd.Context(), v)

				if err != nil {
					return err
				}

				return out(res)

			})
		},
	}

	return cmd
}

func _BlockchainGetAccountCommand(cfg *client.Config) *cobra.Command {
	req := &GetAccountRequest{}

	cmd := &cobra.Command{
		Use:   cfg.CommandNamer("GetAccount"),
		Short: "GetAccount RPC client",
		Long:  "GetAccount retrieves information about an account based on the provided address.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.UseEnvVars {
				if err := flag.SetFlagsFromEnv(cmd.Parent().PersistentFlags(), true, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain"); err != nil {
					return err
				}
				if err := flag.SetFlagsFromEnv(cmd.PersistentFlags(), false, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain", "GetAccount"); err != nil {
					return err
				}
			}
			return client.RoundTrip(cmd.Context(), cfg, func(cc grpc.ClientConnInterface, in iocodec.Decoder, out iocodec.Encoder) error {
				cli := NewBlockchainClient(cc)
				v := &GetAccountRequest{}

				if err := in(v); err != nil {
					return err
				}
				proto.Merge(v, req)

				res, err := cli.GetAccount(cmd.Context(), v)

				if err != nil {
					return err
				}

				return out(res)

			})
		},
	}

	cmd.PersistentFlags().StringVar(&req.Address, cfg.FlagNamer("Address"), "", "The address of the account to retrieve information for.")

	return cmd
}

func _BlockchainGetValidatorCommand(cfg *client.Config) *cobra.Command {
	req := &GetValidatorRequest{}

	cmd := &cobra.Command{
		Use:   cfg.CommandNamer("GetValidator"),
		Short: "GetValidator RPC client",
		Long:  "GetValidator retrieves information about a validator based on the provided address.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.UseEnvVars {
				if err := flag.SetFlagsFromEnv(cmd.Parent().PersistentFlags(), true, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain"); err != nil {
					return err
				}
				if err := flag.SetFlagsFromEnv(cmd.PersistentFlags(), false, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain", "GetValidator"); err != nil {
					return err
				}
			}
			return client.RoundTrip(cmd.Context(), cfg, func(cc grpc.ClientConnInterface, in iocodec.Decoder, out iocodec.Encoder) error {
				cli := NewBlockchainClient(cc)
				v := &GetValidatorRequest{}

				if err := in(v); err != nil {
					return err
				}
				proto.Merge(v, req)

				res, err := cli.GetValidator(cmd.Context(), v)

				if err != nil {
					return err
				}

				return out(res)

			})
		},
	}

	cmd.PersistentFlags().StringVar(&req.Address, cfg.FlagNamer("Address"), "", "The address of the validator to retrieve information for.")

	return cmd
}

func _BlockchainGetValidatorByNumberCommand(cfg *client.Config) *cobra.Command {
	req := &GetValidatorByNumberRequest{}

	cmd := &cobra.Command{
		Use:   cfg.CommandNamer("GetValidatorByNumber"),
		Short: "GetValidatorByNumber RPC client",
		Long:  "GetValidatorByNumber retrieves information about a validator based on the provided number.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.UseEnvVars {
				if err := flag.SetFlagsFromEnv(cmd.Parent().PersistentFlags(), true, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain"); err != nil {
					return err
				}
				if err := flag.SetFlagsFromEnv(cmd.PersistentFlags(), false, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain", "GetValidatorByNumber"); err != nil {
					return err
				}
			}
			return client.RoundTrip(cmd.Context(), cfg, func(cc grpc.ClientConnInterface, in iocodec.Decoder, out iocodec.Encoder) error {
				cli := NewBlockchainClient(cc)
				v := &GetValidatorByNumberRequest{}

				if err := in(v); err != nil {
					return err
				}
				proto.Merge(v, req)

				res, err := cli.GetValidatorByNumber(cmd.Context(), v)

				if err != nil {
					return err
				}

				return out(res)

			})
		},
	}

	cmd.PersistentFlags().Int32Var(&req.Number, cfg.FlagNamer("Number"), 0, "The unique number of the validator to retrieve information for.")

	return cmd
}

func _BlockchainGetValidatorAddressesCommand(cfg *client.Config) *cobra.Command {
	req := &GetValidatorAddressesRequest{}

	cmd := &cobra.Command{
		Use:   cfg.CommandNamer("GetValidatorAddresses"),
		Short: "GetValidatorAddresses RPC client",
		Long:  "GetValidatorAddresses retrieves a list of all validator addresses.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.UseEnvVars {
				if err := flag.SetFlagsFromEnv(cmd.Parent().PersistentFlags(), true, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain"); err != nil {
					return err
				}
				if err := flag.SetFlagsFromEnv(cmd.PersistentFlags(), false, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain", "GetValidatorAddresses"); err != nil {
					return err
				}
			}
			return client.RoundTrip(cmd.Context(), cfg, func(cc grpc.ClientConnInterface, in iocodec.Decoder, out iocodec.Encoder) error {
				cli := NewBlockchainClient(cc)
				v := &GetValidatorAddressesRequest{}

				if err := in(v); err != nil {
					return err
				}
				proto.Merge(v, req)

				res, err := cli.GetValidatorAddresses(cmd.Context(), v)

				if err != nil {
					return err
				}

				return out(res)

			})
		},
	}

	return cmd
}

func _BlockchainGetPublicKeyCommand(cfg *client.Config) *cobra.Command {
	req := &GetPublicKeyRequest{}

	cmd := &cobra.Command{
		Use:   cfg.CommandNamer("GetPublicKey"),
		Short: "GetPublicKey RPC client",
		Long:  "GetPublicKey retrieves the public key of an account based on the provided address.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.UseEnvVars {
				if err := flag.SetFlagsFromEnv(cmd.Parent().PersistentFlags(), true, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain"); err != nil {
					return err
				}
				if err := flag.SetFlagsFromEnv(cmd.PersistentFlags(), false, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain", "GetPublicKey"); err != nil {
					return err
				}
			}
			return client.RoundTrip(cmd.Context(), cfg, func(cc grpc.ClientConnInterface, in iocodec.Decoder, out iocodec.Encoder) error {
				cli := NewBlockchainClient(cc)
				v := &GetPublicKeyRequest{}

				if err := in(v); err != nil {
					return err
				}
				proto.Merge(v, req)

				res, err := cli.GetPublicKey(cmd.Context(), v)

				if err != nil {
					return err
				}

				return out(res)

			})
		},
	}

	cmd.PersistentFlags().StringVar(&req.Address, cfg.FlagNamer("Address"), "", "The address for which to retrieve the public key.")

	return cmd
}

func _BlockchainGetTxPoolContentCommand(cfg *client.Config) *cobra.Command {
	req := &GetTxPoolContentRequest{}

	cmd := &cobra.Command{
		Use:   cfg.CommandNamer("GetTxPoolContent"),
		Short: "GetTxPoolContent RPC client",
		Long:  "GetTxPoolContent retrieves current transactions in the transaction pool.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.UseEnvVars {
				if err := flag.SetFlagsFromEnv(cmd.Parent().PersistentFlags(), true, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain"); err != nil {
					return err
				}
				if err := flag.SetFlagsFromEnv(cmd.PersistentFlags(), false, cfg.EnvVarNamer, cfg.EnvVarPrefix, "Blockchain", "GetTxPoolContent"); err != nil {
					return err
				}
			}
			return client.RoundTrip(cmd.Context(), cfg, func(cc grpc.ClientConnInterface, in iocodec.Decoder, out iocodec.Encoder) error {
				cli := NewBlockchainClient(cc)
				v := &GetTxPoolContentRequest{}

				if err := in(v); err != nil {
					return err
				}
				proto.Merge(v, req)

				res, err := cli.GetTxPoolContent(cmd.Context(), v)

				if err != nil {
					return err
				}

				return out(res)

			})
		},
	}

	flag.EnumVar(cmd.PersistentFlags(), &req.PayloadType, cfg.FlagNamer("PayloadType"), "The type of transactions to retrieve from the transaction pool. 0 means all types.")

	return cmd
}
