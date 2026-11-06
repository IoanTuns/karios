package calcul

// PrioritateAbsenta foloseste categoria configurata, nu denumirea libera a codului.
// Valoarea zero reprezinta programul de lucru; medicalul precede celelalte concedii.
func PrioritateAbsenta(categorie string) int {
	if categorie == "" {
		return 0
	}
	if categorie == "MEDICAL" {
		return 2
	}
	return 1
}
