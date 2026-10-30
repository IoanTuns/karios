"""F02-T22: verifica traduceri si legaturi locale, fara acces la internet."""
from pathlib import Path
from html.parser import HTMLParser
from collections import Counter
import argparse
import json
import re
import sys

RADACINA = Path(__file__).resolve().parents[1]


def traduceri(radacina=RADACINA):
    erori = []
    cataloage = {l: json.loads((radacina / f'internal/i18n/locales/{l}.json').read_text(encoding='utf-8')) for l in ('ro', 'en')}
    ro, en = cataloage['ro'], cataloage['en']
    for k in sorted(ro.keys() ^ en.keys()):
        erori.append(f'Cheie lipsa intr-o limba: {k}')
    for k in ro.keys() & en.keys():
        for l in cataloage:
            if not isinstance(cataloage[l][k], str) or not cataloage[l][k].strip():
                erori.append(f'Text gol sau invalid: {l}:{k}')
        if Counter(re.findall(r'{{\s*([^{}]+?)\s*}}', ro[k])) != Counter(re.findall(r'{{\s*([^{}]+?)\s*}}', en[k])):
            erori.append(f'Parametri diferiti RO/EN: {k}')
    for fisier in (radacina / 'frontend/src').rglob('*'):
        if fisier.suffix not in ('.ts', '.tsx') or '.test.' in fisier.name:
            continue
        for m in re.finditer(r"\bt\(\s*(['\"])([^'\"]+)\1\s*[,)]", fisier.read_text(encoding='utf-8')):
            if m[2] not in ro:
                erori.append(f'Cheie statica inexistenta: {fisier.relative_to(radacina)}: {m[2]}')
    return erori


class LegaturiHTML(HTMLParser):
    def __init__(self):
        super().__init__()
        self.identificatori = []
        self.legaturi = []

    def handle_starttag(self, tag, attrs):
        a = dict(attrs)
        if 'id' in a:
            self.identificatori.append(a['id'])
        if tag == 'a' and 'href' in a:
            self.legaturi.append(a['href'])


def legaturi(cale, radacina=RADACINA):
    parser = LegaturiHTML()
    parser.feed(Path(cale).read_text(encoding='utf-8'))
    ids = Counter(parser.identificatori)
    erori = [f'Ancora duplicata: {k}' for k, n in ids.items() if n > 1]
    for href in parser.legaturi:
        if href.startswith('#') and href[1:] not in ids:
            erori.append(f'Legatura interna inexistenta: {href}')
        elif href and not href.startswith(('#', 'https:', 'http:', 'mailto:')):
            erori.append(f'Legatura neportabila in manual: {href}')
    catalog = json.loads((radacina / 'internal/ajutor/catalog.json').read_text(encoding='utf-8'))
    for articol in catalog:
        if ids['ajutor-' + articol['id']] != 1:
            erori.append(f'Articol fara ancora unica: {articol["id"]}')
    return erori


if __name__ == '__main__':
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--html', type=Path)
    a = p.parse_args()
    erori = traduceri()
    if a.html:
        erori.extend(legaturi(a.html))
    for eroare in erori:
        print(eroare)
    print(f'Verificare documentatie: {len(erori)} erori. Cheile construite dinamic si linkurile externe necesita probe separate.')
    sys.exit(bool(erori))
