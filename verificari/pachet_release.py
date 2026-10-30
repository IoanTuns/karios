"""F02-T22: pachet portabil pe lista explicita, cu verificarea hashurilor."""
from pathlib import Path
import argparse
import hashlib
import json
import re
import zipfile
from documentatie import RADACINA, legaturi, traduceri


def suma(continut):
    return hashlib.sha256(continut).hexdigest().upper()


def impacheteaza(director, radacina=RADACINA):
    director = Path(director).resolve()
    m = json.loads((director / 'build.json').read_text(encoding='utf-8-sig'))
    versiune, commit = m['versiune'], m['commit']
    if m['stare'] != 'complet' or not re.fullmatch(r'\d+\.\d+\.\d+', versiune) or not re.fullmatch(r'[a-f0-9]+(?:-modificat)?', commit):
        raise ValueError('Manifest incomplet sau identitate invalida.')
    manuale = m['manuale']
    if manuale['versiune'] != versiune or manuale['commit'] != commit:
        raise ValueError('Manualele nu corespund executabilului.')
    asteptate = {f'manual-utilizare-ro-{versiune}.{ext}' for ext in ('md', 'html', 'pdf')}
    if {f['nume'] for f in manuale['fisiere']} != asteptate or len(manuale['fisiere']) != 3:
        raise ValueError('Lista manualelor nu este completa.')
    fisiere = {}
    def adauga(nume, hash_asteptat=None):
        if Path(nume).is_absolute() or '..' in Path(nume).parts or '\\' in nume:
            raise ValueError('Cale nepermisa in manifest.')
        cale = director / nume
        if cale.is_symlink() or not cale.resolve().is_relative_to(director):
            raise ValueError('Legatura de fisier nepermisa.')
        continut = cale.read_bytes()
        if hash_asteptat and suma(continut) != hash_asteptat.upper():
            raise ValueError(f'Hash incorect: {nume}')
        fisiere[nume] = continut
    if Path(m['executabil']).name != m['executabil'] or not m['executabil'].endswith('.exe'):
        raise ValueError('Executabil invalid.')
    adauga(m['executabil'], m['sha256'])
    adauga('build.json')
    for f in manuale['fisiere']:
        adauga(f['nume'], f['sha256'])
    erori = traduceri(radacina) + legaturi(director / f'manual-utilizare-ro-{versiune}.html', radacina)
    if erori:
        raise ValueError('\n'.join(erori))
    for sursa, destinatie in [('actualizare.md', 'ACTUALIZARE.md'), ('note-utilizator.md', 'NOTE-VERSIUNE.md'), ('checklist-acceptanta.md', 'CHECKLIST-ACCEPTANTA.md')]:
        text = (radacina / 'docs/livrare' / sursa).read_text(encoding='utf-8')
        text = text.replace('{{VERSIUNE}}', versiune).replace('{{COMMIT}}', commit)
        if re.search(r'{{[^{}]+}}', text):
            raise ValueError('Marcaj necompletat in documentul de livrare.')
        fisiere[destinatie] = text.encode('utf-8')
    fisiere['LICENSE'] = (radacina / 'LICENSE').read_bytes()
    if m.get('llama_server_sha256'):
        adauga('karios/ai/llama-server.exe', m['llama_server_sha256'])
    fisiere['SHA256SUMS.txt'] = ''.join(f'{suma(v)}  {k}\n' for k, v in sorted(fisiere.items())).encode('utf-8')
    iesire = director.parent / (director.name + '.zip')
    if Path(str(iesire) + '.sha256').exists():
        raise FileExistsError('Suma arhivei exista deja.')
    # Modul x refuza suprascrierea. Nu se parcurge recursiv directorul cu date.
    with zipfile.ZipFile(iesire, 'x', compression=zipfile.ZIP_DEFLATED) as z:
        for nume, continut in sorted(fisiere.items()):
            z.writestr(f'Pontix-{versiune}-{commit}/' + nume, continut)
    with zipfile.ZipFile(iesire) as z:
        if z.testzip() or len(z.namelist()) != len(fisiere):
            raise ValueError('Arhiva generata nu este valida.')
        for nume, continut in fisiere.items():
            if z.read(f'Pontix-{versiune}-{commit}/' + nume) != continut:
                raise ValueError('Continut ZIP diferit de sursa verificata.')
    with Path(str(iesire) + '.sha256').open('x', encoding='utf-8', newline='\n') as f:
        f.write(suma(iesire.read_bytes()) + '  ' + iesire.name + '\n')
    return iesire


if __name__ == '__main__':
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--director', required=True, type=Path)
    a = p.parse_args()
    print(impacheteaza(a.director))
