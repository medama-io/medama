import { Combobox, InputBase, useCombobox } from '@mantine/core';
import { useDidUpdate } from '@mantine/hooks';
import { useSearchParams } from '@remix-run/react';
import { useCallback, useMemo, useState } from 'react';

import { IconCalendar } from '@/components/icons/calendar';

import classes from './DateSelector.module.css';

const PRESETS = {
	today: 'Today',
	yesterday: 'Yesterday',
	'12h': 'Previous 12 hours',
	'24h': 'Previous 24 hours',
	'72h': 'Previous 72 hours',
	'7d': 'Previous 7 days',
	'14d': 'Previous 14 days',
	'30d': 'Previous 30 days',
	quarter: 'Previous quarter',
	half: 'Previous half year',
	year: 'Previous year',
	all: 'All time',
} as const;

type PresetKey = keyof typeof PRESETS;

const isPreset = (value: string): value is PresetKey =>
	Object.keys(PRESETS).includes(value);

const GROUP_END_VALUES: PresetKey[] = ['yesterday', '30d', 'year'];

export const DateComboBox = () => {
	const combobox = useCombobox({
		onDropdownClose: () => {
			combobox.resetSelectedOption();
		},
		onDropdownOpen: (eventSource) => {
			if (eventSource === 'keyboard') {
				combobox.selectActiveOption();
			} else {
				combobox.updateSelectedOptionIndex('active');
			}
		},
	});

	const [searchParams, setSearchParams] = useSearchParams();
	const [preset, setPreset] = useState<PresetKey>(
		(searchParams.get('period') as PresetKey) || 'today',
	);

	useDidUpdate(() => {
		setSearchParams((prevParams) => {
			const newParams = new URLSearchParams(prevParams);
			newParams.set('period', preset);
			return newParams;
		});
	}, [preset]);

	const handleOptionSubmit = useCallback((value: string) => {
		if (isPreset(value)) {
			setPreset(value);
		}
	}, []);

	const options = useMemo(
		() =>
			Object.entries(PRESETS).map(([value, label]) => (
				<Combobox.Option
					key={value}
					value={value}
					active={value === preset}
					data-group-end={GROUP_END_VALUES.includes(value as PresetKey)}
					role="option"
					aria-selected={value === preset}
				>
					{label}
				</Combobox.Option>
			)),
		[preset],
	);

	return (
		<Combobox
			classNames={{ dropdown: classes.dropdown, option: classes.option }}
			store={combobox}
			resetSelectionOnOptionHover
			onOptionSubmit={handleOptionSubmit}
		>
			<Combobox.Target>
				<InputBase
					classNames={{ input: classes.target }}
					className={classes.targetWrapper}
					component="button"
					type="button"
					pointer
					rightSection={<Combobox.Chevron />}
					rightSectionPointerEvents="none"
					onClick={() => combobox.toggleDropdown()}
					aria-label="Select date range"
					leftSection={<IconCalendar />}
				>
					{isPreset(preset) ? PRESETS[preset] : 'Custom range'}
				</InputBase>
			</Combobox.Target>
			<Combobox.Dropdown>
				<Combobox.Options>{options}</Combobox.Options>
			</Combobox.Dropdown>
		</Combobox>
	);
};
