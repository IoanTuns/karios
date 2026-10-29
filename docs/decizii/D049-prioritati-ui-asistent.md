# D049 · Priorități UI pentru asistentul AI

Stare: propusă · Data: 2026-10-29 · Propusă de: Codex

Task: F14-T23. Bază: [review-ul UI al versiunii 0.82.1](../review-ui-asistent-2026-10-29.md). Solicitantul a cerut review, avantaje, dezavantaje și îmbunătățiri; nu a aprobat prin aceasta implementarea noului lot.

Propunere: primul lot tratează I072–I074 — editarea fără salt automat, o zonă de planificare cu acțiuni vizibile și delimitarea clară a rezultatului de propunere. Urmează I075–I077 pentru excepții, selecții și erori, apoi I078–I080 pentru ciclul planului, lizibilitate și identificarea modelului activ.

Pentru spațiul de lucru, recomand extinderea planului activ pe lățimea disponibilă și o acțiune persistentă, păstrând conversația ca istoric. Alternativa unui panou lateral permanent reduce lățimea pe ecranele modeste; un dialog complet separat implică mai multă navigare. Alegerea finală trebuie verificată pe WebView2, la un viewport CSS de 1024×640 și la scalările Windows folosite de operator.

Nu se simplifică verificările datelor, clarificările lipsă sau confirmarea prin eliminarea lor. Nu se adaugă biblioteci și nu se schimbă reguli de calcul în acest audit.

Decizia și implementarea rămân pentru evaluare umană. Nicio remediere nu este declarată executată.
