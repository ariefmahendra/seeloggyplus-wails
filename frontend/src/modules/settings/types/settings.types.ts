import type { dto } from '../../../../wailsjs/go/models';

export interface BaseSetting {
    key: string;
    name: string;
    description: string;
    type: string;
}

export interface BooleanSetting extends BaseSetting {
    value: boolean;
    type: 'boolean';
}

export interface StringSetting extends BaseSetting {
    value: string;
    type: 'string';
}

export interface NumberSetting extends BaseSetting {
    value: number;
    type: 'number';
}

export interface SelectSetting extends BaseSetting {
    value: string;
    type: 'select';
    options?: string[];
}

export type Setting = BooleanSetting | StringSetting | NumberSetting | SelectSetting;

export interface SettingsState {
    booleanSettings: BooleanSetting[];
    stringSettings: StringSetting[];
    numberSettings: NumberSetting[];
    selectSettings: SelectSetting[];
    isLoading: boolean;
    error: string;
    hasChanges: boolean;
}