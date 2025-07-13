<script lang="ts">
    export let value: string | number = '';
    export let type: 'text' | 'password' | 'email' | 'number' | 'search' | 'url' = 'text';
    export let placeholder: string = '';
    export let id: string | undefined = undefined;
    export let name: string | undefined = undefined;
    export let disabled: boolean = false;
    export let autocomplete: string = 'off';

    let customClass = $$props.class || '';

    const baseClasses = 'w-full px-3 py-2 text-gray-200 bg-gray-700 border border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors';
</script>

{#if type === 'number'}
    <input
            type="number"
            bind:value
            {id}
            {name}
            {placeholder}
            {disabled}
            {autocomplete}
            class="{baseClasses} {customClass}"
            on:input
            on:blur
            on:focus
            on:keydown
    />
{:else if type === 'password'}
    <input
            type="password"
            bind:value
            {id}
            {name}
            {placeholder}
            {disabled}
            {autocomplete}
            class="{baseClasses} {customClass}"
            on:input
            on:blur
            on:focus
            on:keydown
    />
{:else if type === 'email'}
    <input
            type="email"
            bind:value
            {id}
            {name}
            {placeholder}
            {disabled}
            {autocomplete}
            class="{baseClasses} {customClass}"
            on:input
            on:blur
            on:focus
            on:keydown
    />
{:else if type === 'search'}
    <input
            type="search"
            bind:value
            {id}
            {name}
            {placeholder}
            {disabled}
            {autocomplete}
            class="{baseClasses} {customClass}"
            on:input
            on:blur
            on:focus
            on:keydown
    />
{:else if type === 'url'}
    <input
            type="url"
            bind:value
            {id}
            {name}
            {placeholder}
            {disabled}
            {autocomplete}
            class="{baseClasses} {customClass}"
            on:input
            on:blur
            on:focus
            on:keydown
    />
{:else}
    <input
            type="text"
            bind:value
            {id}
            {name}
            {placeholder}
            {disabled}
            {autocomplete}
            class="{baseClasses} {customClass}"
            on:input
            on:blur
            on:focus
            on:keydown
    />
{/if}