# D086 - Programe multiple pe aceeasi marca si zi

Stare: propusa. Responsabil Codex. Data2026-11-04. Task F11-T21.

Cerinta utilizatorului: intervale distincte pe aceeasi marca/zi, adaugare fara inlocuire, fara suprapuneri. Planificarea ramane normalizata pe pozitii cu ID/revizie. Pontajul zilnic pastreaza o proiectie canonica a codurilor separate prin +, separator interzis in codurile individuale existente. Calculul expandeaza proiectia in intervalele originale; nu transforma totalul in ore manuale. Limita proiectiei1024 caractere, validare pe server, intervale si categorii verificate atomic. Migrare noua elimina unicitatea pe zi a Planificarii si extinde spatiul pentru proiectia Pontajului, pastrand toate valorile/reviziile existente. Publicarea oricarei pozitii dintr-o zi publica intregul set activ al zilei, cu token pentru toate reviziile; nu lasa pozitiile vecine sa fie pierdute la publicare partiala. Nu se adauga reguli salariale sau dependente.
