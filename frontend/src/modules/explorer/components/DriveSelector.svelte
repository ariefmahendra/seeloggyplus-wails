<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { dto } from '../../../../wailsjs/go/models';
    import { Button } from 'carbon-components-svelte';
    import { FileStorage, IbmOpenshiftContainerPlatformOnVpcForRegulatedIndustries  } from 'carbon-icons-svelte';

    export let drives: dto.DriveInfo[];

    const dispatch = createEventDispatcher();

    function selectDrive(drive: dto.DriveInfo) {
        dispatch('driveSelected', drive);
    }

    function getDriveIcon(drive: dto.DriveInfo) {
        return drive.type === 'fixed' ? IbmOpenshiftContainerPlatformOnVpcForRegulatedIndustries : FileStorage;
    }

    function getDriveTypeLabel(drive: dto.DriveInfo) {
        switch (drive.type) {
            case 'fixed': return 'Local Disk';
            case 'removable': return 'Removable';
            case 'network': return 'Network';
            case 'cdrom': return 'CD-ROM';
            default: return 'Drive';
        }
    }
</script>

{#if drives.length > 0}
    <div class="p-4 bg-gray-50 border-b border-gray-200">
        <h3 class="text-sm font-medium text-gray-700 mb-3">Available Drives</h3>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
            {#each drives as drive}
                <Button
                        kind="tertiary"
                        size="small"
                        on:click={() => selectDrive(drive)}
                        class="justify-start h-auto p-3"
                >
                    <div class="flex items-center gap-3 w-full">
                        <svelte:component this={getDriveIcon(drive)} size={20} class="text-blue-600" />
                        <div class="flex-1 text-left">
                            <div class="font-medium text-sm">
                                {drive.label || drive.name}
                            </div>
                            <div class="text-xs text-gray-500 flex items-center gap-2">
                                <span>{drive.path}</span>
                                <span>•</span>
                                <span>{getDriveTypeLabel(drive)}</span>
                            </div>
                        </div>
                    </div>
                </Button>
            {/each}
        </div>
    </div>
{/if}