export namespace datatypes {
	
	export class DateOnly {
	
	
	    static createFrom(source: any = {}) {
	        return new DateOnly(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

export namespace models {
	
	export class AppAggregation {
	    AppName: string;
	    AdditionalData: any;
	    Date: datatypes.DateOnly;
	    TimeElapsed: number;
	
	    static createFrom(source: any = {}) {
	        return new AppAggregation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AppName = source["AppName"];
	        this.AdditionalData = source["AdditionalData"];
	        this.Date = this.convertValues(source["Date"], datatypes.DateOnly);
	        this.TimeElapsed = source["TimeElapsed"];
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

