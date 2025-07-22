export namespace dto {
	
	export class DriveInfo {
	    name: string;
	    path: string;
	    label: string;
	    type: string;
	    isRoot: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DriveInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.label = source["label"];
	        this.type = source["type"];
	        this.isRoot = source["isRoot"];
	    }
	}
	export class FileInfo {
	    sessionID: string;
	    name: string;
	    size: number;
	    path: string;
	    isDir: boolean;
	    modTime: string;
	    mode: string;
	    isDrive: boolean;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionID = source["sessionID"];
	        this.name = source["name"];
	        this.size = source["size"];
	        this.path = source["path"];
	        this.isDir = source["isDir"];
	        this.modTime = source["modTime"];
	        this.mode = source["mode"];
	        this.isDrive = source["isDrive"];
	    }
	}
	export class ServerCreateRequest {
	    name?: string;
	    address?: string;
	    port?: number;
	    user?: string;
	    password?: string;
	
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
	    }
	}
	export class ServerSessionManagement {
	    id: string;
	    name: string;
	    address: string;
	    port: number;
	    user: string;
	    password: string;
	
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
	    }
	}
	export class ServerUpdateRequest {
	    id?: string;
	    name?: string;
	    address?: string;
	    port?: number;
	    user?: string;
	    password?: string;
	
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
	    }
	}
	export class SessionManagerDto {
	    id: string;
	    serverInfo?: ServerSessionManagement;
	
	    static createFrom(source: any = {}) {
	        return new SessionManagerDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.serverInfo = this.convertValues(source["serverInfo"], ServerSessionManagement);
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
	export class SettingsRequestDto {
	    key: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new SettingsRequestDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.value = source["value"];
	    }
	}
	export class SettingsResponseDto {
	    id: number;
	    name: string;
	    key: string;
	    value: string;
	    description: string;
	    type: string;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new SettingsResponseDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.key = source["key"];
	        this.value = source["value"];
	        this.description = source["description"];
	        this.type = source["type"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
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

