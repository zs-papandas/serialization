package types

const (
	unknownStatusStr = "unknown"
	createStatusStr = "create"
	soldStatusStr   = "sold"
	ownerChangeStatusStr = "owner"
)

// StatusType : Status type
type StatusType int

// concrete UserType
const (
	UnKnownStatus StatusType = iota
	CreateStatus
	SoldStatus
	OwnerChangeStatus 
)

// String : Stringer interface
func (t StatusType) String() string {
	switch t {
	case CreateStatus:
		return createStatusStr
	case SoldStatus:
		return soldStatusStr
	case OwnerChangeStatus:
		return ownerChangeStatusStr
	default:
		return unknownStatusStr
	}
}