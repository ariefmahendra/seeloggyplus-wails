<script lang="ts">
    import Button from '../components/common/Button.svelte';

    // Data untuk Frequently Asked Questions (FAQ)
    // Anda bisa dengan mudah menambah atau mengubah pertanyaan di sini
    const faqs = [
        {
            question: 'Bagaimana cara memulai dengan SeeLoggyPlus?',
            answer:
                'Untuk memulai, navigasikan ke halaman "Home". Anda dapat memasukkan nama Anda di kolom yang tersedia dan klik tombol "Greet" untuk mendapatkan respons dari aplikasi. Ini adalah cara sederhana untuk memverifikasi bahwa semuanya berjalan dengan baik.',
        },
        {
            question: 'Di mana saya bisa menemukan file log?',
            answer:
                'Fungsionalitas manajemen file log akan tersedia di halaman "File". Dari sana, Anda dapat membuka, melihat, dan menganalisis file log yang Anda perlukan.',
        },
        {
            question: 'Bisakah saya mengubah tema aplikasi?',
            answer:
                'Tentu saja. Anda dapat mengubah tema antara mode terang (light) dan gelap (dark) dari halaman "Settings". Pengaturan Anda akan disimpan untuk sesi berikutnya.',
        },
        {
            question: 'Apakah aplikasi ini gratis?',
            answer:
                'Ya, SeeLoggyPlus adalah perangkat lunak sumber terbuka (open-source) dan gratis untuk digunakan oleh siapa saja.',
        },
    ];

    // State untuk mengontrol item FAQ mana yang sedang terbuka
    let activeFaq: number | null = null;

    function toggleFaq(index: number) {
        if (activeFaq === index) {
            activeFaq = null; // Tutup jika sudah terbuka
        } else {
            activeFaq = index; // Buka item yang baru
        }
    }
</script>

<!-- Kontainer utama dengan padding dan batas lebar untuk keterbacaan -->
<div class="max-w-4xl mx-auto">
    <h1 class="mb-4 text-4xl font-bold text-white">Help & Support</h1>
    <p class="mb-12 text-lg text-gray-400">
        Kami di sini untuk membantu. Temukan jawaban atas pertanyaan Anda di bawah ini.
    </p>

    <!-- Bagian Frequently Asked Questions (FAQ) -->
    <section>
        <h2 class="mb-6 text-2xl font-semibold text-gray-100">Frequently Asked Questions</h2>
        <div class="space-y-4">
            {#each faqs as faq, index}
                <div class="overflow-hidden border border-gray-700 rounded-lg">
                    <!-- Tombol untuk membuka/menutup akordeon -->
                    <button
                            on:click={() => toggleFaq(index)}
                            class="flex items-center justify-between w-full p-5 text-left transition-colors focus:outline-none"
                            class:bg-gray-800={activeFaq !== index}
                            class:bg-gray-700={activeFaq === index}
                            class:hover:bg-gray-700={activeFaq !== index}
                    >
                        <span class="font-medium text-white">{faq.question}</span>
                        <!-- Ikon panah (chevron) -->
                        <svg
                                class="w-5 h-5 text-gray-400 transform transition-transform duration-300"
                                class:rotate-180={activeFaq === index}
                                xmlns="http://www.w3.org/2000/svg"
                                viewBox="0 0 20 20"
                                fill="currentColor"
                        >
                            <path
                                    fill-rule="evenodd"
                                    d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
                                    clip-rule="evenodd"
                            />
                        </svg>
                    </button>

                    <!-- Konten jawaban yang bisa disembunyikan/ditampilkan -->
                    {#if activeFaq === index}
                        <div class="p-5 bg-gray-800 border-t border-gray-700">
                            <p class="text-gray-300 leading-relaxed">
                                {faq.answer}
                            </p>
                        </div>
                    {/if}
                </div>
            {/each}
        </div>
    </section>

    <!-- Bagian Kontak -->
    <section class="p-8 mt-16 text-center bg-gray-800 rounded-lg shadow-lg">
        <h2 class="mb-4 text-2xl font-semibold text-white">Masih Butuh Bantuan?</h2>
        <p class="max-w-2xl mx-auto mb-6 text-gray-400">
            Jika Anda tidak dapat menemukan jawaban yang Anda cari di FAQ, jangan ragu untuk menghubungi tim
            dukungan kami.
        </p>
        <Button on:click={() => (window.location.href = 'mailto:support@seeloggyplus.com')}>
            Contact Support
        </Button>
    </section>
</div>