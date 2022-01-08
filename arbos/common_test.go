//
// Copyright 2021, Offchain Labs, Inc. All rights reserved.
//

package arbos

import (
	"testing"

	"github.com/offchainlabs/arbstate/util/testhelpers"
)

func Require(t *testing.T, err error, text ...string) {
	t.Helper()
	testhelpers.RequireImpl(t, err, text...)
}

func Fail(t *testing.T, printables ...interface{}) {
	t.Helper()
	testhelpers.FailImpl(t, printables...)
}