import {
	Combobox,
	Group,
	InputBase,
	Popover,
	ScrollArea,
	Text,
	TextInput,
	UnstyledButton,
	useCombobox,
} from '@mantine/core';
import { useSearchParams } from '@remix-run/react';
import { useEffect, useState } from 'react';

import { IconChevronDown } from '@/components/icons/chevrondown';
import { IconPlus } from '@/components/icons/plus';

import classes from './Filter.module.css';

const filterTypes = {
	string: [
		{ label: 'equals', value: 'eq' },
		{ label: 'does not equal', value: 'neq' },
		{ label: 'contains', value: 'contains' },
		{ label: 'does not contain', value: 'not_contains' },
		{ label: 'starts with', value: 'starts_with' },
		{ label: 'does not start with', value: 'not_starts_with' },
		{ label: 'ends with', value: 'ends_with' },
		{ label: 'does not end with', value: 'not_ends_with' },
		{ label: 'is in', value: 'in' },
		{ label: 'is not in', value: 'not_in' },
	] as const,
	fixed: [
		{ label: 'equals', value: 'eq' },
		{ label: 'does not equal', value: 'neq' },
		{ label: 'is in', value: 'in' },
		{ label: 'is not in', value: 'not_in' },
	] as const,
};

interface FilterChoicesString {
	label: string;
	value: string;
	type: 'string';
	placeholder?: string;
}

interface FilterChoicesFixed {
	label: string;
	value: string;
	type: 'fixed';
	choices: string[];
}

type FilterOptions = Array<FilterChoicesString | FilterChoicesFixed>;

const filterOptions: FilterOptions = [
	{
		label: 'Path',
		value: 'path',
		type: 'string',
		placeholder: 'e.g. /blog',
	},
	{
		label: 'Referrer',
		value: 'referrer',
		type: 'string',
		placeholder: 'e.g. /blog',
	},
	{
		label: 'Browser',
		value: 'browser',
		type: 'fixed',
		choices: ['Chrome', 'Firefox', 'Safari', 'Edge', 'Internet Explorer'],
	},
];

interface FilterDropdownProps {
	choices: ReadonlyArray<{ label: string; value: string }>;
	value: string;
	setValue: (value: string) => void;
}

const FilterDropdown = ({ choices, value, setValue }: FilterDropdownProps) => {
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

	const options = choices.map((filter) => (
		<Combobox.Option
			key={filter.value}
			value={filter.label}
			active={filter.value === value}
		>
			{filter.label}
		</Combobox.Option>
	));

	return (
		<Combobox
			onOptionSubmit={(option) => {
				setValue(option);
				combobox.updateSelectedOptionIndex('active');
				combobox.closeDropdown();
			}}
			store={combobox}
			withinPortal={false}
		>
			<Combobox.Target>
				<InputBase
					component="button"
					type="button"
					pointer
					className={classes['dropdown-button']}
					rightSection={<IconChevronDown />}
					rightSectionPointerEvents="none"
					onClick={() => {
						combobox.toggleDropdown();
					}}
				>
					{value}
				</InputBase>
			</Combobox.Target>

			<Combobox.Dropdown>
				<Combobox.Options>
					<ScrollArea.Autosize mah={200} type="scroll">
						{options}
					</ScrollArea.Autosize>
				</Combobox.Options>
			</Combobox.Dropdown>
		</Combobox>
	);
};

export const Filters = () => {
	const [opened, setOpened] = useState(false);

	const [filter, setFilter] = useState('Path');
	const [type, setType] = useState('equals');
	const [value, setValue] = useState('');

	const filterList = filterOptions.map((item) => ({
		label: item.label,
		value: item.value,
	}));
	const chosenFilter = filterOptions.find((item) => item.label === filter);

	// If the filter changes, reset the type
	// We only reset the value if there is a change in filter type
	useEffect(() => {
		setType('equals');
		// If the filter is a string, reset the value
		// Otherwise, reset the value to the first choice
		chosenFilter?.type === 'string'
			? setValue(value)
			: setValue(chosenFilter?.choices[0] ?? '');
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [filter]);

	const [searchParams, setSearchParams] = useSearchParams();

	return (
		<Popover
			width={442}
			trapFocus
			position="bottom-start"
			opened={opened}
			onChange={setOpened}
		>
			<Popover.Target>
				<UnstyledButton
					className={classes.add}
					onClick={() => {
						// Reset all filters on open
						if (!opened) {
							setFilter('Path');
							setType('equals');
							setValue('');
						}
						setOpened(!opened);
					}}
				>
					<Group gap="xs" justify="center">
						<IconPlus />
						<Text fz={14} fw={600}>
							Add filter
						</Text>
					</Group>
				</UnstyledButton>
			</Popover.Target>
			<Popover.Dropdown className={classes.popover}>
				<Text fz={16} fw={600} pb="sm">
					New filter
				</Text>
				<Group grow>
					<FilterDropdown
						choices={filterList}
						value={filter}
						setValue={setFilter}
					/>
					<FilterDropdown
						choices={
							chosenFilter?.type === 'string'
								? filterTypes.string
								: filterTypes.fixed
						}
						value={type}
						setValue={setType}
					/>
					{chosenFilter?.type === 'string' ? (
						<TextInput
							h={40}
							value={value}
							onChange={(event) => {
								setValue(event.currentTarget.value);
							}}
							placeholder={chosenFilter.placeholder}
						/>
					) : (
						<FilterDropdown
							choices={
								chosenFilter?.choices.map((item) => ({
									label: item,
									value: item.toLowerCase(),
								})) ?? []
							}
							value={value === '' ? chosenFilter?.choices[0] ?? '' : value}
							setValue={setValue}
						/>
					)}
				</Group>
				<Group justify="flex-end" className={classes.select}>
					<UnstyledButton
						className={classes.cancel}
						onClick={() => {
							setOpened(false);
						}}
					>
						Cancel
					</UnstyledButton>
					<UnstyledButton
						className={classes.apply}
						type="submit"
						onClick={() => {
							const params = new URLSearchParams(searchParams);
							// The state values are just labels which need to be converted
							const filterKey =
								filterOptions.find((item) => item.label === filter)?.value ??
								'';
							const filterType =
								filterTypes.string.find((item) => item.label === type)?.value ??
								'';
							params.append(`${filterKey}[${filterType}]`, value);
							setSearchParams(params, {
								preventScrollReset: true,
							});
							setOpened(false);
						}}
						disabled={value === ''}
					>
						Apply
					</UnstyledButton>
				</Group>
			</Popover.Dropdown>
		</Popover>
	);
};
