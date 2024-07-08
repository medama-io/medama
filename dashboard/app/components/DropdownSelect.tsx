import { Combobox, InputBase, useCombobox } from '@mantine/core';
import { useDidUpdate } from '@mantine/hooks';
import { useSearchParams } from '@remix-run/react';
import { useCallback, useMemo, useState } from 'react';

import classes from './DropdownSelect.module.css';

interface DropdownSelectProps {
	// Key to update and read from search params
	searchParamKey: string;
	defaultValue: string;
	defaultLabel: string;
	selectAriaLabel: string;

	records: Record<string, string>;
	groupEndValues?: string[];

	leftSection?: React.ReactNode;
}

export const DropdownSelect = ({
	searchParamKey,
	defaultLabel,
	defaultValue,
	selectAriaLabel,
	records,
	groupEndValues = [],
	leftSection,
}: DropdownSelectProps) => {
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
	const [option, setOption] = useState<string>(
		searchParams.get(searchParamKey) || defaultValue,
	);

	useDidUpdate(() => {
		setSearchParams((prevParams) => {
			const newParams = new URLSearchParams(prevParams);
			newParams.set(searchParamKey, option);
			return newParams;
		});
	}, [option]);

	const handleOptionSubmit = useCallback(
		(value: string) => {
			setOption(value);
			combobox.toggleDropdown();
		},
		[combobox.toggleDropdown],
	);

	const options = useMemo(
		() =>
			Object.entries(records).map(([value, label]) => (
				<Combobox.Option
					key={value}
					value={value}
					active={value === option}
					data-group-end={groupEndValues.includes(value)}
					role="option"
					aria-selected={value === option}
				>
					{label}
				</Combobox.Option>
			)),
		[records, groupEndValues, option],
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
					aria-label={selectAriaLabel}
					leftSection={leftSection}
				>
					{records[option] || defaultLabel}
				</InputBase>
			</Combobox.Target>
			<Combobox.Dropdown>
				<Combobox.Options>{options}</Combobox.Options>
			</Combobox.Dropdown>
		</Combobox>
	);
};
