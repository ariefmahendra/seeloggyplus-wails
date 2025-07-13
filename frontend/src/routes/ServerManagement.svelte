<script lang="ts">
    import {onMount} from 'svelte';
    import Button from '../components/common/Button.svelte';
    import Input from '../components/common/Input.svelte';
    import Spinner from '../components/common/Spinner.svelte';
    import {
        DeleteServer,
        AddServer,
        ListServers,
        UpdateServer
    } from '../../wailsjs/go/main/App';
    import { dto } from '../../wailsjs/go/models';

    type serverCreateRequest = dto.ServerCreateRequest;
    type serverResponse = dto.ServerResponse;
    type serverUpdateRequest = dto.ServerUpdateRequest;

    let pageStatus: 'loading' | 'loaded' | 'error' = 'loading';
    let servers: serverResponse[] = [];
    let showModal = false;

    let editingServer: serverCreateRequest | serverUpdateRequest | null = null;
    let serverToDelete: serverResponse | null = null;
    let originalServerState: string | null = null;

    let searchTerm = '';
    type ModalStatus = 'idle' | 'saving' | 'deleting' | 'testing';
    let modalStatus: ModalStatus = 'idle';
    let testMessage = '';

    async function refreshServers(){
        try {
            const fetchedServers = await ListServers();
            servers = fetchedServers === null ? [] : fetchedServers;
        } catch (e) {
            console.error('Failed to refresh servers', e);
        }
    }

    onMount(async () => {
        pageStatus = 'loading';
        await refreshServers();
        pageStatus = 'loaded';
    });

    // Reactive variable form validation
    $: isFormValid =
        !!editingServer &&
        editingServer.name?.trim() !== '' &&
        editingServer.address?.trim() !== '' &&
        editingServer.port !== 0 &&
        editingServer.user?.trim() !== '' &&
        editingServer.password?.trim() !== '';

    $: isEditing = !!(editingServer && 'id' in editingServer && editingServer.id);

    $: hasChanges = (() => {
        if (!editingServer) return false;
        if (!isEditing) return true;
        return originalServerState !== JSON.stringify(editingServer);
    })();

    $: filteredServers = servers.filter((server) => {
        if (!searchTerm) return true;
        const lowerSearchTerm = searchTerm.toLowerCase();
        return (
            server.name.toLowerCase().includes(lowerSearchTerm) ||
            server.address.toLowerCase().includes(lowerSearchTerm) ||
            server.user.toLowerCase().includes(lowerSearchTerm)
        );
    });

    async function testConnection() {
        if (!editingServer || modalStatus !== 'idle') return;

        modalStatus = 'testing';
        testMessage = 'Attempting to connect...';

        try {
            await new Promise((resolve) => setTimeout(resolve, 1500));

            if (Math.random() > 0.3) {
                testMessage = '✅ Connection successful!';
            } else {
                console.error('Test connection failed');
            }
        } catch (e: any) {
            testMessage = `❌ ${e.message || 'An unknown error occurred.'}`;
        } finally {
            modalStatus = 'idle';
        }
    }

    async function saveServer() {
        if (!editingServer || modalStatus !== 'idle') return;

        modalStatus = 'saving';
        try {
            if ('id' in editingServer && editingServer.id) {
                await UpdateServer(editingServer as serverUpdateRequest);
            } else {
                await AddServer(editingServer as serverCreateRequest);
            }

            await refreshServers();
            closeModal();
        } catch (e) {
            console.error('Failed to save server', e);
        } finally {
            modalStatus = 'idle';
            editingServer = null;
        }
    }

    async function confirmDelete() {
        if (!serverToDelete || modalStatus !== 'idle') return;

        modalStatus = 'deleting';
        try {
            await DeleteServer(serverToDelete.id);
            await refreshServers();
            closeModal();
        } catch (e){
            console.error('Failed to delete server', e);
        } finally {
            modalStatus = 'idle';
            serverToDelete = null;
        }
    }

    function openAddModal() {
        editingServer = new dto.ServerCreateRequest();
        editingServer.name = '';
        editingServer.address = '';
        editingServer.port = 22;
        editingServer.user = '';
        editingServer.password = '';
        originalServerState = null;
        showModal = true;
    }

    function openEditModal(server: serverResponse) {
        const serverToEdit = new dto.ServerUpdateRequest();
        Object.assign(serverToEdit, server);
        editingServer = serverToEdit;

        originalServerState = JSON.stringify(serverToEdit);
        showModal = true;
    }

    function openDeleteConfirmation(server: serverResponse) {
        serverToDelete = server;
    }

    function closeModal() {
        showModal = false;
        editingServer = null;
        testMessage = '';
        originalServerState = null;
    }
</script>

{#if pageStatus === 'loading'}
    <div class="mx-auto max-w-7xl animate-pulse">
        <div class="mb-6 flex flex-col items-center justify-between gap-4 md:flex-row">
            <div class="h-10 w-64 rounded-md bg-gray-700"/>
            <div class="flex w-full items-center gap-4 md:w-auto">
                <div class="h-12 w-full flex-grow rounded-md bg-gray-700"/>
                <div class="h-12 w-32 flex-shrink-0 rounded-md bg-gray-700"/>
            </div>
        </div>
        <div class="grid grid-cols-1 gap-5 md:grid-cols-2 xl:grid-cols-3">
            {#each Array(3) as _}
                <div class="rounded-lg bg-gray-800 p-4">
                    <div class="mb-4 h-6 w-3/4 rounded bg-gray-700"/>
                    <div class="space-y-3">
                        <div class="h-4 w-1/4 rounded bg-gray-700"/>
                        <div class="h-5 w-1/2 rounded bg-gray-700"/>
                        <div class="mt-2 h-4 w-1/4 rounded bg-gray-700"/>
                        <div class="h-5 w-1/3 rounded bg-gray-700"/>
                    </div>
                </div>
            {/each}
        </div>
    </div>
{:else if pageStatus === 'loaded'}
    <div class="mx-auto max-w-7xl">
        <div class="mb-6 flex flex-col items-center justify-between gap-4 md:flex-row">
            <h1 class="text-3xl font-bold text-white">Server Management</h1>
            <div class="flex w-full items-center gap-4 md:w-auto">
                <div class="relative flex-grow">
                    <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                        <i class="ri-search-line text-gray-400"/>
                    </div>
                    <Input bind:value={searchTerm} placeholder="Search servers..." class="!pl-10"/>
                </div>
                <Button on:click={openAddModal} class="flex-shrink-0">
                    <i class="ri-add-line -ml-1 mr-2"/>
                    Add Server
                </Button>
            </div>
        </div>

        {#if filteredServers.length > 0}
            <div class="grid grid-cols-1 gap-5 md:grid-cols-2 xl:grid-cols-3">
                {#each filteredServers as server (server.id)}
                    <div
                            class="flex flex-col rounded-lg border border-transparent bg-gray-800 p-4 transition-colors duration-200 hover:border-blue-500"
                    >
                        <div class="mb-4 flex items-start justify-between">
                            <h2 class="pr-4 text-lg font-semibold text-white">{server.name}</h2>
                            <div class="flex flex-shrink-0 items-center space-x-2">
                                <button
                                        on:click={() => openEditModal(server)}
                                        class="rounded-full p-1 text-gray-400 transition-colors hover:bg-gray-700 hover:text-blue-400"
                                        title="Edit Server"
                                >
                                    <i class="ri-pencil-line text-base"/>
                                </button>
                                <button
                                        on:click={() => openDeleteConfirmation(server)}
                                        class="rounded-full p-1 text-gray-400 transition-colors hover:bg-gray-700 hover:text-red-400"
                                        title="Delete Server"
                                >
                                    <i class="ri-delete-bin-line text-base"/>
                                </button>
                            </div>
                        </div>
                        <div class="flex-grow space-y-3 text-sm">
                            <div>
                                <span class="block text-gray-400">Host</span>
                                <span class="font-mono text-gray-200">{server.address}:{server.port}</span>
                            </div>
                            <div>
                                <span class="block text-gray-400">User</span>
                                <span class="font-mono text-gray-200">{server.user}</span>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        {:else}
            <div class="rounded-lg bg-gray-800 px-6 py-20 text-center">
                <i class="ri-server-line mb-4 text-6xl text-gray-600"/>
                <h3 class="text-xl font-semibold text-white">
                    {#if servers.length === 0}
                        No Servers Added
                    {:else}
                        No Servers Found
                    {/if}
                </h3>
                <p class="mx-auto mt-2 max-w-sm text-gray-400">
                    {#if servers.length === 0}
                        Get started by adding your first remote server connection.
                    {:else}
                        Your search for "{searchTerm}" did not match any servers. Try a different search
                        term.
                    {/if}
                </p>
            </div>
        {/if}
    </div>
{/if}

<!--Modal For Edit or Create new Server-->
{#if showModal && editingServer}
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-75">
        <div class="w-full max-w-lg rounded-lg bg-gray-800 p-6 shadow-xl">
            <h2 class="mb-6 text-2xl font-bold text-white">
                {'id' in editingServer && editingServer.id? 'Edit Server' : 'Add Server'}
            </h2>
            <form on:submit|preventDefault={saveServer} class="space-y-4">
                <div>
                    <label for="alias" class="mb-2 block text-sm font-medium text-gray-300">Alias / Name</label>
                    <Input id="alias" bind:value={editingServer.name} placeholder="e.g., Production Server" required/>
                </div>
                <div class="grid grid-cols-3 gap-4">
                    <div class="col-span-2">
                        <label for="host" class="mb-2 block text-sm font-medium text-gray-300">Host / IP</label>
                        <Input id="host" bind:value={editingServer.address} placeholder="e.g., 192.168.1.1" required/>
                    </div>
                    <div>
                        <label for="port" class="mb-2 block text-sm font-medium text-gray-300">Port</label>
                        <Input id="port" type="number" bind:value={editingServer.port} required/>
                    </div>
                </div>
                <div>
                    <label for="user" class="mb-2 block text-sm font-medium text-gray-300">Username</label>
                    <Input id="user" bind:value={editingServer.user} placeholder="e.g., root" required/>
                </div>
                <div>
                    <label for="password" class="mb-2 block text-sm font-medium text-gray-300">Password</label>
                    <Input id="password" type="password" bind:value={editingServer.password}/>
                </div>
                <div class="flex items-center justify-end space-x-3 pt-4">
                    <Button
                            type="button"
                            on:click={testConnection}
                            disabled={modalStatus !== 'idle' || !isFormValid}
                            class="bg-gray-700 hover:bg-gray-600 w-48 flex justify-center items-center gap-2 disabled:background-gray-600/50 disabled:cursor-not-allowed"
                    >
                        {#if modalStatus === 'testing'}
                            <Spinner className="w-5 h-5" />
                            <span>Testing...</span>
                        {:else}
                            <i class="ri-plug-line w-5 h-5"/>
                            <span>Test Connection</span>
                        {/if}
                    </Button>
                    <Button type="button" on:click={closeModal} class="!bg-gray-600 hover:!bg-gray-500"
                            disabled={modalStatus !== 'idle'}>Cancel
                    </Button>
                    <Button type="submit"
                            disabled={modalStatus !== 'idle' || !isFormValid || !hasChanges}
                            class="bg-gray-700 hover:bg-gray-600 w-32 flex justify-center items-center gap-2 disabled:background-gray-600/50 disabled:cursor-not-allowed">
                        {#if modalStatus === 'saving'}
                            <Spinner className="w-5 h-5"/>
                            Saving...
                        {:else}
                            Save Server
                        {/if}
                    </Button>
                </div>
                {#if testMessage}
                    <p class="mt-2 text-center text-sm text-gray-400">{testMessage}</p>
                {/if}
            </form>
        </div>
    </div>
{/if}

<!--Modal For Delete Server-->
{#if serverToDelete}
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-75">
        <div class="w-full max-w-md rounded-lg bg-gray-800 p-6 text-center shadow-xl">
            <h2 class="mb-4 text-xl font-semibold text-white">Are you sure?</h2>
            <p class="mb-6 text-gray-400">
                Do you really want to delete the server "<strong class="text-white">{serverToDelete.name}</strong>"?
                This
                action cannot be undone.
            </p>
            <div class="flex justify-center space-x-4">
                <Button on:click={() => (serverToDelete = null)}
                        class="!bg-gray-600 hover:!bg-gray-500"
                        disabled={modalStatus === 'deleting'}>Cancel
                </Button>
                <Button on:click={confirmDelete}
                        class="!bg-red-600 hover:!bg-red-700 w-32 flex justify-center items-center gap-2"
                        disabled={modalStatus === 'deleting'}>
                    {#if modalStatus === 'deleting'}
                        <Spinner/>
                        Deleting...
                    {:else}
                        Delete
                    {/if}
                </Button>
            </div>
        </div>
    </div>
{/if}