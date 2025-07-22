import {writable} from 'svelte/store';
import type {Server, ServerSession} from '../types/server.types';
import {ServerService} from '../services/server.service';

interface ServerStoreState {
    servers: Server[];
    sessions: ServerSession[];
    loading: boolean;
    sessionLoading: boolean;
    error: string | null;
    sessionError: string | null;
}

const createServerStore = () => {
    const {subscribe, set, update} = writable<ServerStoreState>({
        servers: [],
        sessions: [],
        loading: false,
        sessionLoading: false,
        error: null,
        sessionError: null
    });

    return {
        subscribe,

        async loadServers() {
            update(state => ({...state, loading: true, error: null}));
            try {
                const servers = await ServerService.listServers();
                update(state => ({...state, servers, loading: false}));
            } catch (error) {
                update(state => ({
                    ...state,
                    loading: false,
                    error: error ? error : 'Failed to load servers'
                }));
            }
        },

        async addServer(serverData: any) {
            update(state => ({...state, loading: true, error: null}));
            try {
                const newServer = await ServerService.addServer(serverData);
                update(state => ({
                    ...state,
                    servers: [...state.servers, newServer],
                    loading: false
                }));
                return newServer;
            } catch (error) {
                update(state => ({
                    ...state,
                    loading: false,
                    error: error ? error : 'Failed to add server'
                }));
                throw error;
            }
        },

        async updateServer(serverData: any) {
            update(state => ({...state, loading: true, error: null}));
            try {
                await ServerService.updateServer(serverData);
                update(state => ({
                    ...state,
                    servers: state.servers.map(server =>
                        server.id === serverData.id ? {...server, ...serverData} : server
                    ),
                    loading: false
                }));
            } catch (error) {
                update(state => ({
                    ...state,
                    loading: false,
                    error: error ? error : 'Failed to update server'
                }));
                throw error;
            }
        },

        async deleteServer(id: string) {
            update(state => ({...state, loading: true, error: null}));
            try {
                await ServerService.deleteServer(id);
                update(state => ({
                    ...state,
                    servers: state.servers.filter(server => server.id !== id),
                    loading: false
                }));
            } catch (error) {
                update(state => ({
                    ...state,
                    loading: false,
                    error: error ? error : 'Failed to delete server'
                }));
                throw error;
            }
        },

        async loadSessions() {
            update(state => ({...state, sessionLoading: true, sessionError: null}));
            try {
                const sessions = await ServerService.getListSession();
                update(state => ({...state, sessions, sessionLoading: false}));
            } catch (error) {
                update(state => ({
                    ...state,
                    sessionLoading: false,
                    sessionError: error instanceof Error ? error.message : 'Failed to load sessions'
                }));
            }
        },


        async disconnectSession(sessionId: string): Promise<boolean> {
            update(state => ({...state, sessionLoading: true, sessionError: null}));
            try {
                await ServerService.closeSession(sessionId);
                await this.loadSessions(); // Refresh sessions after disconnecting
                return true;
            } catch (error) {
                update(state => ({
                    ...state,
                    sessionLoading: false,
                    sessionError: error ? error : 'Failed to disconnect session'
                }));
                return false;
            }
        },

        async connectToServer(serverId: string): Promise<string | null> {
            update(state => ({...state, sessionLoading: true, sessionError: null}));
            try {
                const sessionId = await ServerService.connectSession(serverId);
                await this.loadSessions(); // Refresh sessions after connecting
                return sessionId;
            } catch (error) {
                update(state => ({
                    ...state,
                    sessionLoading: false,
                    sessionError: error ? error : 'Failed to connect to server'
                }));
                return null;
            }
        },

        // Utility methods
        clearErrors() {
            update(state => ({...state, error: null, sessionError: null}));
        },

        getSessionByServerId(serverId: string) {
            let currentState: ServerStoreState;
            subscribe(state => currentState = state)();
            return currentState.sessions.find(session => session.serverInfo.id === serverId);
        },

        isServerConnected(serverId: string): boolean {
            return !!this.getSessionByServerId(serverId);
        }
    };
};

export const serverStore = createServerStore();