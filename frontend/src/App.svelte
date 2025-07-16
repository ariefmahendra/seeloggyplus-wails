<script lang="ts">
    import { onDestroy, onMount } from 'svelte';
    import Router, { link, location } from 'svelte-spa-router';
    import { wrap } from 'svelte-spa-router/wrap';
    import { fade } from 'svelte/transition';
    import 'remixicon/fonts/remixicon.css';
    import Footer from './components/layout/Footer.svelte';

    const SIDEBAR_CONFIG = {
        DEFAULT_WIDTH: 224,
        COLLAPSED_WIDTH: 68,
        MIN_WIDTH: 180,
        MAX_WIDTH: 400,
        STORAGE_KEY_WIDTH: 'sidebar_width',
        STORAGE_KEY_COLLAPSED: 'sidebar_collapsed'
    };

    // --- State Management ---
    let isCollapsed = false;
    let sidebarWidth = SIDEBAR_CONFIG.DEFAULT_WIDTH;
    let isResizing = false;
    let useTransition = true;

    // --- Data Navigation ---
    interface NavItem {
        name: string;
        icon: string;
        href: string;
    }
    const navItems: NavItem[] = [
        { name: 'Home', icon: 'ri-home-4-line', href: '#/' },
        { name: 'File Finder', icon: 'ri-folder-3-line', href: '#/files' },
        { name: 'Servers', icon: 'ri-server-line', href: '#/server-management' },
        { name: 'Settings', icon: 'ri-settings-3-line', href: '#/settings' },
        { name: 'Help', icon: 'ri-question-line', href: '#/help' }
    ];

    // --- Lifecycle & Persistence ---
    onMount(() => {
        const savedWidth = localStorage.getItem(SIDEBAR_CONFIG.STORAGE_KEY_WIDTH);
        const savedCollapsed = localStorage.getItem(SIDEBAR_CONFIG.STORAGE_KEY_COLLAPSED);

        if (savedCollapsed !== null) {
            isCollapsed = savedCollapsed === 'true';
        }
        if (savedWidth !== null && !isCollapsed) {
            sidebarWidth = Math.max(SIDEBAR_CONFIG.MIN_WIDTH, parseInt(savedWidth, 10));
        }

        window.addEventListener('mousemove', doResize);
        window.addEventListener('mouseup', stopResize);
    });

    onDestroy(() => {
        window.removeEventListener('mousemove', doResize);
        window.removeEventListener('mouseup', stopResize);
    });

    $: if (typeof window !== 'undefined') {
        localStorage.setItem(SIDEBAR_CONFIG.STORAGE_KEY_WIDTH, String(sidebarWidth));
        localStorage.setItem(SIDEBAR_CONFIG.STORAGE_KEY_COLLAPSED, String(isCollapsed));
    }

    function toggleCollapse() {
        useTransition = true;
        isCollapsed = !isCollapsed;
        if (!isCollapsed && sidebarWidth < SIDEBAR_CONFIG.MIN_WIDTH) {
            sidebarWidth = SIDEBAR_CONFIG.DEFAULT_WIDTH;
        }
    }

    function startResize(event: MouseEvent) {
        event.preventDefault();
        isResizing = true;
        useTransition = false;
        document.body.style.userSelect = 'none';
        document.body.style.cursor = 'col-resize';
    }

    function doResize(event: MouseEvent) {
        if (!isResizing) return;

        if (isCollapsed) {
            isCollapsed = false;
        }

        let newWidth = event.clientX;
        if (newWidth < SIDEBAR_CONFIG.MIN_WIDTH) newWidth = SIDEBAR_CONFIG.MIN_WIDTH;
        if (newWidth > SIDEBAR_CONFIG.MAX_WIDTH) newWidth = SIDEBAR_CONFIG.MAX_WIDTH;
        sidebarWidth = newWidth;
    }

    function stopResize() {
        if (!isResizing) return;
        isResizing = false;
        document.body.style.userSelect = '';
        document.body.style.cursor = '';
    }

    const routes = {
        '/': wrap({ asyncComponent: () => import('./routes/Home.svelte') }),
        '/files': wrap({ asyncComponent: () => import('./routes/Files.svelte') }),
        '/server-management': wrap({ asyncComponent: () => import('./routes/ServerManagement.svelte') }),
        '/settings': wrap({ asyncComponent: () => import('./routes/Settings.svelte') }),
        '/help': wrap({ asyncComponent: () => import('./routes/Helps.svelte') }),
        '*': wrap({ asyncComponent: () => import('./routes/NotFound.svelte') }),
    };
</script>

<div class="flex h-screen overflow-hidden bg-gray-900 text-white font-sans">
    <nav
            class="group relative flex flex-shrink-0 flex-col bg-gray-800 p-3"
            class:transition-width={useTransition}
            class:duration-300={useTransition}
            class:ease-in-out={useTransition}
            style="width: {isCollapsed ? SIDEBAR_CONFIG.COLLAPSED_WIDTH : sidebarWidth}px;"
    >
        <div class="flex items-center mb-4 px-2 py-2" class:justify-center={isCollapsed}>
            {#if !isCollapsed}
                <h1 class="text-lg font-semibold text-white whitespace-nowrap" transition:fade={{ duration: 150 }}>
                    SeeLoggyPlus
                </h1>
            {/if}
        </div>

        <div class="flex flex-col space-y-1">
            {#each navItems as item (item.href)}
                <a
                        href={item.href}
                        use:link
                        class="w-full text-left px-3 py-2 rounded-md text-sm font-medium transition-colors duration-150 flex items-center space-x-3"
                        class:justify-center={isCollapsed}
                        class:bg-blue-600={'#' + $location === item.href}
                        class:text-white={'#' + $location === item.href}
                        class:text-gray-300={'#' + $location !== item.href}
                        class:hover:bg-gray-700={'#' + $location !== item.href}
                        title={isCollapsed ? item.name : ''}
                >
                    <i class="text-xl flex-shrink-0 {item.icon}" />
                    {#if !isCollapsed}
                        <span class="whitespace-nowrap" transition:fade={{ duration: 150 }}>{item.name}</span>
                    {/if}
                </a>
            {/each}
        </div>

        <div class="mt-auto">
            <button
                    class="w-full text-left px-3 py-2 rounded-md text-sm font-medium transition-colors duration-150 flex items-center space-x-3 text-gray-400 hover:bg-gray-700 hover:text-white"
                    class:justify-center={isCollapsed}
                    on:click={toggleCollapse}
                    title={isCollapsed ? 'Expand Sidebar' : 'Collapse Sidebar'}
            >
                <i
                        class="text-xl flex-shrink-0 transition-transform duration-300"
                        class:ri-arrow-right-s-line={isCollapsed}
                        class:ri-arrow-left-s-line={!isCollapsed}
                />
                {#if !isCollapsed}
                    <span class="whitespace-nowrap" transition:fade={{ duration: 150 }}>Collapse</span>
                {/if}
            </button>
        </div>

        <div
                role="separator"
                aria-orientation="vertical"
                class="absolute top-0 right-0 h-full w-1.5 cursor-col-resize select-none touch-none opacity-0 group-hover:opacity-100 transition-opacity"
                on:mousedown={startResize}
        >
            <div class="h-full w-px bg-blue-500/50"></div>
        </div>
    </nav>

    <div class="flex flex-1 flex-col overflow-y-auto">
        <main class="flex-grow p-6">
            <Router {routes} />
        </main>
        <Footer />
    </div>
</div>

<style>
    .transition-width {
        transition-property: width;
    }
</style>