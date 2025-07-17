import type { dto } from '../../wailsjs/go/models';

/**
 * Represents a server available for connection in the frontend.
 * It extends the auto-generated ServerResponse from Wails with frontend-specific properties.
 */
export type Server = dto.ServerResponse & {
    isLocal?: boolean;
    defaultPath: string;
    osType: 'windows' | 'linux' | 'unknown';
};

/**
 * Represents the complete state of a single file explorer pane.
 */
export interface Session {
    // Connection State
    server: Server | null;         // The server this pane is targeting.
    wailsSessionId: string | null; // The backend session ID for remote connections.
    status: 'disconnected' | 'connecting' | 'connected' | 'error';

    // File/Directory State
    currentPath: string;
    files: dto.FileInfo[];
    isLoading: boolean;
    error: string | null;
}