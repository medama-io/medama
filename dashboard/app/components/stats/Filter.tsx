import * as DropdownMenu from '@radix-ui/react-dropdown-menu';
import { ChevronDownIcon, Cross1Icon, PlusIcon } from '@radix-ui/react-icons';
import * as Popover from '@radix-ui/react-popover';
import { useSearchParams } from '@remix-run/react';
import { useCallback, useEffect, useMemo, useState } from 'react';

import { ScrollArea } from '@/components/ScrollArea';
import { useFilter } from '@/hooks/use-filter';
import classes from './Filter.module.css';
import type { Filter, FilterOperator } from './types';

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
	prop_name: {
		label: 'Property Name',
		placeholder: 'e.g. logged_in',
	},
	prop_value: {
		label: 'Property Value',
		placeholder: 'e.g. true',
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
	const label = isFilterType(choices)
		? choices[value as FilterOperator]?.label
		: choices[value as Filter]?.label;

	const options = useMemo(() => {
		return Object.entries(choices).map(([key, filter]) => (
			<DropdownMenu.Item key={key} onSelect={() => setValue(key)} asChild>
				<button type="button" className={classes.item}>
					{filter.label}
				</button>
			</DropdownMenu.Item>
		));
	}, [choices, setValue]);

	return (
		<DropdownMenu.Root>
			<DropdownMenu.Trigger asChild>
				<button type="button" className={classes.trigger}>
					<span className={classes.label}>{label ?? 'Unknown'}</span>
					<ChevronDownIcon />
				</button>
			</DropdownMenu.Trigger>

			<DropdownMenu.Content className={classes.dropdown} sideOffset={8}>
				<ScrollArea vertical>{options}</ScrollArea>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
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
					'Property Name': 'prop_name',
					'Property Value': 'prop_value',
				};
				const filterKey = filterMap[label] ?? label.toLowerCase();
				removeFilter(filterKey as Filter, type, value);
			};
		},
		[removeFilter],
	);

	return (
		<div className={classes.root}>
			<Popover.Root
				// width={454}
				open={opened}
				onOpenChange={setOpened}
			>
				<Popover.Trigger asChild>
					<button
						type="button"
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
						data-m:click="filter=open"
					>
						<PlusIcon />
						<span>Add filter</span>
					</button>
				</Popover.Trigger>
				<Popover.Portal>
					<Popover.Content
						className={classes.popover}
						align="start"
						sideOffset={8}
					>
						<h5>New filter</h5>
						<div className={classes['dropdown-list']}>
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
							<input
								className={classes.input}
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
						</div>
						<div className={classes.select}>
							<button
								type="reset"
								className={classes.cancel}
								onClick={() => {
									setOpened(false);
								}}
								data-m:click="filter=cancel"
							>
								Cancel
							</button>
							<button
								className={classes.apply}
								type="submit"
								onClick={() => {
									handleAddFilters();
								}}
								disabled={value === ''}
								data-m:click="filter=apply"
							>
								Apply
							</button>
						</div>
					</Popover.Content>
				</Popover.Portal>
			</Popover.Root>
			{paramsArr.map(([label, type, value]) => {
				return (
					<div key={label + type + value} className={classes['filter-item']}>
						<span>{label}&nbsp;</span>
						<span style={{ fontWeight: 700 }}>
							{FILTER_TYPES[type as FilterOperator]?.label ?? 'Unknown'}&nbsp;
						</span>
						<span>{value}</span>
						<div>
							<Cross1Icon
								onClick={handleRemoveFilter(
									label,
									type as FilterOperator,
									value,
								)}
							/>
						</div>
					</div>
				);
			})}
		</div>
	);
};
