<script lang="ts">
    import {Button, InlineNotification, NumberInput, TextInput} from 'carbon-components-svelte';
    import {ConnectionSignal} from 'carbon-icons-svelte';
    import type {ServerFormData, ServerValidationErrors} from '../types/server.types';
    import {ServerService} from '../services/serverService';
    import {tick} from "svelte";

    export let formData: ServerFormData = {
        name: '',
        address: '',
        port: 22,
        user: '',
        password: ''
    };
    export let errors: ServerValidationErrors = {};
    export let isLoading = false;

    let isTestingConnection = false;
    let testResult: { success: boolean; message: string } | null = null;
    let requestId: string | null = null;

    $: if (testResult) {
        tick().then(() => {
            scrollToElement('notification');
        })
    }

    function scrollToElement(id: string){
        const element = document.getElementById(id);
        if (element){
            element.scrollIntoView({behavior: 'smooth'});
        } else {
            console.warn(`Element with id ${id} not found`);
        }
    }

    function generateRequestId(): string {
        return crypto.randomUUID();
    }

    function validateForm(): boolean {
        errors = {};
        let isValid = true;

        if (!formData.name.trim()) {
            errors.name = 'Server name is required';
            isValid = false;
        }

        if (!formData.address.trim()) {
            errors.address = 'Server address is required';
            isValid = false;
        }

        if (formData.port < 1 || formData.port > 65535) {
            errors.port = 'Port must be between 1 and 65535';
            isValid = false;
        }

        if (!formData.user.trim()) {
            errors.user = 'Username is required';
            isValid = false;
        }

        if (!formData.password.trim()) {
            errors.password = 'Password is required';
            isValid = false;
        }

        return isValid;
    }

    export function validateAndSubmit(): boolean {
        return validateForm();
    }

    async function cancelTestConnection() {
        if (!requestId) return;

        try {
            await ServerService.cancelTestConnection(requestId);
            testResult = {success: false, message: 'Connection test cancelled'};
        } catch (error) {
            console.error('Failed to cancel test connection:', error);
            testResult = {success: false, message: 'Failed to cancel connection test'};
        } finally {
            isTestingConnection = false;
            requestId = null;
        }
    }

    async function testConnection() {
        if (!validateForm()) return;

        isTestingConnection = true;
        testResult = null;
        requestId = generateRequestId();

        try {
            await ServerService.testConnection(requestId, formData);
            testResult = {success: true, message: 'Connection successful!'};
        } catch (error) {
            testResult = {success: false, message: error};
        } finally {
            isTestingConnection = false;
            requestId = null;
        }
    }

    function clearTestResult() {
        testResult = null;
    }
</script>

{#if testResult}
        <InlineNotification
                id="notification"
                kind={testResult.success ? "success" : "error"}
                title={testResult.success ? "Success" : "Error"}
                subtitle={testResult.message}
                hideCloseButton={false}
                on:close={clearTestResult}
        />
{/if}

<div class="space-y-6">
    <!-- Form fields remain the same -->
    <TextInput
            id="server-name"
            invalid={errors.name !== undefined}
            invalidText={errors.name}
            labelText="Server Name *"
            placeholder="Enter a descriptive server name"
            bind:value={formData.name}
            disabled={isLoading}
    />

    <TextInput
            invalid={errors.address !== undefined}
            invalidText={errors.address}
            labelText="Server Address *"
            placeholder="IP address or hostname (e.g., 192.168.1.100)"
            bind:value={formData.address}
            disabled={isLoading}
    />

    <NumberInput
            invalid={errors.port !== undefined}
            invalidText={errors.port}
            label="Port *"
            placeholder="22"
            bind:value={formData.port}
            min={1}
            max={65535}
            step={1}
            disabled={isLoading}
    />

    <TextInput
            invalid={errors.user !== undefined}
            invalidText={errors.user}
            labelText="SSH Username *"
            placeholder="Enter SSH username"
            bind:value={formData.user}
            disabled={isLoading}
    />

    <TextInput
            invalid={errors.password !== undefined}
            invalidText={errors.password}
            labelText="SSH Password *"
            type="password"
            placeholder="Enter SSH password"
            bind:value={formData.password}
            disabled={isLoading}
    />
</div>

<div class="space-y-4 mt-6">
    <Button
            on:click={isTestingConnection ? cancelTestConnection : testConnection}
            icon={ConnectionSignal}
            kind={isTestingConnection ? "danger" : "primary"}
            disabled={isLoading}
    >
        {#if isTestingConnection}
            Cancel Test
        {:else}
            Test Connection
        {/if}
    </Button>
</div>