<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { Modal, Button } from 'carbon-components-svelte';
    import { TrashCan, ServerProxy } from 'carbon-icons-svelte';
    import type { Server } from '../types/server.types';
    import { ServerService } from '../services/server.service';

    export let isOpen = false;
    export let server: Server;

    let isLoading = false;

    const dispatch = createEventDispatcher<{
        success: void;
    }>();

    function closeModal() {
        isOpen = false;
    }

    async function handleDelete() {
        isLoading = true;
        try {
            await ServerService.deleteServer(server.id);
            dispatch('success');
            closeModal();
        } catch (error) {
            console.error('Failed to delete server:', error);
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
        modalHeading="Delete Server"
        primaryButtonText="Delete"
        secondaryButtonText="Cancel"
        on:click:button--primary={handleDelete}
        on:click:button--secondary={closeModal}
        on:close={closeModal}
        size="xs"
        danger
>
    <div slot="heading">
        <div class="flex items-center space-x-2">
            <TrashCan size={24} />
            <span>Delete Server</span>
        </div>
    </div>

    <div class="space-y-4">
        <p class="text-sm text-gray-700">
            Are you sure you want to delete this server? This action cannot be undone.
        </p>

        <div class="bg-gray-50 rounded-lg p-4">
            <div class="flex items-center space-x-3">
                <ServerProxy size={20} class="text-gray-500" />
                <div>
                    <p class="font-medium text-gray-900">{server.name}</p>
                    <p class="text-sm text-gray-500">{server.address}:{server.port}</p>
                </div>
            </div>
        </div>
    </div>

    <div slot="footer" class="flex justify-end space-x-3">
        <Button
                kind="secondary"
                on:click={handleCancel}
                disabled={isLoading}
        >
            Cancel
        </Button>
        <Button
                kind="danger"
                icon={TrashCan}
                on:click={handleDelete}
                disabled={isLoading}
        >
            {#if isLoading}
                Deleting...
            {:else}
                Delete Server
            {/if}
        </Button>
    </div>
</Modal>