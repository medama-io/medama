import { useSearchParams } from '@remix-run/react';
import { useCallback } from 'react';
import type {
	Filter,
	FilterKey,
	FilterOperator,
} from '@/components/stats/types';

const getKey = (filter: Filter, type: FilterOperator): FilterKey =>
	`${filter}[${type}]`;

const useFilter = () => {
	const [searchParams, setSearchParams] = useSearchParams();

	const addFilter = useCallback(
		(filter: Filter, type: FilterOperator, value: string) => {
			const key = getKey(filter, type);
			// Check if the filter is already in the search params with the same value.
			// It is possible to have multiple with the same filter and type but different values.
			setSearchParams(
				(params) => {
					if (!params.getAll(key).includes(value)) {
						params.append(key, value);
					}
					return params;
				},
				{ preventScrollReset: true },
			);
		},
		[setSearchParams],
	);

	const removeFilter = useCallback(
		(filter: Filter, type: FilterOperator, value: string) => {
			const key = getKey(filter, type);

			// Remove the filter with the same filter, type, and value.
			setSearchParams(
				(params) => {
					params.delete(key, value);
					return params;
				},
				{ preventScrollReset: true },
			);
		},
		[setSearchParams],
	);

	const clearAllFilters = useCallback(() => {
		setSearchParams(
			(params) => {
				// Get all filter keys
				const keysToDelete: string[] = [];
				for (const key of params.keys()) {
					// Check if the key is a filter key (contains [ and ])
					if (key.includes('[') && key.includes(']')) {
						keysToDelete.push(key);
					}
				}
				// Delete all filter keys
				for (const key of keysToDelete) {
					params.delete(key);
				}
				return params;
			},
			{ preventScrollReset: true },
		);
	}, [setSearchParams]);

	const getFilterEq = useCallback(
		(filter: Filter): string | null => {
			const key = getKey(filter, 'eq');
			return searchParams.get(key);
		},
		[searchParams],
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
		clearAllFilters,
		isFilterActive,
		isFilterActiveEq,
		getFilterEq,
	};
};

export { useFilter, getKey };
