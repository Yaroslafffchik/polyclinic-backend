package tests

/*
func TestCreatePatient(t *testing.T) {
	db.Init()
	patient, err := factory.NewPatient("John Doe", "123 Street", "M", 30, "1234567890123456")
	if err != nil {
		t.Fatalf("Failed to create patient: %v", err)
	}
	db.DB.Create(patient)

	var retrieved models.Patient
	if err := db.DB.First(&retrieved, patient.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve patient: %v", err)
	}
	if retrieved.FullName != "John Doe" {
		t.Errorf("Expected full name John Doe, got %s", retrieved.FullName)
	}
}
*/
