// Propunerile repetate ale asistentului (F14-T09, F14-T19): se omit doar duplicatele unor
// propuneri încă neconfirmate; cele identice cu una deja aplicată se arată, marcate.
import type {ai} from '../wailsjs/go/models';

interface MesajCuPropuneri {
    propuneri?: ai.PropunereActiune[];
    propuneriAplicata?: Record<string, boolean>;
}

// cheiePropunere identifică o propunere după conținut, nu după id-ul generat la fiecare răspuns.
export function cheiePropunere(p: ai.PropunereActiune): string {
    const parametri = p.parametri ?? {};
    return p.tip + JSON.stringify(Object.keys(parametri).sort().map(k => [k, parametri[k]]));
}

export function imparteRepetari(propuneri: ai.PropunereActiune[] | undefined, conversatie: MesajCuPropuneri[]):
    {noi: ai.PropunereActiune[]; omise: number; reaplicari: Record<string, boolean>} {
    const neconfirmate = new Set<string>(), aplicate = new Set<string>();
    for (const m of conversatie) for (const p of m.propuneri ?? []) {
        (m.propuneriAplicata?.[p.id] ? aplicate : neconfirmate).add(cheiePropunere(p));
    }
    const reaplicari: Record<string, boolean> = {};
    const noi = (propuneri ?? []).filter(p => {
        const k = cheiePropunere(p);
        if (neconfirmate.has(k) && !aplicate.has(k)) return false;
        if (aplicate.has(k)) reaplicari[p.id] = true;
        return true;
    });
    return {noi, omise: (propuneri?.length ?? 0) - noi.length, reaplicari};
}
