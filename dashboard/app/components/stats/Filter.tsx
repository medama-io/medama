import {
	Center,
	Combobox,
	Group,
	InputBase,
	Popover,
	ScrollArea,
	Text,
	UnstyledButton,
	useCombobox,
} from '@mantine/core';
import { useState } from 'react';

import { IconChevronDown } from '@/components/icons/chevrondown';
import { IconPlus } from '@/components/icons/plus';

import classes from './Filter.module.css';

const filterTypes = {
	string: [
		'equals',
		'does not equal',
		'contains',
		'does not contain',
		'starts with',
		'does not start with',
		'ends with',
		'does not end with',
		'is in',
		'is not in',
	] as const,
	fixed: ['equals', 'does not equal', 'is in', 'is not in'] as const,
};

interface FilterChoicesString {
	label: string;
	type: 'string';
	placeholder?: string;
}

interface FilterChoicesFixed {
	label: string;
	type: 'fixed';
	values: readonly string[];
}

type FilterOptions = Array<FilterChoicesString | FilterChoicesFixed>;

const filterOptions: FilterOptions = [
	{
		label: 'Path',
		type: 'string',
		placeholder: 'e.g. /blog',
	},
	{
		label: 'Referrer',
		type: 'string',
		placeholder: 'e.g. /blog',
	},
	{
		label: 'Browser',
		type: 'fixed',
		values: ['Chrome', 'Firefox', 'Safari', 'Edge', 'Internet Explorer'],
	},
];

interface FilterDropdownProps {
	choices: readonly string[];
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
		<Combobox.Option key={filter} value={filter} active={filter === value}>
			{filter}
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

const FilterGroup = () => {
	const [filter, setFilter] = useState('Path');
	const [type, setType] = useState('equals');
	const [value, setValue] = useState('');

	const filterList = filterOptions.map((item) => item.label);
	const chosenFilter = filterOptions.find((item) => item.label === filter);

	if (!chosenFilter) {
		return <Center>No filters found</Center>;
	}

	return (
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
				<InputBase
					value={value}
					onChange={(event) => {
						setValue(event.currentTarget.value);
					}}
					placeholder={chosenFilter.placeholder}
				/>
			) : (
				<FilterDropdown
					choices={chosenFilter.values}
					value={value}
					setValue={setValue}
				/>
			)}
		</Group>
	);
};

export const Filters = () => {
	return (
		<Popover width={442} trapFocus position="bottom-start">
			<Popover.Target>
				<UnstyledButton className={classes.add}>
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
				<FilterGroup />
				<Group justify="flex-end" className={classes.select}>
					<UnstyledButton className={classes.cancel}>Cancel</UnstyledButton>
					<UnstyledButton className={classes.apply}>Apply</UnstyledButton>
				</Group>
			</Popover.Dropdown>
		</Popover>
	);
};
