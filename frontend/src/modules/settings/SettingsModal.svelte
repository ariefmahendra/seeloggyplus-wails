<script lang="ts">
    import { Modal, NumberInput, TextInput, Toggle, InlineNotification, Loading } from "carbon-components-svelte";
    import { SettingsService } from "./services/settings.service";
    import type { dto } from "../../../wailsjs/go/models";
    import type { BooleanSetting, StringSetting, NumberSetting, SelectSetting } from "./types/settings.types";

    export let isOpen: boolean = false;

    let booleanSettings: BooleanSetting[] = [];
    let stringSettings: StringSetting[] = [];
    let numberSettings: NumberSetting[] = [];
    let selectSettings: SelectSetting[] = [];
    let originalSettings: Map<string, any> = new Map();

    let isLoading = false;
    let isSaving = false;
    let error = '';
    let hasChanges = false;

    async function getAllSettings() {
        try {
            isLoading = true;
            error = '';

            const listAllSettings: dto.SettingsResponseDto[] = await SettingsService.getListSettings();
            if (listAllSettings && listAllSettings.length > 0) {
                booleanSettings = [];
                stringSettings = [];
                numberSettings = [];
                selectSettings = [];
                originalSettings.clear();

                mapSettings(listAllSettings);
                trackOriginalValues();
            }
        } catch (err) {
            error = err instanceof Error ? err.message : 'Failed to load settings';
            console.error("Error retrieving settings:", err);
        } finally {
            isLoading = false;
        }
    }

    // Track original values for change detection
    function trackOriginalValues() {
        [...booleanSettings, ...stringSettings, ...numberSettings, ...selectSettings].forEach(setting => {
            originalSettings.set(setting.key, setting.value);
        });
        hasChanges = false;
    }

    // Check if there are any changes
    function checkForChanges() {
        hasChanges = [...booleanSettings, ...stringSettings, ...numberSettings, ...selectSettings].some(setting => {
            return originalSettings.get(setting.key) !== setting.value;
        });
    }

    // util for mapping settings to their respective types
    function mapSettings(listSettings: dto.SettingsResponseDto[]) {
        listSettings.forEach(setting => {
            switch (setting.type) {
                case 'boolean':
                    const booleanSetting: BooleanSetting = {
                        key: setting.key,
                        name: setting.name,
                        value: JSON.parse(setting.value),
                        description: setting.description || '',
                        type: setting.type
                    }
                    booleanSettings.push(booleanSetting);
                    break;
                case 'string':
                    const stringSetting: StringSetting = {
                        key: setting.key,
                        name: setting.name,
                        value: setting.value,
                        description: setting.description || '',
                        type: setting.type
                    }
                    stringSettings.push(stringSetting);
                    break;
                case 'number':
                    const numberSetting: NumberSetting = {
                        key: setting.key,
                        name: setting.name,
                        value: JSON.parse(setting.value),
                        description: setting.description || '',
                        type: setting.type
                    }
                    numberSettings.push(numberSetting);
                    break;
                case 'select':
                    const selectSetting: SelectSetting = {
                        key: setting.key,
                        name: setting.name,
                        value: setting.value,
                        description: setting.description || '',
                        type: setting.type
                    }
                    selectSettings.push(selectSetting);
                    break;
                default:
                    console.warn(`Unknown setting type: ${setting.type}`);
            }
        });
    }

    // Handle bulk update of all settings
    async function handlePrimaryButtonClick() {
        if (!hasChanges) {
            isOpen = false;
            return;
        }

        try {
            isSaving = true;
            error = '';

            const settingsToUpdate: dto.SettingsRequestDto[] = [];

            // Collect all settings that have changed
            booleanSettings.forEach(setting => {
                if (originalSettings.get(setting.key) !== setting.value) {
                    settingsToUpdate.push({
                        key: setting.key,
                        value: JSON.stringify(setting.value),
                    });
                }
            });

            stringSettings.forEach(setting => {
                if (originalSettings.get(setting.key) !== setting.value) {
                    settingsToUpdate.push({
                        key: setting.key,
                        value: setting.value,
                    });
                }
            });

            numberSettings.forEach(setting => {
                if (originalSettings.get(setting.key) !== setting.value) {
                    settingsToUpdate.push({
                        key: setting.key,
                        value: JSON.stringify(setting.value),
                    });
                }
            });

            selectSettings.forEach(setting => {
                if (originalSettings.get(setting.key) !== setting.value) {
                    settingsToUpdate.push({
                        key: setting.key,
                        value: setting.value,
                    });
                }
            });

            if (settingsToUpdate.length > 0) {
                // Use bulk update method
                await SettingsService.updateAllSettings(settingsToUpdate);
                console.log(`Successfully updated ${settingsToUpdate.length} settings`);
            }

            isOpen = false;
        } catch (err) {
            error = err instanceof Error ? err.message : 'Failed to save settings';
            console.error("Error updating settings:", err);
        } finally {
            isSaving = false;
        }
    }

    function handleCancel() {
        if (hasChanges) {
            const confirmCancel = confirm('You have unsaved changes. Are you sure you want to cancel?');
            if (!confirmCancel) return;
        }
        isOpen = false;
    }

    // Reactive statements
    $: if (isOpen) {
        getAllSettings();
    }

    // Watch for changes in settings
    $: {
        if (originalSettings.size > 0) {
            checkForChanges();
        }
    }
</script>

<Modal
        bind:open={isOpen}
        size="lg"
        modalHeading="Settings"
        primaryButtonText={isSaving ? "Saving..." : hasChanges ? "Save Changes" : "Close"}
        secondaryButtonText="Cancel"
        primaryButtonDisabled={isSaving || isLoading}
        on:click:button--primary={handlePrimaryButtonClick}
        on:click:button--secondary={handleCancel}
        on:close={handleCancel}
>
    <div class="space-y-4 m-2">
        <!-- Loading State -->
        {#if isLoading}
            <div class="flex justify-center items-center h-32">
                <Loading description="Loading settings..." />
            </div>
        {/if}

        <!-- Error Display -->
        {#if error}
            <InlineNotification
                    kind="error"
                    title="Error"
                    subtitle={error}
                    on:close={() => error = ''}
            />
        {/if}

        <!-- Changes Indicator -->
        {#if hasChanges}
            <InlineNotification
                    kind="info"
                    title="Unsaved Changes"
                    subtitle="You have modified settings that haven't been saved yet."
                    hideCloseButton
            />
        {/if}

        <!-- Boolean Settings -->
        {#if booleanSettings.length > 0}
            <div class="space-y-3">
                <h3 class="text-lg font-medium text-gray-900">Boolean Settings</h3>
                {#each booleanSettings as setting}
                    <div class="p-3 border border-gray-200 rounded-lg">
                        <Toggle
                                bind:toggled={setting.value}
                                labelText={setting.name}
                                disabled={isSaving}
                                on:toggle={checkForChanges}
                        />
                        {#if setting.description}
                            <p class="text-sm text-gray-500 mt-1">{setting.description}</p>
                        {/if}
                    </div>
                {/each}
            </div>
        {/if}

        <!-- String Settings -->
        {#if stringSettings.length > 0}
            <div class="space-y-3">
                <h3 class="text-lg font-medium text-gray-900">Text Settings</h3>
                {#each stringSettings as setting}
                    <div class="p-3 border border-gray-200 rounded-lg">
                        <TextInput
                                bind:value={setting.value}
                                labelText={setting.name}
                                helperText={setting.description}
                                disabled={isSaving}
                                on:input={checkForChanges}
                        />
                    </div>
                {/each}
            </div>
        {/if}

        <!-- Number Settings -->
        {#if numberSettings.length > 0}
            <div class="space-y-3">
                <h3 class="text-lg font-medium text-gray-900">Number Settings</h3>
                {#each numberSettings as setting}
                    <div class="p-3 border border-gray-200 rounded-lg">
                        <NumberInput
                                bind:value={setting.value}
                                label={setting.name}
                                helperText={setting.description}
                                disabled={isSaving}
                                on:input={checkForChanges}
                        />
                    </div>
                {/each}
            </div>
        {/if}

        <!-- Select Settings -->
        {#if selectSettings.length > 0}
            <div class="space-y-3">
                <h3 class="text-lg font-medium text-gray-900">Selection Settings</h3>
                {#each selectSettings as setting}
                    <div class="p-3 border border-gray-200 rounded-lg">
                        <TextInput
                                bind:value={setting.value}
                                labelText={setting.name}
                                helperText={setting.description}
                                disabled={isSaving}
                                on:input={checkForChanges}
                        />
                    </div>
                {/each}
            </div>
        {/if}

        <!-- Empty State -->
        {#if !isLoading && booleanSettings.length === 0 && stringSettings.length === 0 && numberSettings.length === 0 && selectSettings.length === 0}
            <div class="text-center py-8 text-gray-500">
                <p>No settings available</p>
            </div>
        {/if}

        <!-- Footer Info -->
        {#if !isLoading}
            <div class="text-sm text-gray-500 p-3 bg-gray-50 rounded-lg">
                <p><strong>Total Settings:</strong> {booleanSettings.length + stringSettings.length + numberSettings.length + selectSettings.length}</p>
                {#if hasChanges}
                    <p class="text-blue-600 mt-1"><strong>Note:</strong> Click "Save Changes" to apply your modifications.</p>
                {/if}
            </div>
        {/if}
    </div>
</Modal>