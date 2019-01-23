package types

const (
	unknownUserStr = "unknown"
	manufacturerUserStr = "manufacturer"
	wholesalerUserStr   = "wholesaler"
	retaileUserStr = "retailer"
	patientUserStr = "patient"
)

// ModelType : model type
type UserType int

// concrete ModelType
const (
	UnKnownUser UserType = iota
	ManufacturerUser
	WholesalerUser
	RetailerUser 
	PatientUser 
)

// String : Stringer interface
func (t UserType) String() string {
	switch t {
	case ManufacturerUser:
		return manufacturerUserStr
	case WholesalerUser:
		return wholesalerUserStr
	case RetailerUser:
		return retaileUserStr
	case PatientUser:
		return patientUserStr
	default:
		return unknownUserStr
	}
}