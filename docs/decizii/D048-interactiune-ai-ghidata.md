# D048 · Interpretare editabilă și plan verificat pentru asistent

Stare: propusă · Data: 2026-10-29 · Propusă de: Codex

Task: F14-T22, cerut de Ioan Tuns, cu implementarea tuturor celor șase propuneri discutate.

Modelul și recunoașterea deterministă produc o intenție editabilă, nu o modificare. Interfața cere explicit filtrul zilelor și comportamentul pentru celulele ocupate dacă acestea lipsesc. Grupurile sunt rezolvate din datele curente: toți activii, secție, departament, categorie sau lista persoanelor; filtrul fără pontaj se aplică pe lunile vizate.

Previzualizarea este calculată în Go, cu lista persoanelor, luni, operații, valori vechi/noi și excluderi. Confirmarea se referă la amprenta previzualizării. La aplicare se reconstruiește planul și se refuză orice diferență; salvarea mai multor luni este tranzacțională. Absențele, lunile închise și programele nepermise rămân protejate. Rezultatul este verificat după scriere și numărat, nu afirmat de model.

Contextele pentru model sunt specifice cererii (planificare, absențe, pontaj, ajutor), cu limite explicite și date relevante. Analizele deterministe folosesc datele complete. Reducerea contextului nu trebuie să autorizeze selecții incomplete: grupurile se rezolvă din baza de date la previzualizare.

Alternativa de a schimba doar modelul sau promptul nu garantează selecția completă și confirmarea exactă. Nu adăugăm biblioteci. Testele de regresie folosesc date anonimizate și comenzi RO/EN. Performanța și precizia modelului real se raportează separat de testele deterministe.

Decizia rămâne propusă până la acceptarea explicită de către un om.
