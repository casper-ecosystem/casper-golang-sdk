package types

type AccessRight byte

const (
	AccessRightRead AccessRight = 1 << iota
	AccessRightWrite
	AccessRightAdd
	AccessRightNone AccessRight = 0
)

type URef struct {
	AccessRight AccessRight
	Address     [32]byte
}
