<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { Modal } from 'carbon-components-svelte';
    import { Edit } from 'carbon-icons-svelte';
    import ServerForm from '../components/ServerForm.svelte';
    import type { Server, ServerFormData, ServerUpdateRequest } from '../types/server.types';
    import { ServerService } from '../services/serverService';

    export let isOpen = false;
    export let server: Server;

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

    $: if (server && isOpen) {
        resetFormData();
    }

    function resetFormData() {
        formData = {
            name: server.name,
            address: server.address,
            port: server.port,
            user: server.user,
            password: server.password
        };
    }

    function closeModal() {
        isOpen = false;
    }

    async function handleSubmit() {
        if (!serverForm) return;

        const isValid = serverForm.validateAndSubmit();
        if (!isValid) return;

        isLoading = true;
        try {
            const updateRequest: ServerUpdateRequest = {
                id: server.id,
                ...formData
            };
            await ServerService.updateServer(updateRequest);
            dispatch('success');
            closeModal();
        } catch (error) {
            console.error('Failed to update server:', error);
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
        modalHeading="Edit Server"
        selectorPrimaryFocus="#server-name"
        primaryButtonText={isLoading ? "Updating..." : "Update Server"}
        secondaryButtonText="Cancel"
        primaryButtonDisabled={isLoading}
        on:click:button--primary={handleSubmit}
        on:click:button--secondary={handleCancel}
        on:close={closeModal}
        size="sm"
        hasScrollingContent
        on:open={resetFormData}
>
    <div slot="heading">
        <div class="flex items-center space-x-2">
            <Edit size={24} />
            <span>Edit Server</span>
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