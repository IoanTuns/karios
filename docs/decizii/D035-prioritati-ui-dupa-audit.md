# D035 · Priorități UI după auditul versiunii 0.27.0

Stare: acceptată · Data: 2026-10-26 · Propusă de: Codex · Decisă de: Ioan Tuns · Taskuri: F05-T18, F05-T19–F05-T23

Auditul cerut de Ioan Tuns este în [analiza UI](../analiza-ui-2026-10-26.md). Versiunea curentă include deja ergonomia D032–D034. Problemele confirmate sunt spațiul insuficient pentru grilă la 1024×640, resetarea contextului la navigare, texte rămase din versiuni vechi, eticheta ambiguă a comparației T/Ore norm, contrast redus în unele stări și acțiuni de salvare departe de zona editată.

Se propune păstrarea structurii și identității actuale și următoarea ordine: claritatea indicatorilor și lizibilitatea, ergonomia pontajului și continuitatea navigării, acțiuni persistente și explicații pentru erori, apoi gruparea cererilor și simplificarea configurării/editorului. Idei I045–I052; criteriile de acceptare sunt în audit.

Nu se schimbă sensul regulilor de calcul. Ore norm rămâne clasificarea orelor lucrate în zile normale confirmată de utilizator; o eventuală normă contractuală lunară este un concept separat și depinde de F01. Nu se prezintă diferența T–Ore norm ca ore suplimentare sau deficit contractual.

Auditul este documentație. Nu introduce implementări, biblioteci, migrări ori versiune nouă. Alegerea lotului de implementare rămâne la utilizator; propunerea nu înlocuiește automat deciziile existente.

Acceptare explicită la 2026-10-26: „OK, Creaza o lista de lucru si incepe imbunatatirile”. Implementarea este împărțită în [loturi UI](../lista-lucru-ui.md); primul lot este F05-T19. Regulile contabile neprecizate rămân neconfigurate.
