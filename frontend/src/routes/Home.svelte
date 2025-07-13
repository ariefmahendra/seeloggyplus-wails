<script lang="ts">
    import logo from '../assets/images/logo-universal.png';

    import Button from '../components/common/Button.svelte';
    import Input from '../components/common/Input.svelte';
    let resultText: string = "Please enter your name below 👇";
    let name: string = '';

</script>

<!-- Kontainer utama untuk halaman Home, dibuat agar konten berada di tengah -->
<div class="flex flex-col items-center justify-center h-full p-4 text-center">
    <img alt="Wails logo" class="w-1/2 max-w-xs my-8" src={logo} />

    <div class="my-6 text-xl text-gray-200">{resultText}</div>

    <!-- Kotak input yang berisi komponen Input dan Button -->
    <div class="flex items-center w-full max-w-md">
        <div class="flex-grow">
            <!-- Menggunakan komponen Input kustom -->
            <Input
                    bind:value={name}
                    placeholder="Enter your name..."
            />
        </div>
        <div class="ml-4">
            <!-- Menggunakan komponen Button kustom -->
            <Button>Greet</Button>
        </div>
    </div>
</div>