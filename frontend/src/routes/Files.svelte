<script lang="ts">
    import {onMount} from 'svelte';
    import {
        CloseSession,
        ConnectSession,
        GetListFiles,
        GetListSession,
        GetUserHomeDir,
        ListFiles,
        ListServers
    } from '../../wailsjs/go/main/App';
    import {dto} from '../../wailsjs/go/models';

    let currentPath = "";
    let files: dto.FileInfo[] = [];
    let servers: dto.ServerResponse[] = [];
    let activeSessionId: string = 'local';
    let isLoading = false;
    let userHomeDir = '';
    let showServerModal = false;
    let viewMode: 'grid' | 'list' = 'list';
    let sessionStatuses: Record<string, boolean> = {};
    let isEditingBreadcrumb: boolean = false;

    let errorMessage = '';
    let errorType: 'error' | 'warning' | 'info' = 'error';
    let showError = false;
    let errorTimeout: number;

    // Safe reactive statement with fallback
    $: pathParts = currentPath ? currentPath.split('/').filter(Boolean) : [];

    // Safe computed values
    $: safeFiles = Array.isArray(files) ? files : [];
    $: safeServers = Array.isArray(servers) ? servers : [];

    const getFileIcon = (file: dto.FileInfo): string => {
        if (file.isDir) return 'ri-folder-line';
        const ext = file.name.split('.').pop()?.toLowerCase() || '';

        const iconMap = new Map([
            ['jpg', 'ri-image-line'], ['jpeg', 'ri-image-line'], ['png', 'ri-image-line'], ['gif', 'ri-image-line'],
            ['mp3', 'ri-music-line'], ['wav', 'ri-music-line'], ['mp4', 'ri-video-line'], ['mov', 'ri-video-line'],
            ['pdf', 'ri-file-pdf-line'], ['zip', 'ri-file-zip-line'], ['rar', 'ri-file-zip-line'],
            ['doc', 'ri-file-text-line'], ['docx', 'ri-file-text-line'], ['txt', 'ri-file-text-line'],
            ['js', 'ri-code-line'], ['ts', 'ri-code-line'], ['py', 'ri-code-line'], ['go', 'ri-code-line']
        ]);

        return iconMap.get(ext) || 'ri-file-line';
    };

    const formatFileSize = (size: number): string => {
        if (size === 0) return '0 B';
        const units = ['B', 'KB', 'MB', 'GB'];
        let index = 0;
        while (size >= 1024 && index < 3) {
            size /= 1024;
            index++;
        }
        return `${size.toFixed(index === 0 ? 0 : 1)} ${units[index]}`;
    };

    const formatDate = (date: Date): string => {
        return new Intl.DateTimeFormat('en-US', {
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        }).format(date);
    };

    onMount(() => {
        const init = async () => {
            try {
                userHomeDir = await GetUserHomeDir();
                currentPath = userHomeDir;
                await loadFiles();
                await loadServers();
                await updateSessionStatuses();
            } catch (error) {
                console.error('Error during initialization:', error);
                // Ensure arrays are still initialized even on error
                files = [];
                servers = [];
            }
        };

        init();
    });

    function showErrorMessage(message: string, type: 'error' | 'warning' | 'info' = 'error', duration: number = 5000) {
        errorMessage = message;
        errorType = type;
        showError = true;

        // Clear existing timeout
        if (errorTimeout) {
            clearTimeout(errorTimeout);
        }

        // Auto-hide after duration
        errorTimeout = setTimeout(() => {
            hideError();
        }, duration);
    }

    function hideError() {
        showError = false;
        if (errorTimeout) {
            clearTimeout(errorTimeout);
        }
    }

    function getErrorIcon(type: 'error' | 'warning' | 'info'): string {
        switch (type) {
            case 'error': return 'ri-error-warning-line';
            case 'warning': return 'ri-alert-line';
            case 'info': return 'ri-information-line';
            default: return 'ri-error-warning-line';
        }
    }

    function getErrorClasses(type: 'error' | 'warning' | 'info'): string {
        switch (type) {
            case 'error': return 'bg-red-900/50 border-red-500 text-red-400';
            case 'warning': return 'bg-yellow-900/50 border-yellow-500 text-yellow-400';
            case 'info': return 'bg-blue-900/50 border-blue-500 text-blue-400';
            default: return 'bg-red-900/50 border-red-500 text-red-400';
        }
    }

    async function navigateToDirectory(file: dto.FileInfo) {
        if (!file.isDir) return;

        currentPath = activeSessionId === 'local'
            ? (currentPath.endsWith('/') ? `${currentPath}${file.name}` : `${currentPath}/${file.name}`)
            : (currentPath === '/' ? `/${file.name}` : `${currentPath}/${file.name}`);

        await loadFiles();
    }

    async function navigateToPath(index: number) {
        const parts = pathParts; // Use safe pathParts

        if (index === -1) {
            currentPath = activeSessionId === 'local' ? userHomeDir : '/';
        } else {
            const selectedParts = parts.slice(0, index + 1);
            currentPath = '/' + selectedParts.join('/');
        }

        await loadFiles();
    }

    async function navigateUp() {
        const isAtRoot = currentPath === (activeSessionId === 'local' ? userHomeDir : '/');
        if (isAtRoot) return;

        const parts = currentPath.split('/').filter(Boolean);
        parts.pop();
        currentPath = parts.length ? '/' + parts.join('/') : '/';

        if (activeSessionId === 'local' && currentPath === '/') {
            currentPath = userHomeDir;
        }

        await loadFiles();
    }

    async function loadFiles() {
        isLoading = true;
        isEditingBreadcrumb = false;
        hideError();
        try {
            const result = activeSessionId === 'local'
                ? await ListFiles(currentPath)
                : await GetListFiles(activeSessionId, currentPath);

            // Ensure a result is always an array
            files = Array.isArray(result) ? result : [];
        } catch (error) {
            console.error('Error loading files:', error);
            files = []; // Always ensure files is an array

            const errorMsg = error?.toString() || 'An unexpected error occurred';
            showErrorMessage(`${errorMsg}`, 'error');

            if (activeSessionId !== 'local' && error.includes('session')) {
                console.log('Session error detected, switching to local');
                activeSessionId = 'local';
                currentPath = userHomeDir;
                showErrorMessage('Session expired. Switched to local files.', 'warning');
                await loadFiles();
                await updateSessionStatuses();
            }
        } finally {
            isLoading = false;
        }
    }

    async function loadServers() {
        try {
            const result = await ListServers();
            servers = Array.isArray(result) ? result : []; // Ensure servers is always an array
            await updateSessionStatuses();
        } catch (error) {
            console.error('Error loading servers:', error);
            servers = []; // Always ensure servers is an array
            showErrorMessage('Failed to load servers. Please try again.', 'error');
        }
    }

    async function connectToServer(server: dto.ServerResponse) {
        try {
            isLoading = true;
            showServerModal = false;
            hideError();

            const serverSession: dto.ServerSessionManagement = {
                id: server.id,
                name: server.name,
                address: server.address,
                port: server.port,
                user: server.user,
                password: server.password,
                type: server.type
            };

            activeSessionId = await ConnectSession(serverSession);
            currentPath = '/';
            await loadFiles();
            await updateSessionStatuses();
            showErrorMessage(`Connected to ${server.name}`, 'info', 3000);
        } catch (error) {
            console.error('Error connecting to server:', error);
            const errorMsg = error?.toString() || 'Connection failed';
            showErrorMessage(`${errorMsg}`, 'error');
        } finally {
            isLoading = false;
        }
    }

    async function disconnectFromServer() {
        try {
            await CloseSession(activeSessionId);
            activeSessionId = 'local';
            currentPath = userHomeDir;
            await loadFiles();
            await updateSessionStatuses();
            showErrorMessage('Disconnected from server', 'info', 3000);
        } catch (error) {
            console.error('Error disconnecting from server:', error);
            showErrorMessage('Failed to disconnect from server', 'error');
        }
    }

    async function updateSessionStatuses() {
        try {
            const newStatuses: Record<string, boolean> = {};

            // Safely iterate over servers
            for (const server of safeServers) {
                newStatuses[server.id] = false;
            }

            const activeSessions = await GetListSession();

            if (activeSessions && Array.isArray(activeSessions)) {
                for (const session of activeSessions) {
                    if (session.server_info) {
                        newStatuses[session.server_info.id] = true;
                    }
                }
            }

            sessionStatuses = newStatuses;
        } catch (error) {
            console.error('Error updating session statuses:', error);
            sessionStatuses = {}; // Reset to empty objects on error
            showErrorMessage('Failed to update connection status', 'warning');
        }
    }

    async function validateActiveSession() {
        if (activeSessionId === 'local') return true;

        try {
            const activeSessions = await GetListSession();
            return Array.isArray(activeSessions) && activeSessions.some(session => session.id === activeSessionId);
        } catch (error) {
            console.error('Error validating active session:', error);
            return false;
        }
    }
</script>

<div class="flex flex-col h-screen bg-slate-900">
    <!-- Simplified Header -->
    <header class="flex items-center justify-between px-4 py-4 bg-slat-900 shadow-md border-b border-slate-900">
        <!-- Left controls -->
        <div class="flex items-center gap-3">
            <button
                    on:click={async () => {
                        await updateSessionStatuses();
                        showServerModal = true;
                    }}
                    class="px-3 py-1 bg-blue-600 hover:bg-blue-700 text-white text-xs rounded-md transition-colors flex items-center gap-1"
            >
                <i class="ri-server-line text-sm"></i>
                Servers
            </button>
            <button
                    on:click={navigateUp}
                    disabled={currentPath === (activeSessionId === 'local' ? userHomeDir : '/')}
                    class="px-3 py-1 bg-slate-600 hover:bg-slate-500 disabled:bg-slate-700 disabled:cursor-not-allowed text-white text-xs rounded-md transition-colors flex items-center gap-1"
            >
                <i class="ri-arrow-up-line text-sm"></i>
                Up
            </button>
            <button
                    on:click={async () => {
                        if (!(await validateActiveSession())) {
                            activeSessionId = 'local';
                            currentPath = userHomeDir;
                            await updateSessionStatuses();
                        }
                        await loadFiles();
                    }}
                    class="px-3 py-1 bg-green-600 hover:bg-green-700 text-white text-xs rounded-md transition-colors flex items-center gap-1"
                    disabled={isLoading}
            >
                <i class="ri-refresh-line {isLoading ? 'animate-spin' : ''} text-sm"></i>
                Refresh
            </button>
        </div>

        <!-- Breadcrumb Navigation -->
        <div class="flex-1 mx-4 overflow-x-auto">
            {#if isEditingBreadcrumb}
                <div class="flex items-center px-3 py-1 bg-slate-700 rounded-md text-xs font-mono">
                    <i class="ri-home-line mr-2 text-slate-400"></i>
                    <input
                            type="text"
                            class="bg-transparent border-none outline-none text-blue-400 flex-1"
                            bind:value={currentPath}
                            on:keydown={(e) => {
                        if (e.key === 'Enter') {
                            loadFiles();
                            isEditingBreadcrumb = false;
                        }
                    }}
                            on:blur={() => isEditingBreadcrumb = false}
                    />
                </div>
            {:else}
                <div
                        class="flex items-center px-3 py-1 bg-slate-700 rounded-md text-xs font-mono whitespace-nowrap cursor-pointer"
                        on:click={() => isEditingBreadcrumb = true}
                        on:keydown={(e) => {
                            if (e.key === 'Enter') {
                                loadFiles();
                            }
                        }}
                >
                    <i class="ri-home-line mr-2 text-slate-400"></i>
                    {#each pathParts as part, index}
                        <button
                                class="text-blue-400 hover:text-blue-300"
                                on:click={() => navigateToPath(index)}
                        >
                            {part}
                        </button>
                        {#if index < pathParts.length - 1}
                            <span class="mx-1 text-slate-400">/</span>
                        {/if}
                    {/each}
                </div>
            {/if}
        </div>

        <!-- Right controls -->
        <div class="flex items-center gap-3">
            <div class="flex bg-slate-700 rounded-md p-0.5">
                <button
                        on:click={() => viewMode = 'list'}
                        class="p-2 rounded text-xs {viewMode === 'list' ? 'bg-blue-600 text-white' : 'text-slate-400'}"
                >
                    <i class="ri-list-check"></i>
                </button>
                <button
                        on:click={() => viewMode = 'grid'}
                        class="p-2 rounded text-xs {viewMode === 'grid' ? 'bg-blue-600 text-white' : 'text-slate-400'}"
                >
                    <i class="ri-grid-line"></i>
                </button>
            </div>
            {#if activeSessionId !== 'local'}
                <div class="flex items-center gap-2 px-3 py-1 bg-green-900/50 text-green-400 text-xs rounded-md">
                    <i class="ri-wifi-line"></i>
                    <span>Connected</span>
                </div>
                <button
                        on:click={disconnectFromServer}
                        class="px-3 py-1 bg-red-600 hover:bg-red-700 text-white text-xs rounded-md transition-colors flex items-center gap-1"
                >
                    <i class="ri-logout-circle-line text-sm"></i>
                    Disconnect
                </button>
            {/if}
        </div>
    </header>

    <!-- Error Message -->
    {#if showError}
        <div class="px-4 py-2">
            <div class="flex items-center justify-between p-3 rounded-md border {getErrorClasses(errorType)} animate-slide-down">
                <div class="flex items-center gap-3">
                    <i class="{getErrorIcon(errorType)} text-lg"></i>
                    <span class="text-sm font-medium">{errorMessage}</span>
                </div>
                <button
                        on:click={hideError}
                        class="text-current hover:opacity-70 transition-opacity"
                >
                    <i class="ri-close-line text-lg"></i>
                </button>
            </div>
        </div>
    {/if}

    <!-- Main content -->
    <main class="flex-1 overflow-hidden px-4">
        {#if isLoading}
            <div class="flex items-center justify-center h-full">
                <div class="flex items-center gap-2 text-slate-400">
                    <i class="ri-loader-4-line animate-spin text-xl"></i>
                    <span>Loading...</span>
                </div>
            </div>
        {:else}
            <div class="h-full bg-slate-800 rounded-lg overflow-hidden border border-slate-700">
                {#if viewMode === 'list'}
                    <!-- List View -->
                    <div class="h-full overflow-auto">
                        <table class="w-full text-sm">
                            <thead class="bg-slate-700 sticky top-0">
                            <tr>
                                <th class="text-left p-2 font-medium text-slate-300">Name</th>
                                <th class="text-left p-2 font-medium text-slate-300 w-20">Size</th>
                                <th class="text-left p-2 font-medium text-slate-300 w-32">Modified</th>
                            </tr>
                            </thead>
                            <tbody>
                            <!-- Navigation Up Row -->
                            {#if currentPath !== (activeSessionId === 'local' ? userHomeDir : '/')}
                                <tr
                                        class="border-b border-slate-700 hover:bg-slate-700/50 cursor-pointer"
                                        on:click={navigateUp}
                                >
                                    <td class="p-2 text-blue-400 hover:text-blue-300 flex items-center gap-2">
                                        <i class="ri-arrow-up-line"></i>
                                        <span>..</span>
                                    </td>
                                    <td class="p-2 text-slate-400"></td>
                                    <td class="p-2 text-slate-400"></td>
                                </tr>
                            {/if}

                            <!-- File Rows -->
                            {#each safeFiles as file}
                                <tr
                                        class="border-b border-slate-700 hover:bg-slate-700/50 cursor-pointer"
                                        on:click={() => navigateToDirectory(file)}
                                >
                                    <td class="p-2 flex items-center gap-2">
                                        <i class="{getFileIcon(file)} {file.isDir ? 'text-blue-400' : 'text-slate-400'}"></i>
                                        <span>{file.name}</span>
                                    </td>
                                    <td class="p-2 text-slate-400 text-left">
                                        {file.isDir ? '-' : formatFileSize(file.size)}
                                    </td>
                                    <td class="p-2 text-slate-400 text-left">
                                        {formatDate(new Date(file.modTime))}
                                    </td>
                                </tr>
                            {/each}
                            </tbody>
                        </table>
                    </div>
                {:else}
                    <!-- Grid View -->
                    <div class="h-full overflow-auto p-4">
                        <div class="grid grid-cols-6 sm:grid-cols-8 md:grid-cols-10 lg:grid-cols-12 gap-3">
                            <!-- Navigation Up Grid -->
                            {#if currentPath !== (activeSessionId === 'local' ? userHomeDir : '/')}
                                <div
                                        class="flex flex-col items-center p-2 rounded-md hover:bg-slate-700/50 cursor-pointer"
                                        on:click={navigateUp}
                                        on:keydown={(e) => e.key === 'Enter' && navigateUp()}
                                        tabindex="0"
                                        role="button"
                                >
                                    <i class="ri-arrow-up-line text-2xl text-blue-400 mb-1"></i>
                                    <span class="text-xs text-center truncate w-full" title="Up">
                                        ..
                                    </span>
                                </div>
                            {/if}

                            <!-- File/Folder Grids -->
                            {#each safeFiles as file}
                                <div
                                        class="flex flex-col items-center p-2 rounded-md hover:bg-slate-700/50 cursor-pointer"
                                        on:click={() => navigateToDirectory(file)}
                                        on:keydown={(e) => e.key === 'Enter' && navigateToDirectory(file)}
                                        tabindex="0"
                                        role="button"
                                >
                                    <i class="{getFileIcon(file)} text-2xl {file.isDir ? 'text-blue-400' : 'text-slate-400'} mb-1"></i>
                                    <span class="text-xs text-center truncate w-full" title={file.name}>
                                        {file.name}
                                    </span>
                                </div>
                            {/each}
                        </div>
                    </div>
                {/if}

                {#if safeFiles.length === 0}
                    <div class="flex flex-col items-center justify-center h-full text-slate-500">
                        <i class="ri-folder-open-line text-4xl mb-2"></i>
                        <p>Empty folder</p>
                    </div>
                {/if}
            </div>
        {/if}
    </main>
</div>

<!-- Server Modal -->
{#if showServerModal}
    <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-slate-800 rounded-lg shadow-xl max-w-md w-full border border-slate-700">
            <div class="flex items-center justify-between p-4 border-b border-slate-700">
                <h2 class="text-lg font-semibold">Select Server</h2>
                <button
                        on:click={() => showServerModal = false}
                        class="text-slate-400 hover:text-slate-200"
                >
                    <i class="ri-close-line text-xl"></i>
                </button>
            </div>

            <div class="p-4 space-y-2 max-h-80 overflow-y-auto">
                <!-- Local option -->
                <button
                        class="w-full p-3 text-left rounded-md border transition-colors {activeSessionId === 'local' ? 'bg-blue-600/20 border-blue-500' : 'border-slate-600 hover:bg-slate-700'}"
                        on:click={() => {
                        activeSessionId = 'local';
                        currentPath = userHomeDir;
                        loadFiles();
                        showServerModal = false;
                    }}
                >
                    <span class="flex items-center gap-3">
                        <i class="ri-computer-line text-slate-400"></i>
                        <span>
                            <span class="font-medium">Local System</span>
                            <span class="text-sm text-slate-400">Browse local files</span>
                        </span>
                    </span>
                </button>

                <!-- Remote servers -->
                {#each safeServers as server}
                    <button
                            class="w-full p-3 text-left rounded-md border transition-colors {sessionStatuses[server.id] ? 'bg-green-600/20 border-green-500' : 'border-slate-600 hover:bg-slate-700'}"
                            on:click={() => connectToServer(server)}
                    >
                        <span class="flex items-center gap-3">
                            <i class="ri-server-line text-slate-400"></i>
                            <span class="flex-1">
                                <span class="flex items-center gap-2">
                                    <span class="font-medium">{server.name}</span>
                                    {#if sessionStatuses[server.id]}
                                        <i class="ri-wifi-line text-green-400 text-sm" title="Connected"></i>
                                    {:else}
                                        <i class="ri-wifi-off-line text-slate-500 text-sm" title="Disconnected"></i>
                                    {/if}
                                </span>
                                <span class="text-sm text-slate-400">{server.address}:{server.port}</span>
                            </span>
                        </span>
                    </button>
                {/each}

                {#if safeServers.length === 0}
                    <div class="text-center py-8 text-slate-400">
                        <i class="ri-server-line text-3xl mb-2"></i>
                        <p>No servers configured</p>
                    </div>
                {/if}
            </div>
        </div>
    </div>
{/if}

<style>
    :global(.overflow-auto::-webkit-scrollbar) {
        width: 6px;
        height: 6px;
    }

    :global(.overflow-auto::-webkit-scrollbar-track) {
        background: transparent;
    }

    :global(.overflow-auto::-webkit-scrollbar-thumb) {
        background: rgba(148, 163, 184, 0.3);
        border-radius: 3px;
    }
</style>