package portal

import "testing"

func TestMockGetPatientConfig(t *testing.T) {
	mock := NewMock()

	// Init mock data
	s1 := "123"
	pc := &PatientConfig{
		ID: &s1,
	}
	mock.SetMockPatientConfig("1", pc, nil)

	s2 := "456"
	pc = &PatientConfig{
		ID: &s2,
	}
	mock.SetMockPatientConfig("2", pc, nil)

	// Test return
	pc, err := mock.GetPatientConfig("1")
	if pc == nil || pc.ID == nil || *pc.ID != "123" || err != nil {
		t.Error("Invalid mock return for token 1")
		return

	}
	pc, err = mock.GetPatientConfig("2")
	if pc == nil || pc.ID == nil || *pc.ID != "456" || err != nil {
		t.Error("Invalid mock return for token 2")
		return
	}

	pc, err = mock.GetPatientConfig("3")
	if pc != nil || err == nil {
		t.Error("Invalid mock return for token 3")
		return
	}
}
