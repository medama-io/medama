import { Combobox, InputBase, useCombobox } from '@mantine/core';
import { useDidUpdate } from '@mantine/hooks';
import { useSearchParams } from '@remix-run/react';
import { useState } from 'react';
import classes from './DateSelector.module.css';

const presets = {
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
};

const isPreset = (value: string): value is keyof typeof presets =>
	Object.keys(presets).includes(value);

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
	const [preset, setPreset] = useState<keyof typeof presets>(
		(searchParams.get('period') as keyof typeof presets) || 'today',
	);

	useDidUpdate(() => {
		setSearchParams((prevParams) => {
			const newParams = new URLSearchParams(prevParams);
			newParams.set('period', preset);
			return newParams;
		});
	}, [preset]);

	const options = Object.entries(presets).map(([value, label]) => {
		const isGroupEnd = ['yesterday', '30d', 'year'].includes(value);

		return (
			<Combobox.Option
				key={value}
				value={value}
				active={value === preset}
				data-group-end={isGroupEnd}
				role="option"
				aria-selected={value === preset}
			>
				{label}
			</Combobox.Option>
		);
	});

	return (
		<Combobox
			classNames={{ dropdown: classes.dropdown, option: classes.option }}
			store={combobox}
			resetSelectionOnOptionHover
			onOptionSubmit={(value) => {
				setPreset(value as keyof typeof presets);
			}}
		>
			<Combobox.Target>
				<InputBase
					classNames={{ input: classes.target }}
					component="button"
					type="button"
					pointer
					rightSection={<Combobox.Chevron />}
					rightSectionPointerEvents="none"
					onClick={() => {
						combobox.toggleDropdown();
					}}
					aria-label="Select date range"
				>
					{isPreset(preset) ? presets[preset] : 'Custom range'}
				</InputBase>
			</Combobox.Target>
			<Combobox.Dropdown>
				<Combobox.Options>{options}</Combobox.Options>
			</Combobox.Dropdown>
		</Combobox>
	);
};
