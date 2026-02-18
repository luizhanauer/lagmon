export namespace domain {
	
	export class Host {
	    id: string;
	    name: string;
	    ip: string;
	    isGateway: boolean;
	    active: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Host(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.ip = source["ip"];
	        this.isGateway = source["isGateway"];
	        this.active = source["active"];
	    }
	}

}

