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
import { useCallback, useEffect, useMemo, useState } from 'react';
import {
	ChevronDownIcon,
	ChevronUpIcon,
	PlusIcon,
} from '@radix-ui/react-icons';

import { useFilter } from '@/hooks/use-filter';

import type { Filter, FilterOperator } from './types';

import classes from './Filter.module.css';

interface FilterChoices {
	label: string;
	placeholder?: string;
}

type FilterOptions = Record<Filter, FilterChoices>;
type FilterType = Record<
	FilterOperator,
	{ label: string; value: FilterOperator }
>;

const FILTER_TYPES: FilterType = {
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

const FILTER_OPTIONS: FilterOptions = {
	path: {
		label: 'Path',
		placeholder: 'e.g. /blog',
	},
	referrer: {
		label: 'Referrer',
		placeholder: 'e.g. example.com',
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

// Add this type guard function
const isFilterType = (obj: FilterOptions | FilterType): obj is FilterType => {
	return 'eq' in obj;
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

	const label = isFilterType(choices)
		? choices[value as FilterOperator]?.label
		: choices[value as Filter]?.label;

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
					rightSection={
						combobox.dropdownOpened ? <ChevronUpIcon /> : <ChevronDownIcon />
					}
					rightSectionPointerEvents="none"
					onClick={() => {
						combobox.toggleDropdown();
					}}
				>
					<span className={classes['dropdown-label']}>
						{label ?? 'Unknown'}
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
	const [searchParams] = useSearchParams();
	const { addFilter, removeFilter } = useFilter();

	const [filter, setFilter] = useState<Filter>('path');
	const [type, setType] = useState<FilterOperator>('eq');
	const [value, setValue] = useState<string>('');

	const chosenFilter = useMemo(() => FILTER_OPTIONS[filter], [filter]);

	// If the filter changes, reset the type
	// We only reset the value if there is a change in filter type
	useEffect(() => {
		setType('eq');
		// Reset the value
		setValue('');
	}, []);

	// Convert the search params to an array of label-type-value tuples
	// This is used to render the current filters.
	const paramsArr = useMemo(() => {
		const arr: Array<[string, string, string]> = [];
		for (const [key, value] of searchParams.entries()) {
			const [filter, type] = key.split('['); // e.g. path[eq]
			// If the filter is not in the filter options, don't render it
			if (filter && FILTER_OPTIONS[filter as Filter]) {
				const { label = 'N/A' } = FILTER_OPTIONS[filter as Filter] ?? {};
				arr.push([label, type?.replace(']', '') ?? 'Unknown', value]);
			}
		}
		return arr;
	}, [searchParams]);

	// Apply the filters to the search params and close the popover
	// This should trigger a reload of data on the page with new
	// data from the server.
	const handleAddFilters = useCallback(() => {
		addFilter(filter, type, value);
		setOpened(false);
	}, [addFilter, filter, type, value]);

	// Remove a filter.
	const handleRemoveFilter = useCallback(
		(label: string, type: FilterOperator, value: string) => {
			return () => {
				const filterMap: Record<string, string> = {
					'UTM Source': 'utm_source',
					'UTM Medium': 'utm_medium',
					'UTM Campaign': 'utm_campaign',
				};
				const filterKey = filterMap[label] ?? label.toLowerCase();
				removeFilter(filterKey as Filter, type, value);
			};
		},
		[removeFilter],
	);

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
						<Group gap={8} justify="center">
							<PlusIcon />
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
							choices={FILTER_OPTIONS}
							value={filter}
							setValue={setFilter as (value: string) => void}
						/>
						<FilterDropdown
							choices={FILTER_TYPES}
							value={type}
							setValue={setType as (value: string) => void}
						/>
						<TextInput
							h={40}
							value={value}
							onChange={(event) => {
								setValue(event.currentTarget.value);
							}}
							onKeyDown={(event) => {
								if (event.key === 'Enter' && value !== '') {
									handleAddFilters();
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
								handleAddFilters();
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
							{FILTER_TYPES[type as FilterOperator]?.label ?? 'Unknown'}&nbsp;
						</Text>
						<Text fz={14}>{value}</Text>
						<CloseButton
							onClick={handleRemoveFilter(label, type as FilterOperator, value)}
						/>
					</Group>
				);
			})}
		</Group>
	);
};
