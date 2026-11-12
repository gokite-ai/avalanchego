// Copyright (C) 2019, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// A subprocess handshake only: no VM Initialize, chain, database or staking identity.
package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/version"
	"github.com/ava-labs/avalanchego/vms/rpcchainvm/runtime/subprocess"
)

const pluginPath = "/plugins-template/subnet-evm"

func main() {
	if len(os.Args) != 2 || os.Args[1] != pluginPath {
		panic("usage: candidate-handshake /plugins-template/subnet-evm")
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	status, stopper, err := subprocess.Bootstrap(ctx, listener, exec.Command(pluginPath), &subprocess.Config{
		Stdout: io.Discard, Stderr: io.Discard,
		HandshakeTimeout: 10 * time.Second, Log: logging.NoLog{},
	})
	if err != nil {
		panic(err)
	}
	defer stopper.Stop(context.Background())
	if status.Addr == "" {
		panic("empty plugin address")
	}
	fmt.Printf("PASS: actual Subnet-EVM subprocess handshake; rpcchainvm=%d\n", version.RPCChainVMProtocol)
}
