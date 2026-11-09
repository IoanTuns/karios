import type {datebaza} from '../wailsjs/go/models';
// programSpecial: gardele, noptile, programele de zi libera si cele de dupa-amiaza se pun doar din Planificare.
export function programSpecial(p:datebaza.Schimb){return p.tip==='GARDA'||p.tip==='NOAPTE'||p.regim==='NELUCRATOARE'||p.regim==='INTRARE_ZI_LIBERA'||(p.ora_start>=720&&p.ora_start<1320);}
// programPlanificat: ce se poate alege in Planificare; din F11-T21 si programele normale explicite.
export function programPlanificat(p:datebaza.Schimb){return p.tip==='NORMAL'||programSpecial(p);}
export function motivProgramExclus(p:datebaza.Schimb,c?:datebaza.CategorieProgram){if(!p.activ)return 'inactiv';if(!programPlanificat(p))return 'clasificare';if(c&&!c.activ)return 'categorie_inactiva';if(c&&!c.schimburi_id.includes(p.id))return 'categorie';return '';}
