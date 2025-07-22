<script lang="ts">
    import {Button, Modal, Search} from 'carbon-components-svelte';
    import {Add, ServerProxy} from 'carbon-icons-svelte';
    import {serverStore} from './stores/server.store';
    import ServerList from './components/ServerList.svelte';
    import AddServerModal from './modals/AddServerModal.svelte';
    import EditServerModal from './modals/EditServerModal.svelte';
    import DeleteServerModal from './modals/DeleteServerModal.svelte';
    import type {Server} from './types/server.types';

    export let isOpen = false;
    export let disableInteraction = false;

    let showAddModal = false;
    let showEditModal = false;
    let showDeleteModal = false;
    let selectedServer: Server | null = null;
    let searchQuery = '';

    // Reactive variables
    let serverListComponent: ServerList;

    $: if (isOpen) {
        serverStore.loadServers();
        serverStore.loadSessions();
    }
    $: filteredServers = $serverStore.servers.filter(server =>
        server.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        server.address.toLowerCase().includes(searchQuery.toLowerCase())
    );
    $: isInteractionDisabled = disableInteraction || $serverStore.loading;
    $: connectedServersCount = $serverStore.sessions.length;

    function closeModal() {
        if (isInteractionDisabled) return;
        isOpen = false;
    }

    function handleAddServer() {
        if (isInteractionDisabled) return;
        showAddModal = true;
    }

    function handleEditServer(event: CustomEvent<{ server: Server }>) {
        if (isInteractionDisabled) return;
        selectedServer = event.detail.server;
        showEditModal = true;
    }

    function handleDeleteServer(event: CustomEvent<{ server: Server }>) {
        if (isInteractionDisabled) return;
        selectedServer = event.detail.server;
        showDeleteModal = true;
    }

    // async function
    async function handleConnectServer(event: CustomEvent<{ server: Server }>) {
        if (isInteractionDisabled) return;

        const server = event.detail.server;

        if (serverStore.isServerConnected(server.id)) {
            console.log('Server already connected');
            return;
        }

        try {
            serverListComponent?.setServerLoading(server.id);

            await serverStore.connectToServer(server.id);
        } catch (error) {
            console.error('Failed to connect to server:', error);
        } finally {
            serverListComponent?.clearServerLoading(server.id);
        }
    }

    async function handleDisconnectServer(event: CustomEvent<{ server: Server }>) {
        if (isInteractionDisabled) return;

        const server = event.detail.server;
        const session = serverStore.getSessionByServerId(server.id);

        if (!session) {
            console.log('No active session found for server');
            return;
        }

        try {
            serverListComponent?.setServerLoading(server.id);
            await serverStore.disconnectSession(session.id);
        } catch (error) {
            console.error('Failed to disconnect from server:', error);
        } finally {
            serverListComponent?.clearServerLoading(server.id);
        }
    }


    function handleServerAdded() {
        showAddModal = false;
        serverStore.loadServers();
    }

    function handleServerUpdated() {
        showEditModal = false;
        selectedServer = null;
        serverStore.loadServers();
    }

    function handleServerDeleted() {
        showDeleteModal = false;
        selectedServer = null;
        serverStore.loadServers();
        serverStore.loadSessions();
    }
</script>

<Modal
        bind:open={isOpen}
        modalHeading="Server Management"
        primaryButtonText="Ok"
        secondaryButtonText="Close"
        selectorPrimaryFocus="#search-input"
        on:click:button--primary={closeModal}
        on:click:button--secondary={closeModal}
        on:close={closeModal}
        size="lg"
        hasScrollingContent
        passiveModal
>
    <div slot="heading">
        <div class="flex items-center space-x-2">
            <ServerProxy size={24}/>
            <span>Server Management</span>
        </div>
    </div>

    <!-- Modal Content -->
    <div class="space-y-6 overflow-hidden">
        <!-- Header Section -->
        <div class="flex items-center justify-between">
            <div>
                <h3 class="text-lg font-medium text-gray-900">Manage Your Servers</h3>
                <p class="text-sm text-gray-500">Add, edit, delete, and connect to your servers</p>
            </div>
            <div class="flex items-center space-x-2">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                    {$serverStore.servers.length} servers
                </span>
            </div>
        </div>

        <!-- Toolbar Section -->
        <div class="flex items-center justify-between space-x-4">
            <!-- Search -->
            <div class="flex-1 max-w-md">
                <Search
                        id="search-input"
                        bind:value={searchQuery}
                        placeholder="Search servers..."
                        size="sm"
                        disabled={isInteractionDisabled}
                />
            </div>

            <!-- Add Server Button -->
            <Button
                    kind="primary"
                    icon={Add}
                    on:click={handleAddServer}
                    disabled={$serverStore.loading}
            >
                Add Server
            </Button>
        </div>

        <!-- Content Section -->
        <div class="min-h-[300px]">
            {#if $serverStore.loading}
                <div class="text-center py-8">Loading servers...</div>
            {:else if $serverStore.error}
                <!-- Error State -->
                <div class="text-red-600 text-center py-4">{$serverStore.error}</div>
            {:else if filteredServers.length === 0}
                <!-- Empty State -->
                <div class="text-center py-8 text-gray-500">
                    {searchQuery ? 'No servers match your search' : 'No servers configured'}
                </div>
            {:else}
                <!-- Server List -->
                <div class="max-h-96 p-1 overflow-y-auto overflow-x-hidden">
                    <ServerList
                            bind:this={serverListComponent}
                            servers={filteredServers}
                            sessions={$serverStore.sessions}
                            disabled={isInteractionDisabled}
                            on:edit={handleEditServer}
                            on:delete={handleDeleteServer}
                            on:connect={handleConnectServer}
                            on:disconnect={handleDisconnectServer}
                    />
                </div>
            {/if}

            {#if $serverStore.sessionError}
                <div class="text-red-600 text-center py-2 text-sm">
                    {$serverStore.sessionError}
                </div>
            {/if}
        </div>
    </div>
</Modal>

<!-- Nested Modals -->
{#if showAddModal}
    <AddServerModal
            bind:isOpen={showAddModal}
            on:success={handleServerAdded}
    />
{/if}

{#if showEditModal && selectedServer}
    <EditServerModal
            bind:isOpen={showEditModal}
            server={selectedServer}
            on:success={handleServerUpdated}
    />
{/if}

{#if showDeleteModal && selectedServer}
    <DeleteServerModal
            bind:isOpen={showDeleteModal}
            server={selectedServer}
            on:success={handleServerDeleted}
    />
{/if}
