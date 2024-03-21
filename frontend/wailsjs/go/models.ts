export namespace app {
	
	export class QuerySetting {
	    keywords: string;
	    cookie: string;
	    threads: number;
	    assetsNumber: number;
	    checkAlive: boolean;
	
	    static createFrom(source: any = {}) {
	        return new QuerySetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.keywords = source["keywords"];
	        this.cookie = source["cookie"];
	        this.threads = source["threads"];
	        this.assetsNumber = source["assetsNumber"];
	        this.checkAlive = source["checkAlive"];
	    }
	}

}

export namespace types {
	
	export class Asset {
	    icon: string;
	    host: string;
	    ip: string;
	    port: number;
	    title: string;
	    server: string;
	    country: string;
	    organization: string;
	    header: string;
	    certificate: string;
	    alive: boolean;
	    lastUpdateTime: string;
	
	    static createFrom(source: any = {}) {
	        return new Asset(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.icon = source["icon"];
	        this.host = source["host"];
	        this.ip = source["ip"];
	        this.port = source["port"];
	        this.title = source["title"];
	        this.server = source["server"];
	        this.country = source["country"];
	        this.organization = source["organization"];
	        this.header = source["header"];
	        this.certificate = source["certificate"];
	        this.alive = source["alive"];
	        this.lastUpdateTime = source["lastUpdateTime"];
	    }
	}

}

