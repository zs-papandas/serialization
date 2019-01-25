package types

const (
	unknownProductStr = "unknown"
	palletProductStr = "pallet"
	boxProductStr   = "box"
	packetProductStr = "packet"
	itemProductStr = "item"
)

// ProductType : product type
type ProductType int

// concrete UserType
const (
	UnKnownProduct ProductType = iota
	PalletProduct
	BoxProduct
	PacketProduct 
	ItemProduct 
)

// String : Stringer interface
func (t ProductType) String() string {
	switch t {
	case PalletProduct:
		return palletProductStr
	case BoxProduct:
		return boxProductStr
	case PacketProduct:
		return packetProductStr
	case ItemProduct:
		return itemProductStr
	default:
		return unknownUserStr
	}
}