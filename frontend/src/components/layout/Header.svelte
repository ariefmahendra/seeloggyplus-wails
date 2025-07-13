<script lang="ts">
    import { createEventDispatcher } from 'svelte';

    export let activePage: 'home' | 'settings' = 'home';

    const dispatch = createEventDispatcher();

    function navigate(page: 'home' | 'settings') {
        dispatch('navigate', { page });
    }
</script>

<header class="flex items-center justify-between p-4 bg-gray-800 shadow-lg">
    <div class="text-xl font-bold text-white">
        SeeLoggyPlus
    </div>

    <nav class="flex space-x-2">
        <button
                on:click={() => navigate('home')}
                class="px-4 py-2 rounded-md text-sm font-medium transition-colors"
                class:bg-blue-600={activePage === 'home'}
                class:text-white={activePage === 'home'}
                class:text-gray-300={activePage !== 'home'}
                class:hover:bg-gray-700={activePage !== 'home'}
        >
            Home
        </button>
        <button
                on:click={() => navigate('settings')}
                class="px-4 py-2 rounded-md text-sm font-medium transition-colors"
                class:bg-blue-600={activePage === 'settings'}
                class:text-white={activePage === 'settings'}
                class:text-gray-300={activePage !== 'settings'}
                class:hover:bg-gray-700={activePage !== 'settings'}
        >
            Settings
        </button>
    </nav>
</header>