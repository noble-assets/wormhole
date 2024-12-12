package types

import "cosmossdk.io/errors"

var (
	ErrInvalidRequest = errors.Register(ModuleName, 0, "invalid request")

	ErrInvalidVAA                  = errors.Register(ModuleName, 1, "invalid vaa")
	ErrAlreadyExecutedVAA          = errors.Register(ModuleName, 2, "vaa already executed")
	ErrNotGovernanceVAA            = errors.Register(ModuleName, 3, "not governance vaa")
	ErrInvalidGovernanceVAA        = errors.Register(ModuleName, 4, "invalid governance vaa")
	ErrUnsupportedGovernanceAction = errors.Register(ModuleName, 5, "unsupported governance action")
	ErrMalformedPayload            = errors.New(ModuleName, 6, "payload is malformed")

	ErrInvalidPort    = errors.New(ModuleName, 101, "invalid port")
	ErrInvalidVersion = errors.New(ModuleName, 102, "invalid version")
	ErrInvalidChannel = errors.New(ModuleName, 103, "channel is not associated with wormchain")
)
