package calcul

import "testing"

func TestPrioritateMedicalConcediuProgram(t *testing.T) {
	for _, c := range []string{"ODIHNA", "SPECIAL"} {
		if PrioritateAbsenta("MEDICAL") <= PrioritateAbsenta(c) || PrioritateAbsenta(c) <= PrioritateAbsenta("") {
			t.Fatal(c)
		}
	}
	if PrioritateAbsenta("ODIHNA") != PrioritateAbsenta("SPECIAL") {
		t.Fatal("prioritate inventata intre concedii")
	}
}
