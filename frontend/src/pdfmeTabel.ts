// F09-T22: stil pe rand si pe celula in tabelele pdfme (bold pentru weekend si
// sarbatori, celulele editate in copia raportului, D085). Tabelul pdfme nu are
// un punct de extindere pentru stiluri pe rand, asa ca la generare inlocuim doar
// desenarea lui (table.pdf). Editorul, panoul de proprietati si impartirea pe
// pagini raman cele din pdfme; un tabel fara stiluri se deseneaza cu codul original.
import type {Font, PDFRenderProps, Plugin, Schema} from '@pdfme/common';

// Stilul unei celule sau al unui rand, in termenii StilText din Go.
export interface StilCelula {
    aldin?: boolean;
    cursiv?: boolean;
    font_pt?: number;
    aliniere?: 'L' | 'C' | 'R';
    fundal?: string;
}

// StiluriTabel vine de la Go (intrarea pontix_stiluri). Indicii randurilor sunt
// in tot corpul tabelului, nu pe pagina; cheile celulelor sunt "rand:coloana".
export interface StiluriTabel {
    cap?: StilCelula;
    corp?: StilCelula;
    randuri?: Record<string, StilCelula>;
    celule?: Record<string, StilCelula>;
}

export const INTRARE_STILURI = 'pontix_stiluri';

type Box = {top: number; right: number; bottom: number; left: number};
type CellStyle = {fontName?: string; alignment: string; verticalAlignment: string; fontSize: number; lineHeight: number; characterSpacing: number; fontColor: string; backgroundColor: string; borderColor: string; borderWidth: Box; padding: Box};
type TabelSchema = Schema & {
    head: string[];
    headWidthPercentages: number[];
    showHead?: boolean;
    repeatHead?: boolean;
    __isSplit?: boolean;
    tableStyles: {borderColor: string; borderWidth: number};
    headStyles: CellStyle;
    bodyStyles: CellStyle & {alternateBackgroundColor?: string};
    columnStyles?: {alignment?: Record<number, string>; verticalAlignment?: Record<number, string>; cellType?: Record<number, string>};
};

const ALINIERI: Record<string, string> = {L: 'left', C: 'center', R: 'right'};
const VARIANTE_GO: Record<string, [string, string, string, string]> = {
    // [normal, aldin, cursiv, aldin cursiv]
    Go: ['Go', 'GoBold', 'GoItalic', 'GoBoldItalic'],
    GoBold: ['Go', 'GoBold', 'GoItalic', 'GoBoldItalic'],
    GoItalic: ['Go', 'GoBold', 'GoItalic', 'GoBoldItalic'],
    GoBoldItalic: ['Go', 'GoBold', 'GoItalic', 'GoBoldItalic'],
};

function bazaFont(nume: string): string {
    return nume.replace(/\s+(Bold Italic|Bold|Italic|Regular)$/i, '');
}

// numeVarianta da numele fontului aldin/cursiv cerut pentru un font dat:
// familia Go are variante fixe; pentru fonturile Windows se cauta "Nume Bold",
// "Nume Italic", "Nume Bold Italic". Daca varianta lipseste, ramane fontul initial.
export function numeVarianta(nume: string | undefined, aldin: boolean, cursiv: boolean, fonturi: Font): string | undefined {
    if (!nume) return nume;
    const go = VARIANTE_GO[nume];
    if (go) return go[(aldin ? 1 : 0) + (cursiv ? 2 : 0)];
    if (!aldin && !cursiv) return nume;
    const baza = bazaFont(nume);
    const sufix = aldin && cursiv ? ' Bold Italic' : aldin ? ' Bold' : ' Italic';
    return fonturi[baza + sufix] ? baza + sufix : nume;
}

// fonturiNecesare: variantele aldin/cursiv de care are nevoie tabelul, ca sa
// poata fi incarcate inainte de generare (fonturile Windows se incarca la cerere).
export function variantaNecesare(fonturiTabel: string[], stiluri: StiluriTabel): string[] {
    const toate = [stiluri.cap, stiluri.corp, ...Object.values(stiluri.randuri || {}), ...Object.values(stiluri.celule || {})].filter(Boolean) as StilCelula[];
    const rez = new Set<string>();
    for (const f of fonturiTabel) {
        if (VARIANTE_GO[f]) continue;
        const baza = bazaFont(f);
        for (const s of toate) {
            if (s.aldin && s.cursiv) rez.add(baza + ' Bold Italic');
            else if (s.aldin) rez.add(baza + ' Bold');
            else if (s.cursiv) rez.add(baza + ' Italic');
        }
    }
    return [...rez];
}

export function areStiluri(s?: StiluriTabel | null): boolean {
    return !!s && !!(s.cap || s.corp || Object.keys(s.randuri || {}).length || Object.keys(s.celule || {}).length);
}

// stilFinal combina stilul pdfme al celulei cu stilurile Pontix, in ordinea
// sectiune -> rand -> celula (cel mai specific castiga).
export function stilFinal(baza: CellStyle, straturi: (StilCelula | undefined)[], fonturi: Font): CellStyle {
    const rez = {...baza};
    let aldin: boolean | undefined, cursiv: boolean | undefined;
    for (const s of straturi) {
        if (!s) continue;
        if (s.aldin !== undefined) aldin = s.aldin;
        if (s.cursiv !== undefined) cursiv = s.cursiv;
        if (s.font_pt && s.font_pt > 0) rez.fontSize = s.font_pt;
        if (s.aliniere && ALINIERI[s.aliniere]) rez.alignment = ALINIERI[s.aliniere];
        if (s.fundal) rez.backgroundColor = s.fundal;
    }
    if (aldin !== undefined || cursiv !== undefined) {
        const go = baza.fontName ? VARIANTE_GO[baza.fontName] : undefined;
        const eraAldin = !!baza.fontName && (/Bold/.test(baza.fontName) || go?.[1] === baza.fontName);
        const eraCursiv = !!baza.fontName && /Italic/.test(baza.fontName);
        rez.fontName = numeVarianta(baza.fontName, aldin ?? eraAldin, cursiv ?? eraCursiv, fonturi);
    }
    return rez;
}

function cutie(v: unknown, implicit = 0): Box {
    if (typeof v === 'number') return {top: v, right: v, bottom: v, left: v};
    const b = (v || {}) as Partial<Box>;
    return {top: b.top ?? implicit, right: b.right ?? implicit, bottom: b.bottom ?? implicit, left: b.left ?? implicit};
}

type Desenatori = {text: Plugin['pdf']; rectangle: Plugin['pdf']; line: Plugin['pdf']};
type Inaltimi = (value: string, args: {schema: Schema; basePdf: unknown; options: unknown; _cache: Map<string | number, unknown>}) => Promise<number[]>;
type Interval = (schema: Schema) => {start: number; end?: number} | undefined;

async function deseneazaCelula(arg: PDFRenderProps<Schema>, d: Desenatori, valoare: string, x: number, y: number, w: number, h: number, st: CellStyle) {
    const bw = cutie(st.borderWidth), pad = cutie(st.padding);
    const lucrari: Promise<void>[] = [];
    if (st.backgroundColor) lucrari.push(Promise.resolve(d.rectangle({...arg, schema: {name: '', type: 'rectangle', position: {x, y}, width: w, height: h, borderWidth: 0, borderColor: '', color: st.backgroundColor, readOnly: true} as Schema})));
    const linie = (lx: number, ly: number, lw: number, lh: number) => {
        if (lw > 0 && lh > 0) lucrari.push(Promise.resolve(d.line({...arg, schema: {name: '', type: 'line', position: {x: lx, y: ly}, width: lw, height: lh, color: st.borderColor, readOnly: true} as Schema})));
    };
    linie(x, y, w, bw.top);
    linie(x + w - bw.right, y, bw.right, h);
    linie(x, y + h - bw.bottom, w, bw.bottom);
    linie(x, y, bw.left, h);
    await Promise.all(lucrari);
    const sx = bw.left + pad.left, sy = bw.top + pad.top;
    await d.text({...arg, value: valoare, schema: {
        name: '', type: 'text', position: {x: x + sx, y: y + sy},
        width: Math.max(0, w - sx - bw.right - pad.right), height: Math.max(0, h - sy - bw.bottom - pad.bottom),
        fontName: st.fontName, alignment: st.alignment, verticalAlignment: st.verticalAlignment, fontSize: st.fontSize,
        lineHeight: st.lineHeight, characterSpacing: st.characterSpacing, fontColor: st.fontColor,
        backgroundColor: '',
    } as Schema});
}

// tabelCuStiluri intoarce pluginul de tabel pentru generare: acelasi editor si
// aceeasi impartire pe pagini, dar desenarea aplica stilurile primite.
export function tabelCuStiluri(original: Plugin, desenatori: Desenatori, inaltimi: Inaltimi, interval: Interval, stiluri: Record<string, StiluriTabel>): Plugin {
    return {...original, pdf: async (arg: PDFRenderProps<Schema>) => {
        const st = stiluri[arg.schema.name];
        const schema = arg.schema as TabelSchema;
        const imagini = Object.values(schema.columnStyles?.cellType || {}).includes('image');
        if (!areStiluri(st) || imagini) return original.pdf(arg);
        const fonturi = (arg.options?.font || {}) as Font;
        const corp: string[][] = (() => { try { const v = JSON.parse(typeof arg.value === 'string' ? arg.value : JSON.stringify(arg.value ?? [])); return Array.isArray(v) ? v.map((r: unknown[]) => (Array.isArray(r) ? r.map(c => String(c ?? '')) : [])) : []; } catch { return []; } })();
        const range = interval(schema);
        const start = range?.start ?? 0, end = range?.end ?? corp.length;
        const randuri = corp.slice(start, end);
        // Ca in pdfme: pe paginile continuate capul apare doar daca se repeta.
        const cap = schema.showHead !== false && (!schema.__isSplit || schema.repeatHead === true);
        const valoare = typeof arg.value === 'string' ? arg.value : JSON.stringify(arg.value ?? []);
        const h = await inaltimi(valoare, {schema: {...schema, showHead: cap, repeatHead: false} as Schema, basePdf: arg.basePdf, options: arg.options, _cache: arg._cache});
        const nCap = cap ? 1 : 0;
        const inaltimiCap = cap ? h.slice(0, 1) : [];
        const inaltimiCorp = cap ? h.slice(1) : h.slice(1);
        const latimi = schema.headWidthPercentages.map(p => (schema.width * p) / 100);
        const coloana = (cs: CellStyle, i: number): CellStyle => ({...cs,
            alignment: schema.columnStyles?.alignment?.[i] ?? cs.alignment,
            verticalAlignment: schema.columnStyles?.verticalAlignment?.[i] ?? cs.verticalAlignment});
        let y = schema.position.y;
        const x0 = schema.position.x;
        const d = desenatori;
        if (nCap) {
            let x = x0;
            for (let c = 0; c < latimi.length; c++) {
                await deseneazaCelula(arg, d, schema.head[c] ?? '', x, y, latimi[c], inaltimiCap[0], stilFinal(coloana(schema.headStyles, c), [st!.cap], fonturi));
                x += latimi[c];
            }
            y += inaltimiCap[0];
        }
        const alternant = schema.bodyStyles.alternateBackgroundColor;
        for (let r = 0; r < randuri.length; r++) {
            const g = start + r;
            const hr = inaltimiCorp[r] ?? 0;
            // Randurile pare (in tot tabelul) primesc fundalul alternant, ca in pdfme.
            const bazaRand = {...schema.bodyStyles, backgroundColor: g % 2 === 0 && alternant ? alternant : schema.bodyStyles.backgroundColor};
            let x = x0;
            for (let c = 0; c < latimi.length; c++) {
                const s = stilFinal(coloana(bazaRand, c), [st!.corp, st!.randuri?.[String(g)], st!.celule?.[g + ':' + c]], fonturi);
                await deseneazaCelula(arg, d, randuri[r][c] ?? '', x, y, latimi[c], hr, s);
                x += latimi[c];
            }
            y += hr;
        }
        const tb = schema.tableStyles;
        if (tb?.borderWidth && tb.borderColor) {
            await d.rectangle({...arg, schema: {name: '', type: 'rectangle', position: {x: x0, y: schema.position.y}, width: latimi.reduce((a, b) => a + b, 0), height: y - schema.position.y, borderWidth: tb.borderWidth, borderColor: tb.borderColor, color: '', readOnly: true} as Schema});
        }
    }};
}
