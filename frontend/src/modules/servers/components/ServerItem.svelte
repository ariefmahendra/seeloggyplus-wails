<script lang="ts">
    import {createEventDispatcher} from 'svelte';
    import type {Server, ServerSession} from '../types/server.types';
    import {Tag, Button} from "carbon-components-svelte";
    import {StopOutline, DirectionRight_01, Edit, TrashCan} from 'carbon-icons-svelte'

    export let server: Server;
    export let session: ServerSession | null = null;
    export let disabled = false;
    export let sessionLoading = false;

    const dispatch = createEventDispatcher<{
        edit: { server: Server };
        delete: { server: Server };
        connect: { server: Server };
        disconnect: { server: Server };
    }>();

    $: isConnected = session !== null;
    $: connectionStatus = isConnected ? 'Connected' : 'Disconnected';

    function handleEdit() {
        if (disabled || sessionLoading) return;
        dispatch('edit', {server});
    }

    function handleDelete() {
        if (disabled || sessionLoading) return;
        dispatch('delete', {server});
    }

    function handleConnect() {
        if (disabled || sessionLoading || isConnected) return;
        dispatch('connect', {server});
    }

    function handleDisconnect() {
        if (disabled || sessionLoading || !isConnected) return;
        dispatch('disconnect', {server});
    }
</script>

<div class="bg-white border border-gray-200 rounded-lg p-3 mx-0 hover:shadow-md transition-shadow duration-200 {disabled ? 'opacity-50' : ''}">
    <!-- Server Header -->
    <div class="flex items-start justify-between mb-1">
        <div class="flex items-center space-x-3">
            <!-- Server Status Indicator -->
            <div class="flex-shrink-0">
                <div class="w-3 h-3 bg-gray-400 rounded-full"></div>
            </div>
            <div>
                <h3 class="text-lg font-semibold text-gray-900">{server.name}</h3>
                <p class="text-sm text-gray-500">{server.address}:{server.port}</p>
            </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center space-x-2">
            <!-- Connection Status-->
            <Tag type={isConnected ? 'green' : 'gray'}>
                {connectionStatus}
            </Tag>

            <!-- Action Buttons -->
            <div class="flex space-x-2">
                <!-- Connect/Disconnect Button -->
                {#if isConnected}
                    <Button
                            kind="danger"
                            size="small"
                            icon={StopOutline}
                            disabled={disabled || sessionLoading}
                            on:click={handleDisconnect}
                    >
                        Disconnect
                    </Button>
                {:else}
                    <Button
                            kind="primary"
                            size="small"
                            icon={DirectionRight_01}
                            disabled={disabled || sessionLoading}
                            on:click={handleConnect}
                    >
                        {sessionLoading ? 'Connecting...' : 'Connect'}
                    </Button>
                {/if}

                <!-- Edit Button -->
                <Button
                        kind="ghost"
                        size="small"
                        iconDescription="Edit server"
                        icon={Edit}
                        disabled={disabled}
                        on:click={handleEdit}
                />

                <!-- Delete Button -->
                <Button
                        kind="ghost"
                        size="small"
                        iconDescription="Delete server"
                        icon={TrashCan}
                        disabled={disabled}
                        on:click={handleDelete}
                />
            </div>
        </div>
    </div>

    <!-- Server Details -->
    <div class="grid grid-cols-2 gap-4 text-sm">
        <div>
            <span class="text-gray-500">User:</span>
            <span class="ml-2 text-gray-900 font-medium">{server.user}</span>
        </div>
    </div>

    <!-- Server Metadata -->
    <div class="mt-1 border-t border-gray-100">
        <div class="flex items-center justify-between text-xs text-gray-500">
            <span>Created: {new Date(server.createdAt).toLocaleDateString()}</span>
            <span>Updated: {new Date(server.updatedAt).toLocaleDateString()}</span>
        </div>
    </div>
</div>