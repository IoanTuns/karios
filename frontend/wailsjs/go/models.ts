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

}

