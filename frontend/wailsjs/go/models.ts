export namespace dto {
	
	export class FileInfo {
	    sessionID: string;
	    name: string;
	    size: number;
	    isDir: boolean;
	    // Go type: time
	    modTime: any;
	    mode: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionID = source["sessionID"];
	        this.name = source["name"];
	        this.size = source["size"];
	        this.isDir = source["isDir"];
	        this.modTime = this.convertValues(source["modTime"], null);
	        this.mode = source["mode"];
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
	export class ServerCreateRequest {
	    name?: string;
	    address?: string;
	    port?: number;
	    user?: string;
	    password?: string;
	    type?: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerCreateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.address = source["address"];
	        this.port = source["port"];
	        this.user = source["user"];
	        this.password = source["password"];
	        this.type = source["type"];
	    }
	}
	export class ServerResponse {
	    id: string;
	    name: string;
	    address: string;
	    port: number;
	    user: string;
	    password: string;
	    createdAt: string;
	    updatedAt: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.address = source["address"];
	        this.port = source["port"];
	        this.user = source["user"];
	        this.password = source["password"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	        this.type = source["type"];
	    }
	}
	export class ServerSessionManagement {
	    id: string;
	    name: string;
	    address: string;
	    port: number;
	    user: string;
	    password: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerSessionManagement(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.address = source["address"];
	        this.port = source["port"];
	        this.user = source["user"];
	        this.password = source["password"];
	        this.type = source["type"];
	    }
	}
	export class ServerUpdateRequest {
	    id?: string;
	    name?: string;
	    address?: string;
	    port?: number;
	    user?: string;
	    password?: string;
	    type?: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerUpdateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.address = source["address"];
	        this.port = source["port"];
	        this.user = source["user"];
	        this.password = source["password"];
	        this.type = source["type"];
	    }
	}

}

export namespace entity {
	
	export class Session {
	    ID: string;
	    ServerInfo?: dto.ServerSessionManagement;
	    SSHClient?: ssh.Client;
	    // Go type: sftp
	    SFTPClient?: any;
	
	    static createFrom(source: any = {}) {
	        return new Session(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.ServerInfo = this.convertValues(source["ServerInfo"], dto.ServerSessionManagement);
	        this.SSHClient = this.convertValues(source["SSHClient"], ssh.Client);
	        this.SFTPClient = this.convertValues(source["SFTPClient"], null);
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

export namespace ssh {
	
	export class Client {
	    Conn: any;
	
	    static createFrom(source: any = {}) {
	        return new Client(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Conn = source["Conn"];
	    }
	}

}

