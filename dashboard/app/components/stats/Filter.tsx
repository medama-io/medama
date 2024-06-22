import {
	CloseButton,
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

type FilterType = Record<string, { label: string; value: string }>;

interface FilterChoices {
	label: string;
	placeholder?: string;
}

type FilterOptions = Record<string, FilterChoices>;

const filterTypes: FilterType = {
	eq: { label: 'equals', value: 'eq' },
	neq: { label: 'does not equal', value: 'neq' },
	contains: { label: 'contains', value: 'contains' },
	not_contains: { label: 'does not contain', value: 'not_contains' },
	starts_with: { label: 'starts with', value: 'starts_with' },
	not_starts_with: {
		label: 'does not start with',
		value: 'not_starts_with',
	},
	ends_with: { label: 'ends with', value: 'ends_with' },
	not_ends_with: { label: 'does not end with', value: 'not_ends_with' },
	// in: { label: 'is in', value: 'in' },
	// not_in: { label: 'is not in', value: 'not_in' },
};

const filterOptions: FilterOptions = {
	path: {
		label: 'Path',
		placeholder: 'e.g. /blog',
	},
	referrer: {
		label: 'Referrer',
		placeholder: 'e.g. /blog',
	},
	utm_source: {
		label: 'UTM Source',
		placeholder: 'e.g. google',
	},
	utm_medium: {
		label: 'UTM Medium',
		placeholder: 'e.g. cpc',
	},
	utm_campaign: {
		label: 'UTM Campaign',
		placeholder: 'e.g. summer_sale',
	},
	browser: {
		label: 'Browser',
		placeholder: 'e.g. Chrome',
	},
	os: {
		label: 'OS',
		placeholder: 'e.g. Windows',
	},
	device: {
		label: 'Device',
		placeholder: 'e.g. Desktop',
	},
	country: {
		label: 'Country',
		placeholder: 'e.g. United States',
	},
	language: {
		label: 'Language',
		placeholder: 'e.g. English',
	},
};

interface FilterDropdownProps {
	choices: FilterOptions | FilterType;
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

	const options = Object.entries(choices).map(([key, filter]) => (
		<Combobox.Option key={key} value={key} active={key === value}>
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
					className={classes.dropdown}
					rightSection={<IconChevronDown />}
					rightSectionPointerEvents="none"
					onClick={() => {
						combobox.toggleDropdown();
					}}
				>
					<span className={classes['dropdown-label']}>
						{choices?.[value]?.label ?? value}
					</span>
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

	const [filter, setFilter] = useState('path');
	const [type, setType] = useState('eq');
	const [value, setValue] = useState('');

	const chosenFilter = filterOptions[filter];

	// If the filter changes, reset the type
	// We only reset the value if there is a change in filter type
	useEffect(() => {
		setType('eq');
		// Reset the value
		setValue('');
	}, []);

	const [searchParams, setSearchParams] = useSearchParams();
	// Convert the search params to an array of label-type-value tuples
	// This is used to render the current filters.
	const paramsArr: Array<[string, string, string]> = [];
	for (const [key, value] of searchParams.entries()) {
		const [filter, type] = key.split('['); // e.g. path[eq]
		// If the filter is not in the filter options, don't render it
		if (!filter || !filterOptions[filter]) {
			continue;
		}

		const { label = 'N/A' } = filterOptions[filter] ?? {};
		paramsArr.push([label, type?.replace(']', '') ?? 'Unknown', value]);
	}

	// Apply the filters to the search params and close the popover
	// This should trigger a reload of data on the page with new
	// data from the server.
	const applyFilters = () => {
		const params = new URLSearchParams(searchParams);
		params.append(`${filter}[${type}]`, value);
		setSearchParams(params, {
			preventScrollReset: true,
		});
		setOpened(false);
	};

	return (
		<Group mt={-40}>
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
								setFilter('path');
								setType('eq');
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
							choices={filterOptions}
							value={filter}
							setValue={setFilter}
						/>
						<FilterDropdown
							choices={filterTypes}
							value={type}
							setValue={setType}
						/>
						<TextInput
							h={40}
							value={value}
							onChange={(event) => {
								setValue(event.currentTarget.value);
							}}
							onKeyDown={(event) => {
								if (event.key === 'Enter' && value !== '') {
									applyFilters();
								}
							}}
							placeholder={chosenFilter?.placeholder}
						/>
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
								applyFilters();
							}}
							disabled={value === ''}
						>
							Apply
						</UnstyledButton>
					</Group>
				</Popover.Dropdown>
			</Popover>
			{paramsArr.map(([label, type, value]) => {
				return (
					<Group
						key={label + type + value}
						className={classes['filter-item']}
						gap={0}
					>
						<Text fz={14}>{label}&nbsp;</Text>
						<Text fz={14} fw={700}>
							{filterTypes[type]?.label ?? 'Unknown'}&nbsp;
						</Text>
						<Text fz={14}>{value}</Text>
						<CloseButton
							onClick={() => {
								const params = new URLSearchParams(searchParams);
								params.delete(`${label.toLowerCase()}[${type}]`, value);
								setSearchParams(params, {
									preventScrollReset: true,
								});
							}}
						/>
					</Group>
				);
			})}
		</Group>
	);
};
