<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { Breadcrumb, BreadcrumbItem } from 'carbon-components-svelte';
    import { Home, ChevronRight } from 'carbon-icons-svelte';

    export let currentPath: string;
    export let isRemote: boolean;

    const dispatch = createEventDispatcher();

    $: pathParts = currentPath ? currentPath.split(/[/\\]/).filter(Boolean) : [];

    function navigateTo(index: number) {
        let newPath: string;

        if (index === -1) {
            newPath = isRemote ? '/' : '';
        } else {
            const separator = currentPath.includes('\\') ? '\\' : '/';
            newPath = separator + pathParts.slice(0, index + 1).join(separator);
        }

        dispatch('navigate', newPath);
    }
</script>

{#if currentPath}
    <Breadcrumb noTrailingSlash>
        <BreadcrumbItem on:click={() => navigateTo(-1)}>
            <div class="flex items-center gap-1 cursor-pointer hover:text-blue-600">
                <Home size={14} />
                <span class="text-sm">{isRemote ? 'Root' : 'Home'}</span>
            </div>
        </BreadcrumbItem>

        {#each pathParts as part, index}
            <BreadcrumbItem
                    on:click={() => navigateTo(index)}
                    class="cursor-pointer hover:text-blue-600"
            >
                <span class="text-sm" title={part}>{part}</span>
            </BreadcrumbItem>
        {/each}
    </Breadcrumb>
{:else}
    <div class="text-sm text-gray-500 italic">
        {isRemote ? 'No remote session selected' : 'Select a drive to browse'}
    </div>
{/if}