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

export namespace dtos {
	
	export class CategoryCreate {
	    name: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new CategoryCreate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	    }
	}
	export class CategoryDetail {
	    id: number;
	    name: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new CategoryDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	    }
	}
	export class CategoryListItem {
	    id: number;
	    name: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new CategoryListItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	    }
	}
	export class CategoryUpdate {
	    id: number;
	    name: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new CategoryUpdate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	    }
	}
	export class RuleCreate {
	    categoryId: number;
	    appName: string;
	    additionalDataKey: string;
	    expression: string;
	    isRegex: boolean;
	    priority: number;
	
	    static createFrom(source: any = {}) {
	        return new RuleCreate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.categoryId = source["categoryId"];
	        this.appName = source["appName"];
	        this.additionalDataKey = source["additionalDataKey"];
	        this.expression = source["expression"];
	        this.isRegex = source["isRegex"];
	        this.priority = source["priority"];
	    }
	}
	export class RuleDetail {
	    id: number;
	    categoryId: number;
	    appName: string;
	    additionalDataKey: string;
	    expression: string;
	    isRegex: boolean;
	    priority: number;
	
	    static createFrom(source: any = {}) {
	        return new RuleDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.categoryId = source["categoryId"];
	        this.appName = source["appName"];
	        this.additionalDataKey = source["additionalDataKey"];
	        this.expression = source["expression"];
	        this.isRegex = source["isRegex"];
	        this.priority = source["priority"];
	    }
	}
	export class RuleListItem {
	    id: number;
	    categoryId: number;
	    appName: string;
	    expression: string;
	    priority: number;
	
	    static createFrom(source: any = {}) {
	        return new RuleListItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.categoryId = source["categoryId"];
	        this.appName = source["appName"];
	        this.expression = source["expression"];
	        this.priority = source["priority"];
	    }
	}
	export class RuleUpdate {
	    id: number;
	    categoryId: number;
	    appName: string;
	    additionalDataKey: string;
	    expression: string;
	    isRegex: boolean;
	    priority: number;
	
	    static createFrom(source: any = {}) {
	        return new RuleUpdate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.categoryId = source["categoryId"];
	        this.appName = source["appName"];
	        this.additionalDataKey = source["additionalDataKey"];
	        this.expression = source["expression"];
	        this.isRegex = source["isRegex"];
	        this.priority = source["priority"];
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

