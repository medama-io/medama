import {
	Combobox,
	Input,
	InputBase,
	ScrollArea,
	Text,
	useCombobox,
} from '@mantine/core';
import { useState } from 'react';
import { useTimezoneSelect, type ITimezone } from 'react-timezone-select';

import classes from './TimezoneSelect.module.css';

export const TimezoneSelect = () => {
	const [search, setSearch] = useState('');
	const { options: tzOptions, parseTimezone } = useTimezoneSelect({
		labelStyle: 'abbrev',
	});

	const [selectedItem, setSelectedItem] = useState<ITimezone>(
		Intl.DateTimeFormat().resolvedOptions().timeZone,
	);
	const combobox = useCombobox({
		onDropdownClose: () => {
			combobox.resetSelectedOption();
			combobox.focusTarget();
			setSearch('');
		},

		onDropdownOpen: () => {
			combobox.focusSearchInput();
		},
	});

	const options = tzOptions
		.filter((item) => item.label.toLowerCase().includes(search.toLowerCase()))
		.map((item) => (
			<Combobox.Option value={item.value} key={item.value}>
				{item.label}
			</Combobox.Option>
		));

	return (
		<Combobox
			store={combobox}
			width={250}
			position="bottom-start"
			onOptionSubmit={(val) => {
				setSelectedItem(val);
				combobox.closeDropdown();
			}}
		>
			<Combobox.Target>
				<Input.Wrapper
					label="Reported Timezone"
					description="Input description"
				>
					<Input
						classNames={{ input: classes.target }}
						component="button"
						type="button"
						pointer
						rightSection={<Combobox.Chevron />}
						rightSectionPointerEvents="none"
						aria-label="Select timezone"
						onClick={() => {
							combobox.toggleDropdown();
						}}
					>
						{String(selectedItem) || 'Select timezone'}{' '}
					</Input>
				</Input.Wrapper>
			</Combobox.Target>

			<Combobox.Dropdown>
				<Combobox.Search
					value={search}
					onChange={(event) => setSearch(event.currentTarget.value)}
					placeholder="Search timezones"
				/>
				<Combobox.Options>
					<ScrollArea.Autosize type="scroll" mah={200}>
						{options.length > 0 ? (
							options
						) : (
							<Combobox.Empty>Nothing found</Combobox.Empty>
						)}
					</ScrollArea.Autosize>
				</Combobox.Options>
			</Combobox.Dropdown>
		</Combobox>
	);
};
