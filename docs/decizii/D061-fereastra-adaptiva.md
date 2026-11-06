# D061 - Fereastra initiala adaptata zonei de lucru

Stare: acceptată. Propusa de Codex, 2026-10-31. Task F16-T24 · Decisă de: Ioan Tuns, 2026-11-06 (F02-T25)

Fereastra pastreaza maximum 1280x800 unitati logice, dar dimensiunea initiala si minima scad pentru a incapea in zona de lucru a monitorului principal, cu o rezerva de 64 unitati pentru margini si bara de titlu. Zona de lucru exclude bara de activitati. Dimensiunile fizice sunt convertite la 96 DPI pentru Wails. Fara biblioteci noi; API Windows prin x/sys existent. Redimensionarea manuala si maximizarea raman disponibile. Nu se modifica rezolutia sau scalarea sistemului. Adaptarea se face la fiecare pornire; mutarea intre monitoare ramane gestionata de Windows/Wails.
