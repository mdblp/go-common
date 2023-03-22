package portal

import "fmt"

type patientConfigCall struct {
	pc  *PatientConfig
	err error
}

// MockClient The mocked interface to portal-api.
//
type MockClient struct {
	nextPatientConfigCall map[string]*patientConfigCall
}

// NewMock create a new portal mock client
func NewMock() *MockClient {
	return &MockClient{
		nextPatientConfigCall: make(map[string]*patientConfigCall),
	}
}

// SetMockPatientConfig Set the result for the next MockPatientConfig calls
//
// - token: The token string for which the response will be
//
// - pc: The PatientConfig to return or nil
//
// - err: The error to return or nil
func (client *MockClient) SetMockPatientConfig(token string, pc *PatientConfig, err error) {
	client.nextPatientConfigCall[token] = &patientConfigCall{
		pc:  pc,
		err: err,
	}
}

// GetPatientConfig mock the GetPatientConfig call
func (client *MockClient) GetPatientConfig(token string) (*PatientConfig, error) {
	pcc, ok := client.nextPatientConfigCall[token]
	if !ok {
		return nil, fmt.Errorf("Unknown response code[404] from service[http://portal/%s]", routeV2GetPatientConfig)
	}
	return pcc.pc, pcc.err
}
