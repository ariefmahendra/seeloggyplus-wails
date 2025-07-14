<script lang="ts">
    import { createEventDispatcher, onMount } from 'svelte';
    import { slide } from 'svelte/transition';

    type Option = {
        value: string;
        label: string;
        icon?: string;
    };

    export let options: Option[] = [];
    export let value: string;
    export let placeholder = 'Select an option';
    export let id: string = '';
    export let name: string = '';
    export let required: boolean = false;

    let isOpen = false;
    let selectedOption: Option | undefined;
    let selectElement: HTMLElement;

    $: selectedOption = options.find((opt) => opt.value === value);

    const dispatch = createEventDispatcher();

    function toggleDropdown() {
        isOpen = !isOpen;
    }

    function selectOption(option: Option) {
        value = option.value;
        isOpen = false;
        dispatch('change', value);
    }

    function handleClickOutside(event: MouseEvent) {
        if (selectElement && !selectElement.contains(event.target as Node)) {
            isOpen = false;
        }
    }

    onMount(() => {
        document.addEventListener('click', handleClickOutside, true);
        return () => {
            document.removeEventListener('click', handleClickOutside, true);
        };
    });
</script>

<div class="relative" bind:this={selectElement}>
    <input type="hidden" name={name || id} bind:value {required} />
    <button
            type="button"
            on:click={toggleDropdown}
            {id}
            class="relative w-full cursor-pointer rounded-md border border-gray-600 bg-gray-700 py-2 pl-3 pr-10 text-left text-white shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            aria-haspopup="listbox"
            aria-expanded={isOpen}
    >
        <span class="flex items-center">
            {#if selectedOption?.icon}
                <i class="{selectedOption.icon} mr-3 text-lg text-gray-400"></i>
            {/if}
            <span class="block truncate">{selectedOption?.label || placeholder}</span>
        </span>
        <span class="pointer-events-none absolute inset-y-0 right-0 ml-3 flex items-center pr-2">
            <i class="ri-arrow-down-s-line text-gray-400" aria-hidden="true"></i>
        </span>
    </button>

    {#if isOpen}
        <ul
                transition:slide={{ duration: 150 }}
                class="absolute z-10 mt-1 max-h-56 w-full overflow-auto rounded-md bg-gray-700 py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm"
                tabindex="-1"
                role="listbox"
        >
            {#each options as option (option.value)}
                <li
                        on:click={() => selectOption(option)}
                        class="relative cursor-pointer select-none py-2 pl-3 pr-9 text-gray-300 hover:bg-gray-600 hover:text-white"
                        role="option"
                >
                    <div class="flex items-center">
                        {#if option.icon}<i class="{option.icon} mr-3 text-lg text-gray-400"></i>{/if}
                        <span class="font-normal block truncate">{option.label}</span>
                    </div>
                </li>
            {/each}
        </ul>
    {/if}
</div>
