// Capturi pentru analiza UX (F16-T36): parcurge ecranele principale cu puntea simulata
// si date fictive, fara asertiuni de continut; doar navigarea trebuie sa reuseasca.
export async function capturiUX({cdp, evalueaza, asteapta, apasa, completeaza, cataloage, limba, captura}) {
    const t = k => cataloage[limba][k];
    const masoara = () => evalueaza(`(()=>({latime:innerWidth,inaltime:innerHeight,derulareVerticala:document.documentElement.scrollHeight,derulareOrizontala:document.documentElement.scrollWidth>innerWidth,antetJos:Math.round(document.querySelector('header.antet')?.getBoundingClientRect().bottom||0),primulRandGrila:Math.round(document.querySelector('.grila-pontaj tbody tr:nth-child(2)')?.getBoundingClientRect().bottom||0)}))()`);
    const masuri = {};
    for (const [latime, inaltime] of [[1024, 640], [1440, 900]]) {
        await cdp('Emulation.setDeviceMetricsOverride', {width: latime, height: inaltime, deviceScaleFactor: 1, mobile: false});
        const sufix = `${latime}`;
        // Aplicatia se incarca o data; la a doua latime se continua din ecranul curent.
        if (latime === 1024) await asteapta(() => evalueaza("Boolean(document.querySelector('.actiuni-antet'))"), 'Aplicatie lipsa');
        const ecrane = [
            ['angajati', t('ui.angajati'), "Boolean(document.querySelector('.lista-angajati, table'))"],
            ['pontaj', t('ui.pontaj'), "Boolean(document.querySelector('.grila-pontaj'))"],
            ['planificare', t('ui.planificare'), "Boolean(document.querySelector('.grila-planificare, .modul-planificare'))"],
            ['absente', t('ui.absente'), "Boolean(document.querySelector('.modul-absente'))"],
            ['rapoarte', t('ui.rapoarte'), "Boolean(document.querySelector('main'))"],
            ['asistent', t('ui.asistent_ai'), "Boolean(document.querySelector('main'))"],
            ['setari', t('ui.setari'), "Boolean(document.querySelector('.continut-setari, main'))"],
        ];
        for (const [nume, buton, gata] of ecrane) {
            await apasa(buton);
            await asteapta(() => evalueaza(gata), 'Ecran lipsa: ' + nume);
            if (nume === 'pontaj') {
                await completeaza('.unelte-perioada input[type=number]', '2026');
                await evalueaza("document.querySelector('.unelte-perioada form').requestSubmit()");
                await completeaza('select[aria-label="' + t('pontaj.luna') + '"]', '10');
                await asteapta(() => evalueaza("Boolean(document.querySelector('.grila-pontaj tbody input'))"), 'Grila goala');
            }
            // Ecranele care incarca datele se fotografiaza dupa incarcare.
            await asteapta(() => evalueaza(`!document.querySelector('main')?.innerText.includes(${JSON.stringify(t('ui.incarcare'))})`), 'Incarcare nesfarsita: ' + nume);
            await evalueaza('window.scrollTo(0,0)');
            await new Promise(r => setTimeout(r, 250));
            masuri[`${nume}-${sufix}`] = await masoara();
            await captura(`${nume}-${sufix}.png`);
            if (nume === 'pontaj' || nume === 'absente') {
                await evalueaza("(document.querySelector('.grila-pontaj, .tabel-absente')||document.body).scrollIntoView({block:'start'})");
                await new Promise(r => setTimeout(r, 250));
                await captura(`${nume}-tabel-${sufix}.png`);
                if (nume === 'pontaj') { await evalueaza("window.scrollTo(0,0);document.querySelector('.meniu-mai-multe')?.setAttribute('open','')"); await new Promise(r => setTimeout(r, 200)); await captura(`pontaj-meniu-${sufix}.png`); await evalueaza("document.querySelector('.meniu-mai-multe')?.removeAttribute('open')"); }
            }
        }
    }
    return masuri;
}
