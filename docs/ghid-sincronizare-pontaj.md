# Sincronizarea Pontaj ? Planificare

Disponibil? ?n corectura F07-T21, versiunea 0.88.8, inclus? ?n buildul 0.88.9. Nu este prezent? ?n executabilul 0.88.7. Integrarea ramurii ?n main r?m?ne pas separat.

## Modific?ri noi

1. ?n **Pontaj**, selecteaz? angajatul ?i ziua.
2. Introdu programul special permis de categoria angajatului, de exemplu `S7`.
3. Apas? **Salveaz? schi?a**. Salvarea actualizeaz? ?i pozi?ia din **Planificare**, ?n sec?iunea corespunz?toare. Pentru S7: **Dup?-amiaz?**.
4. Deschide Planificare la aceea?i lun? ?i acela?i angajat.

Sincronizarea are loc la salvare, nu ?n timpul tast?rii. Revenirea la program normal sau la o celul? goal? retrage numai pozi?ia care corespunde programului anterior. Nu se modific? pontajul efectiv.

## Un program S7 salvat anterior

Selecteaz? celula ?n Pontaj. ?n bara de detalii a celulei, l?ng? codul programului, apare **Actualizeaz? planificarea**. Salveaz? ?nt?i celelalte modific?ri din gril?, apoi apas? acest buton. Se folose?te programul deja salvat, f?r? schimbarea orelor sau rescrierea ?ntregii luni. Repetarea nu creeaz? dubluri.

## Conflicte

Lunile ?nchise, absen?ele, suprapunerile, programele nepermise ?i reviziile expirate blocheaz? opera?ia. Dac? exist? o planificare diferit?, modificat? separat, verific? ?nt?i ziua ?n Planificare. La eroare se anuleaz? ?ntregul lot: nu r?m?n modific?ri par?iale ?ntre Pontaj ?i Planificare.

## Verificare ?i predare

Teste noi: normal ? S7 (inclusiv cod cu litere mici), S7 ? G24 ?ntre sec?iuni, revenire normal/gol, reparare f?r? rescriere ?i f?r? dubluri, plan independent, absen?? ?n ziua urm?toare unei g?rzi, suprapunere cu pontaj efectiv, categorie, lun? ?nchis?, revizie ?i eroare SQL cu rollback. Testul API verific? ruta folosit? de interfa?? p?n? la pozi?ia salvat? ?n baza de date.

Nu s-au folosit date opera?ionale. Proba vizual? ?n aplica?ia Windows cu o baz? real? r?m?ne de f?cut dup? integrare. D051 este propus?.
