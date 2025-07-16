<script lang="ts">
    import { onMount } from 'svelte';
    import Button from '../components/common/Button.svelte';
    import Input from '../components/common/Input.svelte';

    interface Server {
        id: string;
        name: string;
        host: string;
        port: number;
        user: string;
        password?: string;
        defaultPath: string;
    }

    interface FileItem {
        name: string;
        path: string;
        type: 'file' | 'directory';
        size?: string;
        modified?: string;
    }

    interface Session {
        serverId: string;
        status: 'disconnected' | 'connecting' | 'connected' | 'error';
        currentPath: string;
        files: FileItem[];
        error?: string;
        isLoading: boolean;
    }

    let servers: Server[] = [];
    let sessions = new Map<string, Session>();
    let activeSessionId: string | null = null;

    $: activeSession = activeSessionId ? sessions.get(activeSessionId) : null;
    $: activeServer = activeSession ? servers.find(s => s.id === activeSession.serverId) : null;

    onMount(() => {
        servers = [
            { id: 'srv-local', name: 'Local Machine', host: 'localhost', port: 0, user: '', defaultPath: 'C:\\ProgramData' },
            { id: 'srv-prod-web', name: 'Production Web', host: '192.168.1.100', port: 22, user: 'admin', defaultPath: '/var/log/nginx' },
            { id: 'srv-dev-db', name: 'Development DB', host: 'dev.database.server', port: 22, user: 'devuser', defaultPath: '/var/log/mysql' },
        ];
        if (servers.length > 0) {
            selectServer(servers[0].id);
        }
    });

    function selectServer(serverId: string) {
        if (!sessions.has(serverId)) {
            const server = servers.find(s => s.id === serverId);
            if (!server) return;
            sessions.set(serverId, {
                serverId: serverId,
                status: 'disconnected',
                currentPath: server.defaultPath,
                files: [],
                isLoading: false,
            });
        }
        activeSessionId = serverId;
    }

    async function handleConnect() {
        if (!activeSession) return;
        activeSession.status = 'connecting';
        activeSession.error = undefined;
        sessions.set(activeSession.serverId, activeSession);
        await new Promise(resolve => setTimeout(resolve, 1500));
        try {
            console.log('Connecting to:', activeServer);
            await navigateDirectory(activeSession.currentPath, true);
            activeSession.status = 'connected';
        } catch (e: any) {
            activeSession.status = 'error';
            activeSession.error = `Connection failed: ${e.message}`;
            console.error(e);
        } finally {
            sessions.set(activeSession.serverId, activeSession);
        }
    }

    async function navigateDirectory(path: string, isInitialLoad = false) {
        if (!activeSession) return;
        activeSession.isLoading = true;
        if (!isInitialLoad) {
            activeSession.currentPath = path;
        }
        sessions.set(activeSession.serverId, activeSession);
        await new Promise(resolve => setTimeout(resolve, 1000));
        try {
            console.log(`Fetching contents for ${activeSession.serverId} at ${path}`);
            const isRoot = path === '/' || /^[a-zA-Z]:\\?$/.test(path);
            const dummyFiles: FileItem[] = [
                ...(!isRoot ? [{ name: '..', path: getParentPath(path), type: 'directory' as const }] : []),
                { name: 'archive', path: `${path}/archive`, type: 'directory' },
                { name: 'nginx-access.log', path: `${path}/nginx-access.log`, type: 'file', size: '15.8 MB', modified: '2025-07-15 14:10:22' },
                { name: 'old_logs', path: `${path}/old_logs`, type: 'directory' },
                { name: 'nginx-error.log', path: `${path}/nginx-error.log`, type: 'file', size: '2.1 MB', modified: '2025-07-15 09:15:23' },
                { name: 'syslog', path: `${path}/syslog`, type: 'file', size: '5.2 MB', modified: '2025-07-14 18:55:41' },
            ];
            activeSession.files = dummyFiles.sort((a, b) => {
                if (a.type === b.type) return a.name.localeCompare(b.name);
                return a.type === 'directory' ? -1 : 1;
            });
        } catch (e: any) {
            activeSession.status = 'error';
            activeSession.error = `Could not list directory: ${e.message}`;
        } finally {
            activeSession.isLoading = false;
            sessions.set(activeSession.serverId, activeSession);
        }
    }

    function getParentPath(path: string): string {
        const parts = path.replace(/\\/g, '/').split('/').filter(p => p);
        if (parts.length <= 1 && /^[a-zA-Z]:$/.test(parts[0])) return `${parts[0]}\\`;
        if (parts.length <= 1) return '/';
        parts.pop();
        const separator = path.includes('/') ? '/' : '\\';
        if (/^[a-zA-Z]:$/.test(path.split(separator)[0])) {
            return path.split(separator)[0] + separator + parts.slice(1).join(separator);
        }
        return separator + parts.join(separator);
    }

    function getBreadcrumbs(path: string): { name: string, path: string }[] {
        const separator = path.includes('/') ? '/' : '\\';
        const root = separator === '/' ? '/' : path.split(separator)[0] + separator;
        const parts = path.replace(root, '').split(separator).filter(p => p);
        const crumbs = [{ name: 'home', path: root }];
        let currentPath = root;
        for (const part of parts) {
            currentPath += `${part}${separator}`;
            crumbs.push({ name: part, path: currentPath });
        }
        return crumbs;
    }
</script>

<svelte:head>
    <style>
        /* Pola latar belakang halus untuk body */
        body, html {
            background-color: #0f172a; /* bg-slate-900 */
            background-image: radial-gradient(#334155 0.5px, transparent 0.5px);
            background-size: 16px 16px;
        }
    </style>
</svelte:head>

<div class="flex h-screen font-sans text-slate-200">

    <aside class="flex-shrink-0 w-72 p-4 bg-slate-900/70 backdrop-blur-xl border-r border-slate-700/50 shadow-2xl flex flex-col">
        <div class="flex items-center gap-3 px-2 mb-6">
            <svg class="w-8 h-8 text-blue-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M21.75 17.25v-.228a4.5 4.5 0 00-.12-1.03l-2.268-9.64a3.375 3.375 0 00-3.285-2.802H9.12a3.375 3.375 0 00-3.285 2.802l-2.268 9.64a4.5 4.5 0 00-.12 1.03v.228m15.45-1.5H2.25m15.45 0v7.5a2.25 2.25 0 01-2.25 2.25h-11.25a2.25 2.25 0 01-2.25-2.25v-7.5" /></svg>
            <h2 class="text-xl font-bold">Remote Explorer</h2>
        </div>
        <nav class="flex flex-col space-y-1 overflow-y-auto flex-grow">
            {#each servers as server (server.id)}
                <button
                        on:click={() => selectServer(server.id)}
                        class="w-full flex items-center gap-3 px-3 py-2.5 text-left rounded-lg transition-colors duration-200"
                >
                <svg class="w-5 h-5 flex-shrink-0 text-slate-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M21.75 17.25v-.228a4.5 4.5 0 00-.12-1.03l-2.268-9.64a3.375 3.375 0 00-3.285-2.802H9.12a3.375 3.375 0 00-3.285 2.802l-2.268 9.64a4.5 4.5 0 00-.12 1.03v.228m15.45-1.5H2.25m15.45 0v7.5a2.25 2.25 0 01-2.25 2.25h-11.25a2.25 2.25 0 01-2.25-2.25v-7.5" /></svg>
                <div>
                    <span class="font-semibold">{server.name}</span>
                    <span class="block text-xs text-slate-400">{server.user}@{server.host}</span>
                </div>
                </button>
            {/each}
        </nav>
        <Button class="w-full mt-4 border-slate-600 hover:bg-slate-700/50 text-slate-300">Add New Server</Button>
    </aside>

    <main class="flex-1 p-6 lg:p-10 overflow-y-auto">
        {#if activeSession && activeServer}
            <div class="max-w-full mx-auto">
                <h1 class="text-4xl font-bold tracking-tight bg-gradient-to-r from-blue-400 to-cyan-300 bg-clip-text text-transparent mb-2">{activeServer.name}</h1>
                <p class="mb-8 text-slate-400">Manage and explore files on your remote server.</p>

                {#if activeSession.status === 'disconnected' || activeSession.status === 'error'}
                    <div class="p-8 bg-slate-800/50 backdrop-blur-sm rounded-xl border border-slate-700/50 shadow-lg">
                        <h3 class="text-2xl font-bold mb-6">Connect to Server</h3>
                        <dl class="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-6">
                            <div class="p-4 bg-slate-900/50 rounded-lg">
                                <dt class="text-sm font-medium text-slate-400">Host Address</dt>
                                <dd class="text-lg font-mono">{activeServer.host}:{activeServer.port}</dd>
                            </div>
                            <div class="p-4 bg-slate-900/50 rounded-lg">
                                <dt class="text-sm font-medium text-slate-400">Username</dt>
                                <dd class="text-lg font-mono">{activeServer.user}</dd>
                            </div>
                            <div class="md:col-span-2">
                                <label for="start-path" class="text-sm font-medium text-slate-400 mb-2 block">Starting Directory</label>
                                <Input id="start-path" bind:value={activeSession.currentPath} class="font-mono" />
                            </div>
                        </dl>
                        {#if activeSession.status === 'error'}
                            <div class="mt-6 p-4 text-red-300 bg-red-500/10 border border-red-500/30 rounded-lg">
                                <p class="font-bold">Connection Error</p>
                                <p class="text-sm">{activeSession.error}</p>
                            </div>
                        {/if}
                        <div class="flex justify-end mt-8">
                            <Button on:click={handleConnect} disabled={activeSession.status === 'connecting'} class="bg-blue-600 hover:bg-blue-500 text-white shadow-blue-500/20 shadow-lg">
                                Connect Now
                            </Button>
                        </div>
                    </div>

                {:else if activeSession.status === 'connecting'}
                    <div class="flex flex-col items-center justify-center h-96 text-slate-400">
                        <svg class="w-12 h-12 mb-4 text-blue-500 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
                        <span class="text-lg">Establishing connection to {activeServer.host}...</span>
                    </div>

                {:else if activeSession.status === 'connected'}
                    <div class="bg-slate-800/50 backdrop-blur-sm rounded-xl border border-slate-700/50 shadow-lg">
                        <nav class="flex items-center text-sm p-4 text-slate-400 border-b border-slate-700/50">
                            {#each getBreadcrumbs(activeSession.currentPath) as crumb, i}
                                <button on:click={() => navigateDirectory(crumb.path)} class="flex items-center gap-1.5 hover:text-white transition-colors">
                                    {#if crumb.name === 'home'}
                                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5"><path fill-rule="evenodd" d="M9.293 2.293a1 1 0 011.414 0l7 7A1 1 0 0117 11h-1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-3a1 1 0 00-1-1H9a1 1 0 00-1 1v3a1 1 0 01-1 1H5a1 1 0 01-1-1v-6H3a1 1 0 01-.707-1.707l7-7z" clip-rule="evenodd" /></svg>
                                    {:else}
                                        {crumb.name}
                                    {/if}
                                </button>
                                {#if i < getBreadcrumbs(activeSession.currentPath).length - 1}
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor" class="w-4 h-4 mx-1.5 text-slate-600"><path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" /></svg>
                                {/if}
                            {/each}
                        </nav>

                        <div class="relative min-h-[24rem]">
                            <header class="grid grid-cols-[auto,1fr,150px,100px] gap-4 px-4 py-2 text-xs font-semibold text-slate-400 border-b border-slate-700/50 sticky top-0 bg-slate-800/80 backdrop-blur-sm">
                                <div class="col-start-2">Name</div>
                                <div class="text-right">Last Modified</div>
                                <div class="text-right">Size</div>
                            </header>

                            {#if activeSession.isLoading}
                                <div class="absolute inset-0 bg-slate-800/50 flex items-center justify-center z-10">
                                    <span class="text-slate-300">Loading directory...</span>
                                </div>
                            {/if}

                            <ul class="space-y-0.5 p-2">
                                {#each activeSession.files as item (item.path)}
                                    <li
                                            class="grid grid-cols-[32px,1fr,150px,100px] items-center gap-4 p-2 rounded-md transition-colors duration-150"
                                            class:cursor-pointer={item.type === 'directory'}
                                    on:click={() => item.type === 'directory' && navigateDirectory(item.path)}
                                    on:dblclick={() => item.type === 'file' && alert(`Viewing file: ${item.name}`)}
                                    title={item.type === 'file' ? `Double-click to view ${item.name}` : `Open directory ${item.name}`}
                                    >
                                    <div class="text-slate-400">
                                        {#if item.name === '..'}
                                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M9 9l6-6m0 0l6 6m-6-6v12a6 6 0 01-12 0v-3" /></svg>
                                        {:else if item.type === 'directory'}
                                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6 text-yellow-500/80"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12.75V12A2.25 2.25 0 014.5 9.75h15A2.25 2.25 0 0121.75 12v.75m-8.69-6.44l-2.12-2.12a1.5 1.5 0 00-1.061-.44H4.5A2.25 2.25 0 002.25 6v12a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9a2.25 2.25 0 00-2.25-2.25h-5.379a1.5 1.5 0 01-1.06-.44z" /></svg>
                                        {:else}
                                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" /></svg>
                                        {/if}
                                    </div>
                                    <p class="font-medium truncate">{item.name}</p>
                                    <p class="text-sm text-slate-400 text-right font-mono">{item.modified || ''}</p>
                                    <p class="text-sm text-slate-300 text-right font-mono">{item.size || ''}</p>
                                    </li>
                                {:else}
                                    <div class="text-center text-slate-500 pt-16">This directory is empty.</div>
                                {/each}
                            </ul>
                        </div>
                    </div>
                {/if}
            </div>
        {:else}
            <div class="flex items-center justify-center h-full text-center text-slate-600">
                <div>
                    <svg class="w-24 h-24 mx-auto mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M21.75 17.25v-.228a4.5 4.5 0 00-.12-1.03l-2.268-9.64a3.375 3.375 0 00-3.285-2.802H9.12a3.375 3.375 0 00-3.285 2.802l-2.268 9.64a4.5 4.5 0 00-.12 1.03v.228m15.45-1.5H2.25m15.45 0v7.5a2.25 2.25 0 01-2.25 2.25h-11.25a2.25 2.25 0 01-2.25-2.25v-7.5" /></svg>
                    <p class="text-xl font-semibold text-slate-500">Welcome to Remote Explorer</p>
                    <p>Select a server from the sidebar to begin.</p>
                </div>
            </div>
        {/if}
    </main>
</div>