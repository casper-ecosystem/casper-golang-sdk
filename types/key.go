package types

type KeyType byte

const (
	KeyTypeAccount KeyType = iota
	KeyTypeHash
	KeyTypeURef
	KeyTypeTransfer
	KeyTypeDeployInfo
	KeyTypeEraId
	KeyTypeBalance
	KeyTypeBid
	KeyTypeWithdraw
)

// Key represents key structure
type Key struct {
	Type       KeyType
	Account    [32]byte
	Hash       [32]byte
	URef       *URef
	Transfer   [32]byte
	DeployInfo [32]byte
	EraId      *uint64
	Balance    [32]byte
	Bid        [32]byte
	Withdraw   [32]byte
}

func (u Key) SwitchFieldName() string {
	return "Type"
}

func (u Key) ArmForSwitch(sw byte) (string, bool) {
	switch KeyType(sw) {
	case KeyTypeAccount:
		return "Account", true
	case KeyTypeHash:
		return "Hash", true
	case KeyTypeURef:
		return "URef", true
	case KeyTypeTransfer:
		return "Transfer", true
	case KeyTypeDeployInfo:
		return "DeployInfo", true
	case KeyTypeEraId:
		return "EraId", true
	case KeyTypeBalance:
		return "Balance", true
	case KeyTypeBid:
		return "Bid", true
	case KeyTypeWithdraw:
		return "Withdraw", true
	}
	return "-", false
}
