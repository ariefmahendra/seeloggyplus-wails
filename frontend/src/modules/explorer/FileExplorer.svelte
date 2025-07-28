<script lang="ts">
    import {onMount} from 'svelte';
    import {
        ConnectSession,
        GetListSession,
        GetLocalDrives,
        GetLocalUserHomeDir,
        GetRemoteListFiles,
        ListLocalFiles
    } from '../../../wailsjs/go/main/App';
    import {dto} from '../../../wailsjs/go/models';

    import {
        Breadcrumb,
        BreadcrumbItem,
        Button,
        InlineNotification,
        Loading,
        Search,
        Tab,
        Tabs
    } from 'carbon-components-svelte';

    import {ArrowLeft, BareMetalServer, FileStorage, FolderOpen, Home, Restart} from 'carbon-icons-svelte';
    import FileItem from "./FileItem.svelte";

    // State variables
    let currentPath = '';
    let isRemote = false;
    let selectedSessionId = '';
    let files: dto.FileInfo[] = [];
    let drives: dto.DriveInfo[] = [];
    let sessions: dto.SessionManagerDto[] = [];
    let loading = false;
    let error = '';
    let searchValue = '';
    let selectedTab = 0;

    // Reactive statements
    $: filteredFiles = files.filter(file =>
        file.name.toLowerCase().includes(searchValue.toLowerCase())
    );

    $: pathParts = currentPath ? currentPath.split(/[\/\\]/).filter(Boolean) : [];

    // Lifecycle
    onMount(async () => {
        console.log('FileExplorer mounted');
        await initializeExplorer();
    });

    async function initializeExplorer() {
        try {
            await loadSessions();
            if (!isRemote) {
                await loadLocalDrives();
                await loadUserHome();
            }
        } catch (err) {
            console.error('Failed to initialize explorer:', err);
            error = `Initialization failed: ${err}`;
        }
    }

    // API Functions
    async function loadSessions() {
        try {
            sessions = await GetListSession();
            console.log('Sessions loaded:', sessions);
        } catch (err) {
            console.error('Failed to load sessions:', err);
            error = `Failed to load sessions: ${err}`;
        }
    }

    async function loadLocalDrives() {
        try {
            drives = await GetLocalDrives();
            console.log('Drives loaded:', drives);
        } catch (err) {
            console.error('Failed to load drives:', err);
            error = `Failed to load drives: ${err}`;
        }
    }

    async function loadUserHome() {
        try {
            if (!currentPath) {
                currentPath = await GetLocalUserHomeDir();
                console.log('User home loaded:', currentPath);
                await loadFiles();
            }
        } catch (err) {
            console.error('Failed to load user home:', err);
            error = `Failed to load user home: ${err}`;
        }
    }

    async function loadFiles() {
        if (!currentPath && !isRemote) return;

        loading = true;
        error = '';

        try {
            console.log('Loading files for path:', currentPath, 'isRemote:', isRemote);

            if (isRemote && selectedSessionId) {
                files = await GetRemoteListFiles(selectedSessionId, currentPath || '/');
            } else if (!isRemote && currentPath) {
                files = await ListLocalFiles(currentPath);
            }

            console.log('Files loaded:', files);
        } catch (err) {
            console.error('Failed to load files:', err);
            error = `Failed to load files: ${err}`;
            files = [];
        } finally {
            loading = false;
        }
    }

    // Event Handlers
    function handleTabChange(event) {
        const newTab = event.detail.selected;
        console.log('Tab changed to:', newTab);

        if (newTab !== selectedTab) {
            selectedTab = newTab;
            isRemote = selectedTab === 1;
            currentPath = '';
            selectedSessionId = '';
            files = [];
            error = '';
            searchValue = '';

            if (!isRemote) {
                loadUserHome();
            }
        }
    }

    function handleFileClick(file) {
        console.log('File clicked:', file);
        if (file.isDir) {
            currentPath = file.path;
            loadFiles();
        }
    }

    function handleDriveClick(drive) {
        console.log('Drive clicked:', drive);
        currentPath = drive.path;
        loadFiles();
    }

    function navigateUp() {
        if (!currentPath || currentPath === '/') return;

        const separator = currentPath.includes('\\') ? '\\' : '/';
        const pathParts = currentPath.split(separator).filter(Boolean);

        if (pathParts.length > 1) {
            currentPath = separator + pathParts.slice(0, -1).join(separator);
        } else {
            currentPath = isRemote ? '/' : '';
        }

        loadFiles();
    }

    function navigateHome() {
        if (isRemote) {
            currentPath = '/';
        } else {
            loadUserHome();
        }
    }

    function refreshFiles() {
        loadFiles();
    }

    function navigateToBreadcrumb(index) {
        let newPath;

        if (index === -1) {
            newPath = isRemote ? '/' : '';
        } else {
            const separator = currentPath.includes('\\') ? '\\' : '/';
            const partsToJoin = pathParts.slice(0, index + 1);

            if (!isRemote && partsToJoin.length > 0 && /^[a-zA-Z]:$/.test(partsToJoin[0])) {
                newPath = partsToJoin.join(separator);
            } else {
                const pathSuffix = partsToJoin.join(separator);
                newPath = isRemote ? separator + pathSuffix : pathSuffix;
            }
        }

        console.log('Navigating to breadcrumb with new path:', newPath);

        currentPath = newPath;
        loadFiles();
    }

    // Session handling
    async function connectToSession(sessionId) {
        try {
            loading = true;
            await ConnectSession(sessionId);
            selectedSessionId = sessionId;
            currentPath = '/';
            await loadFiles();
        } catch (err) {
            error = `Failed to connect to session: ${err}`;
        } finally {
            loading = false;
        }
    }

    // Utility Functions
    function formatFileSize(bytes) {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
    }

    function formatDate(dateString) {
        if (!dateString) return '—';
        try {
            return new Date(dateString).toLocaleDateString('en-US', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            });
        } catch {
            return dateString;
        }
    }
</script>

<div class="w-full h-[90vh] flex flex-col bg-white border border-gray-200 rounded-lg overflow-hidden gap-2">
    <!-- Header with Tabs -->
    <div class="border-b border-gray-200 flex-shrink-0">
        <Tabs bind:selected={selectedTab} on:change={handleTabChange}>
            <Tab>
                <div class="flex items-center gap-2">
                    <FileStorage size={16}/>
                    Local Storage
                </div>
            </Tab>
            <Tab>
                <div class="flex items-center gap-2">
                    <BareMetalServer size={16}/>
                    Remote Server
                </div>
            </Tab>
        </Tabs>
    </div>

    <!-- Toolbar -->
    <div class="border-b border-gray-200 flex-shrink-0 flex justify-between items-center">
        <div>
            <Button
                    kind="ghost"
                    size="small"
                    icon={ArrowLeft}
                    disabled={!currentPath || currentPath === '/'}
                    on:click={navigateUp}
            />

            <Button
                    kind="ghost"
                    size="small"
                    icon={Home}
                    on:click={navigateHome}
            />

            <Button
                    kind="ghost"
                    size="small"
                    icon={Restart}
                    on:click={refreshFiles}
                    disabled={loading}
            />

            <!-- Breadcrumb -->
            {#if currentPath}
                <Breadcrumb noTrailingSlash class="p-2">
                    {#each pathParts as part, index}
                        <BreadcrumbItem
                                on:click={() => navigateToBreadcrumb(index)}
                                class="cursor-pointer hover:text-blue-600"
                        >
                            <span class="text-sm" title={part}>{part}</span>
                        </BreadcrumbItem>
                    {/each}
                </Breadcrumb>
            {/if}
        </div>
        <div>
            <Search/>
        </div>
    </div>

    <!-- Session Selection for Remote -->
    {#if isRemote}
        <div class="p-3 bg-gray-50 border-b border-gray-200 flex-shrink-0">
            <div class="flex items-center gap-3">
                <span class="text-sm font-medium">Remote Session:</span>
                <div class="flex flex-wrap gap-2">
                    {#if sessions.length === 0}
                        <span class="text-sm text-gray-500">No sessions available</span>
                    {:else}
                        {#each sessions as session}
                            <Button
                                    kind={selectedSessionId === session.id ? "primary" : "tertiary"}
                                    size="small"
                                    on:click={() => connectToSession(session.id)}
                            >
                                {session.serverInfo?.name || session.id}
                            </Button>
                        {/each}
                    {/if}
                </div>
            </div>
        </div>
    {/if}

    <!-- Drives Section (Local only) -->
    {#if !isRemote && drives.length > 0}
        <div class="p-3 bg-gray-50 border-b border-gray-200 flex-shrink-0">
            <div class="text-sm font-medium text-gray-700 mb-2">Available Drives</div>
            <div class="flex flex-wrap gap-2">
                {#each drives as drive}
                    <Button
                            kind="tertiary"
                            size="small"
                            on:click={() => handleDriveClick(drive)}
                    >
                        <div class="flex items-center gap-2">
                            <FileStorage size={16}/>
                            <span>{drive.label || drive.name} ({drive.path})</span>
                        </div>
                    </Button>
                {/each}
            </div>
        </div>
    {/if}

    <!-- Error Message -->
    {#if error}
        <div class="p-3 flex-shrink-0">
            <InlineNotification
                    kind="error"
                    title="Error"
                    subtitle={error}
                    on:close={() => error = ''}
            />
        </div>
    {/if}

    <!-- Files List -->
    <div class="flex-1 overflow-auto">
        {#if loading}
            <div class="h-32 flex items-center justify-center">
                <div class="text-center">
                    <Loading withOverlay={false}/>
                    <p class="mt-2 text-sm text-gray-600">Loading files...</p>
                </div>
            </div>
        {:else if filteredFiles.length === 0}
            <div class="h-32 flex flex-col items-center justify-center text-gray-500">
                <FolderOpen size={32} class="mb-2 text-gray-300"/>
                {#if !currentPath && !selectedSessionId}
                    <p class="text-sm">Select a location to browse files</p>
                {:else if searchValue}
                    <p class="text-sm">No files match your search</p>
                {:else}
                    <p class="text-sm">This folder is empty</p>
                {/if}
            </div>
        {:else}
            <div class="divide-y divide-gray-100">
                {#each filteredFiles as file}
                    <FileItem {file} on:click={(e) => handleFileClick(e.detail)}/>
                {/each}
            </div>
        {/if}
    </div>

    <!-- Status Bar -->
    <div class="border-t border-gray-200 px-3 py-2 bg-gray-50 flex-shrink-0">
        <div class="flex items-center justify-between text-xs text-gray-600">
            <div class="flex items-center gap-4">
                <span>{filteredFiles.length} items</span>
                {#if currentPath}
          <span class="font-mono text-xs bg-white px-2 py-1 rounded border">
            {currentPath}
          </span>
                {/if}
            </div>
            <div class="flex items-center gap-2">
                <span>{isRemote ? 'Remote' : 'Local'}</span>
                {#if selectedSessionId}
                    <span class="text-green-600">Connected</span>
                {/if}
            </div>
        </div>
    </div>
</div>