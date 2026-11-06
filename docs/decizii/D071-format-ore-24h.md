# D071 - Format unitar al orelor

Stare: acceptată. Codex, 2026-10-31. F16-T26 · Decisă de: Ioan Tuns, 2026-11-06 (F02-T25)

Utilizatorul solicita toate orele in format 24h. Formatterul comun impune ciclul h23 (00-23), inclusiv la miezul noptii si in engleza. Afisarile directe dependente de sistem sunt trecute prin formatterul comun. Datele si fusul orar isi pastreaza comportamentul anterior; duratele si totalurile nu sunt ore de ceas si nu se convertesc modulo 24. Rapoartele Go folosesc deja formatul 15:04. Fara biblioteca noua sau schimbari de calcul.

Campul de inceput Schimburi nu mai foloseste controlul nativ time, deoarece WebView2 il poate afisa cu AM/PM dupa setarile sistemului. Introducerea foloseste text HH:MM validat 00:00-23:59, pastrand validarea backend existenta.
