import {
    AddServer,
    DeleteServer,
    GetServerById,
    ListServers,
    TestConnection,
    UpdateServer,
    CancelRequest,
    ConnectSession,
    GetListSession,
    CloseSession
} from '../../../../wailsjs/go/main/App';
import type {
    Server,
    ServerCreateRequest,
    ServerSession,
    ServerUpdateRequest,
} from '../types/server.types';

export class ServerService {
    static async listServers(): Promise<Server[]> {
        try {
            const result = await ListServers();
            return result || [];
        } catch (error) {
            console.error('Failed to list servers:', error);
            throw error;
        }
    }

    static async addServer(payload: ServerCreateRequest): Promise<Server> {
        try {
            return await AddServer(payload);
        } catch (error) {
            console.error('Failed to add server:', error);
            throw error;
        }
    }

    static async updateServer(payload: ServerUpdateRequest): Promise<void> {
        try {
            await UpdateServer(payload);
        } catch (error) {
            console.error('Failed to update server:', error);
            throw error;
        }
    }

    static async deleteServer(id: string): Promise<void> {
        try {
            await DeleteServer(id);
        } catch (error) {
            console.error('Failed to delete server:', error);
            throw error;
        }
    }

    static async getServerById(id: string): Promise<Server> {
        try {
            return await GetServerById(id);
        } catch (error) {
            console.error('Failed to get server by ID:', error);
            throw error;
        }
    }

    static async testConnection(requestId: string, payload: ServerCreateRequest): Promise<string> {
        try {
            await TestConnection(requestId, payload);
            return requestId;
        } catch (error) {
            console.error('Failed to test connection:', error);
            throw error;
        }
    }

    static async cancelTestConnection(requestId: string): Promise<void> {
        try {
            await CancelRequest(requestId);
        } catch (error) {
            console.error('Failed to cancel test connection:', error);
            throw error;
        }
    }

    static async connectSession(serverId: string): Promise<string> {
        try {
            return await ConnectSession(serverId);
        } catch (error) {
            console.error('Failed to connect session:', error);
            throw error;
        }
    }

    static async closeSession(sessionId: string): Promise<void> {
        try {
            await CloseSession(sessionId);
        } catch (error) {
            console.error('Failed to close session:', error);
            throw error;
        }
    }

    static async getListSession(): Promise<Array<ServerSession>> {
        try {
            const sessions = await GetListSession();

            if (!sessions) {
                return [];
            }

            const serverSessions: ServerSession[] = [];
            sessions.forEach(session => {
                serverSessions.push({
                    id: session.id,
                    serverInfo: {
                        id: session.serverInfo.id,
                        name: session.serverInfo.name,
                        address: session.serverInfo.address,
                        port: session.serverInfo.port,
                        user: session.serverInfo.user,
                        password: session.serverInfo.password
                    }
                });
            })

            return serverSessions;
        } catch (error) {
            console.error('Failed to get list of sessions:', error);
            throw error;
        }
    }
}