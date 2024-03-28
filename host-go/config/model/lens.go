// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

/*
The model package contains a lens-file-format agnostic set of types that describe the contents
of a given lens file.
*/
package model

import (
	"crypto/sha256"
	"encoding/hex"
)

type Lens struct {
	// The LensModules that should be applied to the source data, declared in the order
	// in which they should be executed.
	Lenses []LensModule `json:"lenses"`
}

type LensModule struct {
	// Content is the WASM binary
	Content []byte `json:"content"`

	// If true, the module will be inversed.
	//
	// This may result in an error if the module does not provide an inverse function.
	Inverse bool `json:"inverse"`

	// Any additional parameters that you wish to be passed to the lens transform.
	//
	// The lens module must expose a `set_param` function if values are provided here.
	Arguments map[string]any `json:"arguments"`
}

func (m LensModule) Hash() string {
	hash := sha256.Sum256(m.Content)
	return hex.EncodeToString(hash[:])
}
