export namespace ai {

	export class MesajConversatie {
	    role: string;
	    content: string;

	    static createFrom(source: any = {}) {
	        return new MesajConversatie(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.content = source["content"];
	    }
	}
	export class PropunereActiune {
	    id: string;
	    tip: string;
	    descriere: string;
	    parametri: Record<string, any>;
	    necesita_confirmare: boolean;

	    static createFrom(source: any = {}) {
	        return new PropunereActiune(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.tip = source["tip"];
	        this.descriere = source["descriere"];
	        this.parametri = source["parametri"];
	        this.necesita_confirmare = source["necesita_confirmare"];
	    }
	}
	export class RaspunsAsistent {
	    mesaj: string;
	    propuneri: PropunereActiune[];
	    eroare?: string;

	    static createFrom(source: any = {}) {
	        return new RaspunsAsistent(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mesaj = source["mesaj"];
	        this.propuneri = this.convertValues(source["propuneri"], PropunereActiune);
	        this.eroare = source["eroare"];
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
	export class StareMediu {
	    instalat: boolean;
	    cale_folder: string;
	    cale_binar: string;
	    cale_model: string;
	    modele_disponibile: string[];
	    activ: boolean;
	    port: number;
	    mesaj: string;

	    static createFrom(source: any = {}) {
	        return new StareMediu(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.instalat = source["instalat"];
	        this.cale_folder = source["cale_folder"];
	        this.cale_binar = source["cale_binar"];
	        this.cale_model = source["cale_model"];
	        this.modele_disponibile = source["modele_disponibile"];
	        this.activ = source["activ"];
	        this.port = source["port"];
	        this.mesaj = source["mesaj"];
	    }
	}

}

export namespace audit {

	export class Rezultat {
	    stare: string;
	    tabele: string[];
	    rand_rupt: number;
	    verificat_la: string;
	    durata_ms: number;

	    static createFrom(source: any = {}) {
	        return new Rezultat(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stare = source["stare"];
	        this.tabele = source["tabele"];
	        this.rand_rupt = source["rand_rupt"];
	        this.verificat_la = source["verificat_la"];
	        this.durata_ms = source["durata_ms"];
	    }
	}

}

export namespace calcul {

	export class AvertismentEchitate {
	    angajat_id: string;
	    data: string;
	    cod_schimb: string;
	    tip: string;
	    mesaj: string;

	    static createFrom(source: any = {}) {
	        return new AvertismentEchitate(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.angajat_id = source["angajat_id"];
	        this.data = source["data"];
	        this.cod_schimb = source["cod_schimb"];
	        this.tip = source["tip"];
	        this.mesaj = source["mesaj"];
	    }
	}
	export class CelulaPrecompletare {
	    angajat_id: string;
	    data: string;
	    valoare: string;

	    static createFrom(source: any = {}) {
	        return new CelulaPrecompletare(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.angajat_id = source["angajat_id"];
	        this.data = source["data"];
	        this.valoare = source["valoare"];
	    }
	}
	export class RegulaColoana {
	    coloana: string;
	    sursa: string;
	    coduri: string[];
	    factor_sutimi: number;
	    numai_aprobate: boolean;

	    static createFrom(source: any = {}) {
	        return new RegulaColoana(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.coloana = source["coloana"];
	        this.sursa = source["sursa"];
	        this.coduri = source["coduri"];
	        this.factor_sutimi = source["factor_sutimi"];
	        this.numai_aprobate = source["numai_aprobate"];
	    }
	}
	export class ConfiguratiePontaj {
	    atribuire: string;
	    reguli: RegulaColoana[];
	    revizie: string;

	    static createFrom(source: any = {}) {
	        return new ConfiguratiePontaj(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.atribuire = source["atribuire"];
	        this.reguli = this.convertValues(source["reguli"], RegulaColoana);
	        this.revizie = source["revizie"];
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
	export class IncalcareRegula {
	    regula_id: string;
	    regula_cod: string;
	    severitate: string;
	    angajat_id: string;
	    nume: string;
	    data: string;
	    mesaj: string;

	    static createFrom(source: any = {}) {
	        return new IncalcareRegula(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.regula_id = source["regula_id"];
	        this.regula_cod = source["regula_cod"];
	        this.severitate = source["severitate"];
	        this.angajat_id = source["angajat_id"];
	        this.nume = source["nume"];
	        this.data = source["data"];
	        this.mesaj = source["mesaj"];
	    }
	}
	export class OreZi {
	    angajat_id: string;
	    data: string;
	    minute: number;
	    calculabil: boolean;
	    surse: string[];

	    static createFrom(source: any = {}) {
	        return new OreZi(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.angajat_id = source["angajat_id"];
	        this.data = source["data"];
	        this.minute = source["minute"];
	        this.calculabil = source["calculabil"];
	        this.surse = source["surse"];
	    }
	}
	export class PlanPrecompletare {
	    celule: CelulaPrecompletare[];
	    omise: Record<string, number>;

	    static createFrom(source: any = {}) {
	        return new PlanPrecompletare(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.celule = this.convertValues(source["celule"], CelulaPrecompletare);
	        this.omise = source["omise"];
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
	export class ZiComparativ {
	    data: string;
	    zi: number;
	    plan_valoare: string;
	    plan_ore: number;
	    real_valoare: string;
	    real_ore: number;
	    stare: string;

	    static createFrom(source: any = {}) {
	        return new ZiComparativ(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = source["data"];
	        this.zi = source["zi"];
	        this.plan_valoare = source["plan_valoare"];
	        this.plan_ore = source["plan_ore"];
	        this.real_valoare = source["real_valoare"];
	        this.real_ore = source["real_ore"];
	        this.stare = source["stare"];
	    }
	}
	export class RandComparativAngajat {
	    angajat_id: string;
	    nume: string;
	    prenume: string;
	    marca: string;
	    functie: string;
	    total_ore_planificate: number;
	    total_ore_realizate: number;
	    diferenta_ore: number;
	    zile_modificate: number;
	    zile: ZiComparativ[];

	    static createFrom(source: any = {}) {
	        return new RandComparativAngajat(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.angajat_id = source["angajat_id"];
	        this.nume = source["nume"];
	        this.prenume = source["prenume"];
	        this.marca = source["marca"];
	        this.functie = source["functie"];
	        this.total_ore_planificate = source["total_ore_planificate"];
	        this.total_ore_realizate = source["total_ore_realizate"];
	        this.diferenta_ore = source["diferenta_ore"];
	        this.zile_modificate = source["zile_modificate"];
	        this.zile = this.convertValues(source["zile"], ZiComparativ);
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
	export class ZiAcoperire {
	    data: string;
	    ziua_saptamanii: number;
	    este_weekend: boolean;
	    este_sarbatoare: boolean;
	    denumire_sarbatoare: string;
	    posturi: Record<string, Array<string>>;
	    lipsuri: string[];
	    complet: boolean;

	    static createFrom(source: any = {}) {
	        return new ZiAcoperire(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = source["data"];
	        this.ziua_saptamanii = source["ziua_saptamanii"];
	        this.este_weekend = source["este_weekend"];
	        this.este_sarbatoare = source["este_sarbatoare"];
	        this.denumire_sarbatoare = source["denumire_sarbatoare"];
	        this.posturi = source["posturi"];
	        this.lipsuri = source["lipsuri"];
	        this.complet = source["complet"];
	    }
	}
	export class RaportAcoperire {
	    an: number;
	    luna: number;
	    total_zile_nelucratoare: number;
	    total_posturi_necesar: number;
	    total_posturi_acoperite: number;
	    total_posturi_lipsa: number;
	    zile: ZiAcoperire[];

	    static createFrom(source: any = {}) {
	        return new RaportAcoperire(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.an = source["an"];
	        this.luna = source["luna"];
	        this.total_zile_nelucratoare = source["total_zile_nelucratoare"];
	        this.total_posturi_necesar = source["total_posturi_necesar"];
	        this.total_posturi_acoperite = source["total_posturi_acoperite"];
	        this.total_posturi_lipsa = source["total_posturi_lipsa"];
	        this.zile = this.convertValues(source["zile"], ZiAcoperire);
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
	export class RaportPlanificatVsRealizat {
	    an: number;
	    luna: number;
	    total_planificat_ore: number;
	    total_realizat_ore: number;
	    diferenta_globala_ore: number;
	    total_discrepante: number;
	    randuri: RandComparativAngajat[];

	    static createFrom(source: any = {}) {
	        return new RaportPlanificatVsRealizat(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.an = source["an"];
	        this.luna = source["luna"];
	        this.total_planificat_ore = source["total_planificat_ore"];
	        this.total_realizat_ore = source["total_realizat_ore"];
	        this.diferenta_globala_ore = source["diferenta_globala_ore"];
	        this.total_discrepante = source["total_discrepante"];
	        this.randuri = this.convertValues(source["randuri"], RandComparativAngajat);
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

	export class RegulaConfigurabila {
	    id: string;
	    cod: string;
	    denumire: string;
	    descriere: string;
	    tip_regula: string;
	    severitate: string;
	    coduri_schimb: string;
	    valoare_numerica: number;
	    activ: boolean;

	    static createFrom(source: any = {}) {
	        return new RegulaConfigurabila(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.cod = source["cod"];
	        this.denumire = source["denumire"];
	        this.descriere = source["descriere"];
	        this.tip_regula = source["tip_regula"];
	        this.severitate = source["severitate"];
	        this.coduri_schimb = source["coduri_schimb"];
	        this.valoare_numerica = source["valoare_numerica"];
	        this.activ = source["activ"];
	    }
	}
	export class ValoarePontaj {
	    sutimi: number;
	    unitate: string;
	    definita: boolean;

	    static createFrom(source: any = {}) {
	        return new ValoarePontaj(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sutimi = source["sutimi"];
	        this.unitate = source["unitate"];
	        this.definita = source["definita"];
	    }
	}
	export class TotalPontaj {
	    angajat_id: string;
	    valori: Record<string, ValoarePontaj>;

	    static createFrom(source: any = {}) {
	        return new TotalPontaj(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.angajat_id = source["angajat_id"];
	        this.valori = this.convertValues(source["valori"], ValoarePontaj, true);
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
	export class RezultatPontaj {
	    celule: OreZi[];
	    totaluri: TotalPontaj[];

	    static createFrom(source: any = {}) {
	        return new RezultatPontaj(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.celule = this.convertValues(source["celule"], OreZi);
	        this.totaluri = this.convertValues(source["totaluri"], TotalPontaj);
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

export namespace datebaza {

	export class FiltruJurnal {
	    tabel: string;
	    angajat_id: string;
	    luna: string;
	    text: string;
	    limita: number;
	    deplasare: number;

	    static createFrom(source: any = {}) {
	        return new FiltruJurnal(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tabel = source["tabel"];
	        this.angajat_id = source["angajat_id"];
	        this.luna = source["luna"];
	        this.text = source["text"];
	        this.limita = source["limita"];
	        this.deplasare = source["deplasare"];
	    }
	}
	export class IntrareJurnal {
	    id: number;
	    tabel: string;
	    rand_id: string;
	    camp: string;
	    vechi: string;
	    nou: string;
	    cand: string;
	    cine: string;

	    static createFrom(source: any = {}) {
	        return new IntrareJurnal(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.tabel = source["tabel"];
	        this.rand_id = source["rand_id"];
	        this.camp = source["camp"];
	        this.vechi = source["vechi"];
	        this.nou = source["nou"];
	        this.cand = source["cand"];
	        this.cine = source["cine"];
	    }
	}
	export class PaginaJurnal {
	    intrari: IntrareJurnal[];
	    total: number;
	    tabele: string[];

	    static createFrom(source: any = {}) {
	        return new PaginaJurnal(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.intrari = this.convertValues(source["intrari"], IntrareJurnal);
	        this.total = source["total"];
	        this.tabele = source["tabele"];
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
	export class ZiAbsenta {
	    angajat_id: string;
	    data: string;
	    cod: string;
	    activ: boolean;
	    revizie: string;
	    nr_inregistrare: string;
	    inlocuitor_id: string;
	    stare_aprobare: string;
	    aprobat_de: string;
	    document_id: string;
	    document_nume: string;
	    document_mime: string;
	    document_continut_b64?: string;

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
	        this.nr_inregistrare = source["nr_inregistrare"];
	        this.inlocuitor_id = source["inlocuitor_id"];
	        this.stare_aprobare = source["stare_aprobare"];
	        this.aprobat_de = source["aprobat_de"];
	        this.document_id = source["document_id"];
	        this.document_nume = source["document_nume"];
	        this.document_mime = source["document_mime"];
	        this.document_continut_b64 = source["document_continut_b64"];
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
	    categorie_program_id: string;
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
	        this.categorie_program_id = source["categorie_program_id"];
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
	export class CategorieProgram {
	    id: string;
	    denumire: Record<string, string>;
	    schimburi_id: string[];
	    schimb_implicit_id: string;
	    activ: boolean;
	    revizie: string;

	    static createFrom(source: any = {}) {
	        return new CategorieProgram(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.denumire = source["denumire"];
	        this.schimburi_id = source["schimburi_id"];
	        this.schimb_implicit_id = source["schimb_implicit_id"];
	        this.activ = source["activ"];
	        this.revizie = source["revizie"];
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
	export class RandImportAngajat {
	    nume: string;
	    prenume: string;
	    marca: string;
	    functie: string;
	    sectie: string;
	    departament: string;
	    tip_contract: string;
	    norma_minute: number;
	    concediu_zile_an?: number;
	    valid: boolean;
	    mesaj_eroare: string;
	    exista_deja: boolean;

	    static createFrom(source: any = {}) {
	        return new RandImportAngajat(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nume = source["nume"];
	        this.prenume = source["prenume"];
	        this.marca = source["marca"];
	        this.functie = source["functie"];
	        this.sectie = source["sectie"];
	        this.departament = source["departament"];
	        this.tip_contract = source["tip_contract"];
	        this.norma_minute = source["norma_minute"];
	        this.concediu_zile_an = source["concediu_zile_an"];
	        this.valid = source["valid"];
	        this.mesaj_eroare = source["mesaj_eroare"];
	        this.exista_deja = source["exista_deja"];
	    }
	}
	export class CerereExecutaImportAngajati {
	    randuri: RandImportAngajat[];
	    actualizeaza_existente: boolean;

	    static createFrom(source: any = {}) {
	        return new CerereExecutaImportAngajati(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.randuri = this.convertValues(source["randuri"], RandImportAngajat);
	        this.actualizeaza_existente = source["actualizeaza_existente"];
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
	export class RandImportSarbatoare {
	    data: string;
	    denumire_ro: string;
	    denumire_en: string;
	    valid: boolean;
	    mesaj_eroare: string;
	    exista_deja: boolean;

	    static createFrom(source: any = {}) {
	        return new RandImportSarbatoare(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = source["data"];
	        this.denumire_ro = source["denumire_ro"];
	        this.denumire_en = source["denumire_en"];
	        this.valid = source["valid"];
	        this.mesaj_eroare = source["mesaj_eroare"];
	        this.exista_deja = source["exista_deja"];
	    }
	}
	export class CerereExecutaImportSarbatori {
	    randuri: RandImportSarbatoare[];
	    inlocuieste_existente: boolean;

	    static createFrom(source: any = {}) {
	        return new CerereExecutaImportSarbatori(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.randuri = this.convertValues(source["randuri"], RandImportSarbatoare);
	        this.inlocuieste_existente = source["inlocuieste_existente"];
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
	export class CererePublicarePlan {
	    an: number;
	    luna: number;
	    sectiune: string;
	    pozitii_id: string[];
	    inlocuieste: boolean;
	    token: string;

	    static createFrom(source: any = {}) {
	        return new CererePublicarePlan(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.an = source["an"];
	        this.luna = source["luna"];
	        this.sectiune = source["sectiune"];
	        this.pozitii_id = source["pozitii_id"];
	        this.inlocuieste = source["inlocuieste"];
	        this.token = source["token"];
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
	    categorii_program: CategorieProgram[];
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
	        this.categorii_program = this.convertValues(source["categorii_program"], CategorieProgram);
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

	export class DocumentAbsenta {
	    id: string;
	    nume_fisier: string;
	    tip_mime: string;
	    continut_b64: string;

	    static createFrom(source: any = {}) {
	        return new DocumentAbsenta(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.nume_fisier = source["nume_fisier"];
	        this.tip_mime = source["tip_mime"];
	        this.continut_b64 = source["continut_b64"];
	    }
	}
	export class OptiuneRegula {
	    cod: string;
	    denumire: Record<string, string>;
	    tip: string;

	    static createFrom(source: any = {}) {
	        return new OptiuneRegula(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cod = source["cod"];
	        this.denumire = source["denumire"];
	        this.tip = source["tip"];
	    }
	}
	export class PozitiePlanificare {
	    id: string;
	    angajat_id: string;
	    data: string;
	    schimb_id: string;
	    sectiune: string;
	    activ: boolean;
	    revizie: string;
	    stare: string;
	    publicat_la: string;

	    static createFrom(source: any = {}) {
	        return new PozitiePlanificare(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.angajat_id = source["angajat_id"];
	        this.data = source["data"];
	        this.schimb_id = source["schimb_id"];
	        this.sectiune = source["sectiune"];
	        this.activ = source["activ"];
	        this.revizie = source["revizie"];
	        this.stare = source["stare"];
	        this.publicat_la = source["publicat_la"];
	    }
	}
	export class PlanificareLuna {
	    an: number;
	    luna: number;
	    stare: string;
	    pozitii: PozitiePlanificare[];

	    static createFrom(source: any = {}) {
	        return new PlanificareLuna(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.an = source["an"];
	        this.luna = source["luna"];
	        this.stare = source["stare"];
	        this.pozitii = this.convertValues(source["pozitii"], PozitiePlanificare);
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

	export class RandPublicarePlan {
	    pozitie_id: string;
	    angajat_id: string;
	    data: string;
	    schimb_id: string;
	    valoare: string;
	    inainte: string;
	    revizie_plan: string;
	    revizie_pontaj: string;
	    ora_start: number;
	    durata_minute: number;
	    stare: string;
	    motiv: string;

	    static createFrom(source: any = {}) {
	        return new RandPublicarePlan(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pozitie_id = source["pozitie_id"];
	        this.angajat_id = source["angajat_id"];
	        this.data = source["data"];
	        this.schimb_id = source["schimb_id"];
	        this.valoare = source["valoare"];
	        this.inainte = source["inainte"];
	        this.revizie_plan = source["revizie_plan"];
	        this.revizie_pontaj = source["revizie_pontaj"];
	        this.ora_start = source["ora_start"];
	        this.durata_minute = source["durata_minute"];
	        this.stare = source["stare"];
	        this.motiv = source["motiv"];
	    }
	}
	export class PublicarePlan {
	    randuri: RandPublicarePlan[];
	    token: string;
	    publicabile: number;
	    identice: number;
	    blocate: number;

	    static createFrom(source: any = {}) {
	        return new PublicarePlan(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.randuri = this.convertValues(source["randuri"], RandPublicarePlan);
	        this.token = source["token"];
	        this.publicabile = source["publicabile"];
	        this.identice = source["identice"];
	        this.blocate = source["blocate"];
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



	export class RezultatAnalizaImportAngajati {
	    nume_fisier: string;
	    format: string;
	    randuri: RandImportAngajat[];
	    total_randuri: number;
	    total_valide: number;
	    total_invalide: number;
	    total_existente: number;

	    static createFrom(source: any = {}) {
	        return new RezultatAnalizaImportAngajati(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nume_fisier = source["nume_fisier"];
	        this.format = source["format"];
	        this.randuri = this.convertValues(source["randuri"], RandImportAngajat);
	        this.total_randuri = source["total_randuri"];
	        this.total_valide = source["total_valide"];
	        this.total_invalide = source["total_invalide"];
	        this.total_existente = source["total_existente"];
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
	export class RezultatAnalizaImportSarbatori {
	    nume_fisier: string;
	    format: string;
	    randuri: RandImportSarbatoare[];
	    total_randuri: number;
	    total_valide: number;
	    total_invalide: number;
	    total_existente: number;

	    static createFrom(source: any = {}) {
	        return new RezultatAnalizaImportSarbatori(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nume_fisier = source["nume_fisier"];
	        this.format = source["format"];
	        this.randuri = this.convertValues(source["randuri"], RandImportSarbatoare);
	        this.total_randuri = source["total_randuri"];
	        this.total_valide = source["total_valide"];
	        this.total_invalide = source["total_invalide"];
	        this.total_existente = source["total_existente"];
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
	    regim: string;
	    categorie: string;
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
	        this.regim = source["regim"];
	        this.categorie = source["categorie"];
	        this.activ = source["activ"];
	    }
	}
	export class SchitaLuna {
	    an: number;
	    luna: number;
	    stare: string;
	    celule: CelulaSchita[];
	    absente: ZiAbsenta[];
	    calcul: calcul.RezultatPontaj;
	    configuratie: calcul.ConfiguratiePontaj;

	    static createFrom(source: any = {}) {
	        return new SchitaLuna(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.an = source["an"];
	        this.luna = source["luna"];
	        this.stare = source["stare"];
	        this.celule = this.convertValues(source["celule"], CelulaSchita);
	        this.absente = this.convertValues(source["absente"], ZiAbsenta);
	        this.calcul = this.convertValues(source["calcul"], calcul.RezultatPontaj);
	        this.configuratie = this.convertValues(source["configuratie"], calcul.ConfiguratiePontaj);
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

	export class StareLuna {
	    an: number;
	    luna: number;
	    stare: string;
	    modificat_la: string;
	    modificat_de: string;
	    motiv: string;
	    amprenta: string;
	    amprenta_corespunde: boolean;

	    static createFrom(source: any = {}) {
	        return new StareLuna(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.an = source["an"];
	        this.luna = source["luna"];
	        this.stare = source["stare"];
	        this.modificat_la = source["modificat_la"];
	        this.modificat_de = source["modificat_de"];
	        this.motiv = source["motiv"];
	        this.amprenta = source["amprenta"];
	        this.amprenta_corespunde = source["amprenta_corespunde"];
	    }
	}


}

export namespace diagnoza {

	export class StatisticiDiagnoza {
	    numar_erori: number;
	    ultima_eroare: string;
	    total_interogari: number;
	    total_lente: number;
	    lenta_interogare: string;
	    timp_maxim_ms: number;
	    timp_mediu_ms: number;
	    ultima_interogare: string;
	    ultimul_timp_ms: number;
	    cale_folder_jurnale: string;
	    fisier_jurnal_curent: string;
	    dimensiune_jurnal: string;

	    static createFrom(source: any = {}) {
	        return new StatisticiDiagnoza(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.numar_erori = source["numar_erori"];
	        this.ultima_eroare = source["ultima_eroare"];
	        this.total_interogari = source["total_interogari"];
	        this.total_lente = source["total_lente"];
	        this.lenta_interogare = source["lenta_interogare"];
	        this.timp_maxim_ms = source["timp_maxim_ms"];
	        this.timp_mediu_ms = source["timp_mediu_ms"];
	        this.ultima_interogare = source["ultima_interogare"];
	        this.ultimul_timp_ms = source["ultimul_timp_ms"];
	        this.cale_folder_jurnale = source["cale_folder_jurnale"];
	        this.fisier_jurnal_curent = source["fisier_jurnal_curent"];
	        this.dimensiune_jurnal = source["dimensiune_jurnal"];
	    }
	}

}

export namespace main {

	export class CerereAplicaPropunereAI {
	    an: number;
	    luna: number;
	    propunere: ai.PropunereActiune;

	    static createFrom(source: any = {}) {
	        return new CerereAplicaPropunereAI(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.an = source["an"];
	        this.luna = source["luna"];
	        this.propunere = this.convertValues(source["propunere"], ai.PropunereActiune);
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
	export class CerereInterogareAI {
	    istoric: ai.MesajConversatie[];
	    an: number;
	    luna: number;

	    static createFrom(source: any = {}) {
	        return new CerereInterogareAI(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.istoric = this.convertValues(source["istoric"], ai.MesajConversatie);
	        this.an = source["an"];
	        this.luna = source["luna"];
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
	export class CerinteSistem {
	    sistem_operare: string;
	    procesor: string;
	    ram_minim: string;
	    ram_recomandat: string;
	    spatiu_disc: string;
	    dependinte: string;
	    retea: string;

	    static createFrom(source: any = {}) {
	        return new CerinteSistem(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sistem_operare = source["sistem_operare"];
	        this.procesor = source["procesor"];
	        this.ram_minim = source["ram_minim"];
	        this.ram_recomandat = source["ram_recomandat"];
	        this.spatiu_disc = source["spatiu_disc"];
	        this.dependinte = source["dependinte"];
	        this.retea = source["retea"];
	    }
	}
	export class InfoDetaliatLLM {
	    nume_model: string;
	    format: string;
	    motor_inferenta: string;
	    fereastra_context: number;
	    port_local: number;
	    stare: string;
	    dimensiune_model: string;
	    izolare_retea: string;

	    static createFrom(source: any = {}) {
	        return new InfoDetaliatLLM(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nume_model = source["nume_model"];
	        this.format = source["format"];
	        this.motor_inferenta = source["motor_inferenta"];
	        this.fereastra_context = source["fereastra_context"];
	        this.port_local = source["port_local"];
	        this.stare = source["stare"];
	        this.dimensiune_model = source["dimensiune_model"];
	        this.izolare_retea = source["izolare_retea"];
	    }
	}
	export class DetaliiDezvoltare {
	    nume_repo: string;
	    nume_aplicatie: string;
	    versiune: string;
	    commit: string;
	    arhitectura: string;
	    versiune_go: string;
	    versiune_wails: string;
	    cale_executabil: string;
	    checksum_exe: string;
	    autor: string;
	    licenta: string;
	    confidentialitate: string;
	    cerinte_minime: CerinteSistem;
	    detalii_llm: InfoDetaliatLLM;
	    statistici_diagnoza: diagnoza.StatisticiDiagnoza;

	    static createFrom(source: any = {}) {
	        return new DetaliiDezvoltare(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nume_repo = source["nume_repo"];
	        this.nume_aplicatie = source["nume_aplicatie"];
	        this.versiune = source["versiune"];
	        this.commit = source["commit"];
	        this.arhitectura = source["arhitectura"];
	        this.versiune_go = source["versiune_go"];
	        this.versiune_wails = source["versiune_wails"];
	        this.cale_executabil = source["cale_executabil"];
	        this.checksum_exe = source["checksum_exe"];
	        this.autor = source["autor"];
	        this.licenta = source["licenta"];
	        this.confidentialitate = source["confidentialitate"];
	        this.cerinte_minime = this.convertValues(source["cerinte_minime"], CerinteSistem);
	        this.detalii_llm = this.convertValues(source["detalii_llm"], InfoDetaliatLLM);
	        this.statistici_diagnoza = this.convertValues(source["statistici_diagnoza"], diagnoza.StatisticiDiagnoza);
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
	export class RaspunsAcoperirePlan {
	    acoperire?: calcul.RaportAcoperire;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsAcoperirePlan(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.acoperire = this.convertValues(source["acoperire"], calcul.RaportAcoperire);
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
	export class RaspunsAnalizaImportAngajati {
	    rezultat?: datebaza.RezultatAnalizaImportAngajati;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsAnalizaImportAngajati(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rezultat = this.convertValues(source["rezultat"], datebaza.RezultatAnalizaImportAngajati);
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
	export class RaspunsAnalizaImportSarbatori {
	    rezultat?: datebaza.RezultatAnalizaImportSarbatori;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsAnalizaImportSarbatori(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rezultat = this.convertValues(source["rezultat"], datebaza.RezultatAnalizaImportSarbatori);
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
	export class RaspunsAplicaPropunereAI {
	    succes: boolean;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsAplicaPropunereAI(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.succes = source["succes"];
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
	export class RaspunsAudit {
	    rezultat: audit.Rezultat;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsAudit(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rezultat = this.convertValues(source["rezultat"], audit.Rezultat);
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
	export class RaspunsBackup {
	    cale: string;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsBackup(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cale = source["cale"];
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
	export class RaspunsConfiguratiePontaj {
	    configuratie?: calcul.ConfiguratiePontaj;
	    optiuni: datebaza.OptiuneRegula[];
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsConfiguratiePontaj(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.configuratie = this.convertValues(source["configuratie"], calcul.ConfiguratiePontaj);
	        this.optiuni = this.convertValues(source["optiuni"], datebaza.OptiuneRegula);
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
	export class RaspunsDescarcaSablon {
	    continut_base64: string;
	    nume_fisier: string;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsDescarcaSablon(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.continut_base64 = source["continut_base64"];
	        this.nume_fisier = source["nume_fisier"];
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
	export class RaspunsDocumentAbsenta {
	    document?: datebaza.DocumentAbsenta;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsDocumentAbsenta(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.document = this.convertValues(source["document"], datebaza.DocumentAbsenta);
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
	export class RaspunsEchitatePlan {
	    avertismente: calcul.AvertismentEchitate[];
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsEchitatePlan(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.avertismente = this.convertValues(source["avertismente"], calcul.AvertismentEchitate);
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
	export class RaspunsExecutaImportAngajati {
	    date?: datebaza.Date;
	    importati: number;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsExecutaImportAngajati(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = this.convertValues(source["date"], datebaza.Date);
	        this.importati = source["importati"];
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
	export class RaspunsExecutaImportSarbatori {
	    calendar?: datebaza.CalendarAn;
	    importate: number;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsExecutaImportSarbatori(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.calendar = this.convertValues(source["calendar"], datebaza.CalendarAn);
	        this.importate = source["importate"];
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
	export class RaspunsExportExcel {
	    nume_fisier: string;
	    continut: string;
	    cale?: string;
	    anulat?: boolean;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsExportExcel(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nume_fisier = source["nume_fisier"];
	        this.continut = source["continut"];
	        this.cale = source["cale"];
	        this.anulat = source["anulat"];
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
	export class RaspunsExportJurnal {
	    cale?: string;
	    numar: number;
	    anulat?: boolean;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsExportJurnal(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cale = source["cale"];
	        this.numar = source["numar"];
	        this.anulat = source["anulat"];
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
	export class RaspunsFisierRaport {
	    anulat: boolean;
	    cale: string;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsFisierRaport(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.anulat = source["anulat"];
	        this.cale = source["cale"];
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
	export class RaspunsInterogareAI {
	    raspuns?: ai.RaspunsAsistent;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsInterogareAI(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.raspuns = this.convertValues(source["raspuns"], ai.RaspunsAsistent);
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
	export class RaspunsJurnal {
	    pagina: datebaza.PaginaJurnal;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsJurnal(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pagina = this.convertValues(source["pagina"], datebaza.PaginaJurnal);
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
	export class RaspunsListaReguli {
	    reguli: calcul.RegulaConfigurabila[];
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsListaReguli(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.reguli = this.convertValues(source["reguli"], calcul.RegulaConfigurabila);
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
	export class RaspunsModeleRapoarte {
	    configuratie?: rapoarte.Configuratie;
	    salvat_id: string;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsModeleRapoarte(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.configuratie = this.convertValues(source["configuratie"], rapoarte.Configuratie);
	        this.salvat_id = source["salvat_id"];
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
	export class RaspunsPlanificare {
	    planificare?: datebaza.PlanificareLuna;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsPlanificare(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.planificare = this.convertValues(source["planificare"], datebaza.PlanificareLuna);
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
	export class RaspunsPlanificatVsRealizat {
	    raport?: calcul.RaportPlanificatVsRealizat;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsPlanificatVsRealizat(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.raport = this.convertValues(source["raport"], calcul.RaportPlanificatVsRealizat);
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
	export class RaspunsPrecompletare {
	    plan?: calcul.PlanPrecompletare;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsPrecompletare(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.plan = this.convertValues(source["plan"], calcul.PlanPrecompletare);
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
	export class RaspunsPublicarePlan {
	    publicare?: datebaza.PublicarePlan;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsPublicarePlan(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.publicare = this.convertValues(source["publicare"], datebaza.PublicarePlan);
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
	export class RaspunsRaport {
	    document?: rapoarte.Document;
	    pdf: string;
	    token: string;
	    nume_fisier: string;
	    pagini: number;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsRaport(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.document = this.convertValues(source["document"], rapoarte.Document);
	        this.pdf = source["pdf"];
	        this.token = source["token"];
	        this.nume_fisier = source["nume_fisier"];
	        this.pagini = source["pagini"];
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
	export class RaspunsSalvareRegula {
	    regula?: calcul.RegulaConfigurabila;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsSalvareRegula(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.regula = this.convertValues(source["regula"], calcul.RegulaConfigurabila);
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
	export class RaspunsSarcinaAI {
	    id: string;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsSarcinaAI(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
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
	export class RaspunsStareAI {
	    stare: ai.StareMediu;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsStareAI(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stare = this.convertValues(source["stare"], ai.StareMediu);
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
	export class RaspunsStareLuna {
	    stare: datebaza.StareLuna;
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsStareLuna(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stare = this.convertValues(source["stare"], datebaza.StareLuna);
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
	export class RaspunsVerificaReguli {
	    incalcari: calcul.IncalcareRegula[];
	    eroare?: MesajEroare;

	    static createFrom(source: any = {}) {
	        return new RaspunsVerificaReguli(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.incalcari = this.convertValues(source["incalcari"], calcul.IncalcareRegula);
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
	export class RezultatSarcinaAI {
	    id: string;
	    raspuns?: ai.RaspunsAsistent;
	    eroare?: MesajEroare;
	    oprita: boolean;

	    static createFrom(source: any = {}) {
	        return new RezultatSarcinaAI(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.raspuns = this.convertValues(source["raspuns"], ai.RaspunsAsistent);
	        this.eroare = this.convertValues(source["eroare"], MesajEroare);
	        this.oprita = source["oprita"];
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

export namespace rapoarte {

	export class Imagine {
	    nume: string;
	    date: string;
	    inaltime_mm: number;

	    static createFrom(source: any = {}) {
	        return new Imagine(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nume = source["nume"];
	        this.date = source["date"];
	        this.inaltime_mm = source["inaltime_mm"];
	    }
	}
	export class Bloc {
	    stanga: Record<string, string>;
	    centru: Record<string, string>;
	    dreapta: Record<string, string>;
	    imagini?: Record<string, Imagine>;

	    static createFrom(source: any = {}) {
	        return new Bloc(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stanga = source["stanga"];
	        this.centru = source["centru"];
	        this.dreapta = source["dreapta"];
	        this.imagini = this.convertValues(source["imagini"], Imagine, true);
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
	export class Cerere {
	    model_id: string;
	    revizie: string;
	    an: number;
	    luna: number;
	    sectie_id: string;
	    departament_id: string;
	    angajat_id: string;
	    limba: string;

	    static createFrom(source: any = {}) {
	        return new Cerere(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model_id = source["model_id"];
	        this.revizie = source["revizie"];
	        this.an = source["an"];
	        this.luna = source["luna"];
	        this.sectie_id = source["sectie_id"];
	        this.departament_id = source["departament_id"];
	        this.angajat_id = source["angajat_id"];
	        this.limba = source["limba"];
	    }
	}
	export class Coloana {
	    cod: string;
	    titlu: string;
	    pondere: number;
	    culoare: string;

	    static createFrom(source: any = {}) {
	        return new Coloana(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cod = source["cod"];
	        this.titlu = source["titlu"];
	        this.pondere = source["pondere"];
	        this.culoare = source["culoare"];
	    }
	}
	export class Sursa {
	    cod: string;
	    coloane: string[];

	    static createFrom(source: any = {}) {
	        return new Sursa(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cod = source["cod"];
	        this.coloane = source["coloane"];
	    }
	}
	export class Pagina {
	    format: string;
	    orientare: string;
	    margine_mm: number;
	    font_pt: number;

	    static createFrom(source: any = {}) {
	        return new Pagina(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.format = source["format"];
	        this.orientare = source["orientare"];
	        this.margine_mm = source["margine_mm"];
	        this.font_pt = source["font_pt"];
	    }
	}
	export class Continut {
	    sursa: string;
	    coloane: string[];
	    text: Record<string, string>;

	    static createFrom(source: any = {}) {
	        return new Continut(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sursa = source["sursa"];
	        this.coloane = source["coloane"];
	        this.text = source["text"];
	    }
	}
	export class Definitie {
	    id: string;
	    denumire: Record<string, string>;
	    antet: Bloc;
	    continut: Continut;
	    subsol: Bloc;
	    pagina: Pagina;
	    activ: boolean;
	    revizie: string;

	    static createFrom(source: any = {}) {
	        return new Definitie(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.denumire = source["denumire"];
	        this.antet = this.convertValues(source["antet"], Bloc);
	        this.continut = this.convertValues(source["continut"], Continut);
	        this.subsol = this.convertValues(source["subsol"], Bloc);
	        this.pagina = this.convertValues(source["pagina"], Pagina);
	        this.activ = source["activ"];
	        this.revizie = source["revizie"];
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
	export class Configuratie {
	    modele: Definitie[];
	    surse: Sursa[];
	    variabile: string[];

	    static createFrom(source: any = {}) {
	        return new Configuratie(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.modele = this.convertValues(source["modele"], Definitie);
	        this.surse = this.convertValues(source["surse"], Sursa);
	        this.variabile = source["variabile"];
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


	export class Rand {
	    valori: string[];
	    culori: string[];

	    static createFrom(source: any = {}) {
	        return new Rand(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.valori = source["valori"];
	        this.culori = source["culori"];
	    }
	}
	export class Document {
	    titlu: string;
	    sursa: string;
	    versiune_raport: number;
	    revizie_model: string;
	    pagina: Pagina;
	    antet: string[];
	    antet_imagini?: Record<string, Imagine>;
	    introducere: string;
	    coloane: Coloana[];
	    randuri: Rand[];
	    subsol: string[];
	    subsol_imagini?: Record<string, Imagine>;
	    metadate: string;
	    nota: string;
	    fara_date: string;
	    paginare: string;

	    static createFrom(source: any = {}) {
	        return new Document(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.titlu = source["titlu"];
	        this.sursa = source["sursa"];
	        this.versiune_raport = source["versiune_raport"];
	        this.revizie_model = source["revizie_model"];
	        this.pagina = this.convertValues(source["pagina"], Pagina);
	        this.antet = source["antet"];
	        this.antet_imagini = this.convertValues(source["antet_imagini"], Imagine, true);
	        this.introducere = source["introducere"];
	        this.coloane = this.convertValues(source["coloane"], Coloana);
	        this.randuri = this.convertValues(source["randuri"], Rand);
	        this.subsol = source["subsol"];
	        this.subsol_imagini = this.convertValues(source["subsol_imagini"], Imagine, true);
	        this.metadate = source["metadate"];
	        this.nota = source["nota"];
	        this.fara_date = source["fara_date"];
	        this.paginare = source["paginare"];
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

