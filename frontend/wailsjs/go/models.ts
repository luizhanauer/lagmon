export namespace config {
	
	export class DiagramNode {
	    name: string;
	    ip: string;
	
	    static createFrom(source: any = {}) {
	        return new DiagramNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.ip = source["ip"];
	    }
	}
	export class NetworkDiagramConfig {
	    local: DiagramNode;
	    gateway: DiagramNode;
	    internet: DiagramNode;
	
	    static createFrom(source: any = {}) {
	        return new NetworkDiagramConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.local = this.convertValues(source["local"], DiagramNode);
	        this.gateway = this.convertValues(source["gateway"], DiagramNode);
	        this.internet = this.convertValues(source["internet"], DiagramNode);
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
	export class AppConfig {
	    retention_days: number;
	    network_diagram: NetworkDiagramConfig;
	    targets: domain.Host[];
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.retention_days = source["retention_days"];
	        this.network_diagram = this.convertValues(source["network_diagram"], NetworkDiagramConfig);
	        this.targets = this.convertValues(source["targets"], domain.Host);
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

export namespace domain {
	
	export class Host {
	    id: string;
	    name: string;
	    ip: string;
	    isGateway: boolean;
	    active: boolean;
	    showInDiagram: boolean;
	
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
	        this.showInDiagram = source["showInDiagram"];
	    }
	}

}

