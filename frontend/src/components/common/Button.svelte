<script lang="ts">
    import {twMerge} from "tailwind-merge";

    export let type: 'button' | 'submit' | 'reset' = 'button';
    export let disabled: boolean = false;
    let customClass = $$props.class || '';

    const baseClasses = 'px-4 py-2 font-bold text-white bg-blue-600 rounded-md shadow-sm transition-colors duration-200 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-gray-900 disabled:bg-gray-500 disabled:cursor-not-allowed';

    $: finalClass = twMerge(baseClasses, customClass);
</script>

<button {type} {disabled} on:click class={finalClass}>
    <slot />
</button>