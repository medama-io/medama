import type {
	Filter,
	FilterKey,
	FilterOperator,
} from '@/components/stats/types';
import { useSearchParams } from '@remix-run/react';
import { useCallback } from 'react';

const getKey = (filter: Filter, type: FilterOperator): FilterKey =>
	`${filter}[${type}]`;

const useFilter = () => {
	const [searchParams, setSearchParams] = useSearchParams();

	const addFilter = useCallback(
		(filter: Filter, type: FilterOperator, value: string) => {
			const key = getKey(filter, type);
			// Check if the filter is already in the search params with the same value.
			// It is possible to have multiple with the same filter and type but different values.
			if (!searchParams.getAll(key).includes(value)) {
				searchParams.append(key, value);
				setSearchParams(searchParams, { preventScrollReset: true });
			}
		},
		[searchParams, setSearchParams],
	);

	const removeFilter = useCallback(
		(filter: Filter, type: FilterOperator, value: string) => {
			const key = getKey(filter, type);

			// Remove the filter with the same filter, type, and value.
			searchParams.delete(key, value);
			setSearchParams(searchParams, { preventScrollReset: true });
		},
		[searchParams, setSearchParams],
	);

	const isFilterActive = useCallback(
		(filter: Filter) => {
			const keys = searchParams.keys();
			for (const key of keys) {
				if (key.startsWith(filter)) {
					return true;
				}
			}

			return false;
		},
		[searchParams],
	);

	const isFilterActiveEq = useCallback(
		(filter: Filter) => {
			const key = getKey(filter, 'eq');
			return searchParams.has(key);
		},
		[searchParams],
	);

	return {
		searchParams,
		addFilter,
		removeFilter,
		isFilterActive,
		isFilterActiveEq,
	};
};

export { useFilter, getKey };
