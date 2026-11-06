# D069 - Prioritatea absentelor

Stare: acceptată. Responsabil: Codex. Task: F07-T33 · Decisă de: Ioan Tuns, 2026-11-06 (F02-T25)

Regula functionala a fost confirmata explicit de utilizator: toate concediile numai in zile lucratoare; medical > alte concedii > program; anularea intregii ture de noapte suprapuse si a programelor tuturor marcilor asociate explicit persoanei.

Implementare propusa: tranzactie comuna pentru schita, planificare, absenta si jurnal. Nu duplicam medicalul pe marcile asociate. Nu modificam pontajele efective; conflictele si lunile inchise opresc tranzactia. Nu migram retroactiv absente istorice. Programele anulate nu se restaureaza implicit. Fara reguli noi de plata, dependente sau migrari.
