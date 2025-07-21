<script lang="ts">
    import {onMount} from 'svelte';
    import {
        Button,
        CodeSnippet,
        Content,
        DataTable,
        Header,
        HeaderGlobalAction,
        HeaderUtilities,
        InlineNotification,
        Loading,
        Modal,
        Search,
        Select,
        SelectItem,
        SideNav,
        SideNavItems,
        SideNavLink,
        Tag,
        TextInput,
        Toggle
    } from 'carbon-components-svelte';

    import {Add, DocumentView, Filter, Pause, Play, ServerProxy, Settings, TrashCan, ConnectionSignal} from 'carbon-icons-svelte';

    import ServerManagementModal from './modules/servers/ServerManagementModal.svelte'

    // Types
    interface LogEntry {
        id: string;
        timestamp: string;
        level: string;
        source: string;
        message: string;
        requestId?: string;
        stackTrace?: string;
    }

    interface LogFile {
        id: string;
        name: string;
        path: string;
        size: number;
        lastModified: string;
        isActive: boolean;
    }

    // Component state
    let isSideNavOpen = false;
    let logEntries: LogEntry[] = [];
    let searchTerm = '';
    let isLoading = false;
    let error = '';
    let isTailing = false;

    // state for server management
    let serverManagementOpen = false;

    // File management
    let openFiles: LogFile[] = [];
    let activeFileId = '';

    // Modal states
    let showAddFileModal = false;
    let showServerModal = false;
    let showSettingsModal = false;
    let showDetailModal = false;

    // Forms
    let newFilePath = '';
    let selectedLogEntry: LogEntry | null = null;

    // Filter settings
    let filterType = 'text';
    let filterQuery = '';

    // Table configuration
    const headers = [
        {key: 'timestamp', value: 'Time', empty: false},
        {key: 'level', value: 'Level', empty: false},
        {key: 'source', value: 'Source', empty: false},
        {key: 'message', value: 'Message', empty: false},
        {key: 'actions', value: 'Actions', empty: false}
    ];

    // Sample data for UI demo
    onMount(() => {
        loadSampleFiles();
        loadSampleLogs();
    });

    // function for open server management modal
    function openServerManagement() {
        showServerModal = true;
    }

    function loadSampleFiles() {
        openFiles = [
            {
                id: 'app-1',
                name: 'application.log',
                path: '/var/log/app/application.log',
                size: 10485760,
                lastModified: new Date().toISOString(),
                isActive: true
            },
            {
                id: 'error-1',
                name: 'error.log',
                path: '/var/log/app/error.log',
                size: 5242880,
                lastModified: new Date().toISOString(),
                isActive: false
            }
        ];
        activeFileId = 'app-1';
    }

    function loadSampleLogs() {
        logEntries = [
            {
                id: '1',
                timestamp: new Date().toISOString(),
                level: 'ERROR',
                source: 'auth-service',
                message: 'Authentication failed for user ID 12345'
            },
            {
                id: '2',
                timestamp: new Date(Date.now() - 60000).toISOString(),
                level: 'WARN',
                source: 'api-gateway',
                message: 'High memory usage detected: 85%'
            },
            {
                id: '3',
                timestamp: new Date(Date.now() - 120000).toISOString(),
                level: 'INFO',
                source: 'payment-service',
                message: 'Payment processed successfully'
            }
        ];
    }

    // File management functions
    function switchToFile(fileId: string) {
        activeFileId = fileId;
        openFiles = openFiles.map(file => ({
            ...file,
            isActive: file.id === fileId
        }));
        console.log('Switching to file:', fileId);
    }

    function addLogFile() {
        if (!newFilePath.trim()) return;

        const newFile: LogFile = {
            id: Date.now().toString(),
            name: newFilePath.split('/').pop() || 'unknown.log',
            path: newFilePath,
            size: 0,
            lastModified: new Date().toISOString(),
            isActive: false
        };

        openFiles = [...openFiles, newFile];
        newFilePath = '';
        showAddFileModal = false;
        console.log('Adding file:', newFile);
    }

    // UI functions
    function toggleTailing() {
        isTailing = !isTailing;
        console.log('Tailing mode:', isTailing);
    }

    function applyFilter() {
        console.log('Applying filter:', {filterType, filterQuery});
    }

    function showLogDetail(log: LogEntry) {
        selectedLogEntry = log;
        showDetailModal = true;
    }

    function getLevelTagType(level: string): "red" | "warm-gray" | "blue" | "gray" {
        switch (level) {
            case 'ERROR':
                return 'red';
            case 'WARN':
                return 'warm-gray';
            case 'INFO':
                return 'blue';
            case 'DEBUG':
                return 'gray';
            default:
                return 'gray';
        }
    }

    function formatTimestamp(timestamp: string): string {
        return new Date(timestamp).toLocaleString();
    }

    function truncateMessage(message: string, maxLength: number = 80): string {
        return message.length > maxLength ? message.substring(0, maxLength) + '...' : message;
    }

    function formatFileSize(bytes: number): string {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    // Filtered logs for search
    $: filteredLogs = searchTerm
        ? logEntries.filter(log =>
            log.message.toLowerCase().includes(searchTerm.toLowerCase()) ||
            log.source.toLowerCase().includes(searchTerm.toLowerCase()) ||
            log.level.toLowerCase().includes(searchTerm.toLowerCase())
        )
        : logEntries;
</script>

<!-- Header -->
<Header platformName="SeeLoggyPlus" bind:isSideNavOpen persistentHamburgerMenu={true}>
    <HeaderUtilities>
        <HeaderGlobalAction
                icon={Add}
                iconDescription="Add Log File"
                tooltipAlignment="end"
                on:click={() => showAddFileModal = true}
        />
        <HeaderGlobalAction
                iconDescription="Server Management"
                tooltipAlignment="end"
                icon={ServerProxy}
                on:click={() => showServerModal = !showServerModal}
        />
        <HeaderGlobalAction
            icon={Settings}
            iconDescription="Settings"
            tooltipAlignment="end"
            on:click={() => showSettingsModal = !showSettingsModal}
        />
    </HeaderUtilities>
</Header>

<!-- Sidebar for File Management -->
<SideNav bind:isOpen={isSideNavOpen}>
    <SideNavItems>
        <div class="p-4 font-semibold text-gray-600 border-b border-gray-200">
            Open Files ({openFiles.length})
        </div>

        {#each openFiles as file (file.id)}
            <SideNavLink
                    text={file.name}
                    isSelected={file.isActive}
                    on:click={() => switchToFile(file.id)}
            >
                <div class="flex justify-between items-center w-full">
                    <div class="flex flex-col">
                        <span class="text-sm font-medium">{file.name}</span>
                        <span class="text-xs text-gray-500">{formatFileSize(file.size)}</span>
                    </div>
                    <Button
                            kind="ghost"
                            size="small"
                            iconDescription="Remove file"
                            icon={TrashCan}
                    />
                </div>
            </SideNavLink>
        {/each}

        {#if openFiles.length === 0}
            <div class="p-4 text-gray-500 italic">
                No files open
            </div>
        {/if}
    </SideNavItems>
</SideNav>

<!-- Main Content -->
<Content id="main-content">
    <div class="p-6">
        <!-- Control Panel -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 mb-4 p-4">
            <div class="flex gap-4 items-center flex-wrap">
                <!-- Search -->
                <div class="flex-1 min-w-64">
                    <Search
                            placeholder="Search logs..."
                            bind:value={searchTerm}
                            size="sm"
                    />
                </div>

                <!-- Filter -->
                <div class="flex gap-2 items-center">
                    <Select bind:selected={filterType} size="sm" class="w-20">
                        <SelectItem value="text" text="Text"/>
                        <SelectItem value="regex" text="Regex"/>
                    </Select>
                    <TextInput
                            placeholder="Filter query..."
                            bind:value={filterQuery}
                            size="sm"
                            class="w-48"
                    />
                    <Button size="small" on:click={applyFilter}>
                        <Filter class="mr-2"/>
                        Apply
                    </Button>
                </div>

                <!-- Tail Controls -->
                <div class="flex gap-2 items-center">
                    <Button
                            kind={isTailing ? "danger" : "primary"}
                            size="small"
                            on:click={toggleTailing}
                    >
                        {#if isTailing}
                            <Pause class="mr-2"/>
                            Stop Tail
                        {:else}
                            <Play class="mr-2"/>
                            Start Tail
                        {/if}
                    </Button>
                </div>
            </div>

            <!-- Status Info -->
            <div class="mt-4 flex gap-4 items-center text-sm text-gray-600">
                <span>Total: {filteredLogs.length} logs</span>
                {#if activeFileId}
                    <span>File: {openFiles.find(f => f.id === activeFileId)?.name || 'Unknown'}</span>
                {/if}
                <span class:text-green-600={isTailing} class:text-gray-600={!isTailing}>
                    {isTailing ? '● Live' : '○ Paused'}
                </span>
            </div>
        </div>

        <!-- Log Table -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200">
            {#if isLoading}
                <div class="p-8 text-center">
                    <Loading/>
                </div>
            {:else if error}
                <div class="p-4">
                    <InlineNotification
                            kind="error"
                            title="Error loading logs"
                            subtitle={error}
                    />
                </div>
            {:else}
                <DataTable
                        {headers}
                        rows={filteredLogs.map(log => ({
                        id: log.id,
                        timestamp: formatTimestamp(log.timestamp),
                        level: log.level,
                        source: log.source,
                        message: truncateMessage(log.message),
                        actions: log.id
                    }))}
                        size="compact"
                        zebra
                        sortable
                >
                    <svelte:fragment slot="cell" let:row let:cell>
                        {#if cell.key === 'level'}
                            <Tag type={getLevelTagType(cell.value)}>{cell.value}</Tag>
                        {:else if cell.key === 'message'}
                            <span class="font-mono text-sm">{cell.value}</span>
                        {:else if cell.key === 'actions'}
                            <Button
                                    kind="ghost"
                                    size="small"
                                    on:click={() => showLogDetail(logEntries.find(l => l.id === row.id))}
                            >
                                View Details
                            </Button>
                        {:else}
                            {cell.value}
                        {/if}
                    </svelte:fragment>
                </DataTable>
            {/if}

            {#if filteredLogs.length === 0 && !isLoading}
                <div class="p-8 text-center text-gray-500">
                    <DocumentView size="32" class="mx-auto mb-4 opacity-50"/>
                    <p>No logs found</p>
                    <p class="text-sm mt-2">Try adjusting your search or filter criteria</p>
                </div>
            {/if}
        </div>
    </div>
</Content>

<!-- Log Detail Modal -->
<Modal
        bind:open={showDetailModal}
        modalHeading="Log Entry Details"
        primaryButtonText="Close"
        on:click:button--primary={() => showDetailModal = false}
>
    {#if selectedLogEntry}
        <div class="flex flex-col gap-4">
            <div><strong>Timestamp:</strong> {formatTimestamp(selectedLogEntry.timestamp)}</div>
            <div><strong>Level:</strong>
                <Tag type={getLevelTagType(selectedLogEntry.level)}>{selectedLogEntry.level}</Tag>
            </div>
            <div><strong>Source:</strong> {selectedLogEntry.source}</div>
            <div><strong>Message:</strong></div>
            <CodeSnippet type="multi">{selectedLogEntry.message}</CodeSnippet>

            {#if selectedLogEntry.requestId}
                <div><strong>Request ID:</strong> {selectedLogEntry.requestId}</div>
            {/if}

            {#if selectedLogEntry.stackTrace}
                <div><strong>Stack Trace:</strong></div>
                <CodeSnippet type="multi">{selectedLogEntry.stackTrace}</CodeSnippet>
            {/if}
        </div>
    {/if}
</Modal>

<!-- Add File Modal -->
<Modal
        bind:open={showAddFileModal}
        modalHeading="Add Log File"
        primaryButtonText="Add File"
        secondaryButtonText="Cancel"
        on:click:button--primary={addLogFile}
        on:click:button--secondary={() => showAddFileModal = false}
>
    <div class="space-y-4">
        <TextInput
                labelText="File Path"
                placeholder="/var/log/application.log"
                bind:value={newFilePath}
                helperText="Enter the full path to the log file"
        />
    </div>
</Modal>

<!-- Server Management Modal -->
<ServerManagementModal bind:isOpen={showServerModal}/>

<!-- Settings Modal -->
<Modal
        bind:open={showSettingsModal}
        modalHeading="Application Settings"
        primaryButtonText="Save"
        secondaryButtonText="Cancel"
        on:click:button--primary={() => showSettingsModal = false}
        on:click:button--secondary={() => showSettingsModal = false}
>
    <div class="space-y-4">
        <Toggle labelText="Auto-refresh logs"/>
        <Toggle labelText="Show timestamps in relative format"/>
        <Toggle labelText="Enable syntax highlighting"/>
        <Select labelText="Default page size">
            <SelectItem value="25" text="25 logs per page"/>
            <SelectItem value="50" text="50 logs per page"/>
            <SelectItem value="100" text="100 logs per page"/>
        </Select>
        <TextInput
                labelText="Refresh interval (seconds)"
                type="number"
                value="5"
        />
    </div>
</Modal>