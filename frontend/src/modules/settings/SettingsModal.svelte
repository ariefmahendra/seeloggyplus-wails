<script lang="ts">
    import {Modal, NumberInput, TextInput, Toggle} from "carbon-components-svelte";
    import {SettingsService} from "./services/settings.service";
    import type {dto} from "../../../wailsjs/go/models";
    import type {BooleanSetting, StringSetting, NumberSetting, SelectSetting} from "./types/settings.types";

    export let isOpen: boolean = false;

    let booleanSettings: BooleanSetting[] = [];
    let stringSettings: StringSetting[] = [];
    let numberSettings: NumberSetting[] = [];
    let selectSettings: SelectSetting[] = [];

    async function getAllSettings(){
        try {
             const listAllSettings: dto.SettingsResponseDto[] = await SettingsService.getListSettings();
            if (listAllSettings && listAllSettings.length > 0) {
                booleanSettings = [];
                stringSettings = [];
                numberSettings = [];
                selectSettings = [];

                mapSettings(listAllSettings);
            }
        } catch (error) {
            console.error("Error retrieving settings:", error);
        }
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

    async function updateSettings(settings: dto.SettingsRequestDto) {
        try {
            await SettingsService.updateSetting(settings);
            isOpen = false;
        } catch (error) {
            console.error("Error updating settings:", error);
        }
    }

    // Handle the primary button click to save settings
    function handlePrimaryButtonClick() {
        const settingsToUpdate: dto.SettingsRequestDto[] = [];

        booleanSettings.forEach(setting => {
            settingsToUpdate.push({
                key: setting.key,
                value: JSON.stringify(setting.value),
            });
        });

        stringSettings.forEach(setting => {
            settingsToUpdate.push({
                key: setting.key,
                value: setting.value,
            });
        });

        numberSettings.forEach(setting => {
            settingsToUpdate.push({
                key: setting.key,
                value: JSON.stringify(setting.value),
            });
        });

        selectSettings.forEach(setting => {
            settingsToUpdate.push({
                key: setting.key,
                value: setting.value,
            });
        });

        try {

        }catch (error) {
            console.error("Error preparing settings for update:", error);
        }
    }

    $: if (isOpen){
        getAllSettings();
    }
</script>

<Modal
    bind:open={isOpen}
    size="lg"
    modalHeading="Settings"
    primaryButtonText="Save"
    secondaryButtonText="Cancel"
    on:click:button--secondary={() => isOpen = false}
>
    <div class="space-y-4 m-2">
        {#if booleanSettings.length > 0}
            {#each booleanSettings as setting}
                <Toggle
                    toggled={setting.value}
                    labelText={setting.name}
                />
            {/each}
        {/if}

        {#if stringSettings.length > 0}
            {#each stringSettings as setting}
                <TextInput
                    bind:value={setting.value}
                    labelText={setting.name}
                    helperText={setting.description}
                />
            {/each}
        {/if}

        {#if numberSettings.length > 0}
            {#each numberSettings as setting}
                <NumberInput
                    bind:value={setting.value}
                    label={setting.name}
                    helperText={setting.description}
                />
            {/each}
        {/if}
    </div>
</Modal>

