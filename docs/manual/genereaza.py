"""F02-T18: produce manualele versionate MD/HTML/PDF, fără pachete externe.

Suportă convențiile folosite în manual: titluri, paragrafe, liste, tabele,
citate, text aldin și cod inline. Nu este un convertor Markdown general.
"""
from pathlib import Path
import argparse
import hashlib
import html
import json
import os
import re
import subprocess
import sys

radacina = Path(__file__).resolve().parents[1]
sursa = radacina / 'manual-utilizare-ro.md'
optiuni = argparse.ArgumentParser(description='Generează manualele aceleiași livrări Pontix.')
optiuni.add_argument('--versiune', required=True)
optiuni.add_argument('--commit', required=True)
optiuni.add_argument('--data', required=True, help='Data buildului, ISO 8601.')
optiuni.add_argument('--destinatie', type=Path)
optiuni.add_argument('--temporare', type=Path)
optiuni.add_argument('--browser', type=Path)
optiuni.add_argument('--verifica-sursa', action='store_true')
optiuni.add_argument('--fara-pdf', action='store_true', help='Numai pentru previzualizarea editorială.')
args = optiuni.parse_args()
configuratie = json.loads((radacina / 'manual/manual.json').read_text(encoding='utf-8'))
if not re.fullmatch(r'\d+\.\d+\.\d+', args.versiune):
    optiuni.error('Versiunea aplicației este invalidă.')
if configuratie['versiune_revizuita'] != args.versiune:
    optiuni.error('Manualul nu este revizuit pentru această versiune. Actualizează conținutul și docs/manual/manual.json.')
valori = {'VERSIUNE': args.versiune, 'COMMIT': args.commit, 'DATA_BUILD': args.data, 'EDITIE': configuratie['editie']}
continut = sursa.read_text(encoding='utf-8')
for cheie, valoare in valori.items():
    marcaj = '{{' + cheie + '}}'
    if marcaj not in continut:
        optiuni.error('Lipsește marcajul din sursa manualului: ' + marcaj)
    continut = continut.replace(marcaj, valoare)
if re.search(r'\{\{[A-Z_]+\}\}', continut):
    optiuni.error('Manualul conține marcaje nerezolvate.')
if args.verifica_sursa:
    sys.exit(0)
if args.destinatie is None:
    optiuni.error('Este necesar --destinatie; sursa și livrările precedente nu se suprascriu.')
args.destinatie.mkdir(parents=True, exist_ok=True)
nume = f'manual-utilizare-ro-{args.versiune}'
destinatie = args.destinatie / (nume + '.html')
fisiere = [args.destinatie / (nume + extensie) for extensie in ('.md', '.html', '.pdf')]
if any(f.exists() for f in fisiere):
    optiuni.error('Există deja manuale în destinație; folosește un director nou.')


def inline(text):
    text = html.escape(text, quote=True)
    text = re.sub(r'`([^`]+)`', r'<code>\1</code>', text)
    return re.sub(r'\*\*([^*]+)\*\*', r'<strong>\1</strong>', text)


linii = continut.splitlines()
parti, cuprins = [], []
lista, paragraf, sectiune = None, [], False


def inchide_paragraf():
    if paragraf:
        parti.append('<p>' + inline(' '.join(paragraf)) + '</p>')
        paragraf.clear()


def inchide_lista():
    global lista
    if lista:
        parti.append('</' + lista + '>')
        lista = None


i = 0
while i < len(linii):
    linie = linii[i].strip()
    if not linie:
        inchide_paragraf()
        inchide_lista()
        i += 1
        continue
    titlu = re.match(r'^(#{1,3}) (.+)$', linie)
    if titlu:
        inchide_paragraf()
        inchide_lista()
        nivel, text = len(titlu[1]), titlu[2]
        if nivel == 2:
            if sectiune:
                parti.append('</section>')
            identificator = 'capitol-' + str(len(cuprins) + 1)
            cuprins.append((identificator, text))
            parti.append(f'<section class="capitol" id="{identificator}" aria-labelledby="{identificator}-titlu">')
            sectiune = True
            parti.append(f'<h2 id="{identificator}-titlu">{inline(text)}</h2>')
        else:
            parti.append(f'<h{nivel}>{inline(text)}</h{nivel}>')
    elif linie.startswith('|'):
        inchide_paragraf()
        inchide_lista()
        randuri = []
        while i < len(linii) and linii[i].strip().startswith('|'):
            celule = [c.strip() for c in linii[i].strip().strip('|').split('|')]
            if not all(re.fullmatch(r':?-+:?', c) for c in celule):
                randuri.append(celule)
            i += 1
        parti.append('<table><thead><tr>' + ''.join('<th scope="col">' + inline(c) + '</th>' for c in randuri[0]) + '</tr></thead><tbody>')
        for rand in randuri[1:]:
            parti.append('<tr>' + ''.join('<td>' + inline(c) + '</td>' for c in rand) + '</tr>')
        parti.append('</tbody></table>')
        continue
    elif re.match(r'^\d+\. ', linie) or linie.startswith('- '):
        inchide_paragraf()
        tip = 'ul' if linie.startswith('- ') else 'ol'
        if lista != tip:
            inchide_lista()
            parti.append('<' + tip + '>')
            lista = tip
        parti.append('<li>' + inline(re.sub(r'^(?:\d+\.|-) ', '', linie)) + '</li>')
    elif linie.startswith('> '):
        inchide_paragraf()
        inchide_lista()
        parti.append('<blockquote>' + inline(linie[2:]) + '</blockquote>')
    else:
        inchide_lista()
        paragraf.append(linie)
    i += 1
inchide_paragraf()
inchide_lista()
if sectiune:
    parti.append('</section>')

corp = '\n'.join(parti)
primul = corp.index('<section ')
introducere, capitole = corp[:primul], corp[primul:]
legaturi = ''.join(f'<a href="#{identificator}" data-capitol="{identificator}">{html.escape(text)}</a>' for identificator, text in cuprins)
css = """
:root{font-family:'Segoe UI',Arial,sans-serif;color:#203936;background:#f3f7f6;font-size:16px;line-height:1.6}
*{box-sizing:border-box}body{margin:0}a{color:#12685e}button,input{font:inherit}button{cursor:pointer}
.bara{position:fixed;inset:0 auto 0 0;width:295px;padding:24px;background:#123f3a;color:white;overflow:auto}
.bara strong{display:block;font-size:24px}.bara p{font-size:13px;color:#cce1dc}.bara label{display:block;margin:18px 0 5px}
input{width:100%;padding:9px;border:1px solid #8cb5ac;border-radius:6px}button{padding:7px 12px;border:1px solid #9fbfb6;border-radius:6px;background:#fff;color:#194d43;margin:7px 3px 7px 0}
nav a{display:block;color:#eaf5f2;text-decoration:none;padding:7px 4px;border-bottom:1px solid #37615b;font-size:13px;line-height:1.4}
nav a:hover,nav a:focus{background:#285d53;text-decoration:underline}a:focus-visible,button:focus-visible,input:focus-visible{outline:3px solid #efbd56;outline-offset:3px}
main{max-width:1140px;padding:38px 52px;margin-left:295px}.coperta{padding:32px;border:1px solid #b8d2c9;border-radius:12px;background:#e7f1ed}
h1{font-size:38px;line-height:1.15;margin:0 0 20px;color:#124d43}h2{font-size:25px;line-height:1.3;color:#124d43;margin:0 0 20px}h3{font-size:18px;color:#185c50;margin:26px 0 8px}
section{margin:24px 0;padding:32px;background:white;border:1px solid #d6e4de;border-radius:10px;scroll-margin-top:20px}
li{margin:7px 0}p{margin:12px 0}table{width:100%;border-collapse:collapse;margin:20px 0;font-size:14px}th,td{text-align:left;vertical-align:top;padding:10px;border:1px solid #c9dcd4}th{background:#e6f0eb}tr:nth-child(even){background:#f7faf8}
code{font-size:.9em;background:#edf2f0;padding:1px 4px;border-radius:3px;overflow-wrap:anywhere}blockquote{margin:16px 0;border-left:4px solid #398b74;padding:12px 18px;background:#edf6f1}
.cuprins-tipar{display:none}.ascuns,[hidden]{display:none!important}#rezultate{font-size:13px;min-height:22px}#fara-rezultate{padding:24px;border:1px solid #d3aa59;background:#fff7df}
@media(max-width:850px){.bara{position:static;width:auto}nav{max-height:220px;overflow:auto}main{margin:0;padding:20px}.coperta,section{padding:20px}h1{font-size:30px}table{font-size:12px}td,th{padding:6px}}
@page{size:A4;margin:17mm 16mm;@bottom-left{content:'Pontix · 0.88.7 · Manual 1.0';font:8pt 'Segoe UI';color:#50665f}@bottom-right{content:'Pagina ' counter(page) ' / ' counter(pages);font:8pt 'Segoe UI';color:#50665f}}
@media print{html,body{background:white;font-size:10pt;color:#172f2a}.bara,#fara-rezultate{display:none!important}main{margin:0;padding:0;max-width:none}.coperta{padding:0;border:0;background:none;break-after:page}h1{font-size:32pt;padding-top:14mm}h2{font-size:18pt}h3{font-size:12pt}section{border:0;border-radius:0;margin:0 0 10mm;padding:0;display:block!important}h2,h3{break-after:avoid}p,li{orphans:3;widows:3}table{font-size:9pt}th,td{padding:5px}thead{display:table-header-group}tr{break-inside:avoid}blockquote{break-inside:avoid}.cuprins-tipar{display:block;margin:9mm 0}.cuprins-tipar a{display:block;font-size:9pt;line-height:1.6;color:#174d42;text-decoration:none}.cuprins-tipar div{columns:2;column-gap:12mm}a{color:inherit}}
"""
css = css.replace('0.88.7', args.versiune).replace('Manual 1.0', 'Manual ' + configuratie['editie'])
js = """
const cautare=document.getElementById('cautare'),rezultate=document.getElementById('rezultate');
const capitole=[...document.querySelectorAll('.capitol')],legaturi=[...document.querySelectorAll('nav a')];
const normal=s=>s.normalize('NFD').replace(/[\\u0300-\\u036f]/g,'').toLocaleLowerCase('ro');
const texte=capitole.map(c=>normal(c.textContent));
function filtreaza(){const termeni=normal(cautare.value.trim()).split(/\\s+/).filter(Boolean);let numar=0;capitole.forEach((c,i)=>{const gasit=termeni.every(t=>texte[i].includes(t));c.hidden=!gasit;legaturi[i].hidden=!gasit;if(gasit)numar++});rezultate.textContent=termeni.length?`${numar} capitole găsite din ${capitole.length}`:`${capitole.length} capitole disponibile`;document.getElementById('fara-rezultate').hidden=numar!==0;}
cautare.addEventListener('input',filtreaza);
document.getElementById('reseteaza').addEventListener('click',()=>{cautare.value='';filtreaza();cautare.focus()});
document.getElementById('tipareste').addEventListener('click',()=>window.print());
window.addEventListener('beforeprint',()=>capitole.forEach(c=>c.hidden=false));window.addEventListener('afterprint',filtreaza);
filtreaza();
if(new URLSearchParams(location.search).has('verificare')){try{if(capitole.length!==19)throw Error('Număr incorect de capitole');cautare.value='CONCEDIU';filtreaza();if(!capitole.some(c=>!c.hidden)||!capitole.some(c=>c.hidden))throw Error('Filtrare incorectă');cautare.value='termeninexistent987654';filtreaza();if(!capitole.every(c=>c.hidden))throw Error('Rezultate inexistente');document.getElementById('reseteaza').click();if(capitole.some(c=>c.hidden))throw Error('Resetare incorectă');if(legaturi.some(a=>!document.getElementById(a.dataset.capitol)))throw Error('Legătură invalidă');document.body.dataset.verificare='trecut';}catch(e){document.body.dataset.verificare=e.message;}}
"""
js = js.replace('capitole.length!==19', f'capitole.length!=={len(cuprins)}')
document = f'''<!doctype html>
<html lang="ro"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Pontix — Manual de utilizare în română · {html.escape(args.versiune)}</title><style>{css}</style></head>
<body><aside class="bara"><strong>Pontix</strong><p>Manual de utilizare · {html.escape(args.versiune)}<br>Ediția {html.escape(configuratie['editie'])}<br>Build: {html.escape(args.data)}<br>Commit: {html.escape(args.commit)}</p>
<label for="cautare">Caută în manual</label><input id="cautare" type="search" placeholder="De exemplu: concediu, S7, backup">
<button id="reseteaza" type="button">Resetează</button><button id="tipareste" type="button">Tipărește / PDF</button>
<p id="rezultate" role="status" aria-live="polite"></p><nav aria-label="Cuprins">{legaturi}</nav></aside>
<main><header class="coperta">{introducere}<div class="cuprins-tipar"><h2>Cuprins</h2><div>{legaturi}</div></div></header>
<p id="fara-rezultate" hidden>Nu am găsit capitole pentru această căutare. Încearcă un termen mai scurt sau resetează căutarea.</p>{capitole}</main><script>{js}</script></body></html>'''
fisiere[0].write_text(continut, encoding='utf-8', newline='\n')
destinatie.write_text(document, encoding='utf-8', newline='\n')
if not args.fara_pdf:
    browser = args.browser or next((Path(os.environ.get(variabila, '')) / 'Microsoft/Edge/Application/msedge.exe'
        for variabila in ('ProgramFiles(x86)', 'ProgramFiles')
        if (Path(os.environ.get(variabila, '')) / 'Microsoft/Edge/Application/msedge.exe').is_file()), None)
    if browser is None or not browser.is_file():
        optiuni.error('Microsoft Edge nu este disponibil pentru generarea PDF. Instalează-l separat sau indică --browser.')
    if args.temporare is None:
        optiuni.error('Generarea PDF necesită --temporare, într-un director al buildului.')
    profil = args.temporare / 'profil-manual'
    profil.mkdir(parents=True, exist_ok=False)
    try:
        rezultat = subprocess.run([str(browser), '--headless', '--disable-gpu', '--no-first-run', '--no-default-browser-check',
            '--no-pdf-header-footer', '--user-data-dir=' + str(profil.resolve()),
            '--print-to-pdf=' + str(fisiere[2].resolve()), destinatie.resolve().as_uri()],
            capture_output=True, timeout=60, creationflags=getattr(subprocess, 'CREATE_NO_WINDOW', 0))
        if rezultat.returncode != 0:
            optiuni.error('Generarea PDF a eșuat, cod browser ' + str(rezultat.returncode))
    except subprocess.TimeoutExpired:
        optiuni.error('Generarea PDF a depășit 60 de secunde.')
    if not fisiere[2].is_file():
        optiuni.error('Browserul nu a produs manualul PDF.')
    pdf = fisiere[2].read_bytes()
    if not pdf.startswith(b'%PDF-') or b'%%EOF' not in pdf[-1024:]:
        optiuni.error('Manualul PDF produs este incomplet.')
livrate = fisiere[:2] if args.fara_pdf else fisiere
print(json.dumps({'versiune': args.versiune, 'commit': args.commit, 'editie': configuratie['editie'],
    'creat_la': args.data, 'fisiere': [{'nume': f.name, 'sha256': hashlib.sha256(f.read_bytes()).hexdigest().upper()} for f in livrate]}))
