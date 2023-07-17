// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package module

// Module represents a lens module loaded into a runtime.
//
// It may be used to instantiate multiple lens instances.
type Module interface {
	// NewInstance returns a new lens instance from this module, hosted
	// within the parent runtime.
	NewInstance(string, ...map[string]any) (Instance, error)
}
