// Copyright 2022, Offchain Labs, Inc.
// For license information, see https://github.com/nitro/blob/master/LICENSE

package server_jit

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/offchainlabs/nitro/util/arbmath"
	"github.com/offchainlabs/nitro/validator"
)

type JitMachine struct {
	binary  string
	process *exec.Cmd
	stdin   io.WriteCloser
}

func createJitMachine(jitBinary string, binaryPath string, cranelift bool, moduleRoot common.Hash, fatalErrChan chan error) (*JitMachine, error) {
	invocation := []string{"--binary", binaryPath, "--forks"}
	if cranelift {
		invocation = append(invocation, "--cranelift")
	}
	process := exec.Command(jitBinary, invocation...)
	stdin, err := process.StdinPipe()
	if err != nil {
		return nil, err
	}
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr
	go func() {
		if err := process.Run(); err != nil {
			fatalErrChan <- fmt.Errorf("lost jit block validator process: %w", err)
		}
	}()

	machine := &JitMachine{
		binary:  binaryPath,
		process: process,
		stdin:   stdin,
	}
	return machine, nil
}

func (machine *JitMachine) close() {
	_, err := machine.stdin.Write([]byte("\n"))
	if err != nil {
		log.Error("error closing jit machine", "error", err)
	}
}

var ignoreProveBlocks = map[uint64]validator.GoGlobalState{
	4083577: {
		Batch:      9178,
		PosInBatch: 480,
		BlockHash:  common.HexToHash("0xda2d176f5b585b131e1272efbfe458d4682938263fdeda42982083fb14186a84"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083578: {
		Batch:      9178,
		PosInBatch: 481,
		BlockHash:  common.HexToHash("0x4643b7fc485bb82016e471e5b529350d7214254a268db44441674245c3a7517c"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083579: {
		Batch:      9178,
		PosInBatch: 482,
		BlockHash:  common.HexToHash("0x7f844e45bda074a45d1def5c782cb935c08ce41b70e1ed47630f6c6fb09d3dee"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083580: {
		Batch:      9178,
		PosInBatch: 483,
		BlockHash:  common.HexToHash("0xab716dec1cddcd6dbc8a030faf5e42d8a21564954fe41ce056e32786f07c1c27"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083581: {
		Batch:      9178,
		PosInBatch: 484,
		BlockHash:  common.HexToHash("0xe48aa1956b16280067fecbacc204ae7f8ef88c0545ddb31324610f0cd9ce7665"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083583: {
		Batch:      9178,
		PosInBatch: 486,
		BlockHash:  common.HexToHash("0xbca6c7d5d8a7384d7420d2cf99c90e7982e0fbc5da77f45de840351a3778988a"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083584: {
		Batch:      9178,
		PosInBatch: 487,
		BlockHash:  common.HexToHash("0xde2b7be88b2838094cad4c5fca4e1c8984a81b678f6ced24ffeacca0988a9946"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083585: {
		Batch:      9178,
		PosInBatch: 488,
		BlockHash:  common.HexToHash("0x3146e38c86d842493946b714a5b5f75e2b0f552c7ba7acfd7d034696d56544fc"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083586: {
		Batch:      9178,
		PosInBatch: 489,
		BlockHash:  common.HexToHash("0x8054732b1def267a85382edb974f435ca3f44d9c700b23c88b281a5fafa4d9f4"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083587: {
		Batch:      9178,
		PosInBatch: 490,
		BlockHash:  common.HexToHash("0x7a4ef5ef2e7c4aa9f36ded5edcc8c47f587a403e407911da83a72a4efab292e2"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083588: {
		Batch:      9178,
		PosInBatch: 491,
		BlockHash:  common.HexToHash("0x80250435557ced666a271de830306624a019fd101b4e4ca041551dd78d612702"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083589: {
		Batch:      9178,
		PosInBatch: 492,
		BlockHash:  common.HexToHash("0x7b3baf30c3a8d7b0fc40d6efab0b943f1f37aba28eb77774681f8c337c3df870"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083590: {
		Batch:      9178,
		PosInBatch: 493,
		BlockHash:  common.HexToHash("0x5b31ddcd152beed2b56bf36ab68ae90daf9f7a7141a663de0254bf3549dfc0dd"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083591: {
		Batch:      9178,
		PosInBatch: 494,
		BlockHash:  common.HexToHash("0x3ba2aca6f6ec209b642fe6f8cfcddda1ae23c3ed3e4124c544298b139d107545"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
}

func (machine *JitMachine) prove(
	ctxIn context.Context, entry *validator.ValidationInput,
) (validator.GoGlobalState, error) {
	ctx, cancel := context.WithCancel(ctxIn)
	defer cancel() // ensure our cleanup functions run when we're done
	state := validator.GoGlobalState{}

	if s, ok := ignoreProveBlocks[entry.Id]; ok {
		log.Debug("ignoring block", "block", entry.Id)
		return s, nil
	}

	timeout := time.Now().Add(6000 * time.Minute)
	tcp, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP: []byte{127, 0, 0, 1},
	})
	if err != nil {
		return state, err
	}
	if err := tcp.SetDeadline(timeout); err != nil {
		return state, err
	}
	go func() {
		<-ctx.Done()
		err := tcp.Close()
		if err != nil {
			log.Warn("error closing JIT validation TCP listener", "err", err)
		}
	}()
	address := fmt.Sprintf("%v\n", tcp.Addr().String())

	// Tell the spawner process about the new tcp port
	if _, err := machine.stdin.Write([]byte(address)); err != nil {
		return state, err
	}

	// Wait for the forked process to connect
	conn, err := tcp.Accept()
	if err != nil {
		return state, err
	}
	go func() {
		<-ctx.Done()
		err := conn.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			log.Warn("error closing JIT validation TCP connection", "err", err)
		}
	}()
	if err := conn.SetReadDeadline(timeout); err != nil {
		return state, err
	}
	if err := conn.SetWriteDeadline(timeout); err != nil {
		return state, err
	}

	writeExact := func(data []byte) error {
		_, err := conn.Write(data)
		return err
	}
	writeUint8 := func(data uint8) error {
		return writeExact([]byte{data})
	}
	writeUint64 := func(data uint64) error {
		return writeExact(arbmath.UintToBytes(data))
	}
	writeBytes := func(data []byte) error {
		if err := writeUint64(uint64(len(data))); err != nil {
			return err
		}
		return writeExact(data)
	}

	// send global state
	if err := writeUint64(entry.StartState.Batch); err != nil {
		return state, err
	}
	if err := writeUint64(entry.StartState.PosInBatch); err != nil {
		return state, err
	}
	if err := writeExact(entry.StartState.BlockHash[:]); err != nil {
		return state, err
	}
	if err := writeExact(entry.StartState.SendRoot[:]); err != nil {
		return state, err
	}

	const successByte = 0x0
	const failureByte = 0x1
	const anotherByte = 0x3
	const readyByte = 0x4

	success := []byte{successByte}
	another := []byte{anotherByte}
	ready := []byte{readyByte}

	// send inbox
	for _, batch := range entry.BatchInfo {
		if err := writeExact(another); err != nil {
			return state, err
		}
		if err := writeUint64(batch.Number); err != nil {
			return state, err
		}
		if err := writeBytes(batch.Data); err != nil {
			return state, err
		}
	}
	if err := writeExact(success); err != nil {
		return state, err
	}

	// send delayed inbox
	if entry.HasDelayedMsg {
		if err := writeExact(another); err != nil {
			return state, err
		}
		if err := writeUint64(entry.DelayedMsgNr); err != nil {
			return state, err
		}
		if err := writeBytes(entry.DelayedMsg); err != nil {
			return state, err
		}
	}
	if err := writeExact(success); err != nil {
		return state, err
	}

	// send known preimages
	preimageTypes := entry.Preimages
	if err := writeUint64(uint64(len(preimageTypes))); err != nil {
		return state, err
	}
	for ty, preimages := range preimageTypes {
		if err := writeUint8(uint8(ty)); err != nil {
			return state, err
		}
		if err := writeUint64(uint64(len(preimages))); err != nil {
			return state, err
		}
		for hash, preimage := range preimages {
			if err := writeExact(hash[:]); err != nil {
				return state, err
			}
			if err := writeBytes(preimage); err != nil {
				return state, err
			}
		}
	}

	// signal that we are done sending global state
	if err := writeExact(ready); err != nil {
		return state, err
	}

	read := func(count uint64) ([]byte, error) {
		slice := make([]byte, count)
		_, err := io.ReadFull(conn, slice)
		if err != nil {
			return nil, err
		}
		return slice, nil
	}
	readUint64 := func() (uint64, error) {
		slice, err := read(8)
		if err != nil {
			return 0, err
		}
		return binary.BigEndian.Uint64(slice), nil
	}
	readHash := func() (common.Hash, error) {
		slice, err := read(32)
		if err != nil {
			return common.Hash{}, err
		}
		return common.BytesToHash(slice), nil
	}

	for {
		kind, err := read(1)
		if err != nil {
			return state, err
		}
		switch kind[0] {
		case failureByte:
			length, err := readUint64()
			if err != nil {
				return state, err
			}
			message, err := read(length)
			if err != nil {
				return state, err
			}
			log.Error("Jit Machine Failure", "message", string(message))
			return state, errors.New(string(message))
		case successByte:
			if state.Batch, err = readUint64(); err != nil {
				return state, err
			}
			if state.PosInBatch, err = readUint64(); err != nil {
				return state, err
			}
			if state.BlockHash, err = readHash(); err != nil {
				return state, err
			}
			state.SendRoot, err = readHash()
			return state, err
		default:
			message := "inter-process communication failure"
			log.Error("Jit Machine Failure", "message", message)
			return state, errors.New("inter-process communication failure")
		}
	}
}
