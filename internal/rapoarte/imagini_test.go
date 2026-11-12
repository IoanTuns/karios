package rapoarte

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"hash/crc32"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"pontix/internal/bd"
)

func imagineFictiva(t *testing.T, format string) Imagine {
	t.Helper()
	w, h := 80, 32
	if format == "jpeg" {
		w, h = 110, 20
	}
	imagine := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			culoare := color.NRGBA{R: 30, G: 125, B: 100, A: 255}
			if format == "jpeg" {
				culoare = color.NRGBA{R: 60, G: 80, B: 160, A: 255}
			} else if x < 4 || y < 4 {
				culoare.A = 0
			}
			imagine.SetNRGBA(x, y, culoare)
		}
	}
	var buf bytes.Buffer
	var err error
	if format == "jpeg" {
		err = jpeg.Encode(&buf, imagine, nil)
	} else {
		err = png.Encode(&buf, imagine)
	}
	if err != nil {
		t.Fatal(err)
	}
	extensie := "png"
	if format == "jpeg" {
		extensie = "jpg"
	}
	return Imagine{Nume: "imagine-fictiva." + extensie, Date: base64.StdEncoding.EncodeToString(buf.Bytes()), InaltimeMM: 15}
}

func TestImaginiValidarePersistentaSiIndependenta(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	d := initiala(t, m, "raport-pontaj")
	png := imagineFictiva(t, "png")
	jpg := imagineFictiva(t, "jpeg")
	d.Antet.Imagini = map[string]Imagine{"stanga": png}
	d.Subsol.Imagini = map[string]Imagine{"dreapta": jpg}
	if _, err := m.s.Salveaza(ctx, d); err != nil {
		t.Fatal(err)
	}
	if err := m.db.Close(); err != nil {
		t.Fatal(err)
	}
	var err error
	m.db, err = bd.Deschide(m.cale, m.cheie)
	if err != nil {
		t.Fatal(err)
	}
	m.s = Nou(m.db)
	salvat := initiala(t, m, d.ID)
	if salvat.Antet.Imagini["stanga"] != png || salvat.Subsol.Imagini["dreapta"] != jpg {
		t.Fatal("imaginile s-au pierdut la repornire")
	}
	if len(initiala(t, m, "raport-fisa").Antet.Imagini) != 0 {
		t.Fatal("imaginea a ajuns in alt raport")
	}
	copie := salvat
	copie.ID = ""
	copie.Revizie = ""
	id, err := m.s.Salveaza(ctx, copie)
	if err != nil {
		t.Fatal(err)
	}
	delete(salvat.Antet.Imagini, "stanga")
	if _, err := m.s.Salveaza(ctx, salvat); err != nil {
		t.Fatal(err)
	}
	if initiala(t, m, id).Antet.Imagini["stanga"] != png || len(initiala(t, m, d.ID).Antet.Imagini) != 0 {
		t.Fatal("eliminarea imaginii a schimbat copia")
	}
	variante := []struct {
		nume    string
		imagine Imagine
		cod     string
	}{
		{"base64 invalid", Imagine{Nume: "x.png", Date: "!!!", InaltimeMM: 15}, "rapoarte.imagine_invalida"},
		{"svg", Imagine{Nume: "x.png", Date: base64.StdEncoding.EncodeToString([]byte("<svg/>")), InaltimeMM: 15}, "rapoarte.imagine_invalida"},
		{"prea multi octeti", Imagine{Nume: "x.png", Date: strings.Repeat("A", base64.StdEncoding.EncodedLen(MaxOctetiImagine)+4), InaltimeMM: 15}, "rapoarte.imagine_prea_mare"},
	}
	for _, inaltime := range []int64{0, 4, 36} {
		i := png
		i.InaltimeMM = inaltime
		variante = append(variante, struct {
			nume    string
			imagine Imagine
			cod     string
		}{"inaltime", i, "rapoarte.inaltime_imagine_invalida"})
	}
	i := png
	i.Nume = "../imagine.png"
	variante = append(variante, struct {
		nume    string
		imagine Imagine
		cod     string
	}{"cale", i, "rapoarte.imagine_invalida"})
	date, _ := base64.StdEncoding.DecodeString(png.Date)
	i = png
	i.Date = base64.StdEncoding.EncodeToString(date[:len(date)/2])
	variante = append(variante, struct {
		nume    string
		imagine Imagine
		cod     string
	}{"trunchiata", i, "rapoarte.imagine_invalida"})
	for _, dim := range [][2]uint32{{4097, 1}, {2001, 2000}} {
		mare := append([]byte(nil), date...)
		binary.BigEndian.PutUint32(mare[16:20], dim[0])
		binary.BigEndian.PutUint32(mare[20:24], dim[1])
		binary.BigEndian.PutUint32(mare[29:33], crc32.ChecksumIEEE(mare[12:29]))
		i = png
		i.Date = base64.StdEncoding.EncodeToString(mare)
		variante = append(variante, struct {
			nume    string
			imagine Imagine
			cod     string
		}{"pixeli", i, "rapoarte.imagine_prea_mare"})
	}
	for _, v := range variante {
		t.Run(v.nume, func(t *testing.T) {
			model := initiala(t, m, d.ID)
			model.Antet.Imagini = map[string]Imagine{"stanga": v.imagine}
			_, err := m.s.Salveaza(ctx, model)
			codEroare(t, err, v.cod)
		})
	}
	model := initiala(t, m, d.ID)
	model.Antet.Imagini = map[string]Imagine{"zona-necunoscuta": png}
	_, err = m.s.Salveaza(ctx, model)
	codEroare(t, err, "rapoarte.imagine_invalida")
}

func TestMigrareImaginiPastreazaDefinitiile(t *testing.T) {
	m := pregateste(t)
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "modele-vechi.db"), m.cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	surse := fstest.MapFS{}
	intrari, err := fs.ReadDir(bd.Migrari, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, intrare := range intrari {
		if intrare.Name() >= "0014_" {
			continue
		}
		date, err := fs.ReadFile(bd.Migrari, intrare.Name())
		if err != nil {
			t.Fatal(err)
		}
		surse[intrare.Name()] = &fstest.MapFile{Data: date}
	}
	if _, err := bd.Migreaza(context.Background(), db, surse); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`UPDATE definitii_rapoarte SET antet=json_set(antet,'$.stanga.ro','Unitate personalizata'), revizie=7,modificat_la='moment-neschimbat' WHERE id='raport-pontaj'`); err != nil {
		t.Fatal(err)
	}
	if n, err := bd.Migreaza(context.Background(), db, bd.Migrari); err != nil || n != 29 {
		t.Fatalf("migrare imagini: %d %v", n, err)
	}
	d := initiala(t, &mediu{s: Nou(db)}, "raport-pontaj")
	var moment, antet string
	if err := db.QueryRow(`SELECT modificat_la,antet FROM definitii_rapoarte WHERE id='raport-pontaj'`).Scan(&moment, &antet); err != nil {
		t.Fatal(err)
	}
	var campuri map[string]json.RawMessage
	if err := json.Unmarshal([]byte(antet), &campuri); err != nil {
		t.Fatal(err)
	}
	if d.Revizie != "7" || moment != "moment-neschimbat" || d.Antet.Stanga["ro"] != "Unitate personalizata" || string(campuri["imagini"]) != "{}" {
		t.Fatal("migrarea a schimbat definitia veche sau nu a initializat imaginile")
	}
}

func TestActualizareDinProbaRapoarte020(t *testing.T) {
	m := pregateste(t)
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "proba-020.db"), m.cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	surse := fstest.MapFS{}
	intrari, err := fs.ReadDir(bd.Migrari, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, intrare := range intrari {
		if intrare.Name() >= "0014_" || strings.HasPrefix(intrare.Name(), "0012_") {
			continue
		}
		date, err := fs.ReadFile(bd.Migrari, intrare.Name())
		if err != nil {
			t.Fatal(err)
		}
		surse[intrare.Name()] = &fstest.MapFile{Data: date}
	}
	if _, err := bd.Migreaza(context.Background(), db, surse); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`UPDATE definitii_rapoarte SET antet=json_set(antet,'$.centru.ro','Raport pastrat'),revizie=4 WHERE id='raport-pontaj'`); err != nil {
		t.Fatal(err)
	}
	versiune, inAsteptare, err := bd.StareMigrari(context.Background(), db, bd.Migrari)
	if err != nil || versiune != 13 || !inAsteptare {
		t.Fatalf("starea actualizarii 0.20: %d %v %v", versiune, inAsteptare, err)
	}
	if n, err := bd.Migreaza(context.Background(), db, bd.Migrari); err != nil || n != 30 {
		t.Fatalf("actualizarea 0.20: %d %v", n, err)
	}
	d := initiala(t, &mediu{s: Nou(db)}, "raport-pontaj")
	if d.Revizie != "4" || d.Antet.Centru["ro"] != "Raport pastrat" {
		t.Fatal("actualizarea a pierdut definitia existenta")
	}
	var coloane int
	if err := db.QueryRow(`SELECT count(*) FROM pragma_table_info('absente_zile') WHERE name='stare_aprobare'`).Scan(&coloane); err != nil || coloane != 1 {
		t.Fatalf("cererile de concediu lipsesc: %d %v", coloane, err)
	}
}
