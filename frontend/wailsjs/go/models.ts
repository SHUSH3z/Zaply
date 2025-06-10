export namespace main {
	
	export class Ponto {
	    id: number;
	    coordenada: string;
	    lat: number;
	    long: number;
	    endereco: string;
	
	    static createFrom(source: any = {}) {
	        return new Ponto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.coordenada = source["coordenada"];
	        this.lat = source["lat"];
	        this.long = source["long"];
	        this.endereco = source["endereco"];
	    }
	}
	export class KMZRequest {
	    pasta: string;
	    pontos: Ponto[];
	
	    static createFrom(source: any = {}) {
	        return new KMZRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pasta = source["pasta"];
	        this.pontos = this.convertValues(source["pontos"], Ponto);
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

