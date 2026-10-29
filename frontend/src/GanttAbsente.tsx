import {useMemo, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {datebaza} from '../wailsjs/go/models';
import {denumireDepartament, denumireSectie} from './Formulare';

export type IntervalAbsenta = {
	angajatId: string;
	nume: string;
	marca: string;
	departamentId: string;
	sectieId: string;
	cod: string;
	categorie: string;
	stareAprobare: string;
	deLa: string;
	panaLa: string;
	numarZile: number;
};

export type SuprapunereDepartament = {
	departamentId: string;
	departamentNume: string;
	deLa: string;
	panaLa: string;
	angajati: {id: string; nume: string; marca: string; cod: string}[];
	numarZile: number;
};

export function PanouGanttAbsente({
	an,
	date,
	zileAbsenta,
	tipuriAbsenta,
}: {
	an: number;
	date: datebaza.Date;
	zileAbsenta: datebaza.ZiAbsenta[];
	tipuriAbsenta: datebaza.TipAbsenta[];
}) {
	const {t, i18n} = useTranslation();
	const [lunaSelectata, setLunaSelectata] = useState<number>(0); // 0 = Panoramă anuală (12 luni), 1..12 = Lună detaliată
	const [filtruDepartament, setFiltruDepartament] = useState<string>('');
	const [doarCuSuprapuneri, setDoarCuSuprapuneri] = useState<boolean>(false);

	const numeLuni = [
		'Ianuarie', 'Februarie', 'Martie', 'Aprilie', 'Mai', 'Iunie',
		'Iulie', 'August', 'Septembrie', 'Octombrie', 'Noiembrie', 'Decembrie',
	];

	const hartaTipuri = useMemo(() => {
		const m = new Map<string, datebaza.TipAbsenta>();
		for (const tip of tipuriAbsenta) {
			m.set(tip.cod, tip);
		}
		return m;
	}, [tipuriAbsenta]);

	const hartaAngajati = useMemo(() => {
		const m = new Map<string, datebaza.Angajat>();
		for (const a of date.angajati) {
			m.set(a.id, a);
		}
		return m;
	}, [date.angajati]);

	// 1. Grupăm zilele consecutive ale aceluiași angajat și tip în intervale de absență
	const intervale = useMemo(() => {
		const rez: IntervalAbsenta[] = [];
		const active = zileAbsenta.filter(z => z.activ && z.data.startsWith(`${an}-`));

		// Grupăm pe angajat
		const peAngajat = new Map<string, datebaza.ZiAbsenta[]>();
		for (const z of active) {
			const list = peAngajat.get(z.angajat_id) ?? [];
			list.push(z);
			peAngajat.set(z.angajat_id, list);
		}

		for (const [angId, zile] of peAngajat.entries()) {
			const ang = hartaAngajati.get(angId);
			if (!ang) continue;

			zile.sort((a, b) => a.data.localeCompare(b.data));

			let crtInterval: {
				cod: string;
				stare: string;
				deLa: string;
				panaLa: string;
				zile: string[];
			} | null = null;

			for (const z of zile) {
				if (!crtInterval) {
					crtInterval = {cod: z.cod, stare: z.stare_aprobare || 'IN_ASTEPTARE', deLa: z.data, panaLa: z.data, zile: [z.data]};
					continue;
				}

				// Verificăm dacă z.data este ziua imediat următoare celei precedente
				const prevDate: Date = new Date(`${crtInterval.panaLa}T12:00:00Z`);
				const nextExpected: string = new Date(prevDate.getTime() + 86400000).toISOString().slice(0, 10);

				if (z.data === nextExpected && z.cod === crtInterval.cod && (z.stare_aprobare || 'IN_ASTEPTARE') === crtInterval.stare) {
					crtInterval.panaLa = z.data;
					crtInterval.zile.push(z.data);
				} else {
					const tip = hartaTipuri.get(crtInterval.cod);
					rez.push({
						angajatId: angId,
						nume: `${ang.nume} ${ang.prenume}`,
						marca: ang.marca,
						departamentId: ang.departament_id,
						sectieId: ang.sectie_id,
						cod: crtInterval.cod,
						categorie: tip?.categorie || 'ODIHNA',
						stareAprobare: crtInterval.stare,
						deLa: crtInterval.deLa,
						panaLa: crtInterval.panaLa,
						numarZile: crtInterval.zile.length,
					});
					crtInterval = {cod: z.cod, stare: z.stare_aprobare || 'IN_ASTEPTARE', deLa: z.data, panaLa: z.data, zile: [z.data]};
				}
			}

			if (crtInterval) {
				const tip = hartaTipuri.get(crtInterval.cod);
				rez.push({
					angajatId: angId,
					nume: `${ang.nume} ${ang.prenume}`,
					marca: ang.marca,
					departamentId: ang.departament_id,
					sectieId: ang.sectie_id,
					cod: crtInterval.cod,
					categorie: tip?.categorie || 'ODIHNA',
					stareAprobare: crtInterval.stare,
					deLa: crtInterval.deLa,
					panaLa: crtInterval.panaLa,
					numarZile: crtInterval.zile.length,
				});
			}
		}

		return rez;
	}, [an, zileAbsenta, hartaAngajati, hartaTipuri]);

	// 2. Detectăm suprapunerile de absență pe același departament
	const {suprapuneri, zileCuSuprapunere} = useMemo(() => {
		const hartaZileDept = new Map<string, {id: string; nume: string; marca: string; cod: string}[]>();
		const setZileSuprapuse = new Set<string>();

		const active = zileAbsenta.filter(z => z.activ && z.data.startsWith(`${an}-`) && (z.stare_aprobare || 'IN_ASTEPTARE') !== 'RESPINSA');

		for (const z of active) {
			const ang = hartaAngajati.get(z.angajat_id);
			if (!ang || !ang.departament_id) continue;

			const cheie = `${ang.departament_id}|${z.data}`;
			const list = hartaZileDept.get(cheie) ?? [];
			list.push({id: ang.id, nume: `${ang.nume} ${ang.prenume}`, marca: ang.marca, cod: z.cod});
			hartaZileDept.set(cheie, list);
		}

		// Filtrăm zilele în care >= 2 angajați din același departament sunt în concediu simultan
		const suprapusePerDept = new Map<string, {data: string; angajati: {id: string; nume: string; marca: string; cod: string}[]}[]>();

		for (const [cheie, lista] of hartaZileDept.entries()) {
			if (lista.length >= 2) {
				const [deptId, data] = cheie.split('|');
				setZileSuprapuse.add(data);
				const l = suprapusePerDept.get(deptId) ?? [];
				l.push({data, angajati: lista});
				suprapusePerDept.set(deptId, l);
			}
		}

		const listaSuprapuneri: SuprapunereDepartament[] = [];

		for (const [deptId, listaZile] of suprapusePerDept.entries()) {
			listaZile.sort((a, b) => a.data.localeCompare(b.data));
			const dept = date.departamente?.find(d => d.id === deptId);
			const deptNume = dept ? denumireDepartament(dept, i18n.language) : deptId;

			let intervalCurent: {
				deLa: string;
				panaLa: string;
				angajatiMap: Map<string, {id: string; nume: string; marca: string; cod: string}>;
				zile: string[];
			} | null = null;

			for (const item of listaZile) {
				if (!intervalCurent) {
					const m = new Map();
					item.angajati.forEach(a => m.set(a.id, a));
					intervalCurent = {deLa: item.data, panaLa: item.data, angajatiMap: m, zile: [item.data]};
					continue;
				}

				const prevDate: Date = new Date(`${intervalCurent.panaLa}T12:00:00Z`);
				const nextExpected: string = new Date(prevDate.getTime() + 86400000).toISOString().slice(0, 10);

				if (item.data === nextExpected) {
					intervalCurent.panaLa = item.data;
					intervalCurent.zile.push(item.data);
					item.angajati.forEach(a => intervalCurent!.angajatiMap.set(a.id, a));
				} else {
					listaSuprapuneri.push({
						departamentId: deptId,
						departamentNume: deptNume,
						deLa: intervalCurent.deLa,
						panaLa: intervalCurent.panaLa,
						angajati: Array.from(intervalCurent.angajatiMap.values()),
						numarZile: intervalCurent.zile.length,
					});
					const m = new Map();
					item.angajati.forEach(a => m.set(a.id, a));
					intervalCurent = {deLa: item.data, panaLa: item.data, angajatiMap: m, zile: [item.data]};
				}
			}

			if (intervalCurent) {
				listaSuprapuneri.push({
					departamentId: deptId,
					departamentNume: deptNume,
					deLa: intervalCurent.deLa,
					panaLa: intervalCurent.panaLa,
					angajati: Array.from(intervalCurent.angajatiMap.values()),
					numarZile: intervalCurent.zile.length,
				});
			}
		}

		return {suprapuneri: listaSuprapuneri, zileCuSuprapunere: setZileSuprapuse};
	}, [an, zileAbsenta, hartaAngajati, date.departamente, i18n.language]);

	// Angajații filtrați pentru afișare
	const angajatiAfisati = useMemo(() => {
		return date.angajati.filter(a => {
			if (!a.activ) return false;
			if (filtruDepartament && a.departament_id !== filtruDepartament) return false;
			if (doarCuSuprapuneri) {
				const areSuprapunere = suprapuneri.some(s => s.angajati.some(ang => ang.id === a.id));
				if (!areSuprapunere) return false;
			}
			return true;
		});
	}, [date.angajati, filtruDepartament, doarCuSuprapuneri, suprapuneri]);

	const culoareCategorie = (cat: string) => {
		switch (cat) {
			case 'ODIHNA': return '#2e7d32'; // verde
			case 'MEDICAL': return '#d84315'; // portocaliu/roșu
			case 'SPECIAL': return '#1565c0'; // albastru
			default: return '#546e7a';
		}
	};

	// Determinăm numărul de zile din lună pentru vizualizarea detaliată
	const zileLunaDetaliata = useMemo(() => {
		if (lunaSelectata <= 0) return [];
		const totalZile = new Date(an, lunaSelectata, 0).getDate();
		return Array.from({length: totalZile}, (_, i) => {
			const zi = i + 1;
			const data = `${an}-${String(lunaSelectata).padStart(2, '0')}-${String(zi).padStart(2, '0')}`;
			const ziSapt = new Date(`${data}T12:00:00Z`).getDay();
			const esteWeekend = ziSapt === 0 || ziSapt === 6;
			return {zi, data, esteWeekend};
		});
	}, [an, lunaSelectata]);

	return (
		<div className="panou-gantt-absente">
			{/* Metric carduri sintetice */}
			<div className="carduri-sumar-comparativ">
				<div className="card-sumar-comp">
					<span className="eticheta-card-comp">{t('gantt.concedii_an', {an: String(an)})}</span>
					<strong className="valoare-card-comp">{intervale.length}</strong>
					<small className="detaliu-card-comp">{t('gantt.intervale')}</small>
				</div>
				<div className="card-sumar-comp">
					<span className="eticheta-card-comp">{t('gantt.suprapuneri')}</span>
					<strong className={`valoare-card-comp ${suprapuneri.length > 0 ? 'culoare-alerta' : ''}`}>
						{suprapuneri.length}
					</strong>
					<small className="detaliu-card-comp">{t('gantt.suprapuneri_detaliu')}</small>
				</div>
				<div className="card-sumar-comp">
					<span className="eticheta-card-comp">{t('gantt.angajati_activi')}</span>
					<strong className="valoare-card-comp">{angajatiAfisati.length}</strong>
					<small className="detaliu-card-comp">{t('gantt.in_selectie')}</small>
				</div>
			</div>

			{/* Bannere de avertizare pentru suprapuneri critice */}
			{suprapuneri.length > 0 && (
				<div className="caseta-avertismente-gantt" role="alert">
					<div className="antet-avertisment-gantt">
						<span>⚠️</span>
						<strong>{t('gantt.atentie_suprapuneri', {numar: String(suprapuneri.length)})}</strong>
					</div>
					<div className="lista-suprapuneri-detalii">
						{suprapuneri.slice(0, 5).map((sup, idx) => (
							<div key={idx} className="element-suprapunere">
								<span className="insigna-dept">{sup.departamentNume}</span>
								<span className="interval-sup">
									{sup.deLa} – {sup.panaLa} ({t(sup.numarZile === 1 ? 'gantt.o_zi' : 'gantt.zile', {numar: String(sup.numarZile)})})
								</span>
								<span className="colegi-sup">
									{t('gantt.colegi_absenti', {lista: sup.angajati.map(a => `${a.nume} [${a.cod}]`).join(', ')})}
								</span>
							</div>
						))}
						{suprapuneri.length > 5 && (
							<small className="text-discret">{t('gantt.inca_suprapuneri', {numar: String(suprapuneri.length - 5)})}</small>
						)}
					</div>
				</div>
			)}

			{/* Bara de filtrare și comutare vizualizare */}
			<div className="bara-filtre-comparativ">
				<div className="grup-filtre-comp">
					<label>
						{t('gantt.scara')}
						<select
							value={lunaSelectata}
							onChange={e => setLunaSelectata(Number(e.target.value))}
						>
							<option value={0}>📅 {t('gantt.panorama')}</option>
							{numeLuni.map((nume, idx) => (
								<option key={idx} value={idx + 1}>
									{t('gantt.luna', {nume, an: String(an)})}
								</option>
							))}
						</select>
					</label>

					<label>
						{t('gantt.departament')}
						<select
							value={filtruDepartament}
							onChange={e => setFiltruDepartament(e.target.value)}
						>
							<option value="">{t('gantt.toate_departamentele')}</option>
							{(date.departamente ?? []).map(d => (
								<option key={d.id} value={d.id}>
									{denumireDepartament(d, i18n.language)}
								</option>
							))}
						</select>
					</label>

					<label className="bifa-discrepante">
						<input
							type="checkbox"
							checked={doarCuSuprapuneri}
							onChange={e => setDoarCuSuprapuneri(e.target.checked)}
						/>
						<span>{t('gantt.doar_suprapuneri')}</span>
					</label>
				</div>

				<div className="legenda-comparativ">
					<span className="element-legenda"><span className="punct-legenda" style={{background: '#2e7d32'}}/> {t('gantt.legenda_odihna')}</span>
					<span className="element-legenda"><span className="punct-legenda" style={{background: '#d84315'}}/> {t('gantt.legenda_medical')}</span>
					<span className="element-legenda"><span className="punct-legenda" style={{background: '#1565c0'}}/> {t('gantt.legenda_speciale')}</span>
					<span className="element-legenda"><span className="punct-legenda" style={{background: '#f57c00'}}/> ⚠️ {t('gantt.legenda_suprapunere')}</span>
				</div>
			</div>

			{/* Grila Gantt propriu-zisă */}
			<div className="tabel-derulant-comparativ">
				{lunaSelectata === 0 ? (
					/* 1. Vizualizare panoramică anuală (12 coloane) */
					<table className="tabel-comparativ tabel-gantt-anual">
						<thead>
							<tr>
								<th className="coloana-fixa th-angajat" style={{minWidth: 200}}>{t('gantt.angajat')}</th>
								<th className="coloana-fixa th-departament" style={{minWidth: 150}}>{t('gantt.departament_coloana')}</th>
								{numeLuni.map((m, idx) => (
									<th key={idx} className="th-luna-gantt">
										{m.slice(0, 3)}
									</th>
								))}
							</tr>
						</thead>
						<tbody>
							{angajatiAfisati.map(a => {
								const dept = date.departamente?.find(d => d.id === a.departament_id);
								const deptNume = dept ? denumireDepartament(dept, i18n.language) : '—';
								const intervAngajat = intervale.filter(i => i.angajatId === a.id);

								return (
									<tr key={a.id}>
										<td className="coloana-fixa td-angajat">
											<strong>{a.nume} {a.prenume}</strong>
											<small className="text-discret"> · {a.marca}</small>
										</td>
										<td className="coloana-fixa td-departament">
											<small>{deptNume}</small>
										</td>
										{Array.from({length: 12}, (_, idx) => {
											const lunaNr = idx + 1;
											const prefixLuna = `${an}-${String(lunaNr).padStart(2, '0')}-`;
											const absenteLuna = intervAngajat.filter(i =>
												(i.deLa.startsWith(prefixLuna) || i.panaLa.startsWith(prefixLuna) || (i.deLa < prefixLuna && i.panaLa > prefixLuna))
											);

											return (
												<td key={idx} className="td-celula-gantt">
													{absenteLuna.length > 0 && (
														<div className="celula-bara-gantt">
															{absenteLuna.map((abs, aIdx) => (
																<span
																	key={aIdx}
																	className="cip-gantt-interval"
																	style={{backgroundColor: culoareCategorie(abs.categorie)}}
																	title={`${abs.cod} (${t('gantt.zile', {numar: String(abs.numarZile)})}): ${abs.deLa} – ${abs.panaLa}`}
																>
																	{abs.cod} ({abs.numarZile}z)
																</span>
															))}
														</div>
													)}
												</td>
											);
										})}
									</tr>
								);
							})}
						</tbody>
					</table>
				) : (
					/* 2. Vizualizare detaliată pe luna selectată (zilele 1..31) */
					<table className="tabel-comparativ tabel-gantt-luna">
						<thead>
							<tr>
								<th className="coloana-fixa th-angajat" style={{minWidth: 200}}>{t('gantt.angajat')}</th>
								<th className="coloana-fixa th-departament" style={{minWidth: 150}}>{t('gantt.departament_coloana')}</th>
								{zileLunaDetaliata.map(z => (
									<th
										key={z.zi}
										className={`th-zi-gantt ${z.esteWeekend ? 'weekend' : ''} ${zileCuSuprapunere.has(z.data) ? 'suprapunere-alerta' : ''}`}
										title={zileCuSuprapunere.has(z.data) ? `⚠️ ${t('gantt.suprapunere_zi')}` : z.data}
									>
										{z.zi}
									</th>
								))}
							</tr>
						</thead>
						<tbody>
							{angajatiAfisati.map(a => {
								const dept = date.departamente?.find(d => d.id === a.departament_id);
								const deptNume = dept ? denumireDepartament(dept, i18n.language) : '—';
								const zileAng = new Map<string, datebaza.ZiAbsenta>();
								for (const z of zileAbsenta) {
									if (z.activ && z.angajat_id === a.id) {
										zileAng.set(z.data, z);
									}
								}

								return (
									<tr key={a.id}>
										<td className="coloana-fixa td-angajat">
											<strong>{a.nume} {a.prenume}</strong>
											<small className="text-discret"> · {a.marca}</small>
										</td>
										<td className="coloana-fixa td-departament">
											<small>{deptNume}</small>
										</td>
										{zileLunaDetaliata.map(z => {
											const abs = zileAng.get(z.data);
											const tip = abs ? hartaTipuri.get(abs.cod) : undefined;
											const cat = tip?.categorie || 'ODIHNA';
											const esteSuprapunere = zileCuSuprapunere.has(z.data) && !!abs;

											return (
												<td
													key={z.zi}
													className={`td-zi-detaliu ${z.esteWeekend ? 'weekend' : ''} ${esteSuprapunere ? 'celula-suprapunere-activa' : ''}`}
												>
													{abs && (
														<span
															className="insigna-absenta-zi"
															style={{backgroundColor: culoareCategorie(cat)}}
															title={t('gantt.cod_pe_zi', {cod: abs.cod, data: z.data}) + (esteSuprapunere ? ` (⚠️ ${t('gantt.suprapuneri')})` : '')}
														>
															{abs.cod}
														</span>
													)}
												</td>
											);
										})}
									</tr>
								);
							})}
						</tbody>
					</table>
				)}
			</div>
		</div>
	);
}
