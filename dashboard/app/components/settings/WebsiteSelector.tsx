import { Combobox, InputBase, ScrollArea, useCombobox } from '@mantine/core';
import { useCallback, useMemo } from 'react';

import classes from './WebsiteSelector.module.css';

interface WebsiteListComboboxProps {
	websites: string[];
	website: string;
	setWebsite: (website: string) => void;
}

export const WebsiteSelector = ({
	websites,
	website,
	setWebsite,
}: WebsiteListComboboxProps) => {
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

	const handleOptionSubmit = useCallback(
		(value: string) => {
			setWebsite(value);
		},
		[setWebsite],
	);

	const options = useMemo(
		() =>
			websites.map((value) => (
				<Combobox.Option
					key={value}
					value={value}
					active={value === website}
					role="option"
					aria-selected={value === website}
				>
					{value}
				</Combobox.Option>
			)),
		[websites, website],
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
					aria-label="Select website hostname"
					disabled={websites.length === 0}
				>
					{website ?? 'No websites'}
				</InputBase>
			</Combobox.Target>
			<Combobox.Dropdown>
				<ScrollArea.Autosize mah={200}>
					<Combobox.Options>{options}</Combobox.Options>
				</ScrollArea.Autosize>
			</Combobox.Dropdown>
		</Combobox>
	);
};
