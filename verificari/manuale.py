"""F02-T18: livrări distincte, versiune editorială și refuzul suprascrierii."""
from pathlib import Path
import hashlib
import json
import subprocess
import sys
import tempfile
import unittest

radacina = Path(__file__).resolve().parents[1]
generator = radacina / 'docs/manual/genereaza.py'
versiune = json.loads((radacina / 'docs/manual/manual.json').read_text(encoding='utf-8'))['versiune_revizuita']


class Manuale(unittest.TestCase):
    def ruleaza(self, *argumente, versiune_ceruta=versiune):
        return subprocess.run([sys.executable, str(generator), '--versiune', versiune_ceruta,
            '--commit', 'proba-<fictiva>', '--data', '2026-10-30T12:00:00+02:00', *map(str, argumente)], capture_output=True)

    def test_refuza_versiune_nerevizuita(self):
        rezultat = self.ruleaza('--verifica-sursa', versiune_ceruta='999.0.0')
        self.assertNotEqual(rezultat.returncode, 0)
        self.assertIn(b'manual.json', rezultat.stderr)

    def test_identitate_hashuri_si_livrari_pastrate(self):
        with tempfile.TemporaryDirectory() as temporar:
            destinatie = Path(temporar) / 'prima livrare'
            rezultat = self.ruleaza('--destinatie', destinatie, '--fara-pdf')
            self.assertEqual(rezultat.returncode, 0, rezultat.stderr)
            manifest = json.loads(rezultat.stdout)
            self.assertEqual(manifest['versiune'], versiune)
            self.assertEqual(manifest['commit'], 'proba-<fictiva>')
            for fisier in manifest['fisiere']:
                continut = (destinatie / fisier['nume']).read_bytes()
                self.assertEqual(hashlib.sha256(continut).hexdigest().upper(), fisier['sha256'])
                self.assertIn(versiune.encode(), continut)
                self.assertNotIn(b'{{VERSIUNE}}', continut)
                self.assertNotIn(b'4600f30', continut)
            document = (destinatie / f'manual-utilizare-ro-{versiune}.html').read_text(encoding='utf-8')
            self.assertIn('proba-&lt;fictiva&gt;', document)
            self.assertNotIn('proba-<fictiva>', document)
            self.assertIn('Actualizează planificarea', document)
            original = {f.name: f.read_bytes() for f in destinatie.iterdir()}
            repetat = self.ruleaza('--destinatie', destinatie, '--fara-pdf')
            self.assertNotEqual(repetat.returncode, 0)
            self.assertEqual(original, {f.name: f.read_bytes() for f in destinatie.iterdir()})
            a_doua = self.ruleaza('--destinatie', Path(temporar) / 'a doua', '--fara-pdf')
            self.assertEqual(a_doua.returncode, 0, a_doua.stderr)

    def test_pdf_neprodus_refuza_livrarea(self):
        with tempfile.TemporaryDirectory() as temporar:
            rezultat = self.ruleaza('--destinatie', Path(temporar) / 'livrare', '--temporare', temporar,
                '--browser', Path(temporar) / 'browser-inexistent.exe')
            self.assertNotEqual(rezultat.returncode, 0)
            self.assertEqual(rezultat.stdout, b'')


if __name__ == '__main__':
    unittest.main()
