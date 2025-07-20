<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import type {Server, ServerSession} from '../types/server.types';
    import ServerItem from "./ServerItem.svelte";

    export let servers: Server[] = [];
    export let sessions: ServerSession[] = [];
    export let disabled = false;
    let loadingServerIds = new Set<string>();

    const dispatch = createEventDispatcher<{
        edit: { server: Server };
        delete: { server: Server };
        connect: { server: Server };
        disconnect: { server: Server };
    }>();

    function isServerConnected(serverId: string): boolean {
        return sessions.some(session => session.serverInfo.id === serverId);
    }

    function getSessionByServerId(serverId: string): ServerSession | null {
        return sessions.find(session => session.serverInfo.id === serverId) || null;
    }

    function handleConnect(event: CustomEvent) {
        const { server } = event.detail;
        loadingServerIds.add(server.id);
        loadingServerIds = loadingServerIds;

        dispatch('connect', { server });
    }

    function handleDisconnect(event: CustomEvent) {
        const { server } = event.detail;
        loadingServerIds.add(server.id);
        loadingServerIds = loadingServerIds;

        dispatch('disconnect', { server });
    }

    export function clearServerLoading(serverId: string) {
        loadingServerIds.delete(serverId);
        loadingServerIds = loadingServerIds;
    }

    export function setServerLoading(serverId: string) {
        loadingServerIds.add(serverId);
        loadingServerIds = loadingServerIds;
    }
</script>

<div class="space-y-4">
    {#each servers as server (server.id)}
        <ServerItem
                {server}
                session={getSessionByServerId(server.id)}
                {disabled}
                sessionLoading={loadingServerIds.has(server.id)}
                on:edit
                on:delete
                on:connect={handleConnect}
                on:disconnect={handleDisconnect}
        />
    {/each}
</div>