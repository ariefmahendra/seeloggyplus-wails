import {FindSettingsByKey, GetListSettings, UpdateSettings, UpdateAllSettings} from "../../../../wailsjs/go/main/App";
import type {dto} from "../../../../wailsjs/go/models";

export class SettingsService {
    static async getListSettings(): Promise<Array<dto.SettingsResponseDto>>{
        try {
            const result = await GetListSettings();
            return result || [];
        } catch (error) {
            console.error('Failed to get settings:', error);
            throw error;
        }
    }

    static async getSettingByKey(key: string): Promise<dto.SettingsResponseDto> {
        try {
            return await FindSettingsByKey(key);
        }catch (error) {
            console.error('Failed to get setting by key:', error);
            throw error;
        }
    }

    static async updateSetting(setting: dto.SettingsRequestDto): Promise<dto.SettingsResponseDto> {
        try {
            return await UpdateSettings(setting);
        } catch (error) {
            console.error('Failed to update setting:', error);
            throw error;
        }
    }

    static async updateAllSettings(settings: Array<dto.SettingsRequestDto>): Promise<void> {
        try {
            await UpdateAllSettings(settings);
        } catch (error) {
            console.error('Failed to update all settings:', error);
            throw error;
        }
    }
}