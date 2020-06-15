// Package portal is a client module to support server-side use of the Diabeloop
// service called user-api.
//
// This file contains the definitions for the patient config
package portal

// PatientConfig structure returned by the API
type PatientConfig struct {
	ID             *string                  `json:"_id,omitempty"`
	Time           *string                  `json:"time,omitempty"`
	Timezone       *string                  `json:"timezone,omitempty"`
	TimezoneOffset *string                  `json:"timezoneOffset,omitempty"`
	Device         *PatientConfigDevice     `json:"device,omitempty"`
	Parameters     *PatientConfigParameters `json:"parameters,omitempty"`
	CGM            *PatientCGMInfo          `json:"cgm,omitempty"`
	Pump           *PatientPumpInfo         `json:"pump,omitempty"`
}

// PatientConfigDevice sub-structure of PatientConfig which old the device information
type PatientConfigDevice struct {
	// HistoryId The _id of device_parameter_history
	HistoryID *string `json:"historyId,omitempty"`
	// DeviceId like "DBLG-12345"
	DeviceID string `json:"deviceId"`
	// IMEI of the terminal device
	IMEI string `json:"imei"`
	// Device name: DBLG1
	Name string `json:"name"`
	// Device manufacturer: Diabeloop
	Manufacturer string `json:"manufacturer"`
	// Version of the software running on the device
	SWVersion string `json:"swVersion"`
}

// PatientConfigParameters sub-structure of PatientConfig
type PatientConfigParameters struct {
	// HistoryId The _id of patient_parameters_history
	HistoryID *string            `json:"historyId,omitempty"`
	Values    []PatientParameter `json:"values"`
}

// PatientParameter is a single parameter for a patient.
type PatientParameter struct {
	// Name is unique act like a primary key
	Name string `json:"name"`
	// Value of the parameter
	Value string `json:"value"`
	// Unit of the parameter
	Unit *string `json:"unit,omitempty"`
	// Level used for a filter in the UI
	Level    string  `json:"level"`
	MinValue *string `json:"minValue,omitempty"`
	MaxValue *string `json:"maxValue,omitempty"`
	// Processed "yes" | "no"
	Processed     *string   `json:"processed,omitempty"`
	LinkedSubType *[]string `json:"linkedSubType,omitempty"`
	EffectiveDate *string   `json:"effectiveDate,omitempty"`
}

// PatientCGMInfo is the Continuous Glucose Monitoring device information
type PatientCGMInfo struct {
	// Manufacturer ex: Dexcom
	Manufacturer string `json:"manufacturer"`
	// TransmitterID ex: 12345
	TransmitterID string `json:"transmitterId"`
	// Name ex: G6
	Name string `json:"name"`
	// SWVersionTransmitter Software version ex: G6
	SWVersionTransmitter string `json:"swVersionTransmitter"`
	// APIVersion of the software running on the cgm
	APIVersion string `json:"apiVersion"`
	// Estimated end of life of transmitter
	EndOfLifeTransmitterDate *string `json:"endOfLifeTransmitterDate,omitempty"`
	// Expiration Date of the session
	ExpirationDate *string `json:"expirationDate,omitempty"`
}

// PatientPumpInfo Insulin pump device information.
type PatientPumpInfo struct {
	// Pump manufacturer: VICENTRA
	Manufacturer string `json:"manufacturer"`
	// Pump name: blue, red, whatever color is provided
	Name string `json:"name"`
	// Pump serial number
	SerialNumber string `json:"serialNumber"`
	// Pump software version: 0.1.0
	SWVersion string `json:"swVersion"`
	// Expiration date of the session
	ExpirationDate *string `json:"expirationDate,omitempty"`
}
