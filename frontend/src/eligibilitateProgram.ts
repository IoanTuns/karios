import type {datebaza} from '../wailsjs/go/models';
export function programPlanificat(p:datebaza.Schimb){return p.tip==='NORMAL'||p.tip==='GARDA'||p.tip==='NOAPTE'||p.regim==='NELUCRATOARE'||p.regim==='INTRARE_ZI_LIBERA'||(p.ora_start>=720&&p.ora_start<1320);}
export function motivProgramExclus(p:datebaza.Schimb,c?:datebaza.CategorieProgram){if(!p.activ)return 'inactiv';if(!programPlanificat(p))return 'clasificare';if(c&&!c.activ)return 'categorie_inactiva';if(c&&!c.schimburi_id.includes(p.id))return 'categorie';return '';}
