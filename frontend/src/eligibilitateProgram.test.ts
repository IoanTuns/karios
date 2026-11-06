import {it,expect} from 'vitest';
import {datebaza} from '../wailsjs/go/models';
import {motivProgramExclus} from './eligibilitateProgram';
it('explica programul nocturn denumit corect dar clasificat diferit',()=>{const p=new datebaza.Schimb({id:'n',activ:true,tip:'TURA',regim:'LUCRATOARE',ora_start:1320});expect(motivProgramExclus(p)).toBe('clasificare');p.tip='NOAPTE';expect(motivProgramExclus(p)).toBe('');expect(motivProgramExclus(p,new datebaza.CategorieProgram({activ:true,schimburi_id:[]}))).toBe('categorie');p.activ=false;expect(motivProgramExclus(p)).toBe('inactiv');});
