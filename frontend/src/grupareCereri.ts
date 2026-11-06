import type {datebaza} from '../wailsjs/go/models';

// GrupCerere: zilele unei cereri (document sau numar de inregistrare) sau, fara cerere,
// un interval de zile apropiate cu acelasi angajat, cod si stare.
export type GrupCerere = {cheie: string; zile: datebaza.ZiAbsenta[]; faraCerere: boolean};
export function grupeazaPeCereri(zile: datebaza.ZiAbsenta[]): GrupCerere[] {
    const ordonate = [...zile].sort((a, b) => a.angajat_id.localeCompare(b.angajat_id) || a.data.localeCompare(b.data));
    const grupuri = new Map<string, GrupCerere>();
    const ultimulFaraCerere = new Map<string, GrupCerere>();
    for (const z of ordonate) {
        if (z.document_id || z.nr_inregistrare) {
            const cheie = z.document_id ? `d:${z.document_id}` : `n:${z.angajat_id}/${z.nr_inregistrare}`;
            const g = grupuri.get(cheie) ?? {cheie, zile: [], faraCerere: false};
            g.zile.push(z); grupuri.set(cheie, g);
            continue;
        }
        const tip = `${z.angajat_id}/${z.cod}/${z.stare_aprobare || 'IN_ASTEPTARE'}/${z.activ}`;
        const g = ultimulFaraCerere.get(tip);
        const ultima = g?.zile[g.zile.length - 1];
        // Zile apropiate (cel mult 3 zile distanta, adica peste un weekend) raman in acelasi interval.
        if (g && ultima && (Date.parse(z.data) - Date.parse(ultima.data)) <= 3 * 86400000) { g.zile.push(z); continue; }
        const nou = {cheie: `z:${tip}/${z.data}`, zile: [z], faraCerere: true};
        grupuri.set(nou.cheie, nou); ultimulFaraCerere.set(tip, nou);
    }
    return [...grupuri.values()].sort((a, b) => a.zile[0].data.localeCompare(b.zile[0].data) || a.zile[0].angajat_id.localeCompare(b.zile[0].angajat_id));
}
