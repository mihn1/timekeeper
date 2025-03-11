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
	export class Category {
	    Id: string;
	    Name: string;
	    Description: string;
	    CategoryTypeId: number;
	
	    static createFrom(source: any = {}) {
	        return new Category(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	        this.CategoryTypeId = source["CategoryTypeId"];
	    }
	}
	export class CategoryRule {
	    RuleId: number;
	    CategoryId: string;
	    AppName: string;
	    AdditionalDataKey: string;
	    Expression: string;
	    IsRegex: boolean;
	    Priority: number;
	
	    static createFrom(source: any = {}) {
	        return new CategoryRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RuleId = source["RuleId"];
	        this.CategoryId = source["CategoryId"];
	        this.AppName = source["AppName"];
	        this.AdditionalDataKey = source["AdditionalDataKey"];
	        this.Expression = source["Expression"];
	        this.IsRegex = source["IsRegex"];
	        this.Priority = source["Priority"];
	    }
	}

}

