// Package fonturi citeste fonturile instalate in Windows pentru rapoartele
// pdfme (D088). Se citesc doar antetele (tabelele name si OS/2) pentru lista;
// datele unui font se citesc numai la cerere, dupa nume, nu dupa cale.
package fonturi

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"unicode/utf16"

	"pontix/internal/i18n"
)

// Font descrie un font utilizabil in PDF.
type Font struct {
	Nume    string `json:"nume"`
	Familie string `json:"familie"`
	Stil    string `json:"stil"`
	// Subset este fals cand licenta fontului cere includerea lui intreaga.
	Subset bool `json:"subset"`
	fisier string
	index  int
}

const maxFont = 64 << 20

// Biti fsType din tabelul OS/2 (specificatia OpenType).
const (
	fsRestrictionat = 0x0002
	fsFaraSubset    = 0x0100
	fsDoarBitmap    = 0x0200
)

// Numele rezervate fonturilor Go incluse in aplicatie.
var rezervate = map[string]bool{"go": true, "gobold": true, "goitalic": true, "gobolditalic": true}

type Serviciu struct {
	mu       sync.Mutex
	dirs     []string
	lista    []Font
	incarcat bool
}

// Nou foloseste directoarele date sau, fara ele, pe cele ale sistemului.
func Nou(dirs ...string) *Serviciu {
	if len(dirs) == 0 {
		dirs = DirectoareSistem()
	}
	return &Serviciu{dirs: dirs}
}

// DirectoareSistem: fonturile Windows si cele instalate doar pentru utilizator.
func DirectoareSistem() []string {
	var d []string
	if w := os.Getenv("WINDIR"); w != "" {
		d = append(d, filepath.Join(w, "Fonts"))
	}
	if l := os.Getenv("LOCALAPPDATA"); l != "" {
		d = append(d, filepath.Join(l, "Microsoft", "Windows", "Fonts"))
	}
	return d
}

// Lista intoarce fonturile permise, sortate dupa nume, citite o singura data.
func (s *Serviciu) Lista() []Font {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.incarcat {
		s.lista = scaneaza(s.dirs)
		s.incarcat = true
	}
	return append([]Font(nil), s.lista...)
}

// Date intoarce fontul cerut ca fisier TrueType/OpenType independent.
func (s *Serviciu) Date(nume string) ([]byte, Font, error) {
	for _, f := range s.Lista() {
		if f.Nume == nume {
			b, err := citesteFont(f.fisier, f.index)
			if err != nil {
				return nil, f, &i18n.Eroare{Cod: "fonturi.citire", Param: map[string]string{"nume": nume}}
			}
			return b, f, nil
		}
	}
	return nil, Font{}, &i18n.Eroare{Cod: "fonturi.inexistent", Param: map[string]string{"nume": nume}}
}

func scaneaza(dirs []string) []Font {
	vazute := map[string]bool{}
	var rez []Font
	for _, d := range dirs {
		intrari, err := os.ReadDir(d)
		if err != nil {
			continue
		}
		for _, e := range intrari {
			ext := strings.ToLower(filepath.Ext(e.Name()))
			if e.IsDir() || (ext != ".ttf" && ext != ".otf" && ext != ".ttc" && ext != ".otc") {
				continue
			}
			cale := filepath.Join(d, e.Name())
			fonturi, err := antete(cale)
			if err != nil {
				continue
			}
			for _, f := range fonturi {
				cheie := strings.ToLower(f.Nume)
				if f.Nume == "" || vazute[cheie] || rezervate[cheie] {
					continue
				}
				vazute[cheie] = true
				rez = append(rez, f)
			}
		}
	}
	sort.Slice(rez, func(i, j int) bool { return strings.ToLower(rez[i].Nume) < strings.ToLower(rez[j].Nume) })
	return rez
}

var errFormat = errors.New("format font necunoscut")

type tabel struct{ offset, lungime uint32 }

// offseturi intoarce pozitia fiecarui font: unul pentru .ttf/.otf, mai multe pentru colectii.
func offseturi(r io.ReaderAt) ([]uint32, error) {
	var h [12]byte
	if _, err := r.ReadAt(h[:], 0); err != nil {
		return nil, err
	}
	if string(h[:4]) != "ttcf" {
		return []uint32{0}, nil
	}
	n := binary.BigEndian.Uint32(h[8:])
	if n == 0 || n > 256 {
		return nil, errFormat
	}
	b := make([]byte, 4*n)
	if _, err := r.ReadAt(b, 12); err != nil {
		return nil, err
	}
	o := make([]uint32, n)
	for i := range o {
		o[i] = binary.BigEndian.Uint32(b[4*i:])
	}
	return o, nil
}

// tabele citeste directorul de tabele al fontului de la pozitia data.
func tabele(r io.ReaderAt, start uint32) (antet [12]byte, t map[string]tabel, ordine []string, err error) {
	if _, err = r.ReadAt(antet[:], int64(start)); err != nil {
		return
	}
	switch binary.BigEndian.Uint32(antet[:4]) {
	case 0x00010000, 0x4F54544F, 0x74727565: // TrueType, 'OTTO', 'true'
	default:
		err = errFormat
		return
	}
	n := int(binary.BigEndian.Uint16(antet[4:]))
	if n == 0 || n > 512 {
		err = errFormat
		return
	}
	b := make([]byte, 16*n)
	if _, err = r.ReadAt(b, int64(start)+12); err != nil {
		return
	}
	t = map[string]tabel{}
	for i := 0; i < n; i++ {
		rec := b[16*i:]
		tag := string(rec[:4])
		t[tag] = tabel{binary.BigEndian.Uint32(rec[8:]), binary.BigEndian.Uint32(rec[12:])}
		ordine = append(ordine, tag)
	}
	return
}

func citeste(r io.ReaderAt, t tabel, max uint32) ([]byte, error) {
	if t.lungime == 0 || t.lungime > max {
		return nil, errFormat
	}
	b := make([]byte, t.lungime)
	_, err := r.ReadAt(b, int64(t.offset))
	return b, err
}

func antete(cale string) ([]Font, error) {
	f, err := os.Open(cale)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if st, err := f.Stat(); err != nil || st.Size() > maxFont {
		return nil, errFormat
	}
	pozitii, err := offseturi(f)
	if err != nil {
		return nil, err
	}
	var rez []Font
	for i, p := range pozitii {
		_, t, _, err := tabele(f, p)
		if err != nil {
			continue
		}
		fsType := uint16(0)
		if os2, ok := t["OS/2"]; ok {
			b, err := citeste(f, os2, 1<<16)
			if err != nil || len(b) < 10 {
				continue
			}
			fsType = binary.BigEndian.Uint16(b[8:])
		}
		if !permis(fsType) {
			continue
		}
		nt, ok := t["name"]
		if !ok {
			continue
		}
		b, err := citeste(f, nt, 1<<20)
		if err != nil {
			continue
		}
		nm := nume(b)
		familie := alege(nm, 16, 1)
		stil := alege(nm, 17, 2)
		complet := alege(nm, 4)
		if complet == "" {
			complet = strings.TrimSpace(familie + " " + stil)
		}
		rez = append(rez, Font{Nume: complet, Familie: familie, Stil: stil, Subset: fsType&fsFaraSubset == 0, fisier: cale, index: i})
	}
	return rez, nil
}

// permis aplica fsType: fara licenta restrictiva si fara fonturi doar bitmap.
// Bitii 0-3 se exclud reciproc; la fonturile vechi cu mai multi biti, cel mai
// permisiv castiga, deci restrictia se aplica doar cand este singura.
func permis(fsType uint16) bool {
	return fsType&0x000F != fsRestrictionat && fsType&fsDoarBitmap == 0
}

func alege(n map[uint16]string, id ...uint16) string {
	for _, i := range id {
		if v := n[i]; v != "" {
			return v
		}
	}
	return ""
}

// nume citeste tabelul name; prefera Windows Unicode engleza, apoi orice
// Windows Unicode, apoi Macintosh Roman.
func nume(b []byte) map[uint16]string {
	rez := map[uint16]string{}
	prioritate := map[uint16]int{}
	if len(b) < 6 {
		return rez
	}
	n := int(binary.BigEndian.Uint16(b[2:]))
	baza := int(binary.BigEndian.Uint16(b[4:]))
	for i := 0; i < n && 6+12*i+12 <= len(b); i++ {
		r := b[6+12*i:]
		platforma, cod, limba, id := binary.BigEndian.Uint16(r), binary.BigEndian.Uint16(r[2:]), binary.BigEndian.Uint16(r[4:]), binary.BigEndian.Uint16(r[6:])
		lung, off := int(binary.BigEndian.Uint16(r[8:])), int(binary.BigEndian.Uint16(r[10:]))
		if baza+off+lung > len(b) {
			continue
		}
		s := b[baza+off : baza+off+lung]
		var text string
		p := 0
		switch {
		case platforma == 3 && (cod == 1 || cod == 10) && lung%2 == 0:
			u := make([]uint16, lung/2)
			for j := range u {
				u[j] = binary.BigEndian.Uint16(s[2*j:])
			}
			text = string(utf16.Decode(u))
			p = 2
			if limba == 0x0409 {
				p = 3
			}
		case platforma == 1 && cod == 0:
			text = string(s)
			p = 1
		default:
			continue
		}
		if text = strings.TrimSpace(text); text != "" && p > prioritate[id] {
			rez[id], prioritate[id] = text, p
		}
	}
	return rez
}

// citesteFont intoarce fontul de la indexul dat ca fisier independent; din
// colectii se copiaza tabelele fontului intr-un fisier nou.
func citesteFont(cale string, index int) ([]byte, error) {
	f, err := os.Open(cale)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil || st.Size() > maxFont {
		return nil, errFormat
	}
	pozitii, err := offseturi(f)
	if err != nil || index >= len(pozitii) {
		return nil, errFormat
	}
	if len(pozitii) == 1 && pozitii[0] == 0 {
		var h [4]byte
		if _, err := f.ReadAt(h[:], 0); err != nil {
			return nil, err
		}
		if string(h[:]) != "ttcf" {
			return io.ReadAll(io.LimitReader(f, maxFont))
		}
	}
	return extrage(f, pozitii[index])
}

func extrage(r io.ReaderAt, start uint32) ([]byte, error) {
	antet, t, ordine, err := tabele(r, start)
	if err != nil {
		return nil, err
	}
	n := len(ordine)
	out := make([]byte, 12+16*n)
	copy(out, antet[:])
	// Recordurile se pastreaza in ordinea originala (sortate dupa eticheta).
	recorduri := make([]byte, 16*n)
	if _, err := r.ReadAt(recorduri, int64(start)+12); err != nil {
		return nil, err
	}
	for i, tag := range ordine {
		var date []byte
		if t[tag].lungime > 0 {
			if date, err = citeste(r, t[tag], maxFont); err != nil {
				return nil, err
			}
		}
		rec := out[12+16*i:]
		copy(rec, recorduri[16*i:16*i+16])
		binary.BigEndian.PutUint32(rec[8:], uint32(len(out)))
		out = append(out, date...)
		for len(out)%4 != 0 {
			out = append(out, 0)
		}
		if len(out) > maxFont {
			return nil, errFormat
		}
	}
	return out, nil
}
