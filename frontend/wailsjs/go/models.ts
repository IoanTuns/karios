export namespace datebaza {

	export class ZiAbsenta {
	    angajat_id: string;
	    data: string;
	    cod: string;
	    activ: boolean;
	    revizie: string;

	    static createFrom(source: any = {}) {
	        return new ZiAbsenta(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.angajat_id = source["angajat_id"];
	        this.data = source["data"];
	        this.cod = source["cod"];
	        this.activ = source["activ"];
	        this.revizie = source["revizie"];
	    }
	}
	export class TipAbsenta {
	    cod: string;
	    denumire: Record<string, string>;
	    categorie: string;
	    activ: boolean;

	    static createFrom(source: any = {}) {
	        return new TipAbsenta(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cod = source["cod"];
	        this.denumire = source["denumire"];
	        this.categorie = source["categorie"];
	        this.activ = source["activ"];
	    }
	}
	export class AbsenteAn {
	    an: number;
	    tipuri: TipAbsenta[];
	    zile: ZiAbsenta[];

	    static createFrom(source: any = {}) {
	        return new AbsenteAn(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.an = source["an"];
	        this.tipuri = this.convertValues(source["tipuri"], TipAbsenta);
	        this.zile = this.convertValues(source["zile"], ZiAbsenta);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Angajat {
	    id: string;
	    nume: string;
	    prenume: string;
	    marca: string;
	    functie: string;
	    sectie_id: string;
	    departament_id: string;
	    tip_contract: string;
	    norma_partiala_minute?: number;
	    norma_minute: number;
	    concediu_zile_an?: number;
	    activ: boolean;

	    static createFrom(source: any = {}) {
	        return new Angajat(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.nume = source["nume"];
	        this.prenume = source["prenume"];
	        this.marca = source["marca"];
	        this.functie = source["functie"];
	        this.sectie_id = source["sectie_id"];
	        this.departament_id = source["departament_id"];
	        this.tip_contract = source["tip_contract"];
	        this.norma_partiala_minute = source["norma_partiala_minute"];
	        this.norma_minute = source["norma_minute"];
	        this.concediu_zile_an = source["concediu_zile_an"];
	        this.activ = source["activ"];
	    }
	}
	export class Sarbatoare {
	    data: string;
	    denumire: Record<string, string>;

	    static createFrom(source: any = {}) {
	        return new Sarbatoare(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = source["data"];
	        this.denumire = source["denumire"];
	    }
	}
	export class CalendarAn {
	    an: number;
	    sarbatori: Sarbatoare[];
	    import_disponibil: boolean;

	    static createFrom(source: any = {}) {
	        return new CalendarAn(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.an = source["an"];
	        this.sarbatori = this.convertValues(source["sarbatori"], Sarbatoare);
	        this.import_disponibil = source["import_disponibil"];
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CelulaSchita {
	    angajat_id: string;
	    data: string;
	    valoare: string;
	    revizie: string;

	    static createFrom(source: any = {}) {
	        return new CelulaSchita(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.angajat_id = source["angajat_id"];
	        this.data = source["data"];
	        this.valoare = source["valoare"];
	        this.revizie = source["revizie"];
	    }
	}
	export class Departament {
	    id: string;
	    sectie_id: string;
	    denumire: Record<string, string>;

	    static createFrom(source: any = {}) {
	        return new Departament(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.sectie_id = source["sectie_id"];
	        this.denumire = source["denumire"];
	    }
	}
	export class Sectie {
	    id: string;
	    denumire: Record<string, string>;

	    static createFrom(source: any = {}) {
	        return new Sectie(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.denumire = source["denumire"];
	    }
	}
	export class Date {
	    angajati: Angajat[];
	    sectii: Sectie[];
	    departamente: Departament[];
	    norma_intreaga_minute: number;
	    demo_disponibil: boolean;

	    static createFrom(source: any = {}) {
	        return new Date(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.angajati = this.convertValues(source["angajati"], Angajat);
	        this.sectii = this.convertValues(source["sectii"], Sectie);
	        this.departamente = this.convertValues(source["departamente"], Departament);
	        this.norma_intreaga_minute = source["norma_intreaga_minute"];
	        this.demo_disponibil = source["demo_disponibil"];
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}


	export class Schimb {
	    id: string;
	    cod: string;
	    denumire: Record<string, string>;
	    ora_start: number;
	    durata_minute: number;
	    pauza_minute: number;
	    tip: string;
	    activ: boolean;

	    static createFrom(source: any = {}) {
	        return new Schimb(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.cod = source["cod"];
	        this.denumire = source["denumire"];
	        this.ora_start = source["ora_start"];
	        this.durata_minute = source["durata_minute"];
	        this.pauza_minute = source["pauza_minute"];
	        this.tip = source["tip"];
	        this.activ = source["activ"];
	    }
	}
	export class SchitaLuna {
	    an: number;
	    luna: number;
	    celule: CelulaSchita[];
	    absente: ZiAbsenta[];

	    static createFrom(source: any = {}) {
	        return new SchitaLuna(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.an = source["an"];
	        this.luna = source["luna"];
	        this.celule = this.convertValues(source["celule"], CelulaSchita);
	        this.absente = this.convertValues(source["absente"], ZiAbsenta);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}



}

export namespace main {

	export class Informatii {
	    versiune: string;
	    commit: string;

	    static createFrom(source: any = {}) {
	        return new Informatii(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.versiune = source["versiune"];
	        this.commit = source["commit"];
	    }
	}
	export class MesajEroare {
	    cod: string;
	    param: Record<string, string>;

	    static createFrom(source: any = {}) {
	        return new MesajEroare(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cod = source["cod"];
	        this.param = source["param"];
	    }
	}
	export class RaspunsAbsente {
	    absente?: datebaza.AbsenteAn;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsAbsente(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.absente = this.convertValues(source["absente"], datebaza.AbsenteAn);
	        this.eroare = this.convertValues(source["eroare"], MesajEroare);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RaspunsCalendar {
	    calendar?: datebaza.CalendarAn;
	    adaugate: number;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsCalendar(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.calendar = this.convertValues(source["calendar"], datebaza.CalendarAn);
	        this.adaugate = source["adaugate"];
	        this.eroare = this.convertValues(source["eroare"], MesajEroare);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RaspunsDate {
	    date?: datebaza.Date;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsDate(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = this.convertValues(source["date"], datebaza.Date);
	        this.eroare = this.convertValues(source["eroare"], MesajEroare);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RaspunsSchimburi {
	    adaugate: number;
	    schimburi: datebaza.Schimb[];
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsSchimburi(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.adaugate = source["adaugate"];
	        this.schimburi = this.convertValues(source["schimburi"], datebaza.Schimb);
	        this.eroare = this.convertValues(source["eroare"], MesajEroare);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RaspunsSchita {
	    schita?: datebaza.SchitaLuna;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsSchita(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schita = this.convertValues(source["schita"], datebaza.SchitaLuna);
	        this.eroare = this.convertValues(source["eroare"], MesajEroare);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

