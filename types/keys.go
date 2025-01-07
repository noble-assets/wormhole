package types

const ModuleName = "wormhole"

var (
	ConfigKey           = []byte("config")
	WormchainChannelKey = []byte("wormchain_channel")
	GuardianSetPrefix   = []byte("guardian_set")
	SequencePrefix      = []byte("sequence")
	VAAArchivePrefix    = []byte("vaa_archive")
	VAAByIDPrefix       = []byte("vaa_by_id")
)
