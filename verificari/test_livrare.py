"""F02-T22: probe negative pentru integritatea livrarii, pe fisiere fictive."""
from pathlib import Path
import json
import tempfile
import unittest
import zipfile
from documentatie import RADACINA, traduceri, legaturi
from pachet_release import impacheteaza, suma


class Livrare(unittest.TestCase):
    def setUp(self):
        self.temporar = tempfile.TemporaryDirectory()
        self.addCleanup(self.temporar.cleanup)
        self.director = Path(self.temporar.name) / 'livrare'
        self.director.mkdir()
        versiune = '0.0.1'
        catalog = json.loads((RADACINA / 'internal/ajutor/catalog.json').read_text(encoding='utf-8'))
        html = '<html>' + ''.join(f'<h3 id="ajutor-{a["id"]}">Fictiv</h3>' for a in catalog) + '</html>'
        (self.director / 'pontix.exe').write_bytes(b'EXECUTABIL FICTIV - NU SE RULEAZA')
        fisiere = []
        for ext in ('md', 'html', 'pdf'):
            nume = f'manual-utilizare-ro-{versiune}.{ext}'
            continut = html.encode() if ext == 'html' else b'DOCUMENT FICTIV'
            (self.director / nume).write_bytes(continut)
            fisiere.append({'nume': nume, 'sha256': suma(continut)})
        self.manifest = {'versiune': versiune, 'commit': 'abcdef0', 'stare': 'complet', 'executabil': 'pontix.exe', 'sha256': suma((self.director / 'pontix.exe').read_bytes()), 'manuale': {'versiune': versiune, 'commit': 'abcdef0', 'fisiere': fisiere}}
        self.scrie_manifest()

    def scrie_manifest(self):
        (self.director / 'build.json').write_text(json.dumps(self.manifest), encoding='utf-8')

    def test_lista_explicit_nu_include_date(self):
        (self.director / 'karios').mkdir()
        (self.director / 'karios/date-fictive.db').write_bytes(b'DATE FICTIVE')
        (self.director / 'parola-fictiva.txt').write_text('NU DISTRIBUI')
        arhiva = impacheteaza(self.director)
        with zipfile.ZipFile(arhiva) as z:
            nume = z.namelist()
            self.assertFalse(any('date-fictive' in n or 'parola' in n for n in nume))
            self.assertTrue(any(n.endswith('/CHECKLIST-ACCEPTANTA.md') for n in nume))
            self.assertTrue(any(n.endswith('/SHA256SUMS.txt') for n in nume))
            note = next(n for n in nume if n.endswith('/NOTE-VERSIUNE.md'))
            self.assertNotIn('{{VERSIUNE}}', z.read(note).decode())
        self.assertTrue(Path(str(arhiva) + '.sha256').is_file())

    def test_executabil_modificat_refuzat(self):
        (self.director / 'pontix.exe').write_bytes(b'MODIFICAT')
        with self.assertRaisesRegex(ValueError, 'Hash incorect'):
            impacheteaza(self.director)
        self.assertFalse(self.director.with_suffix('.zip').exists())

    def test_suprascriere_refuzata(self):
        arhiva = impacheteaza(self.director)
        vechi = arhiva.read_bytes()
        with self.assertRaises(FileExistsError):
            impacheteaza(self.director)
        self.assertEqual(vechi, arhiva.read_bytes())

    def test_versiune_manual_diferita_refuzata(self):
        self.manifest['manuale']['versiune'] = '9.9.9'
        self.scrie_manifest()
        with self.assertRaisesRegex(ValueError, 'nu corespund'):
            impacheteaza(self.director)

    def test_cale_executabil_refuzata(self):
        self.manifest['executabil'] = '../pontix.exe'
        self.scrie_manifest()
        with self.assertRaisesRegex(ValueError, 'Executabil invalid'):
            impacheteaza(self.director)

    def test_ancora_lipsa_si_duplicata(self):
        p = self.director / 'proba.html'
        p.write_text('<a href="#lipsa">Fictiv</a><p id="dublu"></p><p id="dublu"></p>')
        erori = legaturi(p)
        self.assertTrue(any('inexistenta' in e for e in erori))
        self.assertTrue(any('duplicata' in e for e in erori))
        self.assertTrue(any('Articol fara ancora' in e for e in erori))

    def test_traducere_lipsa_parametri_si_cheie_statica(self):
        root = self.director / 'repo-fictiv'
        (root / 'internal/i18n/locales').mkdir(parents=True)
        (root / 'frontend/src').mkdir(parents=True)
        for limba, continut in [('ro', {'titlu': 'Buna {{nume}}', 'numai_ro': 'Fictiv'}), ('en', {'titlu': 'Hello {{name}}'})]:
            (root / f'internal/i18n/locales/{limba}.json').write_text(json.dumps(continut))
        (root / 'frontend/src/proba.tsx').write_text("t('cheie_inexistenta')")
        erori = traduceri(root)
        self.assertTrue(any('Cheie lipsa' in e for e in erori))
        self.assertTrue(any('Parametri diferiti' in e for e in erori))
        self.assertTrue(any('Cheie statica inexistenta' in e for e in erori))


if __name__ == '__main__':
    unittest.main()
