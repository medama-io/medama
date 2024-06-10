import { Combobox, InputBase, useCombobox } from '@mantine/core';
import { useDidUpdate } from '@mantine/hooks';
import { useSearchParams } from '@remix-run/react';
import { useState } from 'react';

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
		searchParams.get('period')
			? (searchParams.get('period') as keyof typeof presets)
			: 'today',
	);

	// Update search params when preset changes
	useDidUpdate(() => {
		searchParams.set('period', preset);
		setSearchParams(searchParams);

		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [preset]);

	const options = Object.entries(presets).map(([value, label]) => (
		<Combobox.Option key={value} value={value} active={value === preset}>
			{label}
		</Combobox.Option>
	));

	return (
		<Combobox
			store={combobox}
			resetSelectionOnOptionHover
			onOptionSubmit={(value) => {
				setPreset(value as keyof typeof presets);
			}}
		>
			<Combobox.Target>
				<InputBase
					component="button"
					type="button"
					pointer
					rightSection={<Combobox.Chevron />}
					rightSectionPointerEvents="none"
					onClick={() => {
						combobox.toggleDropdown();
					}}
					w={200}
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
