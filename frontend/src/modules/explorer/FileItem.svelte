<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { dto } from '../../../wailsjs/go/models';
    import {
        Folder,
        Document,
        DocumentPdf,
        Image,
        Video,
        Music,
        Archive,
        Code,
        FileStorage,
        DocumentBlank,
        ChevronRight
    } from 'carbon-icons-svelte';

    export let file: dto.FileInfo;

    const dispatch = createEventDispatcher();

    function handleClick() {
        dispatch('click', file);
    }

    function getFileIcon(file: dto.FileInfo) {
        if (file.isDrive) return FileStorage;
        if (file.isDir) return Folder;

        const ext = file.name.split('.').pop()?.toLowerCase();
        switch (ext) {
            case 'pdf': return DocumentPdf;
            case 'txt': case 'md': case 'readme': return Document;
            case 'jpg': case 'jpeg': case 'png': case 'gif': case 'bmp': case 'svg': case 'webp':
                return Image;
            case 'mp4': case 'avi': case 'mov': case 'wmv': case 'flv': case 'webm':
                return Video;
            case 'mp3': case 'wav': case 'flac': case 'aac': case 'ogg':
                return Music;
            case 'zip': case 'rar': case '7z': case 'tar': case 'gz':
                return Archive;
            case 'js': case 'ts': case 'go': case 'py': case 'java': case 'cpp': case 'c': case 'html': case 'css': case 'php': case 'rb': case 'rs':
                return Code;
            default: return DocumentBlank;
        }
    }

    function getFileTypeColor(file: dto.FileInfo) {
        if (file.isDrive) return 'text-blue-600';
        if (file.isDir) return 'text-blue-500';

        const ext = file.name.split('.').pop()?.toLowerCase();
        switch (ext) {
            case 'pdf': return 'text-red-500';
            case 'jpg': case 'jpeg': case 'png': case 'gif': case 'bmp': case 'svg': case 'webp':
                return 'text-green-500';
            case 'mp4': case 'avi': case 'mov': case 'wmv': case 'flv': case 'webm':
                return 'text-purple-500';
            case 'mp3': case 'wav': case 'flac': case 'aac': case 'ogg':
                return 'text-orange-500';
            case 'zip': case 'rar': case '7z': case 'tar': case 'gz':
                return 'text-yellow-600';
            case 'js': case 'ts': case 'go': case 'py': case 'java': case 'cpp': case 'c': case 'html': case 'css': case 'php': case 'rb': case 'rs':
                return 'text-indigo-500';
            default: return 'text-gray-500';
        }
    }

    function formatFileSize(bytes: number): string {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
    }

    function formatDate(dateString: string): string {
        if (!dateString) return '—';
        try {
            return new Date(dateString).toLocaleDateString('en-US', {
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            });
        } catch {
            return dateString;
        }
    }
</script>

<button
        class="flex items-center gap-3 w-full text-left p-3 hover:bg-gray-50 focus:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-inset transition-colors duration-150"
        class:font-medium={file.isDir || file.isDrive}
        on:click={handleClick}
>
    <div class="flex-shrink-0 {getFileTypeColor(file)}">
        <svelte:component this={getFileIcon(file)} size={20} />
    </div>

    <div class="flex-1 min-w-0">
        <div class="text-sm text-gray-900 truncate" title={file.name}>
            {file.name}
        </div>
        <div class="text-xs text-gray-500 flex items-center gap-3 mt-1">
            {#if !file.isDir}
                <span>{formatFileSize(file.size)}</span>
            {/if}
            <span>{formatDate(file.modTime)}</span>
        </div>
    </div>

    {#if file.isDir}
        <div class="flex-shrink-0 text-gray-400">
            <ChevronRight size={16} />
        </div>
    {/if}
</button>