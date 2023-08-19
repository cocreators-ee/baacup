export namespace main {

	export class BackupConfig {
	    keepSaves: number;
	    maxMBPerGame: number;

	    static createFrom(source: any = {}) {
	        return new BackupConfig(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.keepSaves = source["keepSaves"];
	        this.maxMBPerGame = source["maxMBPerGame"];
	    }
	}
	export class CompactionConfig {
	    keepSaves: number;
	    compactAfterDays: number;

	    static createFrom(source: any = {}) {
	        return new CompactionConfig(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.keepSaves = source["keepSaves"];
	        this.compactAfterDays = source["compactAfterDays"];
	    }
	}
	export class Config {
	    disabledRules: string[];
	    pathSeparator: string;
	    backups?: BackupConfig;
	    compaction?: CompactionConfig;
	    // Go type: time
	    rulesLastUpdated: any;
	    rulesAutoUpdate: boolean;

	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.disabledRules = source["disabledRules"];
	        this.pathSeparator = source["pathSeparator"];
	        this.backups = this.convertValues(source["backups"], BackupConfig);
	        this.compaction = this.convertValues(source["compaction"], CompactionConfig);
	        this.rulesLastUpdated = this.convertValues(source["rulesLastUpdated"], null);
	        this.rulesAutoUpdate = source["rulesAutoUpdate"];
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	export class Monitor {
	    path: string;
	    ruleFilename: string;

	    static createFrom(source: any = {}) {
	        return new Monitor(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.ruleFilename = source["ruleFilename"];
	    }
	}

}
