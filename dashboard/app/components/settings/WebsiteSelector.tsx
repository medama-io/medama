import { Combobox, InputBase, useCombobox } from '@mantine/core';
import { useDidUpdate } from '@mantine/hooks';
import { useSearchParams } from '@remix-run/react';
import { useCallback, useMemo, useState } from 'react';

import classes from './WebsiteSelector.module.css';

interface WebsiteListComboboxProps {
	websites: string[];
}

export const WebsiteSelector = ({ websites }: WebsiteListComboboxProps) => {
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
	const [website, setWebsite] = useState<string>(
		searchParams.get('website') ?? websites[0] ?? '',
	);

	useDidUpdate(() => {
		setSearchParams((prevParams) => {
			const newParams = new URLSearchParams(prevParams);
			newParams.set('website', website);
			return newParams;
		});
	}, [website]);

	const handleOptionSubmit = useCallback((value: string) => {
		setWebsite(value);
	}, []);

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
					{website}
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
				<Combobox.Options>{options}</Combobox.Options>
			</Combobox.Dropdown>
		</Combobox>
	);
};
