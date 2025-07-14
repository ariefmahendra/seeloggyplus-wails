export namespace dto {
	
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

