<script lang="ts">
    import { onMount } from 'svelte';
    import { GetListFiles, ListFiles, GetUserHomeDir, GetSession, ConnectSession, CloseSession, ListServers } from '../../wailsjs/go/main/App';
    import { dto } from '../../wailsjs/go/models';

    let currentPath = "";
    let files: dto.FileInfo[] = [];
    let servers: dto.ServerResponse[] = [];
    let activeSessionId: string = 'local';
    let isLoading = false;
    let userHomeDir = '';
    let showServerModal = false;
    let viewMode: 'grid' | 'list' = 'grid';

    const icons = {
        folder: 'ri-folder-3-fill',
        folderOpen: 'ri-folder-open-fill',
        file: 'ri-file-3-line',
        image: 'ri-image-2-fill',
        document: 'ri-file-text-fill',
        code: 'ri-code-s-slash-fill',
        server: 'ri-server-fill',
        audio: 'ri-music-2-fill',
        video: 'ri-video-fill',
        pdf: 'ri-file-pdf-2-fill',
        archive: 'ri-file-zip-fill',
        up: 'ri-arrow-up-line',
        close: 'ri-close-line',
        disconnect: 'ri-logout-circle-line',
        computer: 'ri-computer-line',
        check: 'ri-check-circle-fill',
        path: 'ri-home-4-line',
        connected: 'ri-wifi-line',
        grid: 'ri-grid-fill',
        list: 'ri-list-check-2',
        refresh: 'ri-refresh-line',
        cloud: 'ri-cloud-line'
    };

    onMount(async () => {
        userHomeDir = await GetUserHomeDir();
        currentPath = userHomeDir;
        await loadFiles();
        await loadServers();
    });

    function getFileIcon(file: dto.FileInfo): string {
        if (file.isDir) return icons.folder;
        const ext = file.name.split('.').pop()?.toLowerCase();

        if (['jpg', 'jpeg', 'png', 'gif', 'svg', 'webp', 'bmp', 'ico'].includes(ext)) return icons.image;
        if (['mp3', 'wav', 'ogg', 'm4a', 'flac', 'aac'].includes(ext)) return icons.audio;
        if (['mp4', 'mov', 'avi', 'mkv', 'webm', 'flv', 'm4v'].includes(ext)) return icons.video;
        if (['pdf'].includes(ext)) return icons.pdf;
        if (['zip', 'rar', '7z', 'tar', 'gz', 'bz2', 'xz'].includes(ext)) return icons.archive;
        if (['doc', 'docx', 'txt', 'md', 'rtf', 'odt'].includes(ext)) return icons.document;
        if (['js', 'ts', 'py', 'go', 'java', 'cpp', 'c', 'html', 'css', 'scss', 'json', 'xml', 'sql'].includes(ext)) return icons.code;
        return icons.file;
    }

    function getFileColor(file: dto.FileInfo): string {
        if (file.isDir) return 'text-blue-400';
        const ext = file.name.split('.').pop()?.toLowerCase();

        if (['jpg', 'jpeg', 'png', 'gif', 'svg', 'webp', 'bmp', 'ico'].includes(ext)) return 'text-green-400';
        if (['mp3', 'wav', 'ogg', 'm4a', 'flac', 'aac'].includes(ext)) return 'text-purple-400';
        if (['mp4', 'mov', 'avi', 'mkv', 'webm', 'flv', 'm4v'].includes(ext)) return 'text-red-400';
        if (['pdf'].includes(ext)) return 'text-red-500';
        if (['zip', 'rar', '7z', 'tar', 'gz', 'bz2', 'xz'].includes(ext)) return 'text-yellow-400';
        if (['doc', 'docx', 'txt', 'md', 'rtf', 'odt'].includes(ext)) return 'text-blue-300';
        if (['js', 'ts', 'py', 'go', 'java', 'cpp', 'c', 'html', 'css', 'scss', 'json', 'xml', 'sql'].includes(ext)) return 'text-orange-400';
        return 'text-slate-400';
    }

    function formatFileSize(size: number): string {
        if (size === 0) return '0 B';
        const units = ['B', 'KB', 'MB', 'GB', 'TB'];
        let index = 0;
        while (size >= 1024 && index < units.length - 1) {
            size /= 1024;
            index++;
        }
        return `${size.toFixed(index === 0 ? 0 : 1)} ${units[index]}`;
    }

    function formatDate(date: Date): string {
        const now = new Date();
        const diffMs = now.getTime() - date.getTime();
        const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

        if (diffDays === 0) {
            return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
        } else if (diffDays === 1) {
            return 'Yesterday';
        } else if (diffDays < 7) {
            return `${diffDays} days ago`;
        } else {
            return date.toLocaleDateString();
        }
    }

    async function navigateToDirectory(file: dto.FileInfo) {
        if (!file.isDir) return;

        if (activeSessionId === 'local') {
            currentPath = currentPath.endsWith('/') ?
                `${currentPath}${file.name}` :
                `${currentPath}/${file.name}`;
        } else {
            currentPath = currentPath === '/' ?
                `/${file.name}` :
                `${currentPath}/${file.name}`;
        }

        await loadFiles();
    }

    async function navigateToPath(index: number, parts: string[]) {
        if (activeSessionId === 'local') {
            if (index === -1) {
                currentPath = userHomeDir;
            } else {
                const selectedParts = parts.slice(0, index + 1);
                currentPath = '/' + selectedParts.join('/');
            }
        } else {
            if (index === -1) {
                currentPath = '/';
            } else {
                const selectedParts = parts.slice(0, index + 1);
                currentPath = '/' + selectedParts.join('/');
            }
        }
        await loadFiles();
    }

    async function navigateUp() {
        if (currentPath === (activeSessionId === 'local' ? userHomeDir : '/')) return;

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
        try {
            if (activeSessionId === 'local') {
                files = await ListFiles(currentPath);
            } else {
                files = await GetListFiles(activeSessionId, currentPath);
            }
            files.sort((a, b) => {
                if (a.isDir === b.isDir) {
                    return a.name.localeCompare(b.name);
                }
                return a.isDir ? -1 : 1;
            });
        } catch (error) {
            console.error('Error loading files:', error);
            files = [];
        } finally {
            isLoading = false;
        }
    }

    async function loadServers() {
        try {
            servers = await ListServers();
        } catch (error) {
            console.error('Error loading servers:', error);
            servers = [];
        }
    }

    async function connectToServer(server: dto.ServerResponse) {
        try {
            isLoading = true;
            activeSessionId = await ConnectSession(server);
            currentPath = '/';
            await loadFiles();
            showServerModal = false;
        } catch (error) {
            console.error('Error connecting to server:', error);
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
        } catch (error) {
            console.error('Error disconnecting from server:', error);
        }
    }
</script>

<div class="flex flex-col h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
    <!-- Enhanced Header -->
    <div class="flex items-center justify-between p-4 bg-slate-800/90 backdrop-blur-sm border-b border-slate-700/50 shadow-xl">
        <!-- Left Section -->
        <div class="flex items-center space-x-3">
            <!-- Server Connection Button -->
            <button
                    on:click={() => showServerModal = true}
                    class="group flex items-center gap-2 px-4 py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-sm font-medium text-white rounded-lg hover:from-indigo-700 hover:to-purple-700 transition-all duration-300 shadow-lg hover:shadow-indigo-500/25 hover:scale-105"
            >
                <i class="{icons.server} text-base group-hover:animate-pulse"></i>
                <span>Servers</span>
            </button>

            <!-- Navigation Up Button -->
            <button
                    on:click={navigateUp}
                    disabled={currentPath === (activeSessionId === 'local' ? userHomeDir : '/')}
                    class="flex items-center gap-2 px-4 py-2.5 bg-gradient-to-r from-blue-600 to-cyan-600 text-sm font-medium text-white rounded-lg hover:from-blue-700 hover:to-cyan-700 disabled:from-slate-600 disabled:to-slate-600 disabled:cursor-not-allowed transition-all duration-300 shadow-lg hover:shadow-blue-500/25 hover:scale-105 disabled:hover:scale-100"
            >
                <i class="{icons.up} text-base"></i>
                <span>Up</span>
            </button>

            <!-- Refresh Button -->
            <button
                    on:click={loadFiles}
                    class="flex items-center justify-center w-10 h-10 bg-gradient-to-r from-emerald-600 to-teal-600 text-white rounded-lg hover:from-emerald-700 hover:to-teal-700 transition-all duration-300 shadow-lg hover:shadow-emerald-500/25 hover:scale-105"
            >
                <i class="{icons.refresh} text-base {isLoading ? 'animate-spin' : ''}"></i>
            </button>
        </div>

        <!-- Center - Breadcrumb -->
        <div class="flex-1 mx-6">
            <div class="flex items-center px-4 py-2.5 bg-slate-700/50 backdrop-blur-sm rounded-lg shadow-inner border border-slate-600/50 overflow-x-auto">
                <i class="{icons.path} text-slate-400 mr-3 text-base flex-shrink-0"></i>
                <div class="flex items-center gap-1 font-mono text-sm">
                    {#if activeSessionId === 'local'}
                        <button
                                class="text-blue-400 hover:text-blue-300 transition-colors px-2 py-1 rounded-md hover:bg-slate-600/50 flex items-center gap-1"
                                on:click={() => navigateToPath(-1, [])}
                        >
                            <i class="ri-home-4-line text-xs"></i>
                            Home
                        </button>
                    {:else}
                        <button
                                class="text-blue-400 hover:text-blue-300 transition-colors px-2 py-1 rounded-md hover:bg-slate-600/50"
                                on:click={() => navigateToPath(-1, [])}
                        >
                            /
                        </button>
                    {/if}

                    {#each currentPath.split('/').filter(Boolean) as part, index}
                        <span class="text-slate-500">/</span>
                        <button
                                class="text-blue-400 hover:text-blue-300 transition-colors px-2 py-1 rounded-md hover:bg-slate-600/50 max-w-[150px] truncate"
                                on:click={() => navigateToPath(index, currentPath.split('/').filter(Boolean))}
                        >
                            {part}
                        </button>
                    {/each}
                </div>
            </div>
        </div>

        <!-- Right Section -->
        <div class="flex items-center gap-3">
            <!-- View Mode Toggle -->
            <div class="flex items-center bg-slate-700/50 rounded-lg p-1">
                <button
                        on:click={() => viewMode = 'grid'}
                        class="p-2 rounded-md transition-all duration-200 {viewMode === 'grid' ? 'bg-blue-600 text-white shadow-sm' : 'text-slate-400 hover:text-slate-200'}"
                >
                    <i class="{icons.grid} text-sm"></i>
                </button>
                <button
                        on:click={() => viewMode = 'list'}
                        class="p-2 rounded-md transition-all duration-200 {viewMode === 'list' ? 'bg-blue-600 text-white shadow-sm' : 'text-slate-400 hover:text-slate-200'}"
                >
                    <i class="{icons.list} text-sm"></i>
                </button>
            </div>

            <!-- Connection Status -->
            {#if activeSessionId !== 'local'}
                <div class="flex items-center gap-2 px-3 py-2 bg-gradient-to-r from-green-600/20 to-emerald-600/20 text-green-400 rounded-lg border border-green-500/30">
                    <i class="{icons.connected} animate-pulse text-base"></i>
                    <span class="text-sm font-medium">Connected</span>
                </div>
                <button
                        on:click={disconnectFromServer}
                        class="flex items-center gap-2 px-4 py-2.5 bg-gradient-to-r from-red-600 to-rose-600 text-sm font-medium text-white rounded-lg hover:from-red-700 hover:to-rose-700 transition-all duration-300 shadow-lg hover:shadow-red-500/25 hover:scale-105"
                >
                    <i class="{icons.disconnect} text-base"></i>
                    <span>Disconnect</span>
                </button>
            {/if}
        </div>
    </div>

    <!-- Main Content -->
    <div class="flex-1 overflow-hidden p-6">
        {#if isLoading}
            <div class="flex flex-col justify-center items-center h-full">
                <div class="relative">
                    <div class="w-16 h-16 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin"></div>
                    <div class="absolute inset-0 w-16 h-16 border-4 border-transparent border-t-purple-600 rounded-full animate-spin" style="animation-delay: -0.15s; animation-duration: 1.5s;"></div>
                </div>
                <p class="mt-4 text-slate-400 text-sm">Loading files...</p>
            </div>
        {:else}
            <div class="h-full bg-slate-800/50 backdrop-blur-sm rounded-2xl shadow-2xl border border-slate-700/50 overflow-hidden">
                {#if viewMode === 'grid'}
                    <!-- Grid View -->
                    <div class="h-full overflow-auto p-6">
                        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 gap-4">
                            {#each files as file}
                                <div
                                        class="group flex flex-col items-center p-4 rounded-xl bg-slate-700/30 hover:bg-slate-700/60 backdrop-blur-sm border border-slate-600/30 hover:border-slate-500/50 cursor-pointer transition-all duration-300 hover:scale-105 hover:shadow-lg hover:shadow-slate-900/50"
                                        on:click={() => navigateToDirectory(file)}
                                >
                                    <div class="mb-3 p-3 rounded-full bg-slate-600/50 group-hover:bg-slate-600/70 transition-colors duration-300">
                                        <i class="{getFileIcon(file)} text-3xl {getFileColor(file)}"></i>
                                    </div>
                                    <span class="text-sm text-slate-200 text-center font-medium truncate w-full" title={file.name}>
                                        {file.name}
                                    </span>
                                    <span class="text-xs text-slate-400 mt-1">
                                        {file.isDir ? 'Folder' : formatFileSize(file.size)}
                                    </span>
                                </div>
                            {/each}
                        </div>

                        {#if files.length === 0}
                            <div class="flex flex-col items-center justify-center h-full text-slate-400">
                                <i class="ri-folder-open-line text-6xl mb-4 opacity-50"></i>
                                <p class="text-lg font-medium">Empty folder</p>
                                <p class="text-sm">No files or folders found</p>
                            </div>
                        {/if}
                    </div>
                {:else}
                    <!-- List View -->
                    <div class="h-full overflow-auto">
                        <table class="w-full">
                            <thead class="sticky top-0 bg-slate-700/80 backdrop-blur-sm border-b border-slate-600/50">
                            <tr>
                                <th class="px-6 py-4 text-left text-xs font-semibold text-slate-300 uppercase tracking-wider">Name</th>
                                <th class="px-6 py-4 text-left text-xs font-semibold text-slate-300 uppercase tracking-wider">Size</th>
                                <th class="px-6 py-4 text-left text-xs font-semibold text-slate-300 uppercase tracking-wider">Modified</th>
                                <th class="px-6 py-4 text-left text-xs font-semibold text-slate-300 uppercase tracking-wider">Type</th>
                            </tr>
                            </thead>
                            <tbody class="divide-y divide-slate-700/50">
                            {#each files as file}
                                <tr
                                        class="group hover:bg-slate-700/40 cursor-pointer transition-all duration-200"
                                        on:click={() => navigateToDirectory(file)}
                                >
                                    <td class="px-6 py-4">
                                        <div class="flex items-center">
                                            <div class="mr-4 p-2 rounded-lg bg-slate-600/30 group-hover:bg-slate-600/50 transition-colors duration-200">
                                                <i class="{getFileIcon(file)} text-lg {getFileColor(file)}"></i>
                                            </div>
                                            <span class="text-slate-200 font-medium group-hover:text-white transition-colors duration-200">
                                                    {file.name}
                                                </span>
                                        </div>
                                    </td>
                                    <td class="px-6 py-4 text-sm text-slate-400">
                                        {file.isDir ? '--' : formatFileSize(file.size)}
                                    </td>
                                    <td class="px-6 py-4 text-sm text-slate-400">
                                        {formatDate(new Date(file.modTime))}
                                    </td>
                                    <td class="px-6 py-4">
                                            <span class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium {file.isDir ? 'bg-blue-600/20 text-blue-400 border border-blue-500/30' : 'bg-slate-600/20 text-slate-400 border border-slate-500/30'}">
                                                {file.isDir ? 'Directory' : 'File'}
                                            </span>
                                    </td>
                                </tr>
                            {/each}
                            </tbody>
                        </table>

                        {#if files.length === 0}
                            <div class="flex flex-col items-center justify-center h-64 text-slate-400">
                                <i class="ri-folder-open-line text-6xl mb-4 opacity-50"></i>
                                <p class="text-lg font-medium">Empty folder</p>
                                <p class="text-sm">No files or folders found</p>
                            </div>
                        {/if}
                    </div>
                {/if}
            </div>
        {/if}
    </div>
</div>

<!-- Enhanced Server Modal -->
{#if showServerModal}
    <div class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
        <div class="bg-gradient-to-br from-slate-800 to-slate-900 p-6 rounded-2xl shadow-2xl max-w-2xl w-full border border-slate-700/50">
            <div class="flex justify-between items-center mb-6">
                <div class="flex items-center gap-3">
                    <div class="p-2 bg-gradient-to-r from-indigo-600 to-purple-600 rounded-lg">
                        <i class="{icons.server} text-xl text-white"></i>
                    </div>
                    <h2 class="text-xl font-bold text-slate-100">Server Connections</h2>
                </div>
                <button
                        on:click={() => showServerModal = false}
                        class="text-slate-400 hover:text-slate-200 transition-colors p-2 hover:bg-slate-700/50 rounded-lg"
                >
                    <i class="{icons.close} text-xl"></i>
                </button>
            </div>

            <div class="space-y-3 max-h-96 overflow-y-auto">
                <!-- Local System Option -->
                <button
                        class="w-full p-4 text-left hover:bg-slate-700/50 rounded-xl transition-all duration-200 flex items-center justify-between group {activeSessionId === 'local' ? 'bg-gradient-to-r from-blue-600/20 to-indigo-600/20 border-2 border-blue-500/50' : 'border-2 border-slate-600/30 hover:border-slate-500/50'}"
                        on:click={() => {
                        activeSessionId = 'local';
                        currentPath = userHomeDir;
                        loadFiles();
                        showServerModal = false;
                    }}
                >
                    <div class="flex items-center gap-4">
                        <div class="p-3 bg-slate-600/50 group-hover:bg-slate-600/70 rounded-lg transition-colors duration-200">
                            <i class="{icons.computer} text-xl text-blue-400"></i>
                        </div>
                        <div>
                            <span class="font-semibold text-slate-200 text-lg">Local System</span>
                            <p class="text-sm text-slate-400">Browse files on this computer</p>
                        </div>
                    </div>
                    {#if activeSessionId === 'local'}
                        <i class="{icons.check} text-green-400 text-xl"></i>
                    {/if}
                </button>

                <!-- Remote Servers -->
                {#each servers as server}
                    <button
                            class="w-full p-4 text-left hover:bg-slate-700/50 rounded-xl transition-all duration-200 flex items-center justify-between group {activeSessionId === server.id ? 'bg-gradient-to-r from-green-600/20 to-emerald-600/20 border-2 border-green-500/50' : 'border-2 border-slate-600/30 hover:border-slate-500/50'}"
                            on:click={() => connectToServer(server)}
                    >
                        <div class="flex items-center gap-4">
                            <div class="p-3 bg-slate-600/50 group-hover:bg-slate-600/70 rounded-lg transition-colors duration-200">
                                <i class="{icons.cloud} text-xl text-green-400"></i>
                            </div>
                            <div>
                                <span class="font-semibold text-slate-200 text-lg">{server.name}</span>
                                <p class="text-sm text-slate-400">{server.address}:{server.port}</p>
                            </div>
                        </div>
                        {#if activeSessionId === server.id}
                            <i class="{icons.check} text-green-400 text-xl"></i>
                        {/if}
                    </button>
                {/each}

                {#if servers.length === 0}
                    <div class="text-center py-8 text-slate-400">
                        <i class="ri-server-line text-4xl mb-2 opacity-50"></i>
                        <p>No remote servers configured</p>
                    </div>
                {/if}
            </div>
        </div>
    </div>
{/if}

<style>
    /* Custom scrollbar for the entire component */
    :global(.overflow-auto)::-webkit-scrollbar {
        width: 8px;
        height: 8px;
    }
    :global(.overflow-auto)::-webkit-scrollbar-track {
        background: rgba(51, 65, 85, 0.2);
        border-radius: 4px;
    }
    :global(.overflow-auto)::-webkit-scrollbar-thumb {
        background: rgba(148, 163, 184, 0.3);
        border-radius: 4px;
    }
    :global(.overflow-auto)::-webkit-scrollbar-thumb:hover {
        background: rgba(148, 163, 184, 0.5);
    }

    .grid > div {
        transform-origin: center;
    }
</style>