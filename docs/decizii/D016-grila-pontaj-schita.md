# D016 · Prima grilă de pontaj și schițele persistente

Stare: propusă · Data: 2026-10-25 · Taskuri: F08-T01

## Context

Ioan Tuns a cerut începerea grilei și a furnizat un model de centralizator. Codurile de absență, semnificația coloanelor și legătura rândurilor de gardă cu persoanele/contractele nu sunt specificate încă (F01). Introducerea trebuie să poată fi încercată fără inventarea unor calcule.

## Propunere

- Prototip React cu tabel HTML nativ, fără bibliotecă nouă. Verificăm 50 de rânduri × 31 de zile înainte de o eventuală bibliotecă specializată. Antet și identificarea angajatului rămân vizibile la derulare.
- An și lună, filtru pe secție, weekenduri gri și sărbători marcate din calendarul configurat. Coloanele de total din model sunt afișate ca necalculate.
- Un rând per înregistrare de angajat; nu deducem că un rând GARDĂ aparține altui rând după nume. Gruparea profesională și relația persoană-contract sunt I013/F01.
- Tastatură (săgeți/Enter/Tab), selecție cu Shift, copiere TSV, lipire din Excel și completarea unui dreptunghi. Sugestii din nomenclatorul de schimburi active; textele necunoscute sunt permise în schiță, fără interpretare automată.
- Migrare `0005_pontaj_ciorne.sql`: valori text de maximum 16 caractere, per angajat și dată, păstrate separat de intervalele efective din `pontaj`. O celulă goală rămâne înregistrare actualizată, fără ștergere definitivă.
- Salvare explicită și tranzacțională a celulelor modificate, cu jurnal `local`. Revizia fiecărei celule detectează modificări concurente și refuză suprascrierea unei versiuni mai noi; întregul lot este anulat la conflict.
- Avertizare la schimbarea lunii sau navigare cu modificări nesalvate. Ecranul este etichetat ca schiță; nu pretinde că produce totaluri de salarizare sau raportul final.

## Variante evaluate

- Salvarea codurilor fără intervale în `pontaj`: ar încălca modelul actual și ar inventa semnificații pentru absențe.
- Păstrarea doar în memoria React: ar pierde munca la repornire.
- O bibliotecă de grilă nouă înaintea probei: amânată până la verificarea cerințelor pe prototipul nativ; alternativa rămâne deschisă dacă performanța sau accesibilitatea sunt insuficiente.
