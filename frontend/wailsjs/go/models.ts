export namespace datebaza {

	export class Angajat {
	    id: string;
	    nume: string;
	    prenume: string;
	    marca: string;
	    functie: string;
	    sectie_id: string;
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
	        this.tip_contract = source["tip_contract"];
	        this.norma_partiala_minute = source["norma_partiala_minute"];
	        this.norma_minute = source["norma_minute"];
	        this.concediu_zile_an = source["concediu_zile_an"];
	        this.activ = source["activ"];
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
	    norma_intreaga_minute: number;
	    demo_disponibil: boolean;

	    static createFrom(source: any = {}) {
	        return new Date(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.angajati = this.convertValues(source["angajati"], Angajat);
	        this.sectii = this.convertValues(source["sectii"], Sectie);
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

}

