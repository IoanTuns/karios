# D065 - Zoom global al interfetei

Stare: propusa. Codex, 2026-10-31. Task F16-T25.

Cerinta utilizatorului este existenta optiunilor zoom-in/zoom-out. Implementarea propune niveluri 75/80/90/100/110/125/150/175/200%, butoane in antet si resetare procentuala, scurtaturi Ctrl+-/Ctrl++/Ctrl+0 si memorare localStorage. Lipsa accesului la stocare nu blocheaza folosirea in sesiune.

Zoomul este aplicat radacinii CSS in WebView2; dimensiunile limitate la viewport sunt ajustate prin factor, inclusiv dialogurile. Meniul Pontaj converteste coordonatele mouse/tastatura, iar panoul AI isi recalculeaza inaltimea. Meniul lateral se poate derula la marire mare. Densitatea grilei ramane independenta. PDF/Excel nu isi modifica dimensiunile. Fara dependente noi sau migrare BD. Testare in Edge existent, RO/EN.
