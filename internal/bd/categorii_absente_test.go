package bd

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestMigrareCategoriiAbsentePastreazaEfectele(t *testing.T) {
	db, e := Deschide(filepath.Join(t.TempDir(), "categorii.db"), cheieNoua(t))
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	ctx := context.Background()
	vechi := fstest.MapFS{}
	intrari, _ := fs.ReadDir(Migrari, ".")
	for _, v := range intrari {
		if strings.HasSuffix(v.Name(), ".sql") && v.Name() < "0030" {
			b, _ := fs.ReadFile(Migrari, v.Name())
			vechi[v.Name()] = &fstest.MapFile{Data: b}
		}
	}
	if _, e := Migreaza(ctx, db, vechi); e != nil {
		t.Fatal(e)
	}
	if _, e := db.Exec(`INSERT INTO coduri_absenta(cod,denumire,categorie,efect_ore,activ,creat_la,modificat_la) VALUES('DEPL','{"ro":"Deplasare demo"}','SPECIAL','NORMA',0,'x','y');INSERT INTO angajati(id,nume,marca,creat_la,modificat_la) VALUES('a','Exemplu','TEST','x','x');INSERT INTO absente_zile(angajat_id,data,cod,minute_creditate,creat_la,modificat_la) VALUES('a','2026-10-05','DEPL',420,'x','x')`); e != nil {
		t.Fatal(e)
	}
	if n, e := Migreaza(ctx, db, Migrari); e != nil || n != 3 {
		t.Fatal(n, e)
	}
	var cat, natura, efect, mod string
	var activ, minute int
	if e := db.QueryRow(`SELECT t.categorie_id,c.natura,c.efect_ore,t.activ,t.modificat_la,z.minute_creditate FROM coduri_absenta t JOIN categorii_absente c ON c.cod=t.categorie_id JOIN absente_zile z ON z.cod=t.cod WHERE t.cod='DEPL'`).Scan(&cat, &natura, &efect, &activ, &mod, &minute); e != nil {
		t.Fatal(e)
	}
	if cat != "SPECIAL_NORMA" || natura != "SPECIAL" || efect != "NORMA" || activ != 0 || mod != "y" || minute != 420 {
		t.Fatalf("date schimbate: %s %s %s %d %s %d", cat, natura, efect, activ, mod, minute)
	}
	if _, e := db.Exec(`UPDATE coduri_absenta SET categorie_id='lipsa' WHERE cod='DEPL'`); e == nil {
		t.Fatal("referinta invalida acceptata")
	}
}
