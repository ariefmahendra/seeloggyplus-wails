<script lang="ts">
    import { ListLocalFiles, GetLocalDrives } from '../../../wailsjs/go/main/App';
    import { Button, DataTable, CodeSnippet } from "carbon-components-svelte";
    import { onMount } from "svelte";
    import { dto } from "../../../wailsjs/go/models";
    import { ArrowUp } from "carbon-icons-svelte";

    // Interface untuk memetakan data dari backend Go
    // Menambahkan 'path' untuk menyimpan path lengkap demi kemudahan navigasi
    interface FileInfoMapping {
        id: number;
        name: string;
        path: string; // Menyimpan path lengkap (e.g., "C:\" atau "C:\Users")
        size: number;
        modified: any;
        type: 'Directory' | 'File' | 'Drive';
        isDir: boolean;
        mode: any;
    }

    let files: FileInfoMapping[] = [];
    // Path kosong ("") diinterpretasikan sebagai permintaan untuk menampilkan daftar drive.
    let currentPath: string = '';
    let isLoading: boolean = true;
    let errorMessage: string = '';

    // Statement reaktif Svelte: akan berjalan setiap kali `currentPath` berubah.
    $: if (currentPath !== undefined) {
        loadListFiles(currentPath);
    }

    // Variabel reaktif untuk menonaktifkan tombol "Kembali" saat di level daftar drive.
    $: isAtRoot = currentPath === '';

    /**
     * Memuat daftar file atau drive dari backend Wails.
     * @param path - Jika string kosong (""), fungsi akan memanggil GetLocalDrives().
     * Jika tidak, akan memanggil ListLocalFiles(path).
     */
    async function loadListFiles(path: string) {
        isLoading = true;
        errorMessage = '';
        try {
            if (path === '') {
                // KASUS 1: Memuat daftar drive saat path kosong (tampilan awal)
                const driveInfos: dto.DriveInfo[] = await GetLocalDrives();
                files = driveInfos.map((drive, i) => ({
                    id: i,
                    name: `${drive.label || drive.name} (${drive.name})`, // Tampilkan label dan nama drive
                    path: drive.path, // Gunakan path dari DTO, e.g., "C:\"
                    size: 0,
                    modified: 'N/A',
                    type: 'Drive',
                    isDir: true, // Semua drive dianggap direktori
                    mode: null,
                }));
            } else {
                // KASUS 2: Memuat file dan folder dari path tertentu
                const fileInfos: dto.FileInfo[] = await ListLocalFiles(path);
                files = fileInfos.map((fileInfo, i) => {
                    // Gabungkan path saat ini dengan nama file/folder untuk navigasi selanjutnya
                    let basePath = path;
                    if (basePath.endsWith(':')) {
                        basePath += '/';
                    }
                    const fullPath = [basePath.replace(/\/$/, ''), fileInfo.name].join('/');

                    return {
                        id: i,
                        name: fileInfo.name,
                        path: fullPath, // Simpan path lengkap
                        size: fileInfo.size,
                        modified: fileInfo.modTime ? new Date(fileInfo.modTime).toLocaleString() : 'N/A',
                        type: fileInfo.isDir ? 'Directory' : 'File',
                        isDir: fileInfo.isDir,
                        mode: fileInfo.mode,
                    };
                });
            }
        } catch (error) {
            console.error("Failed to load items:", error);
            errorMessage = `Gagal memuat dari path: '${path || 'Drives'}'. Error: ${error}`;
            // Jika gagal, coba kembali ke direktori induk
            currentPath = getParentDirectory(currentPath);
        } finally {
            isLoading = false;
        }
    }

    /**
     * Menangani event klik pada baris tabel. Navigasi terjadi berdasarkan 'path' yang tersimpan.
     */
    function handleRowClick(event: CustomEvent<FileInfoMapping>) {
        const clickedItem = event.detail;

        // Hanya navigasi jika item yang diklik adalah direktori atau drive
        if (clickedItem.isDir) {
            // Langsung gunakan path yang sudah lengkap dari item yang diklik
            currentPath = clickedItem.path;
        } else {
            // Logika untuk file (misalnya membuka file) bisa ditambahkan di sini
            console.log("File clicked:", clickedItem.path);
        }
    }

    /**
     * Menghitung path direktori induk.
     * Contoh: "C:/Users/Test" -> "C:/Users" -> "C:" -> "" (daftar drive)
     */
    function getParentDirectory(path: string): string {
        if (path === '' || path === null) {
            return ''; // Sudah di level tertinggi
        }

        // Normalisasi backslash menjadi forward slash untuk konsistensi
        const normalizedPath = path.replace(/\\/g, '/').replace(/\/$/, ''); // Hapus slash di akhir

        // Jika path adalah drive (e.g., "C:"), induknya adalah daftar drive
        if (normalizedPath.length === 2 && normalizedPath.endsWith(':')) {
            return '';
        }

        const lastSlashIndex = normalizedPath.lastIndexOf('/');

        // Jika tidak ada slash atau hanya ada satu di root (e.g., "/"), kembali ke drive
        if (lastSlashIndex <= 0) {
            return normalizedPath.substring(0, 2); // "C:"
        }

        // Potong string sampai sebelum slash terakhir
        return normalizedPath.substring(0, lastSlashIndex);
    }

    /**
     * Fungsi untuk tombol "Up" / "Kembali".
     */
    function goUp() {
        if (!isAtRoot) {
            currentPath = getParentDirectory(currentPath);
        }
    }

</script>

<div class="p-4 bg-gray-100 min-h-screen font-sans">
    <div class="container mx-auto bg-white p-6 rounded-lg shadow-md">
        <h1 class="text-2xl font-bold mb-4">File Browser</h1>

        <!-- Kontrol Navigasi -->
        <div class="flex justify-between items-center mb-4 p-2 bg-gray-50 rounded-md">
            <!-- Tombol Kembali -->
            <Button
                    kind="ghost"
                    icon={ArrowUp}
                    on:click={goUp}
                    disabled={isAtRoot || isLoading}
            >
                Up
            </Button>

            <!-- Tampilan Path Saat Ini -->
            <div class="flex-grow mx-4">
                <CodeSnippet type="single" wrapText={true}>
                    {currentPath || 'Daftar Drive'}
                </CodeSnippet>
            </div>
        </div>

        <!-- Pesan Error -->
        {#if errorMessage}
            <div class="p-4 mb-4 text-sm text-red-700 bg-red-100 rounded-lg" role="alert">
                <span class="font-medium">Error!</span> {errorMessage}
            </div>
        {/if}

        <!-- Tabel Data -->
        <div class="overflow-x-auto">
            {#if isLoading}
                <p class="text-center p-8">Loading...</p>
            {:else}
                <DataTable
                        stickyHeader
                        size="compact"
                        headers={[
                        {key: 'name', value: 'Name'},
                        {key: 'size', value: 'Size (bytes)'},
                        {key: 'modified', value: 'Last Modified'},
                        {key: 'type', value: 'Type'},
                    ]}
                        rows={files}
                        on:click:row={handleRowClick}
                        sortable
                />
                {#if files.length === 0 && !isLoading}
                    <p class="text-center text-gray-500 p-8">Direktori ini kosong.</p>
                {/if}
            {/if}
        </div>
    </div>
</div>
