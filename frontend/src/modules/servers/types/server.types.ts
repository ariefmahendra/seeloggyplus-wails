export interface Server {
    id: string;
    name: string;
    address: string;
    port: number;
    user: string;
    password: string;
    createdAt: string;
    updatedAt: string;
}

export interface ServerCreateRequest {
    name: string;
    address: string;
    port: number;
    user: string;
    password: string;
}

export interface ServerUpdateRequest {
    id: string;
    name: string;
    address: string;
    port: number;
    user: string;
    password: string;
}

export interface ServerFormData {
    name: string;
    address: string;
    port: number;
    user: string;
    password: string;
}

export interface ServerValidationErrors {
    name?: string;
    address?: string;
    port?: string;
    user?: string;
    password?: string;
}

export interface ServerSession {
    id: string;
    serverInfo: {
        id: string;
        name: string;
        address: string;
        port: number;
        user: string;
        password: string;
    };
}

export interface SessionStoreState {
    sessions: ServerSession[];
    loading: boolean;
    error: string | null;
}