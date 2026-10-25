package bd

import (
	"context"
	"path/filepath"
	"testing"
)

func TestExceptiaProbeiNuPermiteAlteGoluri(t *testing.T) {
	for _, lipsa := range []int{11, 12} {
		t.Run(map[int]string{11: "lipseste categoria", 12: "istoric ulterior probei"}[lipsa], func(t *testing.T) {
			db, err := Deschide(filepath.Join(t.TempDir(), "istoric.db"), cheieNoua(t))
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			if _, err := Migreaza(context.Background(), db, Migrari); err != nil {
				t.Fatal(err)
			}
			if lipsa == 11 {
				if _, err := db.Exec(`DELETE FROM goose_db_version WHERE version_id=14`); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := db.Exec(`DELETE FROM goose_db_version WHERE version_id=?`, lipsa); err != nil {
				t.Fatal(err)
			}
			if _, _, err := StareMigrari(context.Background(), db, Migrari); err == nil {
				t.Fatal("un istoric diferit de proba 0.20 a primit exceptia")
			}
		})
	}
}
