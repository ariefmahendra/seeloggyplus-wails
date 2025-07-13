<script lang="ts">
    import Button from '../components/common/Button.svelte';

    // Contoh state untuk halaman settings
    let theme: 'dark' | 'light' | 'system' = 'dark';
    let notificationsEnabled: boolean = true;
    let apiKey: string = '********************';

    function saveSettings() {
        // Di sini Anda akan menambahkan logika untuk menyimpan pengaturan,
        // misalnya memanggil fungsi dari Wails atau menyimpan ke localStorage.
        console.log('Settings saved!', { theme, notificationsEnabled, apiKey });
        alert('Settings have been saved!');
    }
</script>

<!-- Kontainer utama untuk halaman Settings -->
<div class="max-w-4xl p-8 mx-auto">
    <h1 class="mb-8 text-4xl font-bold text-white">Settings</h1>

    <!-- Menggunakan grid untuk tata letak yang rapi -->
    <div class="grid gap-8 md:grid-cols-1">
        <!-- Bagian Appearance -->
        <section class="p-6 bg-gray-800 rounded-lg shadow-md">
            <h2 class="mb-4 text-2xl font-semibold text-gray-100">Appearance</h2>
            <div class="flex items-center space-x-4">
                <label for="theme-select" class="text-gray-300">Theme:</label>
                <select
                        id="theme-select"
                        bind:value={theme}
                        class="px-3 py-2 text-white bg-gray-700 border border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                    <option value="dark">Dark</option>
                    <option value="light">Light</option>
                    <option value="system">System Default</option>
                </select>
            </div>
        </section>

        <!-- Bagian Notifications -->
        <section class="p-6 bg-gray-800 rounded-lg shadow-md">
            <h2 class="mb-4 text-2xl font-semibold text-gray-100">Notifications</h2>
            <label class="inline-flex items-center cursor-pointer">
                <input type="checkbox" bind:checked={notificationsEnabled} class="w-5 h-5 text-blue-600 bg-gray-700 border-gray-600 rounded focus:ring-blue-500" />
                <span class="ml-3 text-gray-300">Enable desktop notifications</span>
            </label>
        </section>

        <!-- Tombol Aksi -->
        <div class="flex justify-end pt-4">
            <Button on:click={saveSettings}>Save Changes</Button>
        </div>
    </div>
</div>