import {useEffect, useMemo, useRef, useState, type ClipboardEvent, type KeyboardEvent} from 'react';
import {useTranslation} from 'react-i18next';
import {CitesteCalendar, CitesteSchita, ListaSchimburi, SalveazaSchita, PrevizualizeazaSchita} from '../wailsjs/go/main/App';
import {calcul, datebaza} from '../wailsjs/go/models';
import {denumireSectie, denumireDepartament} from './Formulare';

type Pozitie = {rand: number; zi: number};
type Selectie = {inceput: Pozitie; sfarsit: Pozitie};
const totaluri = ['ore_lucrate', 'T', 'ZNL', 'CO', 'CEV', 'CFS', 'ore_norma', 'form_prof', 'CM', 'L_FIN', 'deleg', 'ore_garda'];
const cheie = (id: string, data: string) => JSON.stringify([id, data]);
const oreText = (minute: number) => String(Number((minute/60).toFixed(4)));
const valoareText = (v?: calcul.ValoarePontaj) => v?.definita ? String(Number((v.sutimi/(v.unitate==='ORE'?6000:100)).toFixed(4))) : '—';
function manual(text: string): string | null {
    if(!text.trim())return '';
    if(!/^\d+(?:[.,]\d+)?$/.test(text.trim()))return null;
    const ore=Number(text.trim().replace(',','.')),minute=Math.round(ore*60);
    return ore<=48&&Math.abs(ore*60-minute)<0.00001?`@M${minute}`:null;
}

export function Pontaj({date, onModificari, onOcupat}: {date: datebaza.Date; onModificari: (modificat: boolean) => void; onOcupat: (ocupat: boolean) => void}) {
    const {t, i18n} = useTranslation();
    const acum = useRef(new Date());
    const [an, setAn] = useState(acum.current.getFullYear());
    const [luna, setLuna] = useState(acum.current.getMonth() + 1);
    const [anText, setAnText] = useState(String(an));
    const [sectie, setSectie] = useState('');
    const [departament, setDepartament] = useState('');
    const [schita, setSchita] = useState<datebaza.SchitaLuna | null>(null);
    const [calendar, setCalendar] = useState<datebaza.CalendarAn | null>(null);
    const [schimburi, setSchimburi] = useState<datebaza.Schimb[]>([]);
    const [modificari, setModificari] = useState<Map<string, datebaza.CelulaSchita>>(new Map());
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
    const grila = useRef<HTMLDivElement>(null);
    const ignoraFocus = useRef(false);
    const blocat = useRef(false);

    useEffect(() => { onModificari(modificari.size > 0); }, [modificari.size, onModificari]);
    useEffect(() => {
        if (!modificari.size) return;
        const avertizeaza = (e: BeforeUnloadEvent) => { e.preventDefault(); e.returnValue = ''; };
        window.addEventListener('beforeunload', avertizeaza);
        return () => window.removeEventListener('beforeunload', avertizeaza);
    }, [modificari.size]);

    useEffect(() => {
        let activ = true;
        setIncarcare(true); setSchita(null); setEroare(''); setMesaj(''); setModificari(new Map());
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
    const randuri = useMemo(() => {
        const cuSchite = new Set((schita?.celule ?? []).map(c => c.angajat_id));
        for (const z of schita?.absente ?? []) if (z.activ) cuSchite.add(z.angajat_id);
        for (const z of calculat?.celule ?? []) cuSchite.add(z.angajat_id);
        return date.angajati.filter(a => (a.activ || cuSchite.has(a.id)) && (!sectie || a.sectie_id === sectie) && (!departament || a.departament_id === departament));
    }, [date.angajati, schita, calculat, sectie, departament]);
    const numarZile = [31, an % 400 === 0 || (an % 4 === 0 && an % 100 !== 0) ? 29 : 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31][luna - 1];
    const zile = Array.from({length: numarZile}, (_, i) => `${String(an).padStart(4, '0')}-${String(luna).padStart(2, '0')}-${String(i + 1).padStart(2, '0')}`);
    const sarbatori = new Map((calendar?.sarbatori ?? []).map(s => [s.data, s.denumire[i18n.language] || s.denumire.ro || Object.values(s.denumire).find(Boolean) || '']));
    const numeSectie = (id: string) => { const s = date.sectii.find(s => s.id === id); return s ? denumireSectie(s, i18n.language) : t('ui.fara_sectie'); };
    const inceputRand = Math.min(selectie.inceput.rand, selectie.sfarsit.rand), sfarsitRand = Math.max(selectie.inceput.rand, selectie.sfarsit.rand);
    const inceputZi = Math.min(selectie.inceput.zi, selectie.sfarsit.zi), sfarsitZi = Math.max(selectie.inceput.zi, selectie.sfarsit.zi);
    const valoare = (rand: number, zi: number) => {
        if (!randuri[rand] || !zile[zi]) return '';
        const k = cheie(randuri[rand].id, zile[zi]);
        const z=oreCalculate.get(k);
        return calculare?'…':z?(z.calculabil?oreText(z.minute):'—'):calculat?'':'—';
    };

    useEffect(() => { setSelectie({inceput: {rand: 0, zi: 0}, sfarsit: {rand: 0, zi: 0}}); }, [sectie, departament, an, luna]);

    function aplicaCelule(celule: {rand: number; zi: number; text: string}[]) {
        if (celule.some(c => !randuri[c.rand] || !zile[c.zi] || Array.from(c.text).length > 16 || /[\0\t\r\n]/.test(c.text))) { setEroare(t('pontaj.lipire_invalida')); return; }
        if (celule.some(c => absente.has(cheie(randuri[c.rand].id,zile[c.zi])))) { setEroare(t('absente.celula_protejata')); return; }
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

    function completeaza(text: string) {
        const celule = [];
        for (let r = inceputRand; r <= sfarsitRand; r++) for (let z = inceputZi; z <= sfarsitZi; z++) celule.push({rand: r, zi: z, text});
        aplicaCelule(celule);
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
    }

    function taste(e: KeyboardEvent<HTMLInputElement>, r: number, z: number) {
        const directii: Record<string, [number, number]> = {ArrowLeft: [0, -1], ArrowRight: [0, 1], ArrowUp: [-1, 0], ArrowDown: [1, 0], Enter: [e.shiftKey ? -1 : 1, 0]};
        if (directii[e.key]) { e.preventDefault(); const [dr, dz] = directii[e.key]; focalizeaza({rand: r + dr, zi: z + dz}, e.shiftKey && e.key !== 'Enter'); }
        else if (e.key === 'Tab') {
            let zi = z + (e.shiftKey ? -1 : 1), rand = r;
            if (zi < 0) { rand--; zi = numarZile - 1; } else if (zi >= numarZile) { rand++; zi = 0; }
            if (rand >= 0 && rand < randuri.length) { e.preventDefault(); focalizeaza({rand, zi}, false); }
        } else if (e.key === 'Delete' || e.key === 'Backspace') { e.preventDefault(); completeaza(''); }
        else if(e.key==='F2'){e.preventDefault();manualInput.current?.focus();}
        else if (e.key === 'Escape') { e.preventDefault(); aplicaCelule([{rand: r, zi: z, text: originale.get(cheie(randuri[r].id, zile[z]))?.valoare ?? ''}]); }
    }

    function poateAbandona() { return !modificari.size || window.confirm(t('ui.renunta')); }
    async function salveaza() {
        if (blocat.current) return;
        blocat.current = true; setOcupat(true); onOcupat(true); setEroare(''); setMesaj('');
        try {
            const r = await SalveazaSchita(an, luna, [...modificari.values()]);
            if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
            else if (r.schita) { setSchita(r.schita); setCalculat(r.schita.calcul); setModificari(new Map()); setMesaj(t('pontaj.salvat')); }
            else setEroare(t('date.eroare'));
        } catch { setEroare(t('date.eroare')); }
        finally { blocat.current = false; setOcupat(false); onOcupat(false); }
    }

    return <section className="pontaj-lunar">
        <div className="unelte-calendar unelte-perioada"><form onSubmit={e => { e.preventDefault(); const nou = Number(anText); if (!Number.isInteger(nou) || nou < 1 || nou > 9999) { setEroare(t('sarbatori.an_invalid')); return; } if (nou !== an) { if (poateAbandona()) setAn(nou); else setAnText(String(an)); } }}>
            <label>{t('sarbatori.an')}<input aria-label={t('sarbatori.an')} type="number" min="1" max="9999" required disabled={ocupat} value={anText} onChange={e => setAnText(e.target.value)}/></label>
            <button className="buton secundar" disabled={ocupat}>{t('sarbatori.arata')}</button>
            <label>{t('pontaj.luna')}<select aria-label={t('pontaj.luna')} value={luna} disabled={ocupat} onChange={e => { if (poateAbandona()) setLuna(Number(e.target.value)); }}>
                {Array.from({length: 12}, (_, i) => <option key={i} value={i + 1}>{new Intl.DateTimeFormat(i18n.language, {month: 'long', timeZone: 'UTC'}).format(new Date(`2026-${String(i + 1).padStart(2, '0')}-01T12:00:00Z`))}</option>)}
            </select></label></form>
            <button className="buton primar" disabled={ocupat || incarcare || !modificari.size} onClick={salveaza}>{t(ocupat ? 'ui.salvare' : 'pontaj.salveaza')}</button></div>
        <p className="nota-calendar eticheta-schita">{t('pontaj.schita_indiciu')}</p>
        {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}{mesaj && <div className="mesaj succes" role="status">{mesaj}</div>}
        {incarcare ? <p role="status">{t('ui.incarcare')}</p> : !schita ? <button className="buton secundar" onClick={() => setReincarcare(x => x + 1)}>{t('ui.reincearca')}</button> : <>
            <div className="unelte-grila"><select aria-label={t('ui.sectie')} value={sectie} disabled={ocupat} onChange={e => { setSectie(e.target.value); setDepartament(''); }}><option value="">{t('ui.toate_sectiile')}</option>
                {date.sectii.map(s => <option key={s.id} value={s.id}>{denumireSectie(s, i18n.language)}</option>)}</select>
                <select aria-label={t('ui.departament')} value={departament} disabled={ocupat} onChange={e => setDepartament(e.target.value)}><option value="">{t('ui.toate_departamentele')}</option>
                    {(date.departamente ?? []).filter(d => !sectie || d.sectie_id === sectie).map(d => <option key={d.id} value={d.id}>{!sectie ? `${numeSectie(d.sectie_id)} / ` : ''}{denumireDepartament(d, i18n.language)}</option>)}</select>
                <select aria-label={t('pontaj.program_completare')} value={programCompletare} disabled={ocupat} onChange={e=>setProgramCompletare(e.target.value)}><option value="">{t('pontaj.program_completare')}</option>{schimburi.filter(s=>s.activ).map(s=><option key={s.id} value={s.cod}>{s.denumire[i18n.language]||s.denumire.ro} {t(`schimburi.regim_${s.regim||'GENERAL'}`)} · {s.ora_start<0?t('pontaj.fara_ora'): `${String(Math.floor(s.ora_start/60)).padStart(2,'0')}:${String(s.ora_start%60).padStart(2,'0')}`} · {t('pontaj.durata_ore',{ore:oreText(s.durata_minute)})}</option>)}</select>
                <button className="buton secundar" disabled={ocupat || !randuri.length || !programCompletare} onClick={()=>completeaza(programCompletare)}>{t('pontaj.aplica_program')}</button>
                <input ref={manualInput} aria-label={t('pontaj.valoare_completare')} placeholder={t('pontaj.valoare_completare')} inputMode="decimal" maxLength={16} value={valoareCompletare} disabled={ocupat} onChange={e => setValoareCompletare(e.target.value)}/>
                <button className="buton secundar" disabled={ocupat || !randuri.length} onClick={() => {const v=manual(valoareCompletare);if(v===null)setEroare(t('pontaj.ore_invalide'));else completeaza(v);}}>{t('pontaj.completeaza_selectia')}</button>
                <button className="buton secundar" disabled={ocupat || calculare || !calculat || !randuri.length} onClick={async () => { try { await navigator.clipboard.writeText(textSelectat()); } catch { setEroare(t('pontaj.clipboard')); } }}>{t('pontaj.copiaza')}</button>
                <button className="legatura" disabled={ocupat} onClick={() => { if (poateAbandona()) setReincarcare(x => x + 1); }}>{t('pontaj.reincarca')}</button></div>
            <p className="indiciu-grila">{t('pontaj.tastatura')} <span>{t('pontaj.modificate', {numar: modificari.size})}</span></p>
            <div className="panou grila-container" ref={grila} onPaste={lipeste} onCopy={e => { e.preventDefault(); e.clipboardData.setData('text/plain', textSelectat()); }}>
                <table className="grila-pontaj" aria-label={t('pontaj.grila', {an, luna})}><thead><tr><th className="fix-nume">{t('ui.angajat')}</th><th className="fix-marca">{t('ui.marca')}</th>
                    {zile.map((data, z) => { const d = new Date(`${data}T12:00:00Z`), weekend = [0, 6].includes(d.getUTCDay()); return <th key={data} className={sarbatori.has(data) ? 'zi-sarbatoare' : weekend ? 'zi-weekend' : ''} title={sarbatori.get(data) || ''} scope="col"><strong>{z + 1}</strong><small>{new Intl.DateTimeFormat(i18n.language, {weekday: 'short', timeZone: 'UTC'}).format(d)}</small></th>; })}
                    {totaluri.map(c => <th key={c}>{t(`pontaj.total_${c}`)}</th>)}</tr></thead><tbody>
                    {randuri.map((a, r) => <tr key={a.id}><th scope="row" className="fix-nume"><strong>{a.nume} {a.prenume}</strong>{!a.activ && <small>{t('ui.inactiv')}</small>}</th><td className="fix-marca">{a.marca}</td>
                        {zile.map((data, z) => {
                            const weekend = [0, 6].includes(new Date(`${data}T12:00:00Z`).getUTCDay());
                            const selectata = r >= inceputRand && r <= sfarsitRand && z >= inceputZi && z <= sfarsitZi;
                            const k = cheie(a.id, data);
                            const absenta = absente.get(k);
                            const stareAbsenta = absenta?.stare_aprobare || 'IN_ASTEPTARE';
                            const inlocuitor = absenta?.inlocuitor_id ? date.angajati.find(x => x.id === absenta.inlocuitor_id) : undefined;
                            const titluAbsenta = absenta ? [t('absente.celula_protejata'), `${t(`absente.aprobare_${stareAbsenta}`)}${absenta.aprobat_de ? ` (${absenta.aprobat_de})` : ''}`, absenta.nr_inregistrare ? `${t('absente.nr_inregistrare')}: ${absenta.nr_inregistrare}` : '', inlocuitor ? `${t('absente.inlocuitor')}: ${inlocuitor.nume} ${inlocuitor.prenume}` : ''].filter(Boolean).join(' · ') : '';
                            return <td key={data} className={`${sarbatori.has(data) ? 'zi-sarbatoare' : weekend ? 'zi-weekend' : ''} ${selectata ? 'celula-selectata' : ''} ${modificari.has(k) ? 'celula-modificata' : ''}`}>
                            <input data-rand={r} data-zi={z} aria-label={t('pontaj.celula', {nume: `${a.nume} ${a.prenume}`, data})} title={[titluAbsenta,sarbatori.get(data)||data,t('pontaj.editare_ore'),oreCalculate.get(k)?.surse.length?t('pontaj.surse_ore',{date:oreCalculate.get(k)!.surse.join(', ')}):''].filter(Boolean).join(' · ')} disabled={ocupat} readOnly className={absenta ? `celula-absenta celula-absenta-${stareAbsenta}` : ''} value={valoare(r, z)}
                                onFocus={() => { if (ignoraFocus.current) ignoraFocus.current = false; else setSelectie({inceput: {rand: r, zi: z}, sfarsit: {rand: r, zi: z}}); }}
                                onMouseDown={e => { if (e.shiftKey) { e.preventDefault(); focalizeaza({rand: r, zi: z}, true); } }}
                                onKeyDown={e => taste(e, r, z)} onDoubleClick={()=>manualInput.current?.focus()}/></td>;
                        })}
                        {totaluri.map(c => {const v=totaluriCalculate.get(a.id)?.[c],regula=schita?.configuratie.reguli.find(r=>r.coloana===c);return <td key={c} className={!v?.definita?'total-necalculat':''} title={!v?.definita?t('pontaj.necalculat'):t(v.unitate==='ZILE'?'pontaj.unitate_zile':'pontaj.unitate_ore')}>{calculare?'…':v?valoareText(v):calculat&&regula?.sursa!=='NEDEFINITA'?'0':'—'}</td>;})}</tr>)}
                </tbody></table>{randuri.length === 0 && <p className="fara-rezultate">{t('pontaj.fara_angajati')}</p>}
            </div>
            <p className="indiciu-grila">{t('pontaj.dimensiuni', {randuri: randuri.length, zile: numarZile})} · {t('pontaj.inactivi_indiciu')}</p>
        </>}
    </section>;
}
