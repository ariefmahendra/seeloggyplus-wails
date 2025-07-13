<script lang="ts">
    import Button from '../components/common/Button.svelte';
    import Input from '../components/common/Input.svelte';

    // --- State Management ---

    // Tipe data untuk hasil pencarian
    interface SearchResult {
        name: string;
        path: string;
        size: string;
        modified: string;
    }

    // State untuk mode pencarian: 'local' atau 'remote'
    let searchMode: 'local' | 'remote' = 'local';

    // State untuk UI: loading, error, dan hasil
    let isLoading: boolean = false;
    let error: string | null = null;
    let searchResults: SearchResult[] = [];

    // State untuk input form Lokal
    let localPath: string = 'C:\\ProgramData\\...';
    let localPattern: string = '*.log';

    // State untuk input form Remote
    let remoteHost: string = '';
    let remotePort: number = 22;
    let remoteUser: string = '';
    let remotePassword: string = '';
    let remoteLogPath: string = '/var/log/';

    // --- Logic ---

    async function handleSearch() {
        isLoading = true;
        error = null;
        searchResults = [];

        // Simulasi pencarian dengan delay untuk menunjukkan state loading
        await new Promise(resolve => setTimeout(resolve, 1500));

        try {
            if (searchMode === 'local') {
                // TODO: Ganti dengan pemanggilan fungsi Wails untuk pencarian lokal
                console.log('Searching locally:', { path: localPath, pattern: localPattern });
                // Contoh hasil dummy
                searchResults = [
                    { name: 'app-2024-07-12.log', path: 'C:\\logs\\app-2024-07-12.log', size: '1.2 MB', modified: '2024-07-12 10:30:15' },
                    { name: 'error.log', path: 'C:\\logs\\error.log', size: '345 KB', modified: '2024-07-11 18:05:00' },
                ];
            } else {
                // TODO: Ganti dengan pemanggilan fungsi Wails untuk pencarian remote (SSH/SFTP)
                console.log('Searching remote:', { host: remoteHost, port: remotePort, user: remoteUser, path: remoteLogPath });
                // Contoh hasil dummy
                searchResults = [
                    { name: 'nginx-access.log', path: '/var/log/nginx/access.log', size: '15.8 MB', modified: '2024-07-12 11:00:00' },
                    { name: 'syslog', path: '/var/log/syslog', size: '5.2 MB', modified: '2024-07-12 10:55:41' },
                ];
            }
            // Jika tidak ada hasil setelah pencarian
            if (searchResults.length === 0) {
                error = "No log files found matching your criteria.";
            }
        } catch (e: any) {
            error = `Failed to search: ${e.message}`;
            console.error(e);
        } finally {
            isLoading = false;
        }
    }

    // Fungsi untuk memilih direktori (placeholder untuk fungsi Wails)
    function chooseDirectory() {
        // TODO: Panggil fungsi Wails untuk membuka dialog pemilihan direktori
        // contoh: wails.runtime.Dialog.Directory().then(path => localPath = path);
        alert("Fungsi pemilihan direktori akan diimplementasikan dengan Wails.");
    }
</script>

<div class="max-w-6xl mx-auto">
    <h1 class="mb-8 text-4xl font-bold text-white">Log File Finder</h1>

    <!-- Tab Switcher -->
    <div class="mb-6 border-b border-gray-700">
        <nav class="flex -mb-px space-x-6" aria-label="Tabs">
            <button
                    on:click={() => searchMode = 'local'}
                    class="px-1 py-3 text-sm font-medium transition-colors border-b-2"
                    class:border-blue-500={searchMode === 'local'}
                    class:text-blue-400={searchMode === 'local'}
                    class:border-transparent={searchMode !== 'local'}
                    class:text-gray-400={searchMode !== 'local'}
                    class:hover:text-gray-200={searchMode !== 'local'}
                    class:hover:border-gray-500={searchMode !== 'local'}
            >
                Local Search
            </button>
            <button
                    on:click={() => searchMode = 'remote'}
                    class="px-1 py-3 text-sm font-medium transition-colors border-b-2"
                    class:border-blue-500={searchMode === 'remote'}
                    class:text-blue-400={searchMode === 'remote'}
                    class:border-transparent={searchMode !== 'remote'}
                    class:text-gray-400={searchMode !== 'remote'}
                    class:hover:text-gray-200={searchMode !== 'remote'}
                    class:hover:border-gray-500={searchMode !== 'remote'}
            >
                Remote Server
            </button>
        </nav>
    </div>

    <!-- Search Forms -->
    <div class="p-6 mb-8 bg-gray-800 rounded-lg shadow-md">
        {#if searchMode === 'local'}
            <div class="grid grid-cols-1 gap-6 md:grid-cols-3">
                <div class="md:col-span-2">
                    <label for="local-path" class="block mb-2 text-sm font-medium text-gray-300">Directory Path</label>
                    <div class="flex">
                        <Input id="local-path" bind:value={localPath} placeholder="e.g., C:\Users\YourUser\Documents" />
                        <Button on:click={chooseDirectory} class="ml-2 flex-shrink-0">Choose...</Button>
                    </div>
                </div>
                <div>
                    <label for="local-pattern" class="block mb-2 text-sm font-medium text-gray-300">File Pattern</label>
                    <Input id="local-pattern" bind:value={localPattern} placeholder="e.g., *.log, access-*.log" />
                </div>
            </div>
        {:else}
            <div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
                <div>
                    <label for="remote-host" class="block mb-2 text-sm font-medium text-gray-300">Host / IP Address</label>
                    <Input id="remote-host" bind:value={remoteHost} placeholder="e.g., 192.168.1.100" />
                </div>
                <div>
                    <label for="remote-port" class="block mb-2 text-sm font-medium text-gray-300">Port</label>
                    <Input id="remote-port" type="number" bind:value={remotePort} />
                </div>
                <div>
                    <label for="remote-user" class="block mb-2 text-sm font-medium text-gray-300">Username</label>
                    <Input id="remote-user" bind:value={remoteUser} placeholder="e.g., root, admin" />
                </div>
                <div>
                    <label for="remote-password" class="block mb-2 text-sm font-medium text-gray-300">Password</label>
                    <Input id="remote-password" type="password" bind:value={remotePassword} />
                </div>
                <div class="md:col-span-2 lg:col-span-4">
                    <label for="remote-log-path" class="block mb-2 text-sm font-medium text-gray-300">Log Directory on Server</label>
                    <Input id="remote-log-path" bind:value={remoteLogPath} placeholder="e.g., /var/log/nginx/" />
                </div>
            </div>
        {/if}
        <div class="flex justify-end mt-6">
            <Button on:click={handleSearch} disabled={isLoading}>
                {#if isLoading}
                    Searching...
                {:else}
                    Search Logs
                {/if}
            </Button>
        </div>
    </div>

    <!-- Results Area -->
    <div class="p-4 bg-gray-800 rounded-lg min-h-[20rem] flex items-center justify-center">
        {#if isLoading}
            <!-- Loading Spinner -->
            <div class="flex flex-col items-center text-gray-400">
                <svg class="w-12 h-12 mb-4 text-blue-500 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>Searching for log files...</span>
            </div>
        {:else if error}
            <!-- Error Message -->
            <div class="text-center text-red-400">
                <p class="text-lg font-semibold">An Error Occurred</p>
                <p>{error}</p>
            </div>
        {:else if searchResults.length > 0}
            <!-- Results List -->
            <div class="w-full">
                <h3 class="mb-4 text-xl font-semibold text-white">Found {searchResults.length} file(s)</h3>
                <ul class="space-y-3">
                    {#each searchResults as result}
                        <li class="flex items-center p-4 transition-colors bg-gray-700 rounded-md hover:bg-gray-600">
                            <!-- File Icon -->
                            <svg class="w-6 h-6 mr-4 text-gray-400 flex-shrink-0" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
                            </svg>
                            <div class="flex-grow overflow-hidden">
                                <p class="font-medium text-white truncate" title={result.name}>{result.name}</p>
                                <p class="text-sm text-gray-400 truncate" title={result.path}>{result.path}</p>
                            </div>
                            <div class="flex-shrink-0 hidden ml-6 text-sm text-right md:block">
                                <p class="text-gray-300">{result.size}</p>
                                <p class="text-gray-500">{result.modified}</p>
                            </div>
                        </li>
                    {/each}
                </ul>
            </div>
        {:else}
            <!-- Initial/No Results State -->
            <div class="text-center text-gray-500">
                <p>Enter search criteria and click "Search Logs" to begin.</p>
            </div>
        {/if}
    </div>
</div>