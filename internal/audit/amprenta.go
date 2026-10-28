package audit

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"
)

// tabeleLuna sunt tabelele ale căror rânduri dintr-o lună intră în amprenta ei (D044 punctul 6).
var tabeleLuna = []string{"absente_zile", "planificare", "pontaj_ciorne"}

// AmprentaLuna calculează amprenta datelor unei luni: pontajul, absențele și planificarea,
// cu rândurile ordonate după cheia primară. Rezultatul are forma XXXX-XXXX-XXXX-XXXX.
func AmprentaLuna(ctx context.Context, q Interogator, an, luna int) (string, error) {
	deLa := fmt.Sprintf("%04d-%02d-01", an, luna)
	panaLa := fmt.Sprintf("%04d-%02d-31", an, luna)
	h := sha256.New()
	h.Write([]byte("pontix/amprenta/v1\x00"))
	var b [8]byte
	binary.BigEndian.PutUint32(b[:4], uint32(an))
	binary.BigEndian.PutUint32(b[4:], uint32(luna))
	h.Write(b[:])
	for _, t := range tabeleLuna {
		col, err := coloane(ctx, q, t)
		if err != nil {
			return "", err
		}
		chei, err := textColoana(ctx, q, `SELECT name FROM pragma_table_info(?) WHERE pk > 0 ORDER BY pk`, t)
		if err != nil {
			return "", err
		}
		if len(col) == 0 || len(chei) == 0 {
			return "", fmt.Errorf("audit: tabelul %s lipsește sau nu are cheie primară", t)
		}
		r, err := q.QueryContext(ctx, "SELECT pontix_h("+literal("luna/"+t)+", "+lista("", col)+") FROM main."+cita(t)+
			" WHERE data BETWEEN ? AND ? ORDER BY "+lista("", chei), deLa, panaLa)
		if err != nil {
			return "", err
		}
		scrieOctetiSimplu(h, []byte(t))
		for r.Next() {
			var x []byte
			if err := r.Scan(&x); err != nil {
				r.Close()
				return "", err
			}
			h.Write(x)
		}
		r.Close()
		if err := r.Err(); err != nil {
			return "", err
		}
	}
	return FormateazaAmprenta(h.Sum(nil)), nil
}

func scrieOctetiSimplu(h hash.Hash, o []byte) {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(len(o)))
	h.Write(b[:])
	h.Write(o)
}

// FormateazaAmprenta păstrează primii 8 octeți, în grupuri de câte 4 caractere hex.
func FormateazaAmprenta(sum []byte) string {
	s := strings.ToUpper(hex.EncodeToString(sum[:8]))
	return s[0:4] + "-" + s[4:8] + "-" + s[8:12] + "-" + s[12:16]
}
