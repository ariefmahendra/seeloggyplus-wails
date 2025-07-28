<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { dto } from '../../../wailsjs/go/models';
    import {
        Button,
        Modal,
        DataTable,
        Tag,
        Loading
    } from 'carbon-components-svelte';
    import {
        Connect,
        CloseFilled,
        Settings,
        BareMetalServer
    } from 'carbon-icons-svelte';
    import { ConnectSession, CloseSession } from '../../../wailsjs/go/main/App';

    export let sessions: dto.SessionManagerDto[];
    export let selectedSessionId: string;

    const dispatch = createEventDispatcher();

    let showSessionModal = false;
    let connectingSessionId = '';
    let disconnectingSessionId = '';

    async function connectToSession(sessionId: string) {
        try {
            connectingSessionId = sessionId;
            await ConnectSession(sessionId);
            dispatch('sessionConnected', sessionId);
        } catch (err) {
            console.error('Failed to connect to session:', err);
            dispatch('error', `Failed to connect: ${err}`);
        } finally {
            connectingSessionId = '';
        }
    }

    async function disconnectSession(sessionId: string) {
        try {
            disconnectingSessionId = sessionId;
            await CloseSession(sessionId);
            if (selectedSessionId === sessionId) {
                dispatch('sessionDisconnected');
            }
        } catch (err) {
            console.error('Failed to disconnect session:', err);
            dispatch('error', `Failed to disconnect: ${err}`);
        } finally {
            disconnectingSessionId = '';
        }
    }

    function getSessionStatus(session: dto.SessionManagerDto) {
        return selectedSessionId === session.id ? 'Connected' : 'Disconnected';
    }

    function getSessionStatusKind(session: dto.SessionManagerDto) {
        return selectedSessionId === session.id ? 'green' : 'gray';
    }

    $: sessionRows = sessions.map((session, index) => ({
        id: index,
        name: session.serverInfo?.name || 'Unknown',
        address: `${session.serverInfo?.address}:${session.serverInfo?.port}`,
        user: session.serverInfo?.user || 'N/A',
        status: getSessionStatus(session),
        session: session,
    }));
</script>

<div class="flex items-center gap-4">
    <Button
            kind="tertiary"
            size="small"
            icon={Settings}
            on:click={() => showSessionModal = true}
    >
        Manage Sessions ({sessions.length})
    </Button>

    {#if selectedSessionId}
        <div class="flex items-center gap-2">
            <BareMetalServer size={16} class="text-green-600" />
            <Tag type="green">
                Connected
            </Tag>
        </div>
    {/if}
</div>

<Modal
        bind:open={showSessionModal}
        modalHeading="Remote Session Management"
        primaryButtonText="Close"
        secondaryButtonText=""
        on:click:button--primary={() => showSessionModal = false}
        size="lg"
>
    <div class="space-y-4">
        {#if sessions.length === 0}
            <div class="text-center py-12 text-gray-500">
                <BareMetalServer size={48} class="mx-auto mb-4 text-gray-300" />
                <p class="text-lg">No remote sessions available</p>
                <p class="text-sm">Please add a server connection first</p>
            </div>
        {:else}
            <DataTable
                    headers={[
                      { key: 'name', value: 'Server Name' },
                      { key: 'address', value: 'Address' },
                      { key: 'user', value: 'User' },
                      { key: 'status', value: 'Status' }
                    ]}
                    rows={sessionRows}
            >
                <svelte:fragment slot="cell" let:row let:cell>
                    {#if cell.key === 'status'}
                        <Tag type={getSessionStatusKind(row.session)}>
                            {row.status}
                        </Tag>
                    {:else if cell.key === 'actions'}
                        <div class="flex gap-2">
                            {#if selectedSessionId === row.session.id}
                                <Button
                                        kind="danger-tertiary"
                                        size="small"
                                        icon={CloseFilled}
                                        disabled={disconnectingSessionId === row.session.id}
                                        on:click={() => disconnectSession(row.session.id)}
                                >
                                    {#if disconnectingSessionId === row.session.id}
                                        <Loading withOverlay={false} small />
                                        Disconnecting...
                                    {:else}
                                        Disconnect
                                    {/if}
                                </Button>
                            {:else}
                                <Button
                                        kind="primary"
                                        size="small"
                                        icon={Connect}
                                        disabled={connectingSessionId === row.session.id}
                                        on:click={() => connectToSession(row.session.id)}
                                >
                                    {#if connectingSessionId === row.session.id}
                                        <Loading withOverlay={false} small />
                                        Connecting...
                                    {:else}
                                        Connect
                                    {/if}
                                </Button>
                            {/if}
                        </div>
                    {:else}
                        {cell.value}
                    {/if}
                </svelte:fragment>
            </DataTable>
        {/if}
    </div>
</Modal>