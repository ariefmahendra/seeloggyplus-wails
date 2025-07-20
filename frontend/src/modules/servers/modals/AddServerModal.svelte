<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { Modal } from 'carbon-components-svelte';
    import { Add } from 'carbon-icons-svelte';
    import ServerForm from '../components/ServerForm.svelte';
    import type { ServerFormData } from '../types/server.types';
    import { ServerService } from '../services/serverService';

    export let isOpen = false;

    let isLoading = false;
    let serverForm: ServerForm;
    let formData: ServerFormData = {
        name: '',
        address: '',
        port: 22,
        user: '',
        password: ''
    };

    const dispatch = createEventDispatcher<{
        success: void;
    }>();

    function closeModal() {
        isOpen = false;
    }

    async function handleSubmit() {
        if (!serverForm) return;

        // Call validation and submit from ServerForm
        const isValid = serverForm.validateAndSubmit();
        if (!isValid) return;

        isLoading = true;
        try {
            await ServerService.addServer(formData);
            dispatch('success');
            closeModal();
        } catch (error) {
            console.error('Failed to create server:', error);
        } finally {
            isLoading = false;
        }
    }

    function handleCancel() {
        closeModal();
    }
</script>

<Modal
        bind:open={isOpen}
        modalHeading="Add Server"
        primaryButtonText={isLoading ? "Adding..." : "Add Server"}
        secondaryButtonText="Cancel"
        primaryButtonDisabled={isLoading}
        selectorPrimaryFocus="#server-name"
        on:click:button--primary={handleSubmit}
        on:click:button--secondary={handleCancel}
        on:close={closeModal}
        size="sm"
        hasScrollingContent
>
    <div slot="heading">
        <div class="flex items-center space-x-2">
            <Add size={24} />
            <span>Add Server</span>
        </div>
    </div>

    <div class="p-4">
        <ServerForm
                bind:this={serverForm}
                bind:formData
                {isLoading}
        />
    </div>
</Modal>