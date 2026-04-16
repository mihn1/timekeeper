export namespace dtos {
	
	export class AppUsageItem {
	    appName: string;
	    timeElapsed: number;
	    categoryId: number;
	    categoryName: string;
	
	    static createFrom(source: any = {}) {
	        return new AppUsageItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appName = source["appName"];
	        this.timeElapsed = source["timeElapsed"];
	        this.categoryId = source["categoryId"];
	        this.categoryName = source["categoryName"];
	    }
	}
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
	export class CategoryUsageItem {
	    id: number;
	    name: string;
	    timeElapsed: number;
	
	    static createFrom(source: any = {}) {
	        return new CategoryUsageItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.timeElapsed = source["timeElapsed"];
	    }
	}
	export class DailyCategorySummary {
	    date: string;
	    categoryId: number;
	    categoryName: string;
	    timeElapsed: number;
	
	    static createFrom(source: any = {}) {
	        return new DailyCategorySummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.categoryId = source["categoryId"];
	        this.categoryName = source["categoryName"];
	        this.timeElapsed = source["timeElapsed"];
	    }
	}
	export class DayActivity {
	    date: string;
	    totalMs: number;
	    topCategoryId: number;
	
	    static createFrom(source: any = {}) {
	        return new DayActivity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.totalMs = source["totalMs"];
	        this.topCategoryId = source["topCategoryId"];
	    }
	}
	export class EventLogItem {
	    id: number;
	    appName: string;
	    startTime: string;
	    endTime: string;
	    durationSecs: number;
	    categoryId: number;
	    urlOrTitle: string;
	
	    static createFrom(source: any = {}) {
	        return new EventLogItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.appName = source["appName"];
	        this.startTime = source["startTime"];
	        this.endTime = source["endTime"];
	        this.durationSecs = source["durationSecs"];
	        this.categoryId = source["categoryId"];
	        this.urlOrTitle = source["urlOrTitle"];
	    }
	}
	export class GoalItem {
	    id: number;
	    name: string;
	    isActive: boolean;
	    categoryIds: number[];
	    categoryNames: string[];
	    frequency: number;
	    targetMs: number;
	
	    static createFrom(source: any = {}) {
	        return new GoalItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.isActive = source["isActive"];
	        this.categoryIds = source["categoryIds"];
	        this.categoryNames = source["categoryNames"];
	        this.frequency = source["frequency"];
	        this.targetMs = source["targetMs"];
	    }
	}
	export class PreferencesDto {
	    timezone: string;
	    minEventDurationMs: number;
	
	    static createFrom(source: any = {}) {
	        return new PreferencesDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timezone = source["timezone"];
	        this.minEventDurationMs = source["minEventDurationMs"];
	    }
	}
	export class RuleCreate {
	    categoryId: number;
	    appName: string;
	    additionalDataKey: string;
	    expression: string;
	    isRegex: boolean;
	    priority: number;
	    isExclusion: boolean;
	
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
	        this.isExclusion = source["isExclusion"];
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
	    isExclusion: boolean;
	
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
	        this.isExclusion = source["isExclusion"];
	    }
	}
	export class RuleListItem {
	    id: number;
	    categoryId: number;
	    appName: string;
	    additionalDataKey: string;
	    expression: string;
	    isRegex: boolean;
	    priority: number;
	    isExclusion: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RuleListItem(source);
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
	        this.isExclusion = source["isExclusion"];
	    }
	}
	export class RuleMatchResult {
	    matched: boolean;
	    categoryId: number;
	    categoryName: string;
	    matchedRule?: RuleDetail;
	
	    static createFrom(source: any = {}) {
	        return new RuleMatchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.matched = source["matched"];
	        this.categoryId = source["categoryId"];
	        this.categoryName = source["categoryName"];
	        this.matchedRule = this.convertValues(source["matchedRule"], RuleDetail);
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
	export class RuleUpdate {
	    id: number;
	    categoryId: number;
	    appName: string;
	    additionalDataKey: string;
	    expression: string;
	    isRegex: boolean;
	    priority: number;
	    isExclusion: boolean;
	
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
	        this.isExclusion = source["isExclusion"];
	    }
	}

}

