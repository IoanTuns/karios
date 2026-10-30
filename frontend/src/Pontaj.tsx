import {memo, useEffect, useLayoutEffect, useMemo, useRef, useState, type ClipboardEvent, type KeyboardEvent, type MouseEvent, type RefObject} from 'react';
import {useTranslation} from 'react-i18next';
import {autorJurnal} from './autorJurnal';
import {semnaleazaModificari} from './modificariNesalvate';
import {UltimaModificareCelula, CitesteCalendar, CitesteSchita, ListaSchimburi, SalveazaSchita, PrevizualizeazaSchita, PrevizualizeazaPrecompletare, ExportaPontajExcel, CitesteStareLuna, SeteazaStareLuna} from '../wailsjs/go/main/App';
import {calcul, datebaza} from '../wailsjs/go/models';
import {AmprentaLuna} from './Audit';
import {Dialog, denumireSectie, denumireDepartament} from './Formulare';
import {programPlanificat} from './Planificare';
import {notifica} from './Toast';
import {formatData} from './formatare';

type Pozitie = {rand: number; zi: number};
type Selectie = {inceput: Pozitie; sfarsit: Pozitie};
export type ContextPontaj = {
    an: number; luna: number; sectie: string; departament: string; categorieProgram: string; cautare: string; filtruRanduri: string;
    selectie: {inceput: {angajatId: string; data: string}; sfarsit: {angajatId: string; data: string}} | null;
    derulareX: number; derulareY: number;
};
const totaluri = ['ore_lucrate', 'T', 'ZNL', 'CO', 'CEV', 'CFS', 'ore_norma', 'form_prof', 'CM', 'L_FIN', 'deleg', 'ore_garda'];
const cheie = (id: string, data: string) => JSON.stringify([id, data]);
const oreText = (minute: number) => String(Number((minute/60).toFixed(4)));
const valoareText = (v?: calcul.ValoarePontaj) => v?.definita ? String(Number((v.sutimi/(v.unitate==='ORE'?6000:100)).toFixed(4))) : '—';
const normalizat = (text: string) => text.normalize('NFD').replace(/\p{Diacritic}/gu, '').toLocaleLowerCase();
function manual(text: string): string | null {
    if(!text.trim())return '';
    if(!/^\d+(?:[.,]\d+)?$/.test(text.trim()))return null;
    const ore=Number(text.trim().replace(',','.')),minute=Math.round(ore*60);
    return ore<=48&&Math.abs(ore*60-minute)<0.00001?`@M${minute}`:null;
}

// ActiuniCelula: handler-ele celulelor, ținute într-un obiect stabil (ref), ca rândurile memorate
// să nu se redeseneze doar pentru că funcțiile din Pontaj sunt noi la fiecare randare (PERF-1).
type ActiuniCelula = {
    focus: (r: number, z: number) => void;
    mouseDown: (e: MouseEvent<HTMLInputElement>, r: number, z: number) => void;
    mouseEnter: (e: MouseEvent<HTMLInputElement>, r: number, z: number) => void;
    keyDown: (e: KeyboardEvent<HTMLInputElement>, r: number, z: number) => void;
    dublu: () => void;
};

type PropsRand = {
    a: datebaza.Angajat; r: number; zile: string[]; weekend: boolean[]; sarbatori: Map<string, string>;
    // selDe/selPana: zilele selectate din acest rând, -1 dacă rândul nu e în selecție; așa, un pas
    // de selecție redesenează doar rândurile care intră sau ies din ea.
    selDe: number; selPana: number;
    modificari: Map<string, datebaza.CelulaSchita>; originale: Map<string, datebaza.CelulaSchita>;
    absente: Map<string, datebaza.ZiAbsenta>; inlocuiri: Map<string, datebaza.ZiAbsenta[]>;
    oreCalculate: Map<string, calcul.OreZi>; valori?: Record<string, calcul.ValoarePontaj>;
    reguli?: calcul.RegulaColoana[]; angajati: datebaza.Angajat[];
    ocupat: boolean; calculare: boolean; areCalcul: boolean; actiuni: RefObject<ActiuniCelula>;
};

const RandPontaj = memo(function RandPontaj({a, r, zile, weekend, sarbatori, selDe, selPana, modificari, originale, absente, inlocuiri, oreCalculate, valori, reguli, angajati, ocupat, calculare, areCalcul, actiuni}: PropsRand) {
    const {t} = useTranslation();
    const numeColeg = (id: string) => { const coleg = angajati.find(x => x.id === id); return coleg ? `${coleg.nume} ${coleg.prenume}` : id; };
    const valoareBruta = (z: number) => { const k = cheie(a.id, zile[z]); return modificari.get(k)?.valoare ?? originale.get(k)?.valoare ?? ''; };
    return <tr><th scope="row" className="fix-nume"><strong>{a.nume} {a.prenume}</strong>{!a.activ && <small>{t('ui.inactiv')}</small>}</th><td className="fix-marca">{a.marca}</td>
        {zile.map((data, z) => {
            const selectata = z >= selDe && z <= selPana;
            const k = cheie(a.id, data);
            const absenta = absente.get(k);
            const stareAbsenta = absenta?.stare_aprobare || 'IN_ASTEPTARE';
            const titluAbsenta = absenta ? [t('absente.celula_protejata'), `${t(`absente.aprobare_${stareAbsenta}`)}${absenta.aprobat_de ? ` (${absenta.aprobat_de})` : ''}`, absenta.nr_inregistrare ? `${t('absente.nr_inregistrare')}: ${absenta.nr_inregistrare}` : '', absenta.inlocuitor_id && angajati.some(x => x.id === absenta.inlocuitor_id) ? `${t('absente.inlocuitor')}: ${numeColeg(absenta.inlocuitor_id)}` : ''].filter(Boolean).join(' · ') : '';
            const colegiInlocuiti = inlocuiri.get(k);
            const titluInlocuire = colegiInlocuiti?.length ? t('pontaj.inlocuieste_pe', {persoane: colegiInlocuiti.map(z => `${numeColeg(z.angajat_id)} (${z.cod})`).join(', ')}) : '';
            const calc = oreCalculate.get(k);
            const minuteZi = calc?.minute ?? 0;
            const areLucruZi = !absenta && (valoareBruta(z).trim() !== '' || minuteZi > 0);
            const dupaG24 = z > 0 && valoareBruta(z - 1).trim().toUpperCase() === 'G24' && areLucruZi;
            const peste24h = minuteZi > 1440;
            const titluAvertisment = [dupaG24 ? t('pontaj.avertisment_dupa_g24') : '', peste24h ? t('pontaj.avertisment_peste_24h', {ore: oreText(minuteZi)}) : ''].filter(Boolean).join(' · ');
            const valoare = calculare ? '…' : calc ? (calc.calculabil ? oreText(calc.minute) : '—') : areCalcul ? '' : '—';
            return <td key={data} className={`${sarbatori.has(data) ? 'zi-sarbatoare' : weekend[z] ? 'zi-weekend' : ''} ${selectata ? 'celula-selectata' : ''} ${modificari.has(k) ? 'celula-modificata' : ''} ${colegiInlocuiti?.length ? 'celula-inlocuitor' : ''} ${titluAvertisment ? 'celula-avertisment' : ''}`}>
            <input data-rand={r} data-zi={z} aria-label={t('pontaj.celula', {nume: `${a.nume} ${a.prenume}`, data})} title={[
                t('ui.tooltip_marca', {nume: `${a.nume} ${a.prenume}`, marca: a.marca || '—'}),
                t('ui.tooltip_data', {data}) + (weekend[z] ? ` [${t('ui.tooltip_weekend')}]` : '') + (sarbatori.has(data) ? ` [${sarbatori.get(data)}]` : ''),
                t('ui.tooltip_schimb', {schimb: valoare || t('ui.tooltip_liber')}),
                titluAbsenta,
                titluInlocuire,
                titluAvertisment,
                calc?.surse.length ? t('pontaj.surse_ore', {date: calc.surse.join(', ')}) : ''
            ].filter(Boolean).join(' · ')} disabled={ocupat} readOnly className={absenta ? `celula-absenta celula-absenta-${stareAbsenta}` : ''} value={valoare}
                onFocus={() => actiuni.current.focus(r, z)}
                onMouseDown={e => actiuni.current.mouseDown(e, r, z)}
                onMouseEnter={e => actiuni.current.mouseEnter(e, r, z)}
                onKeyDown={e => actiuni.current.keyDown(e, r, z)} onDoubleClick={() => actiuni.current.dublu()}/></td>;
        })}
        {totaluri.map(c => {const v=valori?.[c],regula=reguli?.find(r=>r.coloana===c),titluTotal=!v?.definita?t('pontaj.necalculat'):regula?.sursa==='ORE_NORMA'?t('pontaj.ore_norm_explicatie'):t(v.unitate==='ZILE'?'pontaj.unitate_zile':'pontaj.unitate_ore');return <td key={c} className={`${!v?.definita?'total-necalculat':''} ${c==='T'?'fix-total-t':''}`.trim()} title={titluTotal}>{calculare?'…':v?valoareText(v):areCalcul&&regula?.sursa!=='NEDEFINITA'?'0':'—'}</td>;})}</tr>;
});

export function Pontaj({date, context, onModificari, onOcupat, onDeschideAbsente}: {date: datebaza.Date; context: RefObject<ContextPontaj | null>; onModificari: (modificat: boolean) => void; onOcupat: (ocupat: boolean) => void; onDeschideAbsente?: (an: number, angajatId: string) => void}) {
    const {t, i18n} = useTranslation();
    const acum = useRef(new Date());
    const initial = useRef(context.current);
    const restaurare = useRef(Boolean(initial.current));
    const primulFiltru = useRef(true);
    const stareCurenta = useRef(context.current);
    const [an, setAn] = useState(initial.current?.an ?? acum.current.getFullYear());
    const [luna, setLuna] = useState(initial.current?.luna ?? acum.current.getMonth() + 1);
    const [anText, setAnText] = useState(String(an));
    const [sectie, setSectie] = useState(initial.current?.sectie ?? '');
    const [departament, setDepartament] = useState(initial.current?.departament ?? '');
    const [categorieProgram,setCategorieProgram]=useState(initial.current?.categorieProgram??'');
    const [plan,setPlan]=useState<calcul.PlanPrecompletare|null>(null);
    const [numaiSelectiaPlan,setNumaiSelectiaPlan]=useState(false);
    const [cautare, setCautare] = useState(initial.current?.cautare ?? '');
    const [filtruRanduri, setFiltruRanduri] = useState(initial.current?.filtruRanduri ?? '');
    const [panou, setPanou] = useState('');
    const unelte = useRef<HTMLDivElement>(null);
    const butoanePanouri = useRef<Record<string, HTMLButtonElement | null>>({});
    const [schita, setSchita] = useState<datebaza.SchitaLuna | null>(null);
    const [calendar, setCalendar] = useState<datebaza.CalendarAn | null>(null);
    const [schimburi, setSchimburi] = useState<datebaza.Schimb[]>([]);
    const [modificari, setModificari] = useState<Map<string, datebaza.CelulaSchita>>(new Map());
    const [istoricInapoi, setIstoricInapoi] = useState<Map<string, datebaza.CelulaSchita>[]>([]);
    const [istoricInainte, setIstoricInainte] = useState<Map<string, datebaza.CelulaSchita>[]>([]);
    const [selectie, setSelectie] = useState<Selectie>({inceput: {rand: 0, zi: 0}, sfarsit: {rand: 0, zi: 0}});
    const [valoareCompletare, setValoareCompletare] = useState('');
    const [programCompletare,setProgramCompletare] = useState('');
    const [calculat,setCalculat] = useState<calcul.RezultatPontaj | null>(null);
    const [calculare,setCalculare] = useState(false);
    const manualInput = useRef<HTMLInputElement>(null);
    const [incarcare, setIncarcare] = useState(true);
    const [ocupat, setOcupat] = useState(false);
    const [eroare, setEroare] = useState('');
    const [mesaj, setMesaj] = useState('');
    const [reincarcare, setReincarcare] = useState(0);
    const [zoomGrila, setZoomGrila] = useState<'compact' | 'normal' | 'extins'>(() => (localStorage.getItem('pontix_zoom_grila') as any) || 'normal');
    const schimbaZoom = (z: 'compact' | 'normal' | 'extins') => {
        setZoomGrila(z);
        localStorage.setItem('pontix_zoom_grila', z);
    };
    const grila = useRef<HTMLDivElement>(null);
    const ignoraFocus = useRef(false);
    const blocat = useRef(false);

    // Numai contextul de navigare rămâne în sesiune; schița se reîncarcă din baza de date.
    useEffect(() => () => { context.current = stareCurenta.current; }, [context]);
    useEffect(() => {
        if (!panou) return;
        const inchide = (e: PointerEvent) => { if (!unelte.current?.contains(e.target as Node)) setPanou(''); };
        document.addEventListener('pointerdown', inchide);
        return () => document.removeEventListener('pointerdown', inchide);
    }, [panou]);

    useEffect(() => { onModificari(modificari.size > 0); }, [modificari.size, onModificari]);
    useEffect(() => semnaleazaModificari('pontaj', modificari.size > 0), [modificari.size]);

    useEffect(() => {
        let activ = true;
        setIncarcare(true); setSchita(null); setEroare(''); setMesaj(''); setModificari(new Map()); setIstoricInapoi([]); setIstoricInainte([]);
        Promise.all([CitesteSchita(an, luna), CitesteCalendar(an), ListaSchimburi()]).then(([s, c, n]) => {
            if (!activ) return;
            const err = s.eroare || c.eroare || n.eroare;
            if (err) { setEroare(t(err.cod, err.param ?? {})); return; }
            if (!s.schita || !c.calendar) { setEroare(t('date.eroare')); return; }
            setSchita(s.schita); setCalculat(s.schita.calcul); setCalendar(c.calendar); setSchimburi(n.schimburi ?? []);
        }).catch(() => { if (activ) setEroare(t('date.eroare')); }).finally(() => { if (activ) setIncarcare(false); });
        return () => { activ = false; };
    }, [an, luna, reincarcare]);

    useEffect(()=>{
        if(!schita)return;
        if(!modificari.size){setCalculat(schita.calcul);setCalculare(false);return;}
        let activ=true;setCalculare(true);
        const timer=window.setTimeout(()=>PrevizualizeazaSchita(an,luna,[...modificari.values()]).then(r=>{
            if(!activ)return;
            if(r.eroare){setEroare(t(r.eroare.cod,r.eroare.param||{}));setCalculat(null);}
            else if(r.schita)setCalculat(r.schita.calcul);
            else {setCalculat(null);setEroare(t('date.eroare'));}
        }).catch(()=>{if(activ){setEroare(t('date.eroare'));setCalculat(null);}}).finally(()=>{if(activ)setCalculare(false);}),180);
        return()=>{activ=false;window.clearTimeout(timer);};
    },[an,luna,schita,modificari]);
    const oreCalculate=useMemo(()=>new Map((calculat?.celule||[]).map(c=>[cheie(c.angajat_id,c.data),c])),[calculat]);
    const totaluriCalculate=useMemo(()=>new Map((calculat?.totaluri||[]).map(c=>[c.angajat_id,c.valori])),[calculat]);

    const originale = useMemo(() => new Map((schita?.celule ?? []).map(c => [cheie(c.angajat_id, c.data), c])), [schita]);
    const absente = useMemo(() => new Map((schita?.absente ?? []).filter(z => z.activ).map(z => [cheie(z.angajat_id,z.data),z])), [schita]);
    const inlocuiri = useMemo(() => {
        const harta = new Map<string, datebaza.ZiAbsenta[]>();
        for (const z of schita?.absente ?? []) {
            if (z.activ && z.inlocuitor_id && (z.stare_aprobare || 'IN_ASTEPTARE') !== 'RESPINSA') {
                const k = cheie(z.inlocuitor_id, z.data);
                harta.set(k, [...(harta.get(k) ?? []), z]);
            }
        }
        return harta;
    }, [schita]);
    const numarZile = [31, an % 400 === 0 || (an % 4 === 0 && an % 100 !== 0) ? 29 : 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31][luna - 1];
    // zile e memorat: altfel lista de rânduri (cu filtrele ei) se recalcula la orice randare (PERF-1).
    const zile = useMemo(() => Array.from({length: numarZile}, (_, i) => `${String(an).padStart(4, '0')}-${String(luna).padStart(2, '0')}-${String(i + 1).padStart(2, '0')}`), [an, luna, numarZile]);
    const weekend = useMemo(() => zile.map(data => [0, 6].includes(new Date(`${data}T12:00:00Z`).getUTCDay())), [zile]);
    const valBrutaId = (id: string, data: string) => {
        const k = cheie(id, data);
        return modificari.get(k)?.valoare ?? originale.get(k)?.valoare ?? '';
    };
    const randuri = useMemo(() => {
        const cuSchite = new Set((schita?.celule ?? []).map(c => c.angajat_id));
        for (const z of schita?.absente ?? []) if (z.activ) cuSchite.add(z.angajat_id);
        for (const z of calculat?.celule ?? []) cuSchite.add(z.angajat_id);
        const termen = normalizat(cautare.trim());
        return date.angajati.filter(a => {
            if (!(a.activ || cuSchite.has(a.id)) || (sectie && a.sectie_id !== sectie) || (departament && a.departament_id !== departament)) return false;
            if(categorieProgram && (categorieProgram==='__fara__'?Boolean(a.categorie_program_id):a.categorie_program_id!==categorieProgram))return false;
            if (termen && !normalizat(`${a.nume} ${a.prenume} ${a.marca}`).includes(termen)) return false;
            if (filtruRanduri === 'modificate') return zile.some(data => modificari.has(cheie(a.id, data)));
            if (filtruRanduri === 'inlocuitori') return zile.some(data => (inlocuiri.get(cheie(a.id, data))?.length ?? 0) > 0);
            if (filtruRanduri === 'avertismente') {
                return zile.some((data, z) => {
                    const k = cheie(a.id, data), min = oreCalculate.get(k)?.minute ?? 0;
                    if (min > 1440) return true;
                    const areLucru = !absente.has(k) && (valBrutaId(a.id, data).trim() !== '' || min > 0);
                    return z > 0 && valBrutaId(a.id, zile[z - 1]).trim().toUpperCase() === 'G24' && areLucru;
                });
            }
            return true;
        });
    }, [date.angajati, schita, calculat, sectie, departament, categorieProgram, cautare, filtruRanduri, modificari, inlocuiri, totaluriCalculate, oreCalculate, absente, originale, zile]);
    const sarbatori = useMemo(() => new Map((calendar?.sarbatori ?? []).map(s => [s.data, s.denumire[i18n.language] || s.denumire.ro || Object.values(s.denumire).find(Boolean) || ''])), [calendar, i18n.language]);
    const numeSectie = (id: string) => { const s = date.sectii.find(s => s.id === id); return s ? denumireSectie(s, i18n.language) : t('ui.fara_sectie'); };
    const inceputRand = Math.min(selectie.inceput.rand, selectie.sfarsit.rand), sfarsitRand = Math.max(selectie.inceput.rand, selectie.sfarsit.rand);
    const inceputZi = Math.min(selectie.inceput.zi, selectie.sfarsit.zi), sfarsitZi = Math.max(selectie.inceput.zi, selectie.sfarsit.zi);
    const valoareBruta = (rand: number, zi: number) => {
        if (!randuri[rand] || !zile[zi]) return '';
        const k = cheie(randuri[rand].id, zile[zi]);
        return modificari.get(k)?.valoare ?? originale.get(k)?.valoare ?? '';
    };
    const valoare = (rand: number, zi: number) => {
        if (!randuri[rand] || !zile[zi]) return '';
        const k = cheie(randuri[rand].id, zile[zi]);
        const z=oreCalculate.get(k);
        return calculare?'…':z?(z.calculabil?oreText(z.minute):'—'):calculat?'':'—';
    };

    useEffect(() => {
        if (primulFiltru.current) { primulFiltru.current = false; return; }
        restaurare.current = false;
        setSelectie({inceput: {rand: 0, zi: 0}, sfarsit: {rand: 0, zi: 0}});
    }, [sectie, departament, categorieProgram, cautare, filtruRanduri, an, luna]);
    useLayoutEffect(() => {
        if (!schita || !restaurare.current) return;
        restaurare.current = false;
        const salvat = initial.current;
        if (salvat?.selectie) {
            const pozitie = (p: {angajatId: string; data: string}) => ({rand: Math.max(0, randuri.findIndex(a => a.id === p.angajatId)), zi: Math.max(0, zile.indexOf(p.data))});
            setSelectie({inceput: pozitie(salvat.selectie.inceput), sfarsit: pozitie(salvat.selectie.sfarsit)});
        }
        grila.current?.scrollTo(salvat?.derulareX ?? 0, salvat?.derulareY ?? 0);
    }, [schita, randuri, zile]);

    const actiuni = useRef<ActiuniCelula>(null!);
    actiuni.current = {
        focus: (r, z) => { if (ignoraFocus.current) ignoraFocus.current = false; else setSelectie({inceput: {rand: r, zi: z}, sfarsit: {rand: r, zi: z}}); },
        mouseDown: (e, r, z) => { if (e.shiftKey) { e.preventDefault(); focalizeaza({rand: r, zi: z}, true); } },
        mouseEnter: (e, r, z) => { if (e.buttons === 1 && !ocupat) setSelectie(s => ({...s, sfarsit: {rand: r, zi: z}})); },
        keyDown: (e, r, z) => taste(e, r, z),
        dublu: () => deschideOreManuale(),
    };

    function deschideOreManuale() {
        setPanou('completare');
        requestAnimationFrame(() => manualInput.current?.focus());
    }

    function aplicaCelule(celule: {rand: number; zi: number; text: string}[]) {
        if (lunaBlocata) { setEroare(t('pontaj.luna_inchisa')); return; }
        if (celule.some(c => !randuri[c.rand] || !zile[c.zi] || Array.from(c.text).length > 16 || /[\0\t\r\n]/.test(c.text))) { setEroare(t('pontaj.lipire_invalida')); return; }
        if (celule.some(c => absente.has(cheie(randuri[c.rand].id,zile[c.zi])))) { setEroare(t('absente.celula_protejata')); return; }
        for(const c of celule){
            if(c.text===(originale.get(cheie(randuri[c.rand].id,zile[c.zi]))?.valoare??''))continue;
            const program=schimburi.find(s=>s.cod.toUpperCase()===c.text.trim().toUpperCase());
            const categorie=date.categorii_program?.find(p=>p.id===randuri[c.rand].categorie_program_id);
            if(program&&categorie?.activ&&!categorie.schimburi_id.includes(program.id)){setEroare(t('program.nepermis',{nume:`${randuri[c.rand].nume} ${randuri[c.rand].prenume}`,cod:program.cod}));return;}
        }
        setIstoricInapoi(stiva => [...stiva.slice(-49), new Map(modificari)]);
        setIstoricInainte([]);
        setModificari(vechi => {
            const noi = new Map(vechi);
            for (const c of celule) {
                const id = randuri[c.rand].id, data = zile[c.zi], k = cheie(id, data), original = originale.get(k);
                if (c.text === (original?.valoare ?? '')) noi.delete(k);
                else noi.set(k, new datebaza.CelulaSchita({angajat_id: id, data, valoare: c.text, revizie: original?.revizie ?? ''}));
            }
            return noi;
        });
        setCalculare(true);
        setEroare(''); setMesaj('');
    }

    function anuleazaPas() {
        if (!istoricInapoi.length || ocupat) return;
        const anterioara = istoricInapoi[istoricInapoi.length - 1];
        setIstoricInapoi(stiva => stiva.slice(0, -1));
        setIstoricInainte(stiva => [...stiva.slice(-49), new Map(modificari)]);
        setModificari(new Map(anterioara));
        setCalculare(anterioara.size > 0);
        setEroare(''); setMesaj('');
    }

    function refaPas() {
        if (!istoricInainte.length || ocupat) return;
        const urmatoare = istoricInainte[istoricInainte.length - 1];
        setIstoricInainte(stiva => stiva.slice(0, -1));
        setIstoricInapoi(stiva => [...stiva.slice(-49), new Map(modificari)]);
        setModificari(new Map(urmatoare));
        setCalculare(urmatoare.size > 0);
        setEroare(''); setMesaj('');
    }

    function completeaza(text: string) {
        const celule = [];
        for (let r = inceputRand; r <= sfarsitRand; r++) for (let z = inceputZi; z <= sfarsitZi; z++) celule.push({rand: r, zi: z, text});
        aplicaCelule(celule);
    }

    function precompleteazaIntervalRanduri(deLaRand: number, panaLaRand: number) {
        if (!programCompletare) return;
        const ales=schimburi.find(s=>s.cod===programCompletare);
        if(ales&&programPlanificat(ales)){setEroare(t('planificare.foloseste_modulul'));return;}
        const celule: {rand: number; zi: number; text: string}[] = [];
        for (let r = deLaRand; r <= panaLaRand; r++) {
            if (!randuri[r]) continue;
            if(!randuri[r].activ)continue;
            const categorie=date.categorii_program?.find(p=>p.id===randuri[r].categorie_program_id);
            const program=schimburi.find(s=>s.cod===programCompletare);
            if(categorie&&(!categorie.activ||!program||!categorie.schimburi_id.includes(program.id)||program.durata_minute!==randuri[r].norma_minute))continue;
            for (let z = 0; z < zile.length; z++) {
                const data = zile[z], ziSapt = new Date(`${data}T12:00:00Z`).getUTCDay();
                if (ziSapt === 0 || ziSapt === 6 || sarbatori.has(data)) continue;
                const k = cheie(randuri[r].id, data);
                if (absente.has(k) || valoareBruta(r, z).trim() !== '') continue;
                celule.push({rand: r, zi: z, text: programCompletare});
            }
        }
        if (celule.length) aplicaCelule(celule);
    }

    function precompleteazaLV() { precompleteazaIntervalRanduri(0, randuri.length - 1); }
    function precompleteazaSelectieLV() { precompleteazaIntervalRanduri(inceputRand, sfarsitRand); }

    async function previzualizeazaCategorii(numaiSelectia:boolean) {
        if(blocat.current||!randuri.length)return;
        blocat.current=true;setOcupat(true);onOcupat(true);setEroare('');setMesaj('');
        try {
            const ids=(numaiSelectia?randuri.slice(inceputRand,sfarsitRand+1):randuri).map(a=>a.id);
            const r=await PrevizualizeazaPrecompletare(an,luna,ids,[...modificari.values()]);
            if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param??{}));
            else if(r.plan){setPlan(r.plan);setNumaiSelectiaPlan(numaiSelectia);setPanou('');}
            else setEroare(t('date.eroare'));
        } catch {setEroare(t('date.eroare'));}
        finally{blocat.current=false;setOcupat(false);onOcupat(false);}
    }

    function aplicaPlan() {
        if(!plan)return;
        const celule=plan.celule.map(c=>({rand:randuri.findIndex(a=>a.id===c.angajat_id),zi:zile.indexOf(c.data),text:c.valoare}));
        if(celule.some(c=>c.rand<0||c.zi<0||valoareBruta(c.rand,c.zi).trim()!=='')){setPlan(null);setEroare(t('program.plan_schimbat'));return;}
        aplicaCelule(celule);setPlan(null);
    }

    function textSelectat() {
        if(calculare || !calculat)return '';
        return Array.from({length: sfarsitRand - inceputRand + 1}, (_, r) => Array.from({length: sfarsitZi - inceputZi + 1}, (_, z) => {
            const celula=oreCalculate.get(cheie(randuri[inceputRand+r].id,zile[inceputZi+z]));
            return celula?.calculabil?String(Number((celula.minute/60).toFixed(10))):valoare(inceputRand+r,inceputZi+z);
        }).join('\t')).join('\n');
    }

    function lipeste(e: ClipboardEvent) {
        e.preventDefault();
        if (ocupat || !randuri.length) return;
        const text = e.clipboardData.getData('text/plain').replace(/\r\n?/g, '\n').replace(/\n$/, '');
        const matrice = text.split('\n').map(linie => linie.split('\t'));
        if(matrice.flat().some(text=>manual(text)===null)){setEroare(t('pontaj.ore_invalide'));return;}
        if (matrice.length === 1 && matrice[0].length === 1) { completeaza(manual(matrice[0][0])!); return; }
        const celule = matrice.flatMap((linie, r) => linie.map((text, z) => ({rand: inceputRand + r, zi: inceputZi + z, text:manual(text)!})));
        aplicaCelule(celule);
    }

    function focalizeaza(p: Pozitie, extinde: boolean) {
        const urmator = {rand: Math.max(0, Math.min(randuri.length - 1, p.rand)), zi: Math.max(0, Math.min(numarZile - 1, p.zi))};
        setSelectie(s => extinde ? {...s, sfarsit: urmator} : {inceput: urmator, sfarsit: urmator});
        const input = grila.current?.querySelector<HTMLInputElement>(`input[data-rand="${urmator.rand}"][data-zi="${urmator.zi}"]`);
        ignoraFocus.current = Boolean(input && document.activeElement !== input);
        input?.focus();
        input?.scrollIntoView({block: 'nearest', inline: 'nearest'});
    }

    function schimbaLunaRelativ(delta: number) {
        if (ocupat || !poateAbandona()) return;
        let nouAn = an, nouaLuna = luna + delta;
        if (nouaLuna < 1) { nouAn--; nouaLuna = 12; }
        else if (nouaLuna > 12) { nouAn++; nouaLuna = 1; }
        if (nouAn < 1 || nouAn > 9999) return;
        setAn(nouAn); setAnText(String(nouAn)); setLuna(nouaLuna);
    }

    function taste(e: KeyboardEvent<HTMLInputElement>, r: number, z: number) {
        if ((e.ctrlKey || e.metaKey) && !e.altKey && e.key.toLowerCase() === 'z') { e.preventDefault(); if (e.shiftKey) refaPas(); else anuleazaPas(); return; }
        if ((e.ctrlKey || e.metaKey) && !e.altKey && e.key.toLowerCase() === 'y') { e.preventDefault(); refaPas(); return; }
        const directii: Record<string, [number, number]> = {ArrowLeft: [0, -1], ArrowRight: [0, 1], ArrowUp: [-1, 0], ArrowDown: [1, 0], Enter: [e.shiftKey ? -1 : 1, 0]};
        if (directii[e.key]) { e.preventDefault(); const [dr, dz] = directii[e.key]; focalizeaza({rand: r + dr, zi: z + dz}, e.shiftKey && e.key !== 'Enter'); }
        else if (e.key === 'Home') { e.preventDefault(); focalizeaza({rand: r, zi: 0}, e.shiftKey); }
        else if (e.key === 'End') { e.preventDefault(); focalizeaza({rand: r, zi: numarZile - 1}, e.shiftKey); }
        else if (e.key === 'PageUp') { e.preventDefault(); focalizeaza({rand: Math.max(0, r - 5), zi: z}, e.shiftKey); }
        else if (e.key === 'PageDown') { e.preventDefault(); focalizeaza({rand: Math.min(randuri.length - 1, r + 5), zi: z}, e.shiftKey); }
        else if (e.key === 'Tab') {
            let zi = z + (e.shiftKey ? -1 : 1), rand = r;
            if (zi < 0) { rand--; zi = numarZile - 1; } else if (zi >= numarZile) { rand++; zi = 0; }
            if (rand >= 0 && rand < randuri.length) { e.preventDefault(); focalizeaza({rand, zi}, false); }
        } else if (e.key === 'Delete' || e.key === 'Backspace') { e.preventDefault(); completeaza(''); }
        else if (e.key === 'F2') { e.preventDefault(); deschideOreManuale(); }
        else if (e.key === 'Escape') { e.preventDefault(); aplicaCelule([{rand: r, zi: z, text: originale.get(cheie(randuri[r].id, zile[z]))?.valoare ?? ''}]); }
        else if (e.key === ' ' && !e.ctrlKey && !e.altKey && !e.metaKey) { e.preventDefault(); deschideOreManuale(); }
        else if (/^[0-9]$/.test(e.key) && !e.ctrlKey && !e.altKey && !e.metaKey && !ocupat && !lunaBlocata) {
            e.preventDefault();
            const ore = Number(e.key);
            const val = ore === 0 ? '' : `@M${ore * 60}`;
            completeaza(val);
            if (z + 1 < numarZile) {
                focalizeaza({rand: r, zi: z + 1}, false);
            }
        }
    }

    function poateAbandona() { return !modificari.size || window.confirm(t('ui.renunta')); }
    const lunaBlocata = schita?.stare === 'INCHISA' || schita?.stare === 'BLOCATA';
    const [stareLuna, setStareLuna] = useState<datebaza.StareLuna | null>(null);
    useEffect(() => {
        let activ = true;
        if (!lunaBlocata) { setStareLuna(null); return; }
        CitesteStareLuna(an, luna).then(r => { if (activ && !r.eroare) setStareLuna(r.stare); }).catch(() => {});
        return () => { activ = false; };
    }, [an, luna, lunaBlocata, schita]);

    // Închiderea și redeschiderea cer un motiv scris de utilizator, care intră în jurnal (F12-T05).
    const [dialogLuna, setDialogLuna] = useState(false);
    const [motivLuna, setMotivLuna] = useState('');
    const [eroareLuna, setEroareLuna] = useState('');

    function comutaStareLuna() {
        if (!schita || ocupat) return;
        setMotivLuna(lunaBlocata ? '' : t('pontaj.motiv_inchidere_implicit'));
        setEroareLuna('');
        setDialogLuna(true);
    }

    async function confirmaStareLuna() {
        if (!schita || ocupat) return;
        if (!motivLuna.trim()) { setEroareLuna(t('pontaj.motiv_obligatoriu')); return; }
        const esteInchisa = lunaBlocata;
        setOcupat(true); onOcupat(true); setEroare(''); setMesaj(''); setEroareLuna('');
        try {
            const resp = await SeteazaStareLuna(an, luna, esteInchisa ? 'DESCHISA' : 'INCHISA', motivLuna.trim());
            if (resp.eroare) {
                setEroareLuna(t(resp.eroare.cod, resp.eroare.param ?? {}));
            } else {
                setSchita(prev => prev ? new datebaza.SchitaLuna({...prev, stare: resp.stare.stare}) : null);
                setStareLuna(resp.stare);
                setDialogLuna(false);
                setMesaj(t(esteInchisa ? 'pontaj.luna_redeschisa' : 'pontaj.luna_inchisa_succes'));
            }
        } catch {
            setEroareLuna(t('date.eroare'));
        } finally {
            setOcupat(false); onOcupat(false);
        }
    }

    async function salveaza() {
        if (lunaBlocata) { setEroare(t('pontaj.luna_inchisa')); return; }
        if (blocat.current) return;
        blocat.current = true; setOcupat(true); onOcupat(true); setEroare(''); setMesaj('');
        try {
            const r = await SalveazaSchita(an, luna, [...modificari.values()]);
            if (r.eroare) { setEroare(t(r.eroare.cod, r.eroare.param ?? {})); notifica(t(r.eroare.cod, r.eroare.param ?? {}), 'eroare'); }
            else if (r.schita) { setSchita(r.schita); setCalculat(r.schita.calcul); setModificari(new Map()); setIstoricInapoi([]); setIstoricInainte([]); setMesaj(t('pontaj.salvat')); notifica(t('pontaj.salvat'), 'succes'); }
            else setEroare(t('date.eroare'));
        } catch { setEroare(t('date.eroare')); }
        finally { blocat.current = false; setOcupat(false); onOcupat(false); }
    }

    async function exportaExcel() {
        setOcupat(true); onOcupat(true); setEroare(''); setMesaj('');
        try {
            const resp = await ExportaPontajExcel(an, luna);
            if (resp.eroare) {
                setEroare(t(resp.eroare.cod, resp.eroare.param ?? {}));
            } else if (resp.cale) {
                setMesaj(t('rapoarte.fisier_salvat', {cale: resp.cale}));
            } else if (!resp.anulat && resp.continut) {
                const byteCharacters = atob(resp.continut);
                const byteNumbers = new Array(byteCharacters.length);
                for (let i = 0; i < byteCharacters.length; i++) {
                    byteNumbers[i] = byteCharacters.charCodeAt(i);
                }
                const byteArray = new Uint8Array(byteNumbers);
                const blob = new Blob([byteArray], {type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'});
                const url = URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = resp.nume_fisier || `Pontix_Pontaj_${an}_${String(luna).padStart(2, '0')}.xlsx`;
                document.body.appendChild(a);
                a.click();
                document.body.removeChild(a);
                URL.revokeObjectURL(url);
                setMesaj(t('pontaj.export_excel_succes'));
            }
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
            onOcupat(false);
        }
    }

    const randActiv = randuri[selectie.sfarsit.rand], dataActiva = zile[selectie.sfarsit.zi];
    const cheieActiva = randActiv && dataActiva ? cheie(randActiv.id, dataActiva) : '';
    const brutaActiva = randActiv && dataActiva ? valoareBruta(selectie.sfarsit.rand, selectie.sfarsit.zi).trim() : '';
    const absentaActiva = cheieActiva ? absente.get(cheieActiva) : undefined;
    const inlocuiriActive = cheieActiva ? inlocuiri.get(cheieActiva) : undefined;
    const calcActiva = cheieActiva ? oreCalculate.get(cheieActiva) : undefined;
    // Ultima modificare a celulei active, din jurnal (F08-T10); se recitește după salvare.
    const [ultimaModificare, setUltimaModificare] = useState<{cheie: string; intrare?: datebaza.IntrareJurnal} | null>(null);
    useEffect(() => {
        if (!randActiv || !dataActiva) return;
        let activ = true;
        const k = cheieActiva;
        const temporizator = window.setTimeout(() => {
            UltimaModificareCelula(randActiv.id, dataActiva).then(r => { if (activ && !r.eroare) setUltimaModificare({cheie: k, intrare: r.intrare}); }).catch(() => {});
        }, 200);
        return () => { activ = false; window.clearTimeout(temporizator); };
    }, [cheieActiva, schita]);
    const schimbActiv = brutaActiva && !brutaActiva.startsWith('@M') ? schimburi.find(s => s.cod === brutaActiva) : undefined;
    const textCodActiv = absentaActiva
        ? `${absentaActiva.cod} (${t(`absente.aprobare_${absentaActiva.stare_aprobare || 'IN_ASTEPTARE'}`)}${absentaActiva.nr_inregistrare ? ` · ${absentaActiva.nr_inregistrare}` : ''})`
        : brutaActiva.startsWith('@M')
            ? t('pontaj.inspectie_ore_manuale', {ore: oreText(Number(brutaActiva.slice(2)))})
            : schimbActiv
                ? `${schimbActiv.cod} (${schimbActiv.denumire[i18n.language] || schimbActiv.denumire.ro})`
                : brutaActiva
                    ? t('pontaj.inspectie_cod', {valoare: brutaActiva})
                    : t('pontaj.inspectie_gol');
    const dupaG24Activa = Boolean(randActiv && selectie.sfarsit.zi > 0 && valoareBruta(selectie.sfarsit.rand, selectie.sfarsit.zi - 1).trim().toUpperCase() === 'G24' && !absentaActiva && (brutaActiva !== '' || (calcActiva?.minute ?? 0) > 0));
    const peste24hActiva = (calcActiva?.minute ?? 0) > 1440;
    const avertismentActiv = [dupaG24Activa ? t('pontaj.avertisment_dupa_g24') : '', peste24hActiva ? t('pontaj.avertisment_peste_24h', {ore: oreText(calcActiva!.minute)}) : ''].filter(Boolean).join(' · ');

    stareCurenta.current = {
        an, luna, sectie, departament, categorieProgram, cautare, filtruRanduri,
        selectie: schita && !restaurare.current && randuri[selectie.inceput.rand] && randActiv
            ? {inceput: {angajatId: randuri[selectie.inceput.rand].id, data: zile[selectie.inceput.zi]}, sfarsit: {angajatId: randActiv.id, data: dataActiva}}
            : stareCurenta.current?.selectie ?? null,
        derulareX: grila.current?.scrollLeft ?? stareCurenta.current?.derulareX ?? 0,
        derulareY: grila.current?.scrollTop ?? stareCurenta.current?.derulareY ?? 0,
    };
    const numarFiltre = [sectie, departament, categorieProgram, cautare.trim(), filtruRanduri].filter(Boolean).length;

    return <section className="pontaj-lunar">
        <div className="unelte-calendar unelte-perioada"><form onSubmit={e => { e.preventDefault(); const nou = Number(anText); if (!Number.isInteger(nou) || nou < 1 || nou > 9999) { setEroare(t('sarbatori.an_invalid')); return; } if (nou !== an) { if (poateAbandona()) setAn(nou); else setAnText(String(an)); } }}>
            <label>{t('sarbatori.an')}<input aria-label={t('sarbatori.an')} type="number" min="1" max="9999" required disabled={ocupat} value={anText} onChange={e => setAnText(e.target.value)}/></label>
            <button className="buton secundar" disabled={ocupat}>{t('sarbatori.arata')}</button>
            <button type="button" className="buton secundar buton-luna-nav" aria-label={t('pontaj.luna_precedenta')} disabled={ocupat} onClick={() => schimbaLunaRelativ(-1)}>◀</button>
            <label>{t('pontaj.luna')}<select aria-label={t('pontaj.luna')} value={luna} disabled={ocupat} onChange={e => { if (poateAbandona()) setLuna(Number(e.target.value)); }}>
                {Array.from({length: 12}, (_, i) => <option key={i} value={i + 1}>{formatData(i18n.language, {month: 'long', timeZone: 'UTC'}).format(new Date(`2026-${String(i + 1).padStart(2, '0')}-01T12:00:00Z`))}</option>)}
            </select></label>
            <button type="button" className="buton secundar buton-luna-nav" aria-label={t('pontaj.luna_urmatoare')} disabled={ocupat} onClick={() => schimbaLunaRelativ(1)}>▶</button></form>
            <button className="buton primar" disabled={ocupat || incarcare || !modificari.size || lunaBlocata} onClick={salveaza}>{t(ocupat ? 'ui.salvare' : 'pontaj.salveaza')}</button>
            <button type="button" className="buton secundar" title={lunaBlocata ? t('pontaj.redeschide_luna') : t('pontaj.inchide_luna')} disabled={ocupat || incarcare} onClick={comutaStareLuna}>{lunaBlocata ? `🔓 ${t('pontaj.redeschide_luna')}` : `🔒 ${t('pontaj.inchide_luna')}`}</button></div>
        {lunaBlocata && (
            <div className="banner-luna-blocata" role="alert">
                <span className="simbol-lacat">🔒</span>
                <div>
                    <strong>{t('pontaj.stare_luna_inchisa_titlu')}</strong>
                    <p>{t('pontaj.stare_luna_inchisa_descriere')}</p>
                    <AmprentaLuna amprenta={stareLuna?.amprenta ?? ''} corespunde={stareLuna?.amprenta_corespunde ?? false}/>
                </div>
                <button type="button" className="buton secundar mic" disabled={ocupat} onClick={comutaStareLuna}>
                    🔓 {t('pontaj.redeschide_luna')}
                </button>
            </div>
        )}
        {dialogLuna && <Dialog titlu={t(lunaBlocata ? 'pontaj.redeschide_luna' : 'pontaj.inchide_luna')} ocupat={ocupat} inchide={() => setDialogLuna(false)}>
            <div className="corp-dialog">
                <p>{t(lunaBlocata ? 'pontaj.confirmare_redeschidere' : 'pontaj.confirmare_inchidere')}</p>
                {lunaBlocata && stareLuna?.amprenta && <p className="text-discret">{t('pontaj.redeschidere_amprenta', {amprenta: stareLuna.amprenta})}</p>}
                <label>{t('pontaj.motiv_luna')}
                    <textarea name="motiv_luna" rows={3} value={motivLuna} disabled={ocupat} onChange={e => setMotivLuna(e.target.value)} required/>
                </label>
                {eroareLuna && <div className="mesaj eroare" role="alert">{eroareLuna}</div>}
            </div>
            <footer className="subsol-dialog">
                <button className="buton secundar" disabled={ocupat} onClick={() => setDialogLuna(false)}>{t('ui.anuleaza')}</button>
                <button className="buton primar" disabled={ocupat || !motivLuna.trim()} onClick={() => void confirmaStareLuna()}>{t(lunaBlocata ? 'pontaj.redeschide_luna' : 'pontaj.inchide_luna')}</button>
            </footer>
        </Dialog>}
        {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}{mesaj && <div className="mesaj succes" role="status">{mesaj}</div>}
        {incarcare ? <p role="status">{t('ui.incarcare')}</p> : !schita ? <button className="buton secundar" onClick={() => setReincarcare(x => x + 1)}>{t('ui.reincearca')}</button> : <>
            <div className="unelte-grila" ref={unelte} onKeyDown={e => {
                if (e.key === 'Escape' && panou) { e.preventDefault(); e.stopPropagation(); butoanePanouri.current[panou]?.focus(); setPanou(''); }
            }}>
                <input name="cautare_pontaj" className="cautare-grila" aria-label={t('ui.cautare')} placeholder={t('ui.cautare_indiciu')} value={cautare} disabled={ocupat} onChange={e => setCautare(e.target.value)}/>
                {['filtre', 'completare', 'ajutor'].map(p => <button key={p} ref={element => { butoanePanouri.current[p] = element; }} type="button" className={`buton secundar ${panou === p ? 'panou-activ' : ''}`} disabled={ocupat} aria-expanded={panou === p} aria-controls={`pontaj-panou-${p}`} onClick={() => setPanou(panou === p ? '' : p)}>{t(`pontaj.panou_${p}`)}{p === 'filtre' && numarFiltre > 0 ? ` (${numarFiltre})` : ''}</button>)}
                <button type="button" className="buton secundar" title={t('program.completare_program_indiciu')} disabled={ocupat||!randuri.length||!programCompletare} onClick={precompleteazaLV}>{t('pontaj.precompleteaza_lv')}</button>
                <button type="button" className="buton secundar" disabled={ocupat||!randuri.length} onClick={()=>void previzualizeazaCategorii(false)}>{t('program.precompleteaza')}</button>
                <button type="button" className="buton secundar" title={t('pontaj.export_excel_titlu')} disabled={ocupat||!randuri.length} onClick={exportaExcel}>📊 {t('pontaj.export_excel')}</button>
                <div className="grup-zoom-grila" role="group" aria-label={t('ui.densitate_grila')}>
                    <span className="eticheta-zoom">🔎 {t('ui.zoom')}</span>
                    <button type="button" className={`buton-zoom ${zoomGrila === 'compact' ? 'activ' : ''}`} title={t('ui.zoom_compact')} onClick={() => schimbaZoom('compact')}>S</button>
                    <button type="button" className={`buton-zoom ${zoomGrila === 'normal' ? 'activ' : ''}`} title={t('ui.zoom_normal')} onClick={() => schimbaZoom('normal')}>M</button>
                    <button type="button" className={`buton-zoom ${zoomGrila === 'extins' ? 'activ' : ''}`} title={t('ui.zoom_extins')} onClick={() => schimbaZoom('extins')}>L</button>
                </div>
                <span className="stare-modificari" aria-live="polite">{t('pontaj.modificate', {numar: modificari.size})}</span>
                <div id="pontaj-panou-filtre" className="panou-unelte" hidden={panou !== 'filtre'}>
                <select aria-label={t('ui.sectie')} value={sectie} disabled={ocupat} onChange={e => { setSectie(e.target.value); setDepartament(''); }}><option value="">{t('ui.toate_sectiile')}</option>
                {date.sectii.map(s => <option key={s.id} value={s.id}>{denumireSectie(s, i18n.language)}</option>)}</select>
                <select aria-label={t('ui.departament')} value={departament} disabled={ocupat} onChange={e => setDepartament(e.target.value)}><option value="">{t('ui.toate_departamentele')}</option>
                    {(date.departamente ?? []).filter(d => !sectie || d.sectie_id === sectie).map(d => <option key={d.id} value={d.id}>{!sectie ? `${numeSectie(d.sectie_id)} / ` : ''}{denumireDepartament(d, i18n.language)}</option>)}</select>
                <select name="filtru_randuri" aria-label={t('pontaj.filtru_randuri')} value={filtruRanduri} disabled={ocupat} onChange={e => setFiltruRanduri(e.target.value)}>
                    <option value="">{t('pontaj.toti_angajatii')}</option>
                    <option value="avertismente">{t('pontaj.doar_avertismente')}</option>
                    <option value="modificate">{t('pontaj.doar_modificate')}</option>
                    <option value="inlocuitori">{t('pontaj.doar_inlocuitori')}</option>
                </select>
                <select name="filtru_categorie_program" aria-label={t('program.categorie')} value={categorieProgram} disabled={ocupat} onChange={e=>setCategorieProgram(e.target.value)}><option value="">{t('program.toate_categoriile')}</option><option value="__fara__">{t('program.fara_categorie')}</option>{(date.categorii_program??[]).map(c=><option key={c.id} value={c.id}>{c.denumire[i18n.language]||c.denumire.ro}</option>)}</select>
                <button type="button" className="legatura" disabled={ocupat || !numarFiltre} onClick={() => { setSectie(''); setDepartament(''); setCategorieProgram(''); setCautare(''); setFiltruRanduri(''); }}>{t('pontaj.reseteaza_filtre')}</button>
                </div>
                <div id="pontaj-panou-completare" className="panou-unelte" hidden={panou !== 'completare'}>
                <select aria-label={t('pontaj.program_completare')} value={programCompletare} disabled={ocupat} onChange={e=>setProgramCompletare(e.target.value)}><option value="">{t('pontaj.program_completare')}</option>{schimburi.filter(s=>s.activ).map(s=><option key={s.id} value={s.cod}>{s.denumire[i18n.language]||s.denumire.ro} {t(`schimburi.regim_${s.regim||'GENERAL'}`)} · {s.ora_start<0?t('pontaj.fara_ora'): `${String(Math.floor(s.ora_start/60)).padStart(2,'0')}:${String(s.ora_start%60).padStart(2,'0')}`} · {t('pontaj.durata_ore',{ore:oreText(s.durata_minute)})}</option>)}</select>
                <button className="buton secundar" disabled={ocupat || !randuri.length || !programCompletare} onClick={()=>completeaza(programCompletare)}>{t('pontaj.aplica_program')}</button>
                <button type="button" className="buton secundar" disabled={ocupat || !randuri.length || !programCompletare} onClick={precompleteazaSelectieLV}>{t('pontaj.precompleteaza_rand_lv')}</button>
                <input name="ore_manuale" ref={manualInput} aria-label={t('pontaj.valoare_completare')} placeholder={t('pontaj.valoare_completare')} inputMode="decimal" maxLength={16} value={valoareCompletare} disabled={ocupat} onChange={e => setValoareCompletare(e.target.value)}/>
                <button className="buton secundar" disabled={ocupat || !randuri.length} onClick={() => {const v=manual(valoareCompletare);if(v===null)setEroare(t('pontaj.ore_invalide'));else completeaza(v);}}>{t('pontaj.completeaza_selectia')}</button>
                <div className="paleta-schimburi">
                    <span className="text-discret">{t('pontaj.paleta_rapida')}</span>
                    {schimburi.filter(s => s.activ).map(s => <button key={s.id} type="button" className="cip-schimb" disabled={ocupat || !randuri.length} onClick={() => completeaza(s.cod)}>{s.cod}</button>)}
                    <button type="button" className="cip-schimb cip-sterge" disabled={ocupat || !randuri.length} onClick={() => completeaza('')}>{t('pontaj.sterge_selectia')}</button>
                </div>
                <button type="button" className="buton secundar" disabled={ocupat || !istoricInapoi.length} onClick={anuleazaPas}>{t('pontaj.anuleaza_pas')}</button>
                <button type="button" className="buton secundar" disabled={ocupat || !istoricInainte.length} onClick={refaPas}>{t('pontaj.refa_pas')}</button>
                <button className="buton secundar" disabled={ocupat || calculare || !calculat || !randuri.length} onClick={async () => { try { await navigator.clipboard.writeText(textSelectat()); notifica(t('pontaj.selectie_copiata'), 'info'); } catch { setEroare(t('pontaj.clipboard')); } }}>{t('pontaj.copiaza')}</button>
                <button className="legatura" disabled={ocupat} onClick={() => { if (poateAbandona()) setReincarcare(x => x + 1); }}>{t('pontaj.reincarca')}</button>
                </div>
                <div id="pontaj-panou-ajutor" className="panou-unelte ajutor-pontaj" hidden={panou !== 'ajutor'}>
                    <p>{t('pontaj.schita_indiciu')}</p><p>{t('pontaj.tastatura')}</p><p>{t('pontaj.ore_norm_explicatie')}</p><p>{t('pontaj.necalculat')}</p>
                </div>
            </div>
            {randActiv && dataActiva && <div className="inspectie-celula" role="region" aria-label={t('pontaj.inspectie_titlu')}>
                <div className="detalii-celula-activa">
                    <strong>{randActiv.nume} {randActiv.prenume}</strong>
                    <span>{dataActiva}{sarbatori.get(dataActiva) ? ` · ${sarbatori.get(dataActiva)}` : ''}</span>
                    <span className="eticheta-contract">{textCodActiv}</span>
                    {calcActiva?.calculabil && <span>{t('pontaj.inspectie_calculat', {ore: oreText(calcActiva.minute)})}</span>}
                    {inlocuiriActive?.length ? <span className="insigna-inlocuitor">{t('pontaj.inlocuieste_pe', {persoane: inlocuiriActive.map(z => { const coleg = date.angajati.find(x => x.id === z.angajat_id); return `${coleg ? `${coleg.nume} ${coleg.prenume}` : z.angajat_id} (${z.cod})`; }).join(', ')})}</span> : null}
                    {avertismentActiv && <span className="avertisment-text">{avertismentActiv}</span>}
                    {ultimaModificare?.cheie === cheieActiva && ultimaModificare.intrare && <span className="ultima-modificare" title={t('pontaj.ultima_modificare_indiciu')}>
                        {t('pontaj.ultima_modificare', {cine: autorJurnal(t, ultimaModificare.intrare.cine), cand: new Date(ultimaModificare.intrare.cand).toLocaleString(i18n.language)})}
                    </span>}
                    {onDeschideAbsente && <button type="button" className="legatura" disabled={ocupat} onClick={() => { if (poateAbandona()) onDeschideAbsente(an, randActiv.id); }}>{t('pontaj.deschide_absente')}</button>}
                </div>
            </div>}
            <div className="legenda-grila" aria-label={t('pontaj.legenda')}>
                {['weekend', 'sarbatoare', 'selectie', 'modificat', 'avertisment', 'inlocuire'].map(stare => <span key={stare}><i aria-hidden="true" className={`mostra-${stare}`}/>{t(`pontaj.legenda_${stare}`)}</span>)}
            </div>
            <div className={`panou grila-container zoom-${zoomGrila}`} ref={grila} onScroll={e => {
                if (stareCurenta.current) { stareCurenta.current.derulareX = e.currentTarget.scrollLeft; stareCurenta.current.derulareY = e.currentTarget.scrollTop; }
            }} onPaste={lipeste} onCopy={e => { e.preventDefault(); e.clipboardData.setData('text/plain', textSelectat()); }}>
                <table className="grila-pontaj" aria-label={t('pontaj.grila', {an, luna})}><thead><tr><th className="fix-nume">{t('ui.angajat')}</th><th className="fix-marca">{t('ui.marca')}</th>
                    {zile.map((data, z) => { const d = new Date(`${data}T12:00:00Z`); return <th key={data} className={sarbatori.has(data) ? 'zi-sarbatoare' : weekend[z] ? 'zi-weekend' : ''} title={sarbatori.get(data) || ''} scope="col"><strong>{z + 1}</strong><small>{formatData(i18n.language, {weekday: 'short', timeZone: 'UTC'}).format(d)}</small></th>; })}
                    {totaluri.map(c => <th key={c} className={c === 'T' ? 'fix-total-t' : ''}>{t(`pontaj.total_${c}`)}</th>)}</tr></thead><tbody>
                    {randuri.map((a, r) => <RandPontaj key={a.id} a={a} r={r} zile={zile} weekend={weekend} sarbatori={sarbatori}
                        selDe={r >= inceputRand && r <= sfarsitRand ? inceputZi : -1} selPana={r >= inceputRand && r <= sfarsitRand ? sfarsitZi : -1}
                        modificari={modificari} originale={originale} absente={absente} inlocuiri={inlocuiri} oreCalculate={oreCalculate}
                        valori={totaluriCalculate.get(a.id)} reguli={schita?.configuratie.reguli} angajati={date.angajati}
                        ocupat={ocupat} calculare={calculare} areCalcul={Boolean(calculat)} actiuni={actiuni}/>)}
                </tbody></table>
                {randuri.length === 0 && (
                    <div className="card-onboarding-gol" role="region" aria-label={t('ui.ghid_pornire')}>
                        <div className="simbol-gol">📋</div>
                        <h3>{t('pontaj.gol_titlu')}</h3>
                        <p>{t('pontaj.gol_descriere')}</p>
                        <div className="actiuni-onboarding-rapide">
                            <button type="button" className="buton secundar" onClick={() => setReincarcare(x => x + 1)}>
                                🔄 {t('ui.reincarca_datele')}
                            </button>
                            <button type="button" className="buton secundar" onClick={() => { setSectie(''); setDepartament(''); setCautare(''); setFiltruRanduri(''); setCategorieProgram(''); }}>
                                {t('ui.reseteaza_filtrele')}
                            </button>
                        </div>
                    </div>
                )}
            </div>
            <p className="indiciu-grila">{t('pontaj.dimensiuni', {randuri: randuri.length, zile: numarZile})} · {t('pontaj.inactivi_indiciu')}</p>
        </>}
        {plan&&<Dialog titlu={t('program.previzualizare')} ocupat={ocupat} inchide={()=>setPlan(null)}><div className="corp-dialog plan-precompletare">
            <p>{t('program.plan_indiciu')}</p>
            <label className="bifa"><input type="checkbox" checked={numaiSelectiaPlan} disabled={ocupat} onChange={e=>void previzualizeazaCategorii(e.target.checked)}/>{t('program.numai_selectia')}</label>
            <p><strong>{t('program.propuneri',{numar:plan.celule.length})}</strong></p>
            <ul>{Object.entries(plan.omise).sort(([a],[b])=>a.localeCompare(b)).map(([motiv,numar])=><li key={motiv}>{t(`program.omise_${motiv}`,{numar})}</li>)}</ul>
            {plan.celule.length>0&&<div className="tabel-derulant"><table><thead><tr><th>{t('ui.angajat')}</th><th>{t('program.data')}</th><th>{t('program.program')}</th></tr></thead><tbody>{plan.celule.slice(0,100).map(c=><tr key={`${c.angajat_id}/${c.data}`}><td>{date.angajati.find(a=>a.id===c.angajat_id)?.nume} {date.angajati.find(a=>a.id===c.angajat_id)?.prenume}</td><td>{c.data}</td><td>{schimburi.find(s=>s.cod===c.valoare)?.denumire[i18n.language]||c.valoare}</td></tr>)}</tbody></table>{plan.celule.length>100&&<p>{t('program.primele_propuneri')}</p>}</div>}
        </div><footer className="subsol-dialog"><button className="buton secundar" disabled={ocupat} onClick={()=>setPlan(null)}>{t('ui.anuleaza')}</button><button className="buton primar" disabled={ocupat||!plan.celule.length} onClick={aplicaPlan}>{t('program.aplica_plan')}</button></footer></Dialog>}
    </section>;
}
