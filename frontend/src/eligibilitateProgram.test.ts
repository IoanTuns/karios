import {it,expect} from 'vitest';
import {datebaza} from '../wailsjs/go/models';
import {motivProgramExclus, programPlanificat, programSpecial} from './eligibilitateProgram';
it('explica programul nocturn denumit corect dar clasificat diferit',()=>{const p=new datebaza.Schimb({id:'n',activ:true,tip:'TURA',regim:'LUCRATOARE',ora_start:1320});expect(motivProgramExclus(p)).toBe('clasificare');p.tip='NOAPTE';expect(motivProgramExclus(p)).toBe('');expect(motivProgramExclus(p,new datebaza.CategorieProgram({activ:true,schimburi_id:[]}))).toBe('categorie');p.activ=false;expect(motivProgramExclus(p)).toBe('inactiv');});

it('programul normal de zi se poate planifica si ramane precompletabil in Pontaj (F11-T22)',()=>{
 const z7=new datebaza.Schimb({id:'z7',activ:true,tip:'NORMAL',regim:'LUCRATOARE',ora_start:420});
 const g24=new datebaza.Schimb({id:'g',activ:true,tip:'GARDA',regim:'LUCRATOARE',ora_start:420});
 expect(programPlanificat(z7)).toBe(true);expect(programSpecial(z7)).toBe(false);expect(programSpecial(g24)).toBe(true);
});
