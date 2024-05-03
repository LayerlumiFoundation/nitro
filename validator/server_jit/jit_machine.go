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
	4083748: {
		Batch:      9178,
		PosInBatch: 651,
		BlockHash:  common.HexToHash("0x3e09fd518f880a952817561f390f28150875cdfc82e646c8325c6b68fd63ab2f"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083749: {
		Batch:      9178,
		PosInBatch: 652,
		BlockHash:  common.HexToHash("0xd9ffa49df6c0931f0a748b8dcf78e50393f7a10c2cecf3fd3e8d9ece23adbdd2"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083750: {
		Batch:      9178,
		PosInBatch: 653,
		BlockHash:  common.HexToHash("0x1051568779341313cfc41bb6425c99031d8cc847b510e964aed8eac2333aa697"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083751: {
		Batch:      9178,
		PosInBatch: 654,
		BlockHash:  common.HexToHash("0x1d798a555187a0f80d2f38e6b0079d4fe8a3a1d57febe11cf442f17e3149d9b0"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083752: {
		Batch:      9178,
		PosInBatch: 655,
		BlockHash:  common.HexToHash("0x338c0954729d69e259e220d97486e9edeb12f58d2cd17986374bbb075d21309b"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083753: {
		Batch:      9178,
		PosInBatch: 656,
		BlockHash:  common.HexToHash("0x029cf4307f5fa3c4920acd16f4d936b94264f46c3765435e9f3c2a2b2b33f85c"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083754: {
		Batch:      9178,
		PosInBatch: 657,
		BlockHash:  common.HexToHash("0x8d585dbf517707d88c8816df4d3edb853622415bb8f44f6f34a6c788b0eac35c"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083755: {
		Batch:      9178,
		PosInBatch: 658,
		BlockHash:  common.HexToHash("0x409ecf21a914f4341a6ba9ac3ef728a3f45dc028b35df738f2ecdc95a95165c2"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083756: {
		Batch:      9178,
		PosInBatch: 659,
		BlockHash:  common.HexToHash("0x8e23132551a67624d98faacf012074c946560a0be1e2a9648084b2524e411b84"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083757: {
		Batch:      9179,
		PosInBatch: 0,
		BlockHash:  common.HexToHash("0xa4b18acb1e38e1389ad2cebdf01451ce9eee86fa3bb35589d31d93b6bff8b648"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},

	4083758: {
		Batch:      9179,
		PosInBatch: 1,
		BlockHash:  common.HexToHash("0x49660151be5e75586f749f5010deaa20a4c20e4a309506dc9e0a90cfa6438ec2"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083759: {
		Batch:      9179,
		PosInBatch: 2,
		BlockHash:  common.HexToHash("0xea193dbd15a4ade4af250013d19de411baf8aae540b984a12eaad3480f50768c"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083760: {
		Batch:      9179,
		PosInBatch: 3,
		BlockHash:  common.HexToHash("0xd22410f061d9e522c1c6dc4cef724d40d1484eb02c49c63a1ce03d1c2a67e66e"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083761: {
		Batch:      9179,
		PosInBatch: 4,
		BlockHash:  common.HexToHash("0x9d5c886c9706885875afdf3902e09d12201f139232e662e165f43fb1c0d61b5a"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083762: {
		Batch:      9179,
		PosInBatch: 5,
		BlockHash:  common.HexToHash("0x04820347686c578d80495b56e53267ee3d0cade3a03d4a1d9a15ee001ca4dca7"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083763: {
		Batch:      9179,
		PosInBatch: 6,
		BlockHash:  common.HexToHash("0xd7ca7bef74fe43f2d7e6d1845c91cfa98c1a0025191ac631c6dcf0b1b05d31ff"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083764: {
		Batch:      9179,
		PosInBatch: 7,
		BlockHash:  common.HexToHash("0x2fe06217c0dbd14cde1b5c9c28985e2106a70c7d2f28e2ad1e9b6336156353ee"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083765: {
		Batch:      9179,
		PosInBatch: 8,
		BlockHash:  common.HexToHash("0xc3a64bd8a09d277f06c0a5b5156c944554ed93ee8d9419e1593a39e0c9c31c65"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083766: {
		Batch:      9179,
		PosInBatch: 9,
		BlockHash:  common.HexToHash("0x98bff06722a2faeea3d9f4a5b87d23b5d4e351a0a268911a978ccfe2ebaf7557"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083767: {
		Batch:      9179,
		PosInBatch: 10,
		BlockHash:  common.HexToHash("0x41964b2279530661f093695836c814797d1b057fc5bc50c1eb2fecd7748891e4"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083768: {
		Batch:      9179,
		PosInBatch: 11,
		BlockHash:  common.HexToHash("0xcc7906a405cbe757fe18a4b706b31b5de10272b6164dfd4a71c4cf04586a1965"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083769: {
		Batch:      9179,
		PosInBatch: 12,
		BlockHash:  common.HexToHash("0x473e99644243e72a1b24643aa511d155141e3b0fa91fcb2538da50ea9221f029"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083770: {
		Batch:      9179,
		PosInBatch: 13,
		BlockHash:  common.HexToHash("0xd010b7dc2f71266d76762ff9c75e179465f99e60766323ad150a42d71d2f0fc1"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083771: {
		Batch:      9179,
		PosInBatch: 14,
		BlockHash:  common.HexToHash("0x8be5ec7d772181e086cccbf18c0619b9332b5291248d8940380e9e4a5a4775c9"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083772: {
		Batch:      9179,
		PosInBatch: 15,
		BlockHash:  common.HexToHash("0x48d13f72a3af6a6e7377be4f46059d9d3648a32273d942977033bf95f34d97aa"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083773: {
		Batch:      9179,
		PosInBatch: 16,
		BlockHash:  common.HexToHash("0x36007d81d2dd023743c3c6988840edd26297b363b9691f97e8779ad0e2c56dc2"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083774: {
		Batch:      9179,
		PosInBatch: 17,
		BlockHash:  common.HexToHash("0x2158c56d28237190bebaa93d5b514789d16986383d8834fa888d0efc29ca8dc7"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083775: {
		Batch:      9179,
		PosInBatch: 18,
		BlockHash:  common.HexToHash("0x873e66a3521d74e92ab638a6626eada0cfd0c95cff23bcd596e4529379128721"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083776: {
		Batch:      9179,
		PosInBatch: 19,
		BlockHash:  common.HexToHash("0x70d86d565fe2a388c59a4fa8c3f2667cc7f41b137f689e03e8ad06bd2ee865f2"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083777: {
		Batch:      9179,
		PosInBatch: 20,
		BlockHash:  common.HexToHash("0xc586b5b34836563a4c65c3cb801d5363e45301a1bad36b8c603e74b6de4859bc"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083778: {
		Batch:      9179,
		PosInBatch: 21,
		BlockHash:  common.HexToHash("0x2a107b42c03dff558bae068a877f215f57d72782bcaeba36e3a7d49425477e2a"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083779: {
		Batch:      9179,
		PosInBatch: 22,
		BlockHash:  common.HexToHash("0xab2fe94ad82ec7a2c18fef9ce5b08114ee332ccd4201663d9f514943ce503cd1"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083780: {
		Batch:      9179,
		PosInBatch: 23,
		BlockHash:  common.HexToHash("0xc0a14248fe4399503d2a85be30f40ef3f5e0d91fa35559a22cb043c87f275b4b"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083781: {
		Batch:      9179,
		PosInBatch: 24,
		BlockHash:  common.HexToHash("0xe19b0f611e545a630d4f1fcb90de59393bb62971578c7aae025e5843e9ac847c"),
		SendRoot:   common.HexToHash("0x73e3fd5339538bb97fb2ab3f1affc7971c8ea591c1f324e22cbc5b4d01b7b4ee"),
	},
	4083782: {
		Batch:      9179,
		PosInBatch: 25,
		BlockHash:  common.HexToHash("0xe3da69cf80a6ffa9c9c6b92abc60484018a7b0892c9dd6c0ab3d7e8eb0528f9b"),
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
